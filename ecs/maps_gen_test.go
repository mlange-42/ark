package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap2(t *testing.T) {
	w := NewWorld(4)

	mapper := NewMap2[Position, Velocity](&w)

	entities := []Entity{}
	for i := range 12 {
		v := float64(i)
		e := w.NewEntity()
		mapper.Add(e, &Position{v + 1, v + 2}, &Velocity{v + 3, v + 4})
		entities = append(entities, e)
	}

	for i, e := range entities {
		v := float64(i)
		pos, vel := mapper.Get(e)
		assert.Equal(t, Position{v + 1, v + 2}, *pos)
		assert.Equal(t, Velocity{v + 3, v + 4}, *vel)
	}

	for _, e := range entities {
		mapper.Remove(e)
	}
}

func BenchmarkMap2(b *testing.B) {
	n := 1000
	world := NewWorld(1024)

	mapper := NewMap2[Position, Velocity](&world)

	entities := make([]Entity, 0, n)
	for range n {
		e := world.NewEntity()
		mapper.Add(e, &Position{}, &Velocity{X: 1, Y: 0})
		entities = append(entities, e)
	}

	for b.Loop() {
		for _, e := range entities {
			pos, vel := mapper.Get(e)
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}

func BenchmarkMap2Unchecked(b *testing.B) {
	n := 1000
	world := NewWorld(1024)

	mapper := NewMap2[Position, Velocity](&world)

	entities := make([]Entity, 0, n)
	for range n {
		e := world.NewEntity()
		mapper.Add(e, &Position{}, &Velocity{X: 1, Y: 0})
		entities = append(entities, e)
	}

	for b.Loop() {
		for _, e := range entities {
			pos, vel := mapper.GetUnchecked(e)
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}
