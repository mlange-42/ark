package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	w := NewWorld(1024)

	posMap := NewMap[Position](&w)
	velMap := NewMap[Velocity](&w)

	e1 := w.NewEntity()

	posMap.Add(e1, &Position{})
	velMap.Add(e1, &Velocity{})

	assert.True(t, posMap.Has(e1))
	assert.True(t, velMap.Has(e1))

	pos := posMap.Get(e1)
	pos.X = 100

	pos = posMap.Get(e1)
	assert.Equal(t, 100.0, pos.X)

	posMap.Remove(e1)
	assert.False(t, posMap.Has(e1))
	assert.Panics(t, func() { posMap.Get(e1) })

	e2 := posMap.NewEntityFn(func(a *Position) {
		a.X = 100
	})
	assert.Equal(t, 100.0, posMap.Get(e2).X)

	posMap.Remove(e2)
	posMap.AddFn(e2, func(a *Position) {
		a.X = 200
	})
	assert.Equal(t, 200.0, posMap.Get(e2).X)

	assert.Panics(t, func() {
		posMap.Get(Entity{})
	})
	assert.Panics(t, func() {
		posMap.Has(Entity{})
	})
	assert.Panics(t, func() {
		posMap.Add(Entity{}, &Position{})
	})
	assert.Panics(t, func() {
		posMap.AddFn(Entity{}, func(a *Position) {})
	})
	assert.Panics(t, func() {
		posMap.Remove(Entity{})
	})
}

func TestMapNewEntity(t *testing.T) {
	w := NewWorld(1024)

	posMap := NewMap[Position](&w)

	e := posMap.NewEntity(&Position{X: 1, Y: 2})

	pos := posMap.Get(e)
	assert.Equal(t, Position{X: 1, Y: 2}, *pos)
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
	assert.True(t, mapper.Has(lastEntity))
	assert.Equal(t, n*3, cnt)
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
	assert.True(t, mapper.Has(lastEntity))
	assert.Equal(t, 3*n, cnt)

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
	assert.Equal(t, 2*n+1, cnt)

	filter := NewFilter1[Position](&w)
	mapper.AddBatch(filter.Batch(), &CompA{})

	filter2 := NewFilter1[CompA](&w)
	query := filter2.Query()
	cnt = 0
	for query.Next() {
		pos := posMap.Get(query.Entity())
		assert.Greater(t, pos.X, 0.0)
		cnt++
	}
	assert.Equal(t, 2*n, cnt)

	mapper.RemoveBatch(filter2.Batch(), nil)

	query = filter2.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)
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
	assert.Equal(t, 2*n+1, cnt)

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
		assert.EqualValues(t, cnt, a.X)
		pos := posMap.Get(query.Entity())
		assert.Greater(t, pos.X, 0.0)
		cnt++
	}
	assert.Equal(t, 2*n, cnt)

	cnt = 0
	mapper.RemoveBatch(filter2.Batch(), func(entity Entity) {
		cnt++
	})
	assert.Equal(t, 2*n, cnt)

	query = filter2.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)
}

func TestMapRelation(t *testing.T) {
	w := NewWorld(32)

	childMap := NewMap[ChildOf](&w)

	parent1 := w.NewEntity()
	parent2 := w.NewEntity()

	e := w.NewEntity()

	childMap.Add(e, &ChildOf{}, parent1)
	assert.Equal(t, parent1, childMap.GetRelation(e))
	assert.Equal(t, parent1, childMap.GetRelationUnchecked(e))

	childMap.SetRelation(e, parent2)
	assert.Equal(t, parent2, childMap.GetRelation(e))
	assert.Equal(t, parent2, childMap.GetRelationUnchecked(e))

	assert.Panics(t, func() {
		childMap.GetRelation(Entity{})
	})

	childMap.SetRelation(e, Entity{})
	assert.Equal(t, Entity{}, childMap.GetRelation(e))
	assert.Equal(t, Entity{}, childMap.GetRelationUnchecked(e))

	deadParent := w.NewEntity()
	w.RemoveEntity(deadParent)
	assert.PanicsWithValue(t,
		"can't use a dead entity as relation target, except for the zero entity",
		func() {
			childMap.SetRelation(e, deadParent)
		})
	assert.PanicsWithValue(t,
		"relations must be fully specified",
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
		assert.Equal(t, parent3, childMap.GetRelation(entity))
	})

	query := filter.Query(RelIdx(0, parent2))
	cnt := 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)

	query = filter.Query(RelIdx(0, parent3))
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, n, cnt)
}
