package ecs

import "unsafe"

// Map2 is a mapper to access 2 components of an entity.
type Map2[A any, B any] struct {
	world    *World
	ids      []ID
	storageA *componentStorage
	storageB *componentStorage
}

// NewMap2 creates a new [Map2].
func NewMap2[A any, B any](w *World) Map2[A, B] {
	ids := []ID{
		ComponentID[A](w),
		ComponentID[B](w),
	}
	return Map2[A, B]{
		world:    w,
		ids:      ids,
		storageA: &w.storage.components[ids[0].id],
		storageB: &w.storage.components[ids[1].id],
	}
}

func (m *Map2[A, B]) Get(entity Entity) (*A, *B) {
	if !m.world.Alive(entity) {
		panic("can't get components of a dead entity")
	}
	index := m.world.entities[entity.id]
	return (*A)(m.storageA.columns[index.table].Get(uintptr(index.row))),
		(*B)(m.storageB.columns[index.table].Get(uintptr(index.row)))
}

func (m *Map2[A, B]) Add(entity Entity, a *A, b *B) {
	if !m.world.Alive(entity) {
		panic("can't add components to a dead entity")
	}
	m.world.exchange(entity, m.ids, nil, []unsafe.Pointer{
		unsafe.Pointer(a),
		unsafe.Pointer(b),
	})
}

func (m *Map2[A, B]) Remove(entity Entity) {
	if !m.world.Alive(entity) {
		panic("can't remove components from a dead entity")
	}
	m.world.exchange(entity, nil, m.ids, nil)
}
