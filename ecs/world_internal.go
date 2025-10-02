package ecs

import (
	"reflect"
)

func (w *World) newEntity(ids []ID, relations []relationID) (Entity, *bitMask) {
	w.checkLocked()
	mask := bitMask{}
	newTable, _ := w.storage.findOrCreateTable(&w.storage.tables[0], ids, nil, relations, &mask)
	entity, _ := w.storage.createEntity(newTable.id)
	w.storage.registerTargets(relations)

	return entity, &w.storage.archetypes[newTable.archetype].mask
}

func (w *World) newEntities(count int, ids []ID, relations []relationID) (tableID, int) {
	mask := bitMask{}
	newTable, _ := w.storage.findOrCreateTable(&w.storage.tables[0], ids, nil, relations, &mask)
	startIdx := newTable.Len()
	w.storage.createEntities(newTable, count)
	w.storage.registerTargets(relations)
	return newTable.id, startIdx
}

func (w *World) exchange(entity Entity, add []ID, rem []ID, relations []relationID) (*bitMask, *bitMask) {
	w.checkLocked()

	if !w.Alive(entity) {
		panic("can't exchange components on a dead entity")
	}
	if len(add) == 0 && len(rem) == 0 {
		if len(relations) > 0 {
			panic("exchange operation has no effect, but relations were specified. Use SetRelation(s) instead")
		}
		return nil, nil
	}

	index := w.storage.entities[entity.id]
	oldTable := &w.storage.tables[index.table]
	oldArchetype := &w.storage.archetypes[oldTable.archetype]

	mask := oldArchetype.mask
	newTable, relRemoved := w.storage.findOrCreateTable(oldTable, add, rem, relations, &mask)
	newIndex := newTable.Add(entity)

	if len(rem) > 0 {
		hasCompObs := w.storage.observers.HasObservers(OnRemoveComponents)
		hasRelObs := relRemoved && w.storage.observers.HasObservers(OnRemoveRelations)
		if hasCompObs || hasRelObs {
			l := w.lock()
			if hasCompObs {
				w.storage.observers.doFireRemove(OnRemoveComponents, entity, &oldArchetype.mask, &mask, true)
			}
			if hasRelObs {
				w.storage.observers.doFireRemove(OnRemoveRelations, entity, &oldArchetype.mask, &mask, true)
			}
			w.unlock(l)
		}
	}

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

	return &oldArchetype.mask, &w.storage.archetypes[newTable.archetype].mask
}

type batchTable struct {
	oldTable tableID
	newTable tableID
	start    uint32
	len      uint32
}

