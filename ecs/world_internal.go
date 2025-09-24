package ecs

import (
	"reflect"
)

func (w *World) newEntity(ids []ID, relations []relationID) Entity {
	w.checkLocked()
	mask := bitMask{}
	newTable := w.storage.findOrCreateTable(&w.storage.tables[0], ids, nil, relations, &mask)
	entity, _ := w.storage.createEntity(newTable.id)
	w.storage.registerTargets(relations)
	return entity
}

func (w *World) newEntities(count int, ids []ID, relations []relationID) (tableID, int) {
	w.checkLocked()
	mask := bitMask{}
	newTable := w.storage.findOrCreateTable(&w.storage.tables[0], ids, nil, relations, &mask)
	startIdx := newTable.Len()
	w.storage.createEntities(newTable, count)
	w.storage.registerTargets(relations)
	return newTable.id, startIdx
}

func (w *World) exchange(entity Entity, add []ID, rem []ID, relations []relationID) {
	w.checkLocked()

	if !w.Alive(entity) {
		panic("can't exchange components on a dead entity")
	}
	if len(add) == 0 && len(rem) == 0 {
		if len(relations) > 0 {
			panic("exchange operation has no effect, but relations were specified. Use SetRelation(s) instead")
		}
		return
	}

	index := w.storage.entities[entity.id]
	oldTable := &w.storage.tables[index.table]
	oldArchetype := &w.storage.archetypes[oldTable.archetype]

	mask := oldArchetype.mask
	newTable := w.storage.findOrCreateTable(oldTable, add, rem, relations, &mask)
	newIndex := newTable.Add(entity)

	// Get the old table and archetype again, as the pointer may have changed.
	oldTable = &w.storage.tables[oldTable.id]
	oldArchetype = &w.storage.archetypes[oldTable.archetype]

	for _, id := range oldArchetype.components {
		if mask.Get(id.id) {
			newTable.Set(id, newIndex, oldTable.Column(id), int(index.row))
		}
	}

	swapped := oldTable.Remove(index.row)

	if swapped {
		swapEntity := oldTable.GetEntity(uintptr(index.row))
		w.storage.entities[swapEntity.id].row = index.row
	}
	w.storage.entities[entity.id] = entityIndex{table: newTable.id, row: newIndex}

	w.storage.registerTargets(relations)
}

func (w *World) exchangeBatch(batch *Batch, add []ID, rem []ID,
	relations []relationID, fn func(table tableID, start, len int)) {
	w.checkLocked()

	if len(add) == 0 && len(rem) == 0 {
		if len(relations) > 0 {
			panic("exchange operation has no effect, but relations were specified. Use SetRelationBatch instead")
		}
		return
	}

	tables := w.storage.getTables(batch)
	lengths := make([]uint32, len(tables))
	var totalEntities uint32 = 0
	for i, tableID := range tables {
		table := &w.storage.tables[tableID]
		lengths[i] = uint32(table.Len())
		totalEntities += uint32(table.Len())
	}

	for i, tableID := range tables {
		tableLen := lengths[i]

		if tableLen == 0 {
			continue
		}
		table := &w.storage.tables[tableID]
		t, start, len := w.exchangeTable(table, int(tableLen), add, rem, relations)
		if fn != nil {
			fn(t, start, len)
		}
	}
}

func (w *World) exchangeTable(oldTable *table, oldLen int, add []ID, rem []ID, relations []relationID) (tableID, int, int) {
	oldArchetype := &w.storage.archetypes[oldTable.archetype]

	oldIDs := oldArchetype.components

	mask := oldArchetype.mask
	newTable := w.storage.findOrCreateTable(oldTable, add, rem, relations, &mask)
	// Get the old table again, as pointers may have changed.
	oldTable = &w.storage.tables[oldTable.id]

	startIdx := uintptr(newTable.Len())
	count := uintptr(oldLen)

	var i uintptr
	for i = range count {
		idx := startIdx + i
		entity := oldTable.GetEntity(i)
		index := &w.storage.entities[entity.id]
		index.table = newTable.id
		index.row = uint32(idx)
	}

	newTable.AddAllEntities(oldTable, uint32(oldLen))
	for _, id := range oldIDs {
		if mask.Get(id.id) {
			oldCol := oldTable.GetColumn(id)
			newCol := newTable.GetColumn(id)
			newCol.CopyToEnd(oldCol, newTable.len, uint32(oldLen))
		}
	}

	oldTable.Reset()
	w.storage.registerTargets(relations)

	return newTable.id, int(startIdx), int(count)
}

