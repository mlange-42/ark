package ecs

import (
	"testing"
)

func TestMap(t *testing.T) {
	w := NewWorld(1024)

	posMap := NewMap[Position](&w)
	velMap := NewMap[Velocity](&w)

	e1 := w.NewEntity()

	posMap.Add(e1, &Position{})
	velMap.Add(e1, &Velocity{})

	if !posMap.Has(e1) {
		t.Error("posMap should have e1")
	}
	if !velMap.Has(e1) {
		t.Error("velMap should have e1")
	}

	pos := posMap.Get(e1)
	pos.X = 100

	pos = posMap.Get(e1)
	if pos.X != 100.0 {
		t.Errorf("expected pos.X = 100, got %v", pos.X)
	}

	posMap.Set(e1, &Position{X: -1, Y: -2})
	pos = posMap.Get(e1)
	if pos.X != -1.0 {
		t.Errorf("expected pos.X = -1, got %v", pos.X)
	}

	posMap.Remove(e1)
	if posMap.Has(e1) {
		t.Error("posMap should not have e1 after removal")
	}

	if posMap.Get(e1) != nil {
		t.Error("expected nil from posMap.Get after removal")
	}
	if posMap.GetUnchecked(e1) != nil {
		t.Error("expected nil from posMap.GetUnchecked after removal")
	}

	e2 := posMap.NewEntityFn(func(a *Position) {
		a.X = 100
	})
	if posMap.Get(e2).X != 100.0 {
		t.Errorf("expected pos.X = 100, got %v", posMap.Get(e2).X)
	}

	posMap.Remove(e2)
	posMap.AddFn(e2, func(a *Position) {
		a.X = 200
	})
	if posMap.Get(e2).X != 200.0 {
		t.Errorf("expected pos.X = 200, got %v", posMap.Get(e2).X)
	}

	expectPanic(t, func() {
		posMap.Get(Entity{})
	})
	expectPanic(t, func() {
		posMap.Has(Entity{})
	})
	expectPanic(t, func() {
		posMap.Add(Entity{}, &Position{})
	})
	expectPanic(t, func() {
		posMap.Set(Entity{}, &Position{})
	})
	expectPanic(t, func() {
		posMap.AddFn(Entity{}, func(a *Position) {})
	})
	expectPanic(t, func() {
		posMap.Remove(Entity{})
	})
}

func TestMapNewEntity(t *testing.T) {
	w := NewWorld(1024)
	posMap := NewMap[Position](&w)

	e := posMap.NewEntity(&Position{X: 1, Y: 2})
	pos := posMap.Get(e)

	if pos == nil {
		t.Fatal("expected non-nil position")
	}
	if *pos != (Position{X: 1, Y: 2}) {
		t.Errorf("expected Position{1,2}, got %+v", *pos)
	}
}

func TestMapNewBatch(t *testing.T) {
	n := 12
	w := NewWorld(8)
	mapper := NewMap[CompA](&w)

	for range n {
		_ = mapper.NewEntity(&CompA{})
	}

	w.RemoveEntity(w.NewEntity())
	mapper.NewBatch(n*2, &CompA{})

	filter := NewFilter1[CompA](&w)
	query := filter.Query()
	cnt := 0
	var lastEntity Entity

	for query.Next() {
		_ = query.Get()
		lastEntity = query.Entity()
		cnt++
	}

	if !mapper.Has(lastEntity) {
		t.Error("mapper should have lastEntity")
	}
	if cnt != n*3 {
		t.Errorf("expected %d entities, got %d", n*3, cnt)
	}
}

func TestMapNewBatchFn(t *testing.T) {
	n := 12
	w := NewWorld(8)
	mapper := NewMap[CompA](&w)

	for range n {
		_ = mapper.NewEntity(&CompA{})
	}

	w.RemoveEntity(w.NewEntity())
	mapper.NewBatchFn(2*n, func(entity Entity, a *CompA) {
		a.X = 5
		a.Y = 6
	})

	filter := NewFilter1[CompA](&w)
	query := filter.Query()
	cnt := 0
	var lastEntity Entity

	for query.Next() {
		_ = query.Get()
		lastEntity = query.Entity()
		cnt++
	}

	if !mapper.Has(lastEntity) {
		t.Error("mapper should have lastEntity")
	}
	if cnt != 3*n {
		t.Errorf("expected %d entities, got %d", 3*n, cnt)
	}

	mapper.NewBatchFn(5, nil)
}

