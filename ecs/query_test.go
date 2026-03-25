package ecs

import (
	"sync"
	"testing"
)

func TestQuery(t *testing.T) {
	n := 10
	w := NewWorld(4)
	u := w.Unsafe()

	compA := ComponentID[CompA](w)
	compB := ComponentID[CompB](w)
	compC := ComponentID[CompC](w)
	posID := ComponentID[Position](w)

	for range n {
		_ = u.NewEntity(compA, compB, compC)

		e := u.NewEntity(compA, compB, compC)
		u.Remove(e, compA)

		e = u.NewEntity(compA, compB, compC)
		u.Add(e, posID)
	}

	// normal filter
	filter := NewUnsafeFilter(w, compA, compB, compC)
	query := filter.Query()
	count := query.Count()
	expectEqual(t, 2*n, count)

	e := query.EntityAt(count - 1)
	expectEqual(t, 31, int(e.ID()))

	cnt := 0
	for query.Next() {
		_ = query.Entity()
		_ = query.Get(compA)
		expectTrue(t, query.Has(compA))
		cnt++
	}
	expectEqual(t, cnt, 2*n)
	query.Close() // should not panic anymore

	// filter without
	filter = NewUnsafeFilter(w, compA, compB, compC).Without(posID)
	query = filter.Query()
	count = query.Count()
	expectEqual(t, n, count)

	e = query.EntityAt(count - 1)
	expectEqual(t, 29, int(e.ID()))

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
	filter = NewUnsafeFilter(w, compA, compB, compC).Exclusive()
	query = filter.Query()
	count = query.Count()
	expectEqual(t, n, count)

	e = query.EntityAt(count - 1)
	expectEqual(t, 29, int(e.ID()))

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

	compA := ComponentID[CompA](w)
	compB := ComponentID[CompB](w)
	compC := ComponentID[CompC](w)
	posID := ComponentID[Position](w)

	for range 10 {
		e1 := w.NewEntity()
		u.Add(e1, posID)
	}

	filter := NewUnsafeFilter(w, compA, compB, compC)
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

	childID := ComponentID[ChildOf](w)
	childID2 := ComponentID[ChildOf2](w)
	compB := ComponentID[CompB](w)
	compC := ComponentID[CompC](w)

	for range n {
		_ = u.NewEntityRel([]ID{childID, childID2, compB, compC}, RelID(childID, parent1), RelID(childID2, parent3))
		_ = u.NewEntityRel([]ID{childID, childID2, compB, compC}, RelID(childID, parent2), RelID(childID2, parent3))
		e := u.NewEntityRel([]ID{childID, childID2, compB, compC}, RelID(childID, parent3), RelID(childID2, parent3))
		w.RemoveEntity(e)
	}

	// normal filter
	filter := NewUnsafeFilter(w, childID, compB, compC)
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
	filter = NewUnsafeFilter(w, childID, compB, compC)

	expectPanicsWithValue(t, "relations created with RelIdx can't be used in the unsafe API, use RelID or Rel instead",
		func() {
			filter.Query(RelIdx(0, parent2))
		})

	query = filter.Query(Rel[ChildOf](parent2))
	expectEqual(t, n, query.Count())
	query.Close()

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

	// multi relation filter
	filter = NewUnsafeFilter(w, childID2, childID, compB, compC)
	query = filter.Query(RelID(childID2, parent3), RelID(childID, parent2))
	count := query.Count()
	expectEqual(t, n, count)

	e := query.EntityAt(count - 1)
	expectEqual(t, 24, int(e.ID()))

	cnt = 0
	for query.Next() {
		_ = query.Entity()
		_ = query.Get(childID2)
		expectEqual(t, parent3, query.GetRelation(childID2))
		cnt++
	}
	expectEqual(t, n, cnt)
}

func TestQueryCount(t *testing.T) {
	world := NewWorld()

	parentMap := NewMap1[Position](world)
	childMap := NewMap2[Velocity, ChildOf](world)
	dummyMap := NewMap1[Velocity](world)

	parent1 := parentMap.NewEntity(&Position{})
	parent2 := parentMap.NewEntity(&Position{})

	childMap.NewEntity(&Velocity{}, &ChildOf{}, Rel[ChildOf](parent1))
	childMap.NewEntity(&Velocity{}, &ChildOf{}, Rel[ChildOf](parent2))

	dummyMap.NewEntity(&Velocity{})

	filter := NewFilter0(world)
	query := filter.Query()

	count := query.Count()
	counter := 0
	for query.Next() {
		counter++
	}

	expectEqual(t, count, counter, "Number of entities should match count")

	filter2 := NewUnsafeFilter(world)
	query2 := filter2.Query()

	count = query2.Count()
	counter = 0
	for query2.Next() {
		counter++
	}

	expectEqual(t, count, counter, "Number of entities should match count")
}

func TestQueryParallel(t *testing.T) {
	n := 10_000
	threads := 4
	perThread := n / threads
	world := NewWorld(1024)

	parents := make([]Entity, 0, threads)
	for range threads {
		parent := world.NewEntity()
		parents = append(parents, parent)
	}

	mapper := NewMap3[Position, Velocity, ChildOf](world)
	for i, p := range parents {
		cnt := 0
		mapper.NewBatchFn(perThread, func(e Entity, p *Position, v *Velocity, co *ChildOf) {
			p.X = float64(i*1_000_000 + cnt)
			expectEqual(t, i*perThread+threads+cnt+2, int(e.id))
			cnt++
		}, RelIdx(2, p))
	}
	filter := NewFilter2[Position, Velocity](world).
		With(C[ChildOf]())

	task := func(i int, par Entity, wg *sync.WaitGroup) {
		defer wg.Done()
		query := filter.Query(RelIdx(2, par))

		cnt := 0
		for query.Next() {
			pos, _ := query.Get()
			expectEqual(t, i*perThread+threads+cnt+2, int(query.Entity().id))
			expectEqual(t, i*1_000_000+cnt, int(pos.X))
			cnt++
		}
	}

	for range 100 {
		var wg sync.WaitGroup
		wg.Add(threads)
		for i, t := range parents {
			go task(i, t, &wg)
		}
		wg.Wait()
	}
}
