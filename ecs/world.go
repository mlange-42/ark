package ecs

import (
	"unsafe"
)

// World is the central type holding entity and component data, as well as resources.
type World struct {
	storage    storage
	entities   entities
	entityPool entityPool
}

// NewWorld creates a new [World].
func NewWorld(initialCapacity uint32) World {
	entities := make([]entityIndex, 1, initialCapacity)
	entities[0] = entityIndex{table: 0, row: 0}
	return World{
		storage:    newStorage(initialCapacity),
		entities:   entities,
		entityPool: newEntityPool(initialCapacity),
	}
}

// NewEntity creates a new [Entity].
func (w *World) NewEntity() Entity {
	return w.createEntity(0)
}

// Alive return whether the given entity is alive.
func (w *World) Alive(entity Entity) bool {
	return w.entityPool.Alive(entity)
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

func (w *World) createEntity(table tableID) Entity {
	entity := w.entityPool.Get()

	idx := w.storage.tables[table].Add(entity)
	len := len(w.entities)
	if int(entity.id) == len {
		w.entities = append(w.entities, entityIndex{table: table, row: idx})
	} else {
		w.entities[entity.id] = entityIndex{table: table, row: idx}
	}
	return entity
}

func (w *World) exchange(entity Entity, add []ID, rem []ID, addComps []unsafe.Pointer) {
	// TODO: check lock!
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