func TestMapAddBatch(t *testing.T) {
	n := 12
	w := NewWorld(8)

	mapper := NewMap[CompA](&w)
	posMap := NewMap[Position](&w)
	posVelMap := NewMap2[Position, Velocity](&w)

	cnt := 1
	posMap.NewBatchFn(n, func(entity Entity, pos *Position) {
		pos.X = float64(cnt)
		cnt++
	})
	posVelMap.NewBatchFn(n, func(entity Entity, pos *Position, _ *Velocity) {
		pos.X = float64(cnt)
		cnt++
	})

	if cnt != 2*n+1 {
		t.Errorf("expected cnt = %d, got %d", 2*n+1, cnt)
	}

	filter := NewFilter1[Position](&w)
	mapper.AddBatch(filter.Batch(), &CompA{})

	filter2 := NewFilter1[CompA](&w)
	query := filter2.Query()
	cnt = 0
	for query.Next() {
		pos := posMap.Get(query.Entity())
		if pos.X <= 0.0 {
			t.Errorf("expected pos.X > 0, got %v", pos.X)
		}
		cnt++
	}

	if cnt != 2*n {
		t.Errorf("expected %d entities, got %d", 2*n, cnt)
	}

	mapper.RemoveBatch(filter2.Batch(), nil)

	query = filter2.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}

	if cnt != 0 {
		t.Errorf("expected 0 entities after removal, got %d", cnt)
	}
}

func TestMapAddBatchFn(t *testing.T) {
	n := 12
	w := NewWorld(8)

	mapper := NewMap[CompA](&w)
	posMap := NewMap[Position](&w)
	posVelMap := NewMap2[Position, Velocity](&w)

	cnt := 1
	posMap.NewBatchFn(n, func(entity Entity, pos *Position) {
		pos.X = float64(cnt)
		cnt++
	})
	posVelMap.NewBatchFn(n, func(entity Entity, pos *Position, _ *Velocity) {
		pos.X = float64(cnt)
		cnt++
	})

	if cnt != 2*n+1 {
		t.Errorf("expected cnt = %d, got %d", 2*n+1, cnt)
	}

	filter := NewFilter1[Position](&w)
	cnt = 0
	mapper.AddBatchFn(filter.Batch(), func(entity Entity, a *CompA) {
		a.X = float64(cnt)
		cnt++
	})

	filter2 := NewFilter1[CompA](&w)
	query := filter2.Query()
	cnt = 0
	for query.Next() {
		a := query.Get()
		if a.X != float64(cnt) {
			t.Errorf("expected a.X = %v, got %v", float64(cnt), a.X)
		}
		pos := posMap.Get(query.Entity())
		if pos.X <= 0.0 {
			t.Errorf("expected pos.X > 0, got %v", pos.X)
		}
		cnt++
	}

	if cnt != 2*n {
		t.Errorf("expected %d entities, got %d", 2*n, cnt)
	}

	cnt = 0
	mapper.RemoveBatch(filter2.Batch(), func(entity Entity) {
		cnt++
	})

	if cnt != 2*n {
		t.Errorf("expected %d removals, got %d", 2*n, cnt)
	}

	query = filter2.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}

	if cnt != 0 {
		t.Errorf("expected 0 entities after removal, got %d", cnt)
	}
}

