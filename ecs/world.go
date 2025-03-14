package ecs

import (
	"reflect"

	"github.com/mlange-42/ark/ecs/stats"
)

// World is the central type holding entity and component data, as well as resources.
type World struct {
	storage   storage
	resources Resources
	locks     lock
	stats     *stats.World
}

// NewWorld creates a new [World].
//
// Accepts zero, one or two arguments.
// The first argument is the initial capacity of the world, and of normal archetypes.
// The second argument is the initial capacity of relation archetypes.
// If only one argument is provided, it is used for both capacities.
// If no arguments are provided, the defaults are 1024 and 128, respectively.
func NewWorld(initialCapacity ...int) World {
	return World{
		storage:   newStorage(initialCapacity...),
		resources: newResources(),
		locks:     newLock(),
		stats:     &stats.World{},
	}
}

// NewEntity creates a new [Entity] without any components.
func (w *World) NewEntity() Entity {
	w.checkLocked()

	entity, _ := w.storage.createEntity(0)
	return entity
}

// NewEntities creates a batch of new entities without any components, running the given callback function on each.
// The callback function can be nil.
func (w *World) NewEntities(count int, fn func(entity Entity)) {
	tableID, start := w.newEntities(count, nil, nil)
	if fn == nil {
		return
	}
	table := &w.storage.tables[tableID]
	lock := w.lock()
	for i := range count {
		index := uintptr(start + i)
		fn(
			table.GetEntity(index),
		)
	}
	w.unlock(lock)
}

// Alive return whether the given entity is alive.
//
// In Ark, entities are returned to a pool when they are removed from the world.
// These entities can be recycled, with the same ID ([Entity.ID]), but an incremented generation ([Entity.Gen]).
// This allows to determine whether an entity held by the user is still alive, despite it was potentially recycled.
func (w *World) Alive(entity Entity) bool {
	return w.storage.entityPool.Alive(entity)
}

// RemoveEntity removes the given entity from the world.
func (w *World) RemoveEntity(entity Entity) {
	w.checkLocked()
	w.storage.RemoveEntity(entity)
}

// RemoveEntities removes all entities matching the given batch filter,
// running the given function on each. The function can be nil.
func (w *World) RemoveEntities(batch *Batch, fn func(entity Entity)) {
	w.checkLocked()

	tables := w.storage.getTables(batch)
	cleanup := []Entity{}
	for _, table := range tables {
		len := uintptr(table.Len())
		var i uintptr
		if fn != nil {
			l := w.lock()
			for i = range len {
				fn(table.GetEntity(i))
			}
			w.unlock(l)
		}
		for i = range len {
			entity := table.GetEntity(i)
			if w.storage.isTarget[entity.id] {
				cleanup = append(cleanup, entity)
			}
			w.storage.entities[entity.id].table = maxTableID
			w.storage.entityPool.Recycle(entity)
		}
		table.Reset()
	}

	for _, entity := range cleanup {
		w.storage.cleanupArchetypes(entity)
		w.storage.isTarget[entity.id] = false
	}
}

// IsLocked returns whether the world is locked by any queries.
func (w *World) IsLocked() bool {
	return w.locks.IsLocked()
}

// Resources of the world.
// Resources are component-like data that is not associated to an entity, but unique to the world.
//
// See also [Resource], [AddResource] and [GetResource].
func (w *World) Resources() *Resources {
	return &w.resources
}

// Unsafe provides access to Ark's unsafe, ID-based API.
// For details, see [Unsafe].
func (w *World) Unsafe() Unsafe {
	return Unsafe{
		world: w,
	}
}

// Reset removes all entities and resources from the world, and clears the filter cache.
//
// Does NOT free reserved memory, remove archetypes, or clear the registry.
// However, it removes archetypes with a relation component.
//
// Can be used to run systematic simulations without the need to re-allocate memory for each run.
// Accelerates re-populating the world by a factor of 2-3.
func (w *World) Reset() {
	w.checkLocked()

	w.storage.Reset()
	w.locks.Reset()
	w.resources.reset()
}

// Stats reports statistics for inspecting the World.
//
// The underlying [stats.World] object is re-used and updated between calls.
// The returned pointer should thus not be stored for later analysis.
// Rather, the required data should be extracted immediately.
func (w *World) Stats() *stats.World {
	w.stats.Entities = stats.Entities{
		Used:     w.storage.entityPool.Len(),
		Total:    w.storage.entityPool.Cap(),
		Recycled: w.storage.entityPool.Available(),
		Capacity: w.storage.entityPool.TotalCap(),
	}

	compCount := len(w.storage.registry.Components)
	types := append([]reflect.Type{}, w.storage.registry.Types[:compCount]...)

	memory := cap(w.storage.entities)*int(entityIndexSize) + w.storage.entityPool.TotalCap()*int(entitySize)
	memoryUsed := w.storage.entityPool.Len() * int(entityIndexSize+entitySize)

	cntOld := int32(len(w.stats.Archetypes))
	cntNew := int32(len(w.storage.archetypes))
	var i int32
	for i = 0; i < cntOld; i++ {
		arch := &w.storage.archetypes[i]
		archStats := &w.stats.Archetypes[i]
		arch.UpdateStats(archStats, &w.storage)
		memory += archStats.Memory
		memoryUsed += archStats.MemoryUsed
	}
	for i = cntOld; i < cntNew; i++ {
		arch := &w.storage.archetypes[i]
		w.stats.Archetypes = append(w.stats.Archetypes, arch.Stats(&w.storage))
		archStats := &w.stats.Archetypes[i]
		memory += archStats.Memory
		memoryUsed += archStats.MemoryUsed
	}

	w.stats.ComponentTypes = types
	w.stats.Locked = w.IsLocked()
	w.stats.Memory = memory
	w.stats.MemoryUsed = memoryUsed
	w.stats.CachedFilters = len(w.storage.cache.filters)

	return w.stats
}
