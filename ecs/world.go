package ecs

import (
	"unsafe"
)

// World is the central type holding entity and component data, as well as resources.
type World struct {
	storage   storage
	resources Resources
	locks     lock
}

// NewWorld creates a new [World].
func NewWorld(initialCapacity uint32) World {
	return World{
		storage:   newStorage(initialCapacity),
		resources: newResources(),
		locks:     lock{},
	}
}

// NewEntity creates a new [Entity].
func (w *World) NewEntity() Entity {
	w.checkLocked()

	entity, _ := w.storage.createEntity(0)
	return entity
}

// Alive return whether the given entity is alive.
func (w *World) Alive(entity Entity) bool {
	return w.storage.entityPool.Alive(entity)
}

// RemoveEntity removes the given entity from the world.
func (w *World) RemoveEntity(entity Entity) {
	w.checkLocked()
	w.storage.RemoveEntity(entity)
}

// IsLocked returns whether the world is locked by any queries.
func (w *World) IsLocked() bool {
	return w.locks.IsLocked()
}

// Resources of the world.
//
// Resources are component-like data that is not associated to an entity, but unique to the world.
func (w *World) Resources() *Resources {
	return &w.resources
}

// Unsafe provides access to Ark's unsafe, ID-based API.
func (w *World) Unsafe() Unsafe {
	return Unsafe{
		world: w,
	}
}

// DumpEntities dumps entity information into an [EntityDump] object.
// This dump can be used with [World.LoadEntities] to set the World's entity state.
//
// For world serialization with components and resources, see module [github.com/mlange-42/arche-serde].
func (w *World) DumpEntities() EntityDump {
	u := w.Unsafe()
	alive := []uint32{}

	filter := NewFilter()
	query := u.Query(filter)
	for query.Next() {
		alive = append(alive, uint32(query.Entity().id))
	}

	data := EntityDump{
		Entities:  append([]Entity{}, w.storage.entityPool.entities...),
		Alive:     alive,
		Next:      uint32(w.storage.entityPool.next),
		Available: w.storage.entityPool.available,
	}

	return data
}

// LoadEntities resets all entities to the state saved with [World.DumpEntities].
//
// Use this only on an empty world! Can be used after [World.Reset].
//
// The resulting world will have the same entities (in terms of ID, generation and alive state)
// as the original world. This is necessary for proper serialization of entity relations.
// However, the entities will not have any components.
//
// Panics if the world has any dead or alive entities.
//
// For world serialization with components and resources, see module [github.com/mlange-42/ark-serde].
func (w *World) LoadEntities(data *EntityDump) {
	w.checkLocked()

	if len(w.storage.entityPool.entities) > 2 || w.storage.entityPool.available > 0 {
		panic("can set entity data only on a fresh or reset world")
	}

	capacity := len(data.Entities)

	entities := make([]Entity, 0, capacity)
	entities = append(entities, data.Entities...)

	if len(data.Entities) > 0 {
		w.storage.entityPool.entities = entities
		w.storage.entityPool.next = entityID(data.Next)
		w.storage.entityPool.available = data.Available
		w.storage.entityPool.pointer = unsafe.Pointer(&w.storage.entityPool.entities[0])
		w.storage.entityPool.reserved = entityID(reservedEntities)
	}

	w.storage.entities = make([]entityIndex, len(data.Entities), capacity)
	w.storage.isTarget = make([]bool, len(data.Entities), capacity)

	table := &w.storage.tables[0]
	for _, idx := range data.Alive {
		entity := w.storage.entityPool.entities[idx]
		tableIdx := table.Add(entity)
		w.storage.entities[entity.id] = entityIndex{table: table.id, row: tableIdx}
	}
}
