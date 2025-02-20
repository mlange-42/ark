package ecs

// Query2 is a filter for two components.
type Query2[A any, B any] struct {
	world      *World
	mask       Mask
	componentA *componentStorage
	componentB *componentStorage
	columnA    *column
	columnB    *column

	table    int
	index    uint32
	maxIndex int
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

		table:    -1,
		index:    0,
		maxIndex: -1,
	}
}

func (q *Query2[A, B]) Next() bool {
	if int(q.index) < q.maxIndex {
		q.index++
		return true
	}
	return q.nextTable()
}

func (q *Query2[A, B]) nextTable() bool {
	maxTableIndex := len(q.world.storage.tables) - 1
	for q.table < maxTableIndex {
		q.table++
		table := &q.world.storage.tables[q.table]
		archetype := &q.world.storage.archetypes[table.archetype]
		if !archetype.mask.Contains(&q.mask) || table.entities.Len() == 0 {
			continue
		}

		q.columnA = q.componentA.columns[q.table]
		q.columnB = q.componentB.columns[q.table]
		q.index = 0
		q.maxIndex = table.entities.Len() - 1
		return true
	}
	q.table = -1
	q.index = 0
	q.maxIndex = -1
	return false
}

func (q *Query2[A, B]) Get() (*A, *B) {
	return (*A)(q.columnA.Get(uint32(q.index))),
		(*B)(q.columnB.Get(uint32(q.index)))
}