// setRelations sets the target entities for an entity relations.
func (w *World) setRelations(entity Entity, relations []relationID) {
	w.checkLocked()

	if !w.storage.entityPool.Alive(entity) {
		panic("can't set relation for a dead entity")
	}
	if len(relations) == 0 {
		panic("no relations specified")
	}

	index := &w.storage.entities[entity.id]
	oldTable := &w.storage.tables[index.table]

	newRelations, changed := w.storage.getExchangeTargets(oldTable, relations)
	if !changed {
		return
	}

	oldArch := &w.storage.archetypes[oldTable.archetype]
	newTable, ok := oldArch.GetTable(&w.storage, newRelations)
	if !ok {
		newTable = w.storage.createTable(oldArch, newRelations)
		// Get the old table again, as pointers may have changed.
		oldTable = &w.storage.tables[oldTable.id]
	}
	newIndex := newTable.Add(entity)

	for _, id := range oldArch.components {
		newTable.Set(id, newIndex, oldTable.Column(id), int(index.row))
	}

	swapped := oldTable.Remove(index.row)

	if swapped {
		swapEntity := oldTable.GetEntity(uintptr(index.row))
		w.storage.entities[swapEntity.id].row = index.row
	}
	w.storage.entities[entity.id] = entityIndex{table: newTable.id, row: newIndex}

	w.storage.registerTargets(relations)
}

func (w *World) setRelationsBatch(batch *Batch, relations []relationID, fn func(table tableID, start, len int)) {
	w.checkLocked()

	if len(relations) == 0 {
		panic("no relations specified")
	}

	tables := w.storage.getTables(batch)
	lengths := make([]uint32, len(tables))
	var totalEntities uint32 = 0
	for i, tableID := range tables {
		table := &w.storage.tables[tableID]
		lengths[i] = uint32(table.Len())
		totalEntities += uint32(table.Len())
	}

	for i, tableID := range tables {
		tableLen := lengths[i]

		if tableLen == 0 {
			continue
		}
		table := &w.storage.tables[tableID]
		t, start, len := w.setRelationsTable(table, int(tableLen), relations)
		if fn != nil {
			fn(t, start, len)
		}
	}

	w.storage.registerTargets(relations)
}

func (w *World) setRelationsTable(oldTable *table, oldLen int, relations []relationID) (tableID, int, int) {
	newRelations, changed := w.storage.getExchangeTargets(oldTable, relations)

	if !changed {
		return oldTable.id, 0, oldLen
	}
	oldArch := &w.storage.archetypes[oldTable.archetype]
	newTable, ok := oldArch.GetTable(&w.storage, newRelations)
	if !ok {
		newTable = w.storage.createTable(oldArch, newRelations)
		// Get the old table again, as pointers may have changed.
		oldTable = &w.storage.tables[oldTable.id]
	}
	startIdx := newTable.Len()
	w.storage.moveEntities(oldTable, newTable, uint32(oldLen))

	return newTable.id, startIdx, oldLen

}

func (w *World) componentID(tp reflect.Type) ID {
	id, newID := w.storage.registry.ComponentID(tp)
	if newID {
		if w.IsLocked() {
			w.storage.registry.unregisterLastComponent()
			panic("attempt to register a new component in a locked world")
		}
		w.storage.AddComponent(id)
	}
	return ID{id: id}
}

func (w *World) resourceID(tp reflect.Type) ResID {
	id, _ := w.resources.registry.ComponentID(tp)
	return ResID{id: id}
}

func (w *World) registerObserver(obs *Observer) {
	w.storage.observers.AddObserver(obs, &w.storage.registry)
}

// lock the world and get the lock bit for later unlocking.
func (w *World) lock() uint8 {
	return w.locks.Lock()
}

// unlock unlocks the given lock bit.
func (w *World) unlock(l uint8) {
	w.locks.Unlock(l)
}

// checkLocked checks if the world is locked, and panics if so.
func (w *World) checkLocked() {
	if w.IsLocked() {
		panic("cannot modify a locked world: collect entities into a slice and apply changes after query iteration has completed")
	}
}
