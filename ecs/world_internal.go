package ecs

import (
	"reflect"
)

// batchTable is a helper struct for collecting tables for batch processing.
type batchTable struct {
	oldTable tableID
	newTable tableID
	start    uint32
	len      uint32
}

// newEntity creates a new entity.
// Returns the entity and its bit.mask.
func (w *World) newEntity(ids []ID, relations []relationID) (Entity, *bitMask) {
	w.checkLocked()
	mask := bitMask{}
	newTable, newArch := w.storage.findOrCreateTableAdd(&w.storage.tables[0], ids, relations, &mask)
	entity, _ := w.storage.createEntity(newTable.id)
	w.storage.registerTargets(relations)

	return entity, &newArch.mask
}

// newEntities creates multiple new entities.
// Returns the table containing the entities, and their start index in the table.
func (w *World) newEntities(count int, ids []ID, relations []relationID) (tableID, int) {
	mask := bitMask{}
	newTable, _ := w.storage.findOrCreateTableAdd(&w.storage.tables[0], ids, relations, &mask)
	startIdx := newTable.Len()
	w.storage.createEntities(newTable, count)
	w.storage.registerTargets(relations)
	return newTable.id, startIdx
}

// add components to an entity.
// Returns the entity's old and new bit-mask.
func (w *World) add(entity Entity, add []ID, relations []relationID) (*bitMask, *bitMask) {
	w.checkLocked()

	if !w.Alive(entity) {
		panic("can't add components to a dead entity")
	}
	if len(add) == 0 {
		panic("at least one component required to add")
	}

	index := w.storage.entities[entity.id]
	oldTable := &w.storage.tables[index.table]
	oldArchetype := &w.storage.archetypes[oldTable.archetype]

	mask := oldArchetype.mask
	newTable, newArch := w.storage.findOrCreateTableAdd(oldTable, add, relations, &mask)
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

	return &oldArchetype.mask, &newArch.mask
}

// remove components from an entity.
func (w *World) remove(entity Entity, rem []ID) {
	w.checkLocked()

	if !w.Alive(entity) {
		panic("can't remove components from a dead entity")
	}
	if len(rem) == 0 {
		panic("at least one component required to remove")
	}

	index := w.storage.entities[entity.id]
	oldTable := &w.storage.tables[index.table]
	oldArchetype := &w.storage.archetypes[oldTable.archetype]

	mask := oldArchetype.mask
	newTable, _, relRemoved := w.storage.findOrCreateTableRemove(oldTable, rem, &mask)
	newIndex := newTable.Add(entity)

	// Get the old table and archetype again, as the pointer may have changed.
	oldTable = &w.storage.tables[oldTable.id]
	oldArchetype = &w.storage.archetypes[oldTable.archetype]

	hasCompObs := w.storage.observers.HasObservers(OnRemoveComponents)
	hasRelObs := relRemoved && w.storage.observers.HasObservers(OnRemoveRelations)
	if hasCompObs || hasRelObs {
		l := w.lock()
		if hasCompObs {
			w.storage.observers.FireRemove(OnRemoveComponents, entity, &oldArchetype.mask, &mask, true)
		}
		if hasRelObs {
			w.storage.observers.FireRemove(OnRemoveRelations, entity, &oldArchetype.mask, &mask, true)
		}
		w.unlock(l)
	}

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
}

// remove components on an entity.
// Returns the entity's old and new bit-mask.
func (w *World) exchange(entity Entity, add []ID, rem []ID, relations []relationID) (*bitMask, *bitMask) {
	w.checkLocked()

	if !w.Alive(entity) {
		panic("can't exchange components on a dead entity")
	}
	if len(add) == 0 && len(rem) == 0 {
		panic("at least one component required to add or remove")
	}

	index := w.storage.entities[entity.id]
	oldTable := &w.storage.tables[index.table]
	oldArchetype := &w.storage.archetypes[oldTable.archetype]

	mask := oldArchetype.mask
	newTable, newArch, relRemoved := w.storage.findOrCreateTable(oldTable, add, rem, relations, &mask)
	newIndex := newTable.Add(entity)

	// Get the old table and archetype again, as the pointer may have changed.
	oldTable = &w.storage.tables[oldTable.id]
	oldArchetype = &w.storage.archetypes[oldTable.archetype]

	if len(rem) > 0 {
		hasCompObs := w.storage.observers.HasObservers(OnRemoveComponents)
		hasRelObs := relRemoved && w.storage.observers.HasObservers(OnRemoveRelations)
		if hasCompObs || hasRelObs {
			l := w.lock()
			if hasCompObs {
				w.storage.observers.FireRemove(OnRemoveComponents, entity, &oldArchetype.mask, &mask, true)
			}
			if hasRelObs {
				w.storage.observers.FireRemove(OnRemoveRelations, entity, &oldArchetype.mask, &mask, true)
			}
			w.unlock(l)
		}
	}

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

	return &oldArchetype.mask, &newArch.mask
}

