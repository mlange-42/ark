package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	n := 10
	w := NewWorld(4)
	u := w.Unsafe()

	compA := ComponentID[CompA](&w)
	compB := ComponentID[CompB](&w)
	compC := ComponentID[CompC](&w)
	posID := ComponentID[Position](&w)

	for range n {
		_ = u.NewEntity(compA, compB, compC)

		e := u.NewEntity(compA, compB, compC)
		u.Remove(e, compA)

		e = u.NewEntity(compA, compB, compC)
		u.Add(e, posID)
	}

	// normal filter
	filter := NewFilter(compA, compB, compC)
	query := filter.Query(&w)
	assert.Equal(t, 2*n, query.Count())

	cnt := 0
	for query.Next() {
		_ = query.Entity()
		_ = query.Get(compA)
		assert.True(t, query.Has(compA))
		cnt++
	}
	assert.Equal(t, cnt, 2*n)

	// filter without
	filter = NewFilter(compA, compB, compC).Without(posID)
	query = filter.Query(&w)
	assert.Equal(t, n, query.Count())

	cnt = 0
	for query.Next() {
		_ = query.Entity()
		_ = query.Get(compA)
		assert.True(t, query.Has(compA))
		assert.False(t, query.Has(posID))
		assert.Equal(t, []ID{{0}, {1}, {2}}, query.IDs())
		cnt++
	}
	assert.Equal(t, cnt, n)

	// filter exclusive
	filter = NewFilter(compA, compB, compC).Exclusive()
	query = u.Query(filter)
	assert.Equal(t, n, query.Count())

	cnt = 0
	for query.Next() {
		_ = query.Entity()
		_ = query.Get(compA)
		cnt++
	}
	assert.Equal(t, cnt, n)
}

func TestQueryEmpty(t *testing.T) {
	w := NewWorld(4)
	u := w.Unsafe()

	compA := ComponentID[CompA](&w)
	compB := ComponentID[CompB](&w)
	compC := ComponentID[CompC](&w)
	posID := ComponentID[Position](&w)

	for range 10 {
		e1 := w.NewEntity()
		u.Add(e1, posID)
	}

	filter := NewFilter(compA, compB, compC)
	query := filter.Query(&w)
	assert.Equal(t, 0, query.Count())

	assert.Panics(t, func() { query.Get(compA) })
	assert.Panics(t, func() { query.Entity() })

	cnt := 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)

	assert.Panics(t, func() { query.Get(compA) })
	assert.Panics(t, func() { query.Entity() })
	assert.Panics(t, func() { query.Next() })
}

func TestQueryRelations(t *testing.T) {
	n := 10
	w := NewWorld(4)
	u := w.Unsafe()

	parent1 := w.NewEntity()
	parent2 := w.NewEntity()
	parent3 := w.NewEntity()

	childID := ComponentID[ChildOf](&w)
	compB := ComponentID[CompB](&w)
	compC := ComponentID[CompC](&w)

	for range n {
		_ = u.NewEntityRel([]ID{childID, compB, compC}, RelID(childID, parent1))
		_ = u.NewEntityRel([]ID{childID, compB, compC}, RelID(childID, parent2))
		e := u.NewEntityRel([]ID{childID, compB, compC}, RelID(childID, parent3))
		w.RemoveEntity(e)
	}

	// normal filter
	filter := NewFilter(childID, compB, compC)
	query := filter.Query(&w)
	assert.Equal(t, 2*n, query.Count())

	cnt := 0
	for query.Next() {
		_ = query.Entity()
		_ = query.Get(childID)
		cnt++
	}
	assert.Equal(t, cnt, 2*n)

	// relation filter
	filter = NewFilter(childID, compB, compC)
	query = filter.Query(&w, RelID(childID, parent2))
	assert.Equal(t, n, query.Count())

	cnt = 0
	for query.Next() {
		_ = query.Entity()
		_ = query.Get(childID)
		assert.Equal(t, parent2, query.GetRelation(childID))
		cnt++
	}
	assert.Equal(t, cnt, n)
}

func BenchmarkQueryUnsafePosVel_1000(b *testing.B) {
	n := 1000
	world := NewWorld(128)

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	mapper := NewMap2[Position, Velocity](&world)
	mapper.NewBatch(n, &Position{}, &Velocity{X: 1, Y: 0})

	filter := NewFilter(posID, velID)
	for b.Loop() {
		query := filter.Query(&world)
		for query.Next() {
			pos := (*Position)(query.Get(posID))
			vel := (*Velocity)(query.Get(velID))
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}
