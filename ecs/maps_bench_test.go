package ecs

import "testing"

func BenchmarkMapPosVel_1000(b *testing.B) {
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

func BenchmarkMapPosVel_1000_Unchecked(b *testing.B) {
	n := 1000
	world := NewWorld(1024)

	mapper := NewMap2[Position, Velocity](&world)

	entities := make([]Entity, 0, n)
	for range n {
		e := mapper.NewEntity(&Position{}, &Velocity{X: 1, Y: 0})
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
