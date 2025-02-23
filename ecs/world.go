package ecs

// World is the central type holding entity and component data, as well as resources.
type World struct {
	storage    storage
	entities   entities
	isTarget   []bool
	entityPool entityPool
	resources  Resources
	locks      lock
}

// NewWorld creates a new [World].
func NewWorld(initialCapacity uint32) World {
	entities := make([]entityIndex, reservedEntities, initialCapacity+reservedEntities)
	isTarget := make([]bool, reservedEntities, initialCapacity+reservedEntities)
	// Reserved zero and wildcard entities
	for i := range reservedEntities {
		entities[i] = entityIndex{table: maxTableID, row: 0}
	}
	return World{
		storage:    newStorage(initialCapacity),
		entities:   entities,
		isTarget:   isTarget,
		entityPool: newEntityPool(initialCapacity, reservedEntities),
		resources:  newResources(),
		locks:      lock{},
	}
}

// NewEntity creates a new [Entity].
func (w *World) NewEntity() Entity {
	w.checkLocked()

	entity, _ := w.createEntity(0)
	return entity
}

// Alive return whether the given entity is alive.
func (w *World) Alive(entity Entity) bool {
	return w.entityPool.Alive(entity)
}

// RemoveEntity removes the given entity from the world.
func (w *World) RemoveEntity(entity Entity) {
	// TODO: check lock.
	if !w.entityPool.Alive(entity) {
		panic("can't remove a dead entity")
	}
	index := &w.entities[entity.id]
	table := &w.storage.tables[index.table]

	swapped := table.Remove(index.row)

	w.entityPool.Recycle(entity)

	if swapped {
		swapEntity := table.GetEntity(uintptr(index.row))
		w.entities[swapEntity.id].row = index.row
	}
	index.table = maxTableID
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