func TestMapSliceComponent(t *testing.T) {
	n := 12
	w := NewWorld(8)

	mapper := NewMap[CompA](&w)
	sliceMap := NewMap[SliceComp](&w)

	cnt := 0
	sliceMap.NewBatchFn(n, func(entity Entity, sl *SliceComp) {
		sl.Slice = []int{cnt + 1, cnt + 2, cnt + 3, cnt + 4}
		cnt++
	})

	if cnt != n {
		t.Errorf("expected cnt = %d, got %d", n, cnt)
	}

	filter := NewFilter1[SliceComp](&w)
	cnt = 0
	mapper.AddBatchFn(filter.Batch(), func(entity Entity, a *CompA) {
		a.X = float64(cnt)
		cnt++
	})

	filter2 := NewFilter1[SliceComp](&w)
	query := filter2.Query()
	cnt = 0
	for query.Next() {
		sl := query.Get()
		expected := []int{cnt + 1, cnt + 2, cnt + 3, cnt + 4}
		actual := sl.Slice
		if len(actual) != len(expected) {
			t.Errorf("expected slice length %d, got %d", len(expected), len(actual))
		} else {
			for i := range expected {
				if actual[i] != expected[i] {
					t.Errorf("expected slice[%d] = %d, got %d", i, expected[i], actual[i])
				}
			}
		}
		cnt++
	}

	if cnt != n {
		t.Errorf("expected %d entities, got %d", n, cnt)
	}
}

func TestMapRelation(t *testing.T) {
	w := NewWorld(32)
	childMap := NewMap[ChildOf](&w)

	parent1 := w.NewEntity()
	parent2 := w.NewEntity()
	e := w.NewEntity()

	childMap.Add(e, &ChildOf{}, parent1)

	if got := childMap.GetRelation(e); got != parent1 {
		t.Errorf("expected relation to parent1, got %v", got)
	}
	if got := childMap.GetRelationUnchecked(e); got != parent1 {
		t.Errorf("expected unchecked relation to parent1, got %v", got)
	}

	childMap.SetRelation(e, parent2)
	if got := childMap.GetRelation(e); got != parent2 {
		t.Errorf("expected relation to parent2, got %v", got)
	}
	if got := childMap.GetRelationUnchecked(e); got != parent2 {
		t.Errorf("expected unchecked relation to parent2, got %v", got)
	}

	expectPanic(t, func() {
		childMap.GetRelation(Entity{})
	})

	childMap.SetRelation(e, Entity{})
	if got := childMap.GetRelation(e); got != (Entity{}) {
		t.Errorf("expected zero entity relation, got %v", got)
	}
	if got := childMap.GetRelationUnchecked(e); got != (Entity{}) {
		t.Errorf("expected unchecked zero entity relation, got %v", got)
	}

	deadParent := w.NewEntity()
	w.RemoveEntity(deadParent)

	expectPanicWithValue(t,
		"can't use a dead entity as relation target, except for the zero entity",
		func() {
			childMap.SetRelation(e, deadParent)
		})

	expectPanicWithValue(t,
		"relation targets must be fully specified",
		func() {
			childMap.NewEntity(&ChildOf{})
		})
}

func TestMapRelationBatch(t *testing.T) {
	n := 24
	w := NewWorld(16)

	parent1 := w.NewEntity()
	parent2 := w.NewEntity()
	parent3 := w.NewEntity()

	mapper := NewMap3[Position, Velocity, ChildOf](&w)
	childMap := NewMap[ChildOf](&w)

	mapper.NewBatch(n, &Position{}, &Velocity{}, &ChildOf{}, RelIdx(2, parent1))
	mapper.NewBatch(n, &Position{}, &Velocity{}, &ChildOf{}, RelIdx(2, parent2))

	filter := NewFilter1[ChildOf](&w)

	childMap.SetRelationBatch(filter.Batch(RelIdx(0, parent2)), parent3, func(entity Entity) {
		if got := childMap.GetRelation(entity); got != parent3 {
			t.Errorf("expected relation to parent3, got %v", got)
		}
	})

	query := filter.Query(RelIdx(0, parent2))
	cnt := 0
	for query.Next() {
		cnt++
	}
	if cnt != 0 {
		t.Errorf("expected 0 entities with relation to parent2, got %d", cnt)
	}

	query = filter.Query(RelIdx(0, parent3))
	cnt = 0
	for query.Next() {
		cnt++
	}
	if cnt != n {
		t.Errorf("expected %d entities with relation to parent3, got %d", n, cnt)
	}
}
