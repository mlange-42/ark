package ecs

import (
	"reflect"
	"unsafe"
)

func (w *World) newEntityWith(ids []ID, comps []unsafe.Pointer, relations []relationID) Entity {
	w.checkLocked()

	mask := All(ids...)
	newTable := w.storage.findOrCreateTable(&w.storage.tables[0], &mask, relations)
	entity, idx := w.storage.createEntity(newTable.id)

	if comps != nil {
		if len(ids) != len(comps) {
			panic("lengths of IDs and components to add do not match")
		}
		for i, id := range ids {
			newTable.Set(id, idx, comps[i])
		}
	}
	w.storage.registerTargets(relations)
	return entity
}

func (w *World) newEntitiesWith(count int, ids []ID, comps []unsafe.Pointer, relations []relationID) {
	w.checkLocked()

	mask := All(ids...)
	newTable := w.storage.findOrCreateTable(&w.storage.tables[0], &mask, relations)

	startIdx := newTable.Len()
	w.storage.createEntities(newTable, count)

	if comps != nil {
		if len(ids) != len(comps) {
			panic("lengths of IDs and components to add do not match")
		}
		for i := range count {
			for j, id := range ids {
				newTable.Set(id, uint32(startIdx+i), comps[j])
			}
		}
	}
	w.storage.registerTargets(relations)
}

func (w *World) get(entity Entity, component ID) unsafe.Pointer {
	if !w.storage.entityPool.Alive(entity) {
		panic("can't get component of a dead entity")
	}
	index := w.storage.entities[entity.id]
	return w.storage.tables[index.table].Get(component, uintptr(index.row))
}

func (w *World) has(entity Entity, component ID) bool {
	if !w.storage.entityPool.Alive(entity) {
		panic("can't get component of a dead entity")
	}
	index := w.storage.entities[entity.id]
	return w.storage.tables[index.table].Has(component)
}

func (w *World) exchange(entity Entity, add []ID, rem []ID, addComps []unsafe.Pointer, relations []relationID) {
	w.checkLocked()

	if !w.Alive(entity) {
		panic("can't exchange components on a dead entity")
	}
	if len(add) == 0 && len(rem) == 0 {
		return
	}

	index := &w.storage.entities[entity.id]
	oldTable := &w.storage.tables[index.table]
	oldArchetype := &w.storage.archetypes[oldTable.archetype]

	mask := oldArchetype.mask
	w.storage.getExchangeMask(&mask, add, rem)

	oldIDs := oldArchetype.components

	newTable := w.storage.findOrCreateTable(oldTable, &mask, relations)
	newIndex := newTable.Add(entity)

	for _, id := range oldIDs {
		if mask.Get(id) {
			comp := oldTable.Get(id, uintptr(index.row))
			newTable.Set(id, newIndex, comp)
		}
	}
	if addComps != nil {
		if len(add) != len(addComps) {
			panic("lengths of IDs and components to add do not match")
		}
		for i, id := range add {
			newTable.Set(id, newIndex, addComps[i])
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

// setRelations sets the target entities for an entity relations.
func (w *World) setRelations(entity Entity, relations []relationID) {
	w.checkLocked()

	if !w.storage.entityPool.Alive(entity) {
		panic("can't set relation for a dead entity")
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
	}
	newIndex := newTable.Add(entity)

	for _, id := range oldArch.components {
		comp := oldTable.Get(id, uintptr(index.row))
		newTable.Set(id, newIndex, comp)
	}

	swapped := oldTable.Remove(index.row)

	if swapped {
		swapEntity := oldTable.GetEntity(uintptr(index.row))
		w.storage.entities[swapEntity.id].row = index.row
	}
	w.storage.entities[entity.id] = entityIndex{table: newTable.id, row: newIndex}

	w.storage.registerTargets(relations)
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
		panic("attempt to modify a locked world")
	}
}
