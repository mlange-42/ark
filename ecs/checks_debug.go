//go:build debug

package ecs

func (c *cursor) checkQueryNext() {
	if c.table < -1 {
		panic("query iteration already finished. Create a new query to iterate again")
	}
}

func (c *cursor) checkQueryGet() {
	if c.table < 0 {
		panic("query already iterated or iteration not started yet")
	}
}

func (s *storage) checkHasComponent(entity Entity, comp ID) {
	index := s.entities[entity.id]
	if s.components[comp.id].columns[index.table] == nil {
		panic("entity does not have the requested component")
	}
}
