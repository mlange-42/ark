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
	assert.Panics(t, func() {
		childMap.SetRelation(e, deadParent)
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
