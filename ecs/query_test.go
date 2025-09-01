package ecs

import (
	"testing"
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
	filter := NewUnsafeFilter(&w, compA, compB, compC)
	query := filter.Query()
	expectEqual(t, 2*n, query.Count())

	cnt := 0
	for query.Next() {
		_ = query.Entity()
		_ = query.Get(compA)
		expectTrue(t, query.Has(compA))
		cnt++
	}
	expectEqual(t, cnt, 2*n)

	// filter without
	filter = NewUnsafeFilter(&w, compA, compB, compC).Without(posID)
	query = filter.Query()
	expectEqual(t, n, query.Count())

	cnt = 0
	for query.Next() {
		_ = query.Entity()
		_ = query.Get(compA)
		expectTrue(t, query.Has(compA))
		expectFalse(t, query.Has(posID))
		expectSlicesEqual(t, []ID{{0}, {1}, {2}}, query.IDs().data)
		cnt++
	}
	expectEqual(t, cnt, n)

	// filter exclusive
	filter = NewUnsafeFilter(&w, compA, compB, compC).Exclusive()
	query = filter.Query()
	expectEqual(t, n, query.Count())

	cnt = 0
	for query.Next() {
		_ = query.Entity()
		_ = query.Get(compA)
		cnt++
	}
	expectEqual(t, cnt, n)
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

	filter := NewUnsafeFilter(&w, compA, compB, compC)
	query := filter.Query()
	expectEqual(t, 0, query.Count())

	expectPanics(t, func() { query.Get(compA) })
	expectPanics(t, func() { query.Entity() })

	cnt := 0
	for query.Next() {
		cnt++
	}
	expectEqual(t, 0, cnt)

	expectPanics(t, func() { query.Get(compA) })
	expectPanics(t, func() { query.Entity() })
	expectPanics(t, func() { query.Next() })
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
	filter := NewUnsafeFilter(&w, childID, compB, compC)
	query := filter.Query()
	expectEqual(t, 2*n, query.Count())

	cnt := 0
	for query.Next() {
		_ = query.Entity()
		_ = query.Get(childID)
		cnt++
	}
	expectEqual(t, cnt, 2*n)

	// relation filter
	filter = NewUnsafeFilter(&w, childID, compB, compC)
	query = filter.Query(RelID(childID, parent2))
	expectEqual(t, n, query.Count())

	cnt = 0
	for query.Next() {
		_ = query.Entity()
		_ = query.Get(childID)
		expectEqual(t, parent2, query.GetRelation(childID))
		cnt++
	}
	expectEqual(t, cnt, n)
}

func TestQueryCount(t *testing.T) {
	world := NewWorld()

	parentMap := NewMap1[Position](&world)
	childMap := NewMap2[Velocity, ChildOf](&world)
	dummyMap := NewMap1[Velocity](&world)

	parent1 := parentMap.NewEntity(&Position{})
	parent2 := parentMap.NewEntity(&Position{})

	childMap.NewEntity(&Velocity{}, &ChildOf{}, Rel[ChildOf](parent1))
	childMap.NewEntity(&Velocity{}, &ChildOf{}, Rel[ChildOf](parent2))

	dummyMap.NewEntity(&Velocity{})

	filter := NewFilter0(&world)
	query := filter.Query()

	count := query.Count()
	counter := 0
	for query.Next() {
		counter++
	}

	expectEqual(t, count, counter, "Number of entities should match count")

	filter2 := NewUnsafeFilter(&world)
	query2 := filter2.Query()

	count = query2.Count()
	counter = 0
	for query2.Next() {
		counter++
	}

	expectEqual(t, count, counter, "Number of entities should match count")
}
