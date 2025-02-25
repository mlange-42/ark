package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery2Old(t *testing.T) {
	w := NewWorld(4)

	posMap := NewMap[Position](&w)
	velMap := NewMap[Velocity](&w)
	headMap := NewMap[Heading](&w)

	for range 10 {
		e1 := w.NewEntity()
		posMap.Add(e1, &Position{})

		e2 := w.NewEntity()
		posMap.Add(e2, &Position{X: 100, Y: 0})
		velMap.Add(e2, &Velocity{})

		e3 := w.NewEntity()
		posMap.Add(e3, &Position{X: 100, Y: 0})
		velMap.Add(e3, &Velocity{})
		headMap.Add(e3, &Heading{})
	}

	filter := NewFilter2[Position, Velocity](&w)
	query := filter.Query()

	cnt := 0
	for query.Next() {
		pos, vel := query.Get()
		assert.Equal(t, 100.0, pos.X)
		vel.X = float64(cnt)
		cnt++
	}

	assert.Equal(t, 20, cnt)

	query = filter.Query()
	cnt = 0
	for query.Next() {
		pos, vel := query.Get()
		pos.X += vel.X * 2
		cnt++
	}

	query = filter.Query()
	cnt = 0
	for query.Next() {
		pos, vel := query.Get()
		assert.Equal(t, float64(cnt), vel.X)
		assert.Equal(t, float64(cnt)*2+100, pos.X)
		cnt++
	}
}

func TestQuery2Empty(t *testing.T) {
	w := NewWorld(4)

	posMap := NewMap[Position](&w)

	for range 10 {
		e1 := w.NewEntity()
		posMap.Add(e1, &Position{})
	}

	filter := NewFilter2[Position, Velocity](&w)
	query := filter.Query()

	cnt := 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)

	assert.Panics(t, func() { query.Get() })
	assert.Panics(t, func() { query.Entity() })
	assert.Panics(t, func() { query.Next() })
}

func TestQuery2Advanced(t *testing.T) {
	w := NewWorld(4)

	mapPosVel := NewMap2[Position, Velocity](&w)
	mapPosVelHead := NewMap3[Position, Velocity, Heading](&w)
	mapAll := NewMap4[Position, Velocity, Heading, CompA](&w)

	for range 10 {
		_ = mapPosVel.NewEntity(&Position{}, &Velocity{})
		_ = mapPosVelHead.NewEntity(&Position{}, &Velocity{}, &Heading{})
		_ = mapAll.NewEntity(&Position{}, &Velocity{}, &Heading{}, &CompA{})
	}

	filter := NewFilter2[Position, Velocity](&w).
		With(C[Heading]()).
		Without(C[CompA]())

	query := filter.Query()
	cnt := 0
	for query.Next() {
		cnt++
		assert.True(t, mapPosVelHead.HasAll(query.Entity()))
		assert.False(t, mapAll.HasAll(query.Entity()))
	}
	assert.Equal(t, 10, cnt)

	filter = NewFilter2[Position, Velocity](&w).Exclusive()
	query = filter.Query()
	cnt = 0
	for query.Next() {
		cnt++
		assert.True(t, mapPosVel.HasAll(query.Entity()))
		assert.False(t, mapPosVelHead.HasAll(query.Entity()))
		assert.False(t, mapAll.HasAll(query.Entity()))
	}
	assert.Equal(t, 10, cnt)
}

func TestQuery2Closed(t *testing.T) {
	w := NewWorld(4)
	mapper := NewMap2[Position, Velocity](&w)
	for range 10 {
		e1 := w.NewEntity()
		mapper.Add(e1, &Position{}, &Velocity{})
	}

	filter := NewFilter2[Position, Velocity](&w)
	query := filter.Query()

	cnt := 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 10, cnt)

	assert.Panics(t, func() { query.Get() })
	assert.Panics(t, func() { query.Entity() })
	assert.Panics(t, func() { query.Next() })
}
