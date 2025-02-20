package ecs

import "unsafe"

// Map is a mapper to access components of an entity.
type Map[T any] struct {
	world   *World
	id      ID
	storage *componentStorage
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

func (m *Map[T]) Get(entity Entity) *T {
	if !m.world.Alive(entity) {
		panic("can't get a component of a dead entity")
	}
	index := m.world.entities[entity.id]
	return (*T)(m.storage.columns[index.table].Get(uintptr(index.row)))
}

func (m *Map[T]) Has(entity Entity) bool {
	if !m.world.Alive(entity) {
		panic("can't get a component of a dead entity")
	}
	index := m.world.entities[entity.id]
	return m.storage.columns[index.table] != nil
}

func (m *Map[T]) Add(entity Entity, comp *T) {
	if !m.world.Alive(entity) {
		panic("can't add a component to a dead entity")
	}
	m.world.exchange(entity, []ID{m.id}, nil, []unsafe.Pointer{unsafe.Pointer(comp)})
}

func (m *Map[T]) Remove(entity Entity) {
	if !m.world.Alive(entity) {
		panic("can't remove a component from a dead entity")
	}
	m.world.exchange(entity, nil, []ID{m.id}, nil)
}
