package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery2(t *testing.T) {
	w := NewWorld(1024)

	posMap := NewMap[Position](&w)
	velMap := NewMap[Velocity](&w)
	headMap := NewMap[Heading](&w)

	for range 10 {
		e1 := w.NewEntity()
		posMap.Add(e1)

		e2 := w.NewEntity()
		posMap.Add(e2)
		velMap.Add(e2)

		e3 := w.NewEntity()
		posMap.Add(e3)
		velMap.Add(e3)
		headMap.Add(e3)
	}

	query := NewQuery2[Position, Velocity](&w)

	cnt := 0
	for query.Next() {
		_, vel := query.Get()
		vel.X = float64(cnt)
		cnt++
	}

	assert.Equal(t, 20, cnt)

	cnt = 0
	for query.Next() {
		pos, vel := query.Get()
		pos.X += vel.X * 2
		cnt++
	}

	cnt = 0
	for query.Next() {
		pos, vel := query.Get()
		assert.Equal(t, float64(cnt), vel.X)
		assert.Equal(t, float64(cnt)*2, pos.X)
		cnt++
	}
}

func BenchmarkQuery2(b *testing.B) {
	n := 100000
	world := NewWorld(1024)

	posMap := NewMap[Position](&world)
	velMap := NewMap[Velocity](&world)

	for range n {
		e := world.NewEntity()
		posMap.Add(e)
		velMap.Add(e)
	}

	query := NewQuery2[Position, Velocity](&world)
	for b.Loop() {
		for query.Next() {
			pos, vel := query.Get()
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}
