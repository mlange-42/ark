package ecs

import (
	"reflect"
	"time"

	"github.com/mlange-42/ark/ecs/stats"
)

// World is the central type holding entity and component data, as well as resources.
type World struct {
	stats     *stats.World // World statistics, for re-use
	resources Resources    // Registered resources
	storage   storage      // The world's storage
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
		storage:   newStorage(16, initialCapacity...),
		resources: newResources(),
		stats:     &stats.World{},
	}
}

// NewEntity creates a new [Entity] without any components.
func (w *World) NewEntity() Entity {
	w.checkLocked()

	entity, _ := w.storage.createEntity(0)
	w.storage.observers.FireCreateEntityIfHas(entity, &w.storage.archetypes[0].mask)
	return entity
}

// NewEntities creates a batch of new entities without any components, running the given callback function on each.
// The callback function can be nil.
func (w *World) NewEntities(count int, fn func(entity Entity)) {
	w.checkLocked()
	tableID, start := w.newEntities(count, nil, nil)

	hasObs := w.storage.observers.HasObservers(OnCreateEntity)
	shouldLock := hasObs || fn != nil
	var lock uint8
	if shouldLock {
		lock = w.lock()
	}

	if fn != nil {
		table := &w.storage.tables[tableID]
		for i := range count {
			index := uintptr(start + i)
			fn(
				table.GetEntity(index),
			)
		}
	}

	if hasObs {
		table := &w.storage.tables[tableID]
		mask := &w.storage.archetypes[table.archetype].mask
		earlyOut := true
		for i := range count {
			index := uintptr(start + i)
			if !w.storage.observers.FireCreateEntity(table.GetEntity(index), mask, earlyOut) {
				break
			}
			earlyOut = false
		}
	}
	if shouldLock {
		w.unlock(lock)
	}
}

