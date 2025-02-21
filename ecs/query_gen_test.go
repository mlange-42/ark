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

func BenchmarkQuery2(b *testing.B) {
	n := 1000
	world := NewWorld(1024)

	posMap := NewMap[Position](&world)
	velMap := NewMap[Velocity](&world)

	for range n {
		e := world.NewEntity()
		posMap.Add(e, &Position{})
		velMap.Add(e, &Velocity{X: 1, Y: 0})
	}

	filter := NewFilter2[Position, Velocity](&world).Build()
	query := filter.Query()
	for b.Loop() {
		for query.Next() {
			pos, vel := query.Get()
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}

type iQuery2[A any, B any] interface {
	Next() bool
	Get() (*A, *B)
}

func newIQuery2[A any, B any](world *World) iQuery2[A, B] {
	filter := NewFilter2[A, B](world).Build()
	q := filter.Query()
	return &q
}

func BenchmarkQuery2Interface(b *testing.B) {
	n := 1000
	world := NewWorld(1024)

	posMap := NewMap[Position](&world)
	velMap := NewMap[Velocity](&world)

	for range n {
		e := world.NewEntity()
		posMap.Add(e, &Position{})
		velMap.Add(e, &Velocity{X: 1, Y: 0})
	}

	query := newIQuery2[Position, Velocity](&world)
	for b.Loop() {
		for query.Next() {
			pos, vel := query.Get()
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}
