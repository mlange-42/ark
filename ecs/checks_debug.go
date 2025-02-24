//go:build debug

package ecs

func checkQueryNext(cursor *cursor) {
	if cursor.archetype < -1 {
		panic("query iteration already finished. Create a new query to iterate again")
	}
}

func checkQueryGet(cursor *cursor) {
	if cursor.archetype < 0 {
		panic("query already iterated or iteration not started yet")
	}
}
