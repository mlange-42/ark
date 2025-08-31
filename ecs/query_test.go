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
	if query.Count() != 2*n {
		t.Errorf("expected %d, got %d", 2*n, query.Count())
	}
	cnt := 0
	for query.Next() {
		_ = query.Entity()
		_ = query.Get(compA)
		if !query.Has(compA) {
			t.Errorf("expected true, got false")
		}
		cnt++
	}
	if cnt != 2*n {
		t.Errorf("expected %d, got %d", 2*n, cnt)
	}
	// filter without
	filter = NewUnsafeFilter(&w, compA, compB, compC).Without(posID)
	query = filter.Query()
	if query.Count() != n {
		t.Errorf("expected %d, got %d", n, query.Count())
	}
	cnt = 0
	for query.Next() {
		_ = query.Entity()
		_ = query.Get(compA)
		if !query.Has(compA) {
			t.Errorf("expected true, got false")
		}
		if query.Has(posID) {
			t.Errorf("expected false, got true")
		}
		if !equalIDs(query.IDs().data, []ID{{0}, {1}, {2}}) {
			t.Errorf("expected %v, got %v", []ID{{0}, {1}, {2}}, query.IDs().data)
		}
		cnt++
	}
	if cnt != n {
		t.Errorf("expected %d, got %d", n, cnt)
	}
	// filter exclusive
	filter = NewUnsafeFilter(&w, compA, compB, compC).Exclusive()
	query = filter.Query()
	if query.Count() != n {
		t.Errorf("expected %d, got %d", n, query.Count())
	}
	cnt = 0
	for query.Next() {
		_ = query.Entity()
		_ = query.Get(compA)
		cnt++
	}
	if cnt != n {
		t.Errorf("expected %d, got %d", n, cnt)
	}
}

func equalIDs(a, b []ID) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
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
	if query.Count() != 0 {
		t.Errorf("expected 0, got %d", query.Count())
	}
	expectPanic(t, func() { query.Get(compA) })
	expectPanic(t, func() { query.Entity() })
	cnt := 0
	for query.Next() {
		cnt++
	}
	if cnt != 0 {
		t.Errorf("expected 0, got %d", cnt)
	}
	expectPanic(t, func() { query.Get(compA) })
	expectPanic(t, func() { query.Entity() })
	expectPanic(t, func() { query.Next() })
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
	if query.Count() != 2*n {
		t.Errorf("expected %d, got %d", 2*n, query.Count())
	}
	cnt := 0
	for query.Next() {
		_ = query.Entity()
		_ = query.Get(childID)
		cnt++
	}
	if cnt != 2*n {
		t.Errorf("expected %d, got %d", 2*n, cnt)
	}
	// relation filter
	filter = NewUnsafeFilter(&w, childID, compB, compC)
	query = filter.Query(RelID(childID, parent2))
	if query.Count() != n {
		t.Errorf("expected %d, got %d", n, query.Count())
	}
	cnt = 0
	for query.Next() {
		_ = query.Entity()
		_ = query.Get(childID)
		if query.GetRelation(childID) != parent2 {
			t.Errorf("expected %v, got %v", parent2, query.GetRelation(childID))
		}
		cnt++
	}
	if cnt != n {
		t.Errorf("expected %d, got %d", n, cnt)
	}
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
	if count != counter {
		t.Errorf("Number of entities (%d) should match count (%d)", counter, count)
	}
	filter2 := NewUnsafeFilter(&world)
	query2 := filter2.Query()
	count = query2.Count()
	counter = 0
	for query2.Next() {
		counter++
	}
	if count != counter {
		t.Errorf("Number of entities (%d) should match count (%d)", counter, count)
	}
}