// CopyEntity copies an entity with all its components.
//
// Creates a new entity and copies the memory of all components.
// Note that pointer-like fields in components (incl. slices and maps)
// are copied shallow. I.e. they will point to the same address as the original.
func (w *World) CopyEntity(e Entity) Entity {
	w.checkLocked()

	s := &w.storage
	entity := s.entityPool.Get()

	index := s.entities[e.id]
	table := &s.tables[index.table]

	idx := table.Add(&w.storage, entity)
	if int(entity.id) == len(s.entities) {
		s.entities = append(s.entities, entityIndex{table: index.table, row: idx})
		s.isTarget = append(s.isTarget, false)
	} else {
		s.entities[entity.id] = entityIndex{table: index.table, row: idx}
	}

	archetype := &s.archetypes[table.archetype]

	table.CopyAll(table, idx, index.row)

	w.storage.observers.FireCreateEntityIfHas(entity, &archetype.mask)
	return entity
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
func (w *World) RemoveEntities(batch Batch, fn func(entity Entity)) {
	w.checkLocked()

	hasEntityObs := w.storage.observers.HasObservers(OnRemoveEntity)
	hasRelationObs := w.storage.observers.HasObservers(OnRemoveRelations)
	shouldLock := hasEntityObs || hasRelationObs || fn != nil
	var lock uint8
	if shouldLock {
		lock = w.lock()
	}

	tables := w.storage.getBatchTables(&batch)

	if fn != nil {
		for _, tableID := range tables {
			table := &w.storage.tables[tableID]
			len := uintptr(table.Len())
			var i uintptr
			for i = range len {
				fn(table.GetEntity(i))
			}
		}
	}

	if hasEntityObs || hasRelationObs {
		if hasEntityObs {
			for _, tableID := range tables {
				table := &w.storage.tables[tableID]
				mask := &w.storage.archetypes[table.archetype].mask
				len := uintptr(table.Len())
				var i uintptr
				earlyOut := true
				for i = range len {
					if !w.storage.observers.FireRemoveEntity(table.GetEntity(i), mask, earlyOut) {
						break
					}
					earlyOut = false
				}
			}
		}
		if hasRelationObs {
			for _, tableID := range tables {
				table := &w.storage.tables[tableID]
				if !table.HasRelations() {
					continue
				}
				mask := &w.storage.archetypes[table.archetype].mask
				len := uintptr(table.Len())
				var i uintptr
				earlyOut := true
				for i = range len {
					if !w.storage.observers.FireRemoveEntityRel(table.GetEntity(i), mask, earlyOut) {
						break
					}
					earlyOut = false
				}
			}
		}
	}

	cleanup := w.storage.slices.entitiesCleanup
	for _, tableID := range tables {
		table := &w.storage.tables[tableID]
		len := uintptr(table.Len())
		var i uintptr
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

	w.storage.slices.tables = tables[:0]

	for _, entity := range cleanup {
		w.storage.cleanupArchetypes(entity)
		w.storage.isTarget[entity.id] = false
	}
	w.storage.slices.entitiesCleanup = cleanup[:0]

	if shouldLock {
		w.unlock(lock)
	}
}

// IsLocked returns whether the world is locked by any queries.
func (w *World) IsLocked() bool {
	return w.storage.locks.IsLocked()
}

// Resources of the world.
// Resources are component-like data that is not associated to an entity, but unique to the world.
//
// This is only required for unsafe/id-based usage of resources.
// For the safe API, see [Resource], [AddResource] and [GetResource].
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

// Event creates a new event of the given type.
//
// The event can be further configured using [Event.For].
// It must be emitted using [Event.Emit] to have an effect.
//
// See [Event] and [Observer] for details.
func (w *World) Event(tp EventType) Event {
	if tp > customEvent {
		panic("only custom events can be emitted manually")
	}
	return Event{
		world:     w,
		eventType: tp,
	}
}

// Reset removes all entities and resources from the world,
// and un-registers all cached filters and observers.
//
// Does NOT free reserved memory, remove archetypes, or clear the registry.
// However, it removes archetypes with a relation component.
//
// Can be used to run systematic simulations without the need to re-allocate memory for each run.
// Accelerates re-populating the world by a factor of 2-3.
func (w *World) Reset() {
	w.checkLocked()

	w.storage.Reset()
	w.resources.reset()
}

// Stats reports statistics for inspecting the World.
//
// The underlying [stats.World] object is re-used and updated between calls.
// The returned pointer should thus not be stored for later analysis.
// Rather, the required data should be extracted immediately.
func (w *World) Stats(flags ...stats.Option) *stats.World {
	var mask bitMask64
	if len(flags) == 0 {
		mask = newMaskFlags(stats.Archetypes, stats.Tables, stats.Filters, stats.Observers)
	} else {
		mask = newMaskFlags(flags...)
	}

	w.stats.Entities = stats.Entities{
		Used:     w.storage.entityPool.Len(),
		Total:    w.storage.entityPool.Cap(),
		Recycled: w.storage.entityPool.Available(),
		Capacity: w.storage.entityPool.TotalCap(),
	}
	prevCount := len(w.stats.ComponentTypes)
	compCount := len(w.storage.registry.Components)
	if compCount != prevCount {
		types := append([]reflect.Type{}, w.storage.registry.Types[:compCount]...)
		typeNames := make([]string, len(types))
		for i, t := range types {
			typeNames[i] = t.Name()
		}
		w.stats.ComponentTypes = types
		w.stats.ComponentTypeNames = typeNames
	}

	memory := cap(w.storage.entities)*int(entityIndexSize) + w.storage.entityPool.TotalCap()*int(entitySize)
	memoryUsed := w.storage.entityPool.Len() * int(entityIndexSize+entitySize)

	includeTables := mask.Get(uint8(stats.Tables))
	if mask.Get(uint8(stats.Archetypes)) {
		cntOld := int32(len(w.stats.Archetypes))
		cntNew := int32(len(w.storage.archetypes))
		var i int32
		for i = range cntOld {
			arch := &w.storage.archetypes[i]
			archStats := &w.stats.Archetypes[i]
			arch.UpdateStats(archStats, &w.storage, includeTables)
			memory += archStats.Memory
			memoryUsed += archStats.MemoryUsed
		}
		for i = cntOld; i < cntNew; i++ {
			arch := &w.storage.archetypes[i]
			w.stats.Archetypes = append(w.stats.Archetypes, arch.Stats(&w.storage, includeTables))
			archStats := &w.stats.Archetypes[i]
			memory += archStats.Memory
			memoryUsed += archStats.MemoryUsed
		}
	} else {
		w.stats.Archetypes = w.stats.Archetypes[:0]
		for i := range len(w.storage.archetypes) {
			mem, used := w.storage.archetypes[i].CountMemory(&w.storage)
			memory += mem
			memoryUsed += used
		}
	}

	w.stats.NumArchetypes = len(w.storage.archetypes)
	w.stats.Locked = w.IsLocked()
	w.stats.Memory = memory
	w.stats.MemoryUsed = memoryUsed
	w.stats.CachedFilters = len(w.storage.cache.filters)
	w.stats.Observers = int(w.storage.observers.totalCount)

	return w.stats
}

// Shrink reduces memory usage by shrinking the capacity of archetype tables.
// Capacity is reduced to the next power-of-2 of what is occupied,
// but never below the initial capacities specified during world initialization.
// Further, it frees empty tables of archetypes with relations.
//
// The optional stopAfter argument sets a time limit for the shrink operation.
// This allows the method to terminate early if the time budget is exceeded,
// which is especially useful in real-time environments (e.g., games or simulations)
// where maintaining frame rate is critical.
//
// If the time limit is reached before all shrink operations are completed,
// the method returns true to indicate that further shrink work remains.
// This enables incremental cleanup across multiple frames or update cycles.
//
// Note that timer resolution is limited, particularly on Windows.
// On Windows, the shortest effective time limit should be around 0.5ms,
// while it is in the range of microseconds on Unix systems.
//
// This method should not be used regularly!
// Usually, memory should stay allocated for reuse when new entities are created or
// moved between archetypes when adding or removing components.
// However, it might be useful in memory-constrained environments e.g. after initialization.
func (w *World) Shrink(stopAfter ...time.Duration) bool {
	if len(stopAfter) > 1 {
		panic("no more than one time limit stopAfter can be given")
	}
	limit := time.Hour
	if len(stopAfter) > 0 {
		limit = stopAfter[0]
	}
	return w.storage.Shrink(limit)
}