// exchangeBatch batch-exchanges components on entities.
//
//nolint:gocyclo
func (w *World) exchangeBatch(batch *Batch, add []ID, rem []ID,
	relations []relationID, fn func(table tableID, start, len uint32)) {
	w.checkLocked()

	if len(add) == 0 && len(rem) == 0 {
		panic("at least one component required to add or remove")
	}
	lock := w.lock()

	relRemoved := false
	tables := w.storage.getBatchTables(batch)
	batchTables := w.storage.slices.batches
	for _, tableID := range tables {
		table := &w.storage.tables[tableID]

		if table.Len() == 0 {
			continue
		}
		oldArchetype := &w.storage.archetypes[table.archetype]
		mask := oldArchetype.mask
		newTable, _, relRemovedTable := w.storage.findOrCreateTable(table, add, rem, relations, &mask)
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

	if len(rem) > 0 {
		if w.storage.observers.HasObservers(OnRemoveComponents) {
			for _, batch := range batchTables {
				table := &w.storage.tables[batch.oldTable]
				oldMask := &w.storage.archetypes[table.archetype].mask
				newMask := &w.storage.archetypes[w.storage.tables[batch.newTable].archetype].mask
				len := uintptr(batch.len)
				earlyOut := true
				for i := uintptr(0); i < len; i++ {
					if !w.storage.observers.FireRemove(OnRemoveComponents, table.GetEntity(i), oldMask, newMask, earlyOut) {
						break
					}
					earlyOut = false
				}
			}
		}
		if relRemoved && w.storage.observers.HasObservers(OnRemoveRelations) {
			for _, batch := range batchTables {
				table := &w.storage.tables[batch.oldTable]
				oldMask := &w.storage.archetypes[table.archetype].mask
				newMask := &w.storage.archetypes[w.storage.tables[batch.newTable].archetype].mask
				len := uintptr(batch.len)
				earlyOut := true
				for i := uintptr(0); i < len; i++ {
					if !w.storage.observers.FireRemove(OnRemoveRelations, table.GetEntity(i), oldMask, newMask, earlyOut) {
						break
					}
					earlyOut = false
				}
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

	if len(add) > 0 {
		if w.storage.observers.HasObservers(OnAddComponents) {
			for _, batch := range batchTables {
				table := &w.storage.tables[batch.newTable]
				oldMask := &w.storage.archetypes[w.storage.tables[batch.oldTable].archetype].mask
				newMask := &w.storage.archetypes[table.archetype].mask
				len := uintptr(batch.start + batch.len)
				earlyOut := true
				for i := uintptr(batch.start); i < len; i++ {
					if !w.storage.observers.FireAdd(OnAddComponents, table.GetEntity(i), oldMask, newMask, earlyOut) {
						break
					}
					earlyOut = false
				}
			}
		}
		if len(relations) > 0 && w.storage.observers.HasObservers(OnAddRelations) {
			for _, batch := range batchTables {
				table := &w.storage.tables[batch.newTable]
				oldMask := &w.storage.archetypes[w.storage.tables[batch.oldTable].archetype].mask
				newMask := &w.storage.archetypes[table.archetype].mask
				len := uintptr(batch.start + batch.len)
				earlyOut := true
				for i := uintptr(batch.start); i < len; i++ {
					if !w.storage.observers.FireAdd(OnAddRelations, table.GetEntity(i), oldMask, newMask, earlyOut) {
						break
					}
					earlyOut = false
				}
			}
		}
	}
	w.storage.slices.batches = batchTables[:0]
	w.unlock(lock)
}

// exchangeTable performs batch-exchange on a single table.
// Returns the start index of the entities in the new table and number of entities.
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
			oldCol := oldTable.Column(id)
			newCol := newTable.Column(id)
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
		w.storage.observers.FireSetRelations(OnRemoveRelations, entity, &changeMask, newMask, true)
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
		w.storage.observers.FireSetRelations(OnAddRelations, entity, &changeMask, newMask, true)
	}
}

// setRelationsBatch batch-changes entity relations.
func (w *World) setRelationsBatch(batch *Batch, relations []relationID, fn func(table tableID, start, len int)) {
	w.checkLocked()

	if len(relations) == 0 {
		panic("no relations specified")
	}
	lock := w.lock()
	hasObserver := w.storage.observers.HasObservers(OnAddRelations) || w.storage.observers.HasObservers(OnRemoveRelations)

	tables := w.storage.getBatchTables(batch)
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

// setRelationsTable batch-changes entity relations for a single table.
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
			if !w.storage.observers.FireSetRelations(OnRemoveRelations, oldTable.GetEntity(i), &changeMask, newMask, earlyOut) {
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
			if !w.storage.observers.FireSetRelations(OnAddRelations, oldTable.GetEntity(index), &changeMask, newMask, earlyOut) {
				break
			}
			earlyOut = false
		}
	}
}

// componentID returns the component ID for a runtime component type.
// Registers the type if necessary, and adds it to the storage.
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

// resourceID returns the resource ID for a runtime resource type.
// Registers the resource of necessary.
func (w *World) resourceID(tp reflect.Type) ResID {
	id, _ := w.resources.registry.ComponentID(tp)
	return ResID{id: id}
}

// registerObserver adds an observer to the [observerManager].
func (w *World) registerObserver(obs *Observer) {
	w.storage.observers.AddObserver(obs, w)
}

// unregisterObserver removes an observer from the [observerManager].
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

// checkLocked checks if the world is locked, and panics if so.
func (w *World) checkLocked() {
	if w.IsLocked() {
		panic("cannot modify a locked world: collect entities into a slice and apply changes after query iteration has completed")
	}
}

// emitEvent distributes an event to the [observerManager].
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
