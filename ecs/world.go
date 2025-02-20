package ecs

// World is the central type holding entity and component data, as well as resources.
type World struct {
	registry   registry
	storage    storage
	entities   entities
	entityPool entityPool
}

// NewWorld creates a new [World].
func NewWorld(initialCapacity uint32) World {
	registry := newRegistry()

	entities := make([]entityIndex, 1, initialCapacity)
	entities[0] = entityIndex{table: 0, row: 0}
	return World{
		registry:   registry,
		storage:    newStorage(initialCapacity, &registry),
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
