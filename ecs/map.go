package ecs

import "unsafe"

// Map is a mapper to access and manipulate components of an entity.
type Map[T any] struct {
	world     *World
	id        ID
	storage   *componentStorage
	relations []RelationID
}

// NewMap creates a new [Map].
func NewMap[T any](w *World) Map[T] {
	id := ComponentID[T](w)
	return Map[T]{
		world:   w,
		id:      id,
		storage: &w.storage.components[id.id],
	}
}

// NewEntity creates a new entity with the mapped component.
func (m *Map[T]) NewEntity(comp *T, rel ...Entity) Entity {
	m.relations = relationEntities(rel).toRelation(m.id, m.relations)
	return m.world.newEntityWith([]ID{m.id}, []unsafe.Pointer{unsafe.Pointer(comp)}, m.relations)
}

// Get returns the mapped component for the given entity.
func (m *Map[T]) Get(entity Entity) *T {
	if !m.world.Alive(entity) {
		panic("can't get a component of a dead entity")
	}
	return m.GetUnchecked(entity)
}

// GetUnchecked returns the mapped component for the given entity.
// In contrast to [Map.Get], it does not check whether the entity is alive.
// Can be used as an optimization when it is certain that the entity is alive.
func (m *Map[T]) GetUnchecked(entity Entity) *T {
	index := m.world.storage.entities[entity.id]
	checkMapHasComponent(m.storage, index.table)
	return (*T)(m.storage.columns[index.table].Get(uintptr(index.row)))
}

// Has return whether the given entity has the mapped component.
func (m *Map[T]) Has(entity Entity) bool {
	if !m.world.Alive(entity) {
		panic("can't get a component of a dead entity")
	}
	return m.HasUnchecked(entity)
}

// HasUnchecked return whether the given entity has the mapped component.
// In contrast to [Map.Has], it does not check whether the entity is alive.
// Can be used as an optimization when it is certain that the entity is alive.
func (m *Map[T]) HasUnchecked(entity Entity) bool {
	index := m.world.storage.entities[entity.id]
	return m.storage.columns[index.table] != nil
}

// Add the mapped component to the given entity.
func (m *Map[T]) Add(entity Entity, comp *T, rel ...Entity) {
	if !m.world.Alive(entity) {
		panic("can't add a component to a dead entity")
	}
	m.relations = relationEntities(rel).toRelation(m.id, m.relations)
	m.world.exchange(entity, []ID{m.id}, nil, []unsafe.Pointer{unsafe.Pointer(comp)}, m.relations)
}

// Remove the mapped component from the given entity.
func (m *Map[T]) Remove(entity Entity) {
	if !m.world.Alive(entity) {
		panic("can't remove a component from a dead entity")
	}
	m.world.exchange(entity, nil, []ID{m.id}, nil, nil)
}

// SetRelation sets the relation target for the entity and the mapped component.
func (m *Map[T]) SetRelation(entity Entity, target Entity) {
	m.relations = target.toRelation(m.id, m.relations)
	m.world.setRelations(entity, m.relations)
}

// GetRelation returns the relation target for the entity and the mapped component.
func (m *Map[T]) GetRelation(entity Entity) Entity {
	return m.world.storage.getRelation(entity, m.id)
}

// GetRelationUnchecked returns the relation target for the entity and the mapped component.
// In contrast to [Map.GetRelation], it does not check whether the entity is alive.
// Can be used as an optimization when it is certain that the entity is alive.
func (m *Map[T]) GetRelationUnchecked(entity Entity) Entity {
	return m.world.storage.getRelationUnchecked(entity, m.id)
}
