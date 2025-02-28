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

func (w *World) checkMapHasComponent(comp *componentStorage, table tableID) {
	if comp.columns[table] == nil {
		panic("entity does not have the requested component")
	}
}
