package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery2(t *testing.T) {
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

	filter := NewFilter2[Position, Velocity](&w).Build()
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

	filter := NewFilter2[Position, Velocity](&w).Build()
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

func TestQuery2Closed(t *testing.T) {
	w := NewWorld(4)
	mapper := NewMap2[Position, Velocity](&w)
	for range 10 {
		e1 := w.NewEntity()
		mapper.Add(e1, &Position{}, &Velocity{})
	}

	filter := NewFilter2[Position, Velocity](&w).Build()
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

func BenchmarkQueryPosVel_1000(b *testing.B) {
	n := 1000
	world := NewWorld(128)

	mapper := NewMap2[Position, Velocity](&world)

	for range n {
		_ = mapper.NewEntity(&Position{}, &Velocity{X: 1, Y: 0})
	}

	filter := NewFilter2[Position, Velocity](&world).Build()
	for b.Loop() {
		query := filter.Query()
		for query.Next() {
			pos, vel := query.Get()
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}
