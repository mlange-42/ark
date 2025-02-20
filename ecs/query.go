package ecs

type cursor struct {
	table    int
	index    uintptr
	maxIndex int64
}

// Query2 is a filter for two components.
type Query2[A any, B any] struct {
	world      *World
	mask       Mask
	cursor     cursor
	componentA *componentStorage
	componentB *componentStorage
	columnA    *column
	columnB    *column
}

// NewQuery2 creates a new [Query2].
func NewQuery2[A any, B any](world *World) Query2[A, B] {
	idA := ComponentID[A](world)
	idB := ComponentID[B](world)

	return Query2[A, B]{
		world:      world,
		mask:       All(idA, idB),
		componentA: &world.storage.components[idA.id],
		componentB: &world.storage.components[idB.id],
		cursor: cursor{
			table:    -1,
			index:    0,
			maxIndex: -1,
		},
	}
}

func (q *Query2[A, B]) Next() bool {
	if int64(q.cursor.index) < q.cursor.maxIndex {
		q.cursor.index++
		return true
	}
	return q.nextTable()
}

func (q *Query2[A, B]) nextTable() bool {
	maxTableIndex := len(q.world.storage.tables) - 1
	for q.cursor.table < maxTableIndex {
		q.cursor.table++
		table := &q.world.storage.tables[q.cursor.table]
		archetype := &q.world.storage.archetypes[table.archetype]
		if !archetype.mask.Contains(&q.mask) || table.entities.Len() == 0 {
			continue
		}

		q.columnA = q.componentA.columns[q.cursor.table]
		q.columnB = q.componentB.columns[q.cursor.table]
		q.cursor.index = 0
		q.cursor.maxIndex = int64(table.entities.Len() - 1)
		return true
	}
	q.cursor.table = -1
	q.cursor.index = 0
	q.cursor.maxIndex = -1
	return false
}

func (q *Query2[A, B]) Get() (*A, *B) {
	return (*A)(q.columnA.Get(q.cursor.index)),
		(*B)(q.columnB.Get(q.cursor.index))
}
