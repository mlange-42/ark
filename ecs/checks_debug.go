//go:build debug

package ecs

func (w *World) checkQueryNext(cursor *cursor) {
	if cursor.archetype < -1 {
		panic("query iteration already finished. Create a new query to iterate again")
	}
}

func (w *World) checkQueryGet(cursor *cursor) {
	if cursor.archetype < 0 {
		panic("query already iterated or iteration not started yet")
	}
}

func (s *storage) checkHasComponent(entity Entity, comp ID) {
	index := s.entities[entity.id]
	if s.components[comp.id].columns[index.table] == nil {
		panic("entity does not have the requested component")
	}
}
