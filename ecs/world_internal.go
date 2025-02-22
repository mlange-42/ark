package ecs

import (
	"reflect"
	"unsafe"
)

func (w *World) newEntityWith(ids []ID, comps []unsafe.Pointer) Entity {
	// TODO: check lock.
	mask := All(ids...)
	newTable := w.storage.findOrCreateTable(&mask)
	entity, idx := w.createEntity(newTable.id)

	if comps != nil {
		if len(ids) != len(comps) {
			panic("lengths of IDs and components to add do not match")
		}
		for i, id := range ids {
			newTable.Set(id, idx, comps[i])
		}
	}
	return entity
}

func (w *World) get(entity Entity, component ID) unsafe.Pointer {
	if !w.entityPool.Alive(entity) {
		panic("can't get component of a dead entity")
	}
	index := w.entities[entity.id]
	return w.storage.tables[index.table].Get(component, uintptr(index.row))
}

func (w *World) has(entity Entity, component ID) bool {
	if !w.entityPool.Alive(entity) {
		panic("can't get component of a dead entity")
	}
	index := w.entities[entity.id]
	return w.storage.tables[index.table].Has(component)
}

func (w *World) getEntityIndex(entity Entity) entityIndex {
	if !w.entityPool.Alive(entity) {
		panic("can't get component of a dead entity")
	}
	return w.entities[entity.id]
}

func (w *World) createEntity(table tableID) (Entity, uint32) {
	entity := w.entityPool.Get()

	idx := w.storage.tables[table].Add(entity)
	len := len(w.entities)
	if int(entity.id) == len {
		w.entities = append(w.entities, entityIndex{table: table, row: idx})
	} else {
		w.entities[entity.id] = entityIndex{table: table, row: idx}
	}
	return entity, idx
}

func (w *World) exchange(entity Entity, add []ID, rem []ID, addComps []unsafe.Pointer) {
	// TODO: check lock.
	if !w.Alive(entity) {
		panic("can't exchange components on a dead entity")
	}
	if len(add) == 0 && len(rem) == 0 {
		return
	}

	index := &w.entities[entity.id]
	oldTable := &w.storage.tables[index.table]
	oldArchetype := &w.storage.archetypes[oldTable.archetype]

	mask := oldArchetype.mask
	w.storage.getExchangeMask(&mask, add, rem)

	oldIDs := oldArchetype.components

	newTable := w.storage.findOrCreateTable(&mask)
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
		w.entities[swapEntity.id].row = index.row
	}
	w.entities[entity.id] = entityIndex{table: newTable.id, row: newIndex}
}

func (w *World) componentID(tp reflect.Type) ID {
	id, newID := w.storage.registry.ComponentID(tp)
	if newID {
		//	TODO: check lock and unroll
		//	if w.IsLocked() {
		//		w.registry.unregisterLastComponent()
		//		panic("attempt to register a new component in a locked world")
		//	}
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
