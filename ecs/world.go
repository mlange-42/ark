package ecs

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