//nolint:gocyclo
func (w *World) exchangeBatch(batch *Batch, add []ID, rem []ID,
	relations []relationID, fn func(table tableID, start, len uint32)) {
	w.checkLocked()

	if len(add) == 0 && len(rem) == 0 {
		if len(relations) > 0 {
			panic("exchange operation has no effect, but relations were specified. Use SetRelationBatch instead")
		}
		return
	}
	lock := w.lock()

	relRemoved := false
	tables := w.storage.getTables(batch)
	batchTables := w.storage.slices.batches
	for _, tableID := range tables {
		table := &w.storage.tables[tableID]

		if table.Len() == 0 {
			continue
		}
		oldArchetype := &w.storage.archetypes[table.archetype]
		mask := oldArchetype.mask
		newTable, relRemovedTable := w.storage.findOrCreateTable(table, add, rem, relations, &mask)
		if relRemovedTable {
			relRemoved = true
		}
		batchTables = append(batchTables, batchTable{
			oldTable: table.id,
			newTable: newTable.id,
			len:      uint32(table.Len()),
		})
	}
	w.storage.slices.tables = tables[:0]

	if len(rem) > 0 && w.storage.observers.HasObservers(OnRemoveComponents) {
		for _, batch := range batchTables {
			table := &w.storage.tables[batch.oldTable]
			oldMask := &w.storage.archetypes[table.archetype].mask
			newMask := &w.storage.archetypes[w.storage.tables[batch.newTable].archetype].mask
			len := uintptr(batch.len)
			earlyOut := true
			for i := uintptr(0); i < len; i++ {
				if !w.storage.observers.doFireRemove(OnRemoveComponents, table.GetEntity(i), oldMask, newMask, earlyOut) {
					break
				}
				earlyOut = false
			}
		}
	}
	if len(rem) > 0 && relRemoved && w.storage.observers.HasObservers(OnRemoveRelations) {
		for _, batch := range batchTables {
			table := &w.storage.tables[batch.oldTable]
			oldMask := &w.storage.archetypes[table.archetype].mask
			newMask := &w.storage.archetypes[w.storage.tables[batch.newTable].archetype].mask
			len := uintptr(batch.len)
			earlyOut := true
			for i := uintptr(0); i < len; i++ {
				if !w.storage.observers.doFireRemove(OnRemoveRelations, table.GetEntity(i), oldMask, newMask, earlyOut) {
					break
				}
				earlyOut = false
			}
		}
	}

	for i := range batchTables {
		batch := &batchTables[i]

		start, len := w.exchangeTable(batch.oldTable, batch.newTable, relations)
		if fn != nil {
			fn(batch.newTable, start, len)
		}
		batch.start = start
		batch.len = len
	}

	if len(add) > 0 && w.storage.observers.HasObservers(OnAddComponents) {
		for _, batch := range batchTables {
			table := &w.storage.tables[batch.newTable]
			oldMask := &w.storage.archetypes[w.storage.tables[batch.oldTable].archetype].mask
			newMask := &w.storage.archetypes[table.archetype].mask
			len := uintptr(batch.start + batch.len)
			earlyOut := true
			for i := uintptr(batch.start); i < len; i++ {
				if !w.storage.observers.doFireAdd(OnAddComponents, table.GetEntity(i), oldMask, newMask, earlyOut) {
					break
				}
				earlyOut = false
			}
		}
	}
	if len(add) > 0 && len(relations) > 0 && w.storage.observers.HasObservers(OnAddRelations) {
		for _, batch := range batchTables {
			table := &w.storage.tables[batch.newTable]
			oldMask := &w.storage.archetypes[w.storage.tables[batch.oldTable].archetype].mask
			newMask := &w.storage.archetypes[table.archetype].mask
			len := uintptr(batch.start + batch.len)
			earlyOut := true
			for i := uintptr(batch.start); i < len; i++ {
				if !w.storage.observers.doFireAdd(OnAddRelations, table.GetEntity(i), oldMask, newMask, earlyOut) {
					break
				}
				earlyOut = false
			}
		}
	}
	w.storage.slices.batches = batchTables[:0]
	w.unlock(lock)
}

func (w *World) exchangeTable(oldTableID, newTableID tableID, relations []relationID) (uint32, uint32) {
	oldTable := &w.storage.tables[oldTableID]

	oldArchetype := &w.storage.archetypes[oldTable.archetype]
	oldIDs := oldArchetype.components

	newTable := &w.storage.tables[newTableID]
	newArchetype := &w.storage.archetypes[newTable.archetype]

	mask := &newArchetype.mask

	startIdx := uint32(newTable.Len())
	count := oldTable.len

	var i uint32
	for i = range count {
		idx := startIdx + i
		entity := oldTable.GetEntity(uintptr(i))
		index := &w.storage.entities[entity.id]
		index.table = newTable.id
		index.row = idx
	}

	newTable.AddAllEntities(oldTable, count)
	for _, id := range oldIDs {
		if mask.Get(id.id) {
			oldCol := oldTable.GetColumn(id)
			newCol := newTable.GetColumn(id)
			newCol.CopyToEnd(oldCol, newTable.len, count)
		}
	}

	oldTable.Reset()
	w.storage.registerTargets(relations)

	return startIdx, count
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
	hasObserver := w.storage.observers.HasObservers(OnAddRelations) || w.storage.observers.HasObservers(OnRemoveRelations)

	index := &w.storage.entities[entity.id]
	oldTable := &w.storage.tables[index.table]

	var changeMask bitMask
	var maskPointer *bitMask
	if hasObserver {
		maskPointer = &changeMask
	}
	newRelations, changed := w.storage.getExchangeTargets(oldTable, relations, maskPointer)
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

	if w.storage.observers.HasObservers(OnRemoveRelations) {
		lock := w.lock()
		newMask := &w.storage.archetypes[newTable.archetype].mask
		w.storage.observers.doFireSetRelations(OnRemoveRelations, entity, &changeMask, newMask, true)
		w.unlock(lock)
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

	if w.storage.observers.HasObservers(OnAddRelations) {
		newMask := &w.storage.archetypes[newTable.archetype].mask
		w.storage.observers.doFireSetRelations(OnAddRelations, entity, &changeMask, newMask, true)
	}
}

func (w *World) setRelationsBatch(batch *Batch, relations []relationID, fn func(table tableID, start, len int)) {
	w.checkLocked()

	if len(relations) == 0 {
		panic("no relations specified")
	}
	lock := w.lock()
	hasObserver := w.storage.observers.HasObservers(OnAddRelations) || w.storage.observers.HasObservers(OnRemoveRelations)

	tables := w.storage.getTables(batch)
	lengths := w.storage.slices.ints
	var totalEntities uint32 = 0
	for _, tableID := range tables {
		table := &w.storage.tables[tableID]
		lengths = append(lengths, uint32(table.Len()))
		totalEntities += uint32(table.Len())
	}

	for i, tableID := range tables {
		tableLen := lengths[i]
		if tableLen == 0 {
			continue
		}
		table := &w.storage.tables[tableID]
		w.setRelationsTable(table, int(tableLen), relations, fn, hasObserver)
	}

	w.storage.slices.ints = lengths[:0]
	w.storage.slices.tables = tables[:0]

	w.storage.registerTargets(relations)

	w.unlock(lock)
}

func (w *World) setRelationsTable(oldTable *table, oldLen int, relations []relationID, fn func(table tableID, start, len int), hasObserver bool) {
	var changeMask bitMask
	var maskPointer *bitMask
	if hasObserver {
		maskPointer = &changeMask
	}
	newRelations, changed := w.storage.getExchangeTargets(oldTable, relations, maskPointer)

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

	// TODO: move this before the entire batch?
	if w.storage.observers.HasObservers(OnRemoveRelations) {
		newMask := &w.storage.archetypes[newTable.archetype].mask
		len := uintptr(oldTable.len)
		earlyOut := true
		for i := uintptr(0); i < len; i++ {
			if !w.storage.observers.doFireSetRelations(OnRemoveRelations, oldTable.GetEntity(i), &changeMask, newMask, earlyOut) {
				break
			}
			earlyOut = false
		}
	}

	startIdx := newTable.Len()
	w.storage.moveEntities(oldTable, newTable, uint32(oldLen))

	if fn != nil {
		fn(newTable.id, startIdx, oldLen)
	}

	// TODO: move this after the entire batch?
	if w.storage.observers.HasObservers(OnAddRelations) {
		newMask := &w.storage.archetypes[newTable.archetype].mask
		earlyOut := true
		for i := range oldLen {
			index := uintptr(startIdx + i)
			if !w.storage.observers.doFireSetRelations(OnAddRelations, oldTable.GetEntity(index), &changeMask, newMask, earlyOut) {
				break
			}
			earlyOut = false
		}
	}
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
	w.storage.observers.AddObserver(obs, w)
}

func (w *World) unregisterObserver(obs *Observer) {
	w.storage.observers.RemoveObserver(obs)
}

// lock the world and get the lock bit for later unlocking.
func (w *World) lock() uint8 {
	return w.storage.locks.Lock()
}

// unlock unlocks the given lock bit.
func (w *World) unlock(l uint8) {
	w.storage.locks.Unlock(l)
}

// lockSafe locks the world and get the lock bit for later unlocking.
func (w *World) lockSafe() uint8 {
	return w.storage.locks.LockSafe()
}

// unlockSafe unlocks the given lock bit.
func (w *World) unlockSafe(l uint8) {
	w.storage.locks.UnlockSafe(l)
}

// checkLocked checks if the world is locked, and panics if so.
func (w *World) checkLocked() {
	if w.IsLocked() {
		panic("cannot modify a locked world: collect entities into a slice and apply changes after query iteration has completed")
	}
}

func (w *World) emitEvent(e *Event, entity Entity) {
	var mask *bitMask
	if entity.IsZero() {
		mask = &w.storage.archetypes[0].mask
	} else {
		if !w.Alive(entity) {
			panic("can't emit an event for a dead entity")
		}
		table := w.storage.entities[entity.id].table
		mask = &w.storage.archetypes[w.storage.tables[table].archetype].mask
	}
	if !mask.Contains(&e.mask) {
		panic("entity does not have the required event components")
	}
	w.storage.observers.FireCustom(e.eventType, entity, &e.mask, mask)
}
