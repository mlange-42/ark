package ecs

import "testing"

func BenchmarkQueryPosVel_1000(b *testing.B) {
	n := 1000
	world := NewWorld(128)

	mapper := NewMap2[Position, Velocity](&world)

	for range n {
		_ = mapper.NewEntity(&Position{}, &Velocity{X: 1, Y: 0})
	}

	filter := NewFilter2[Position, Velocity](&world)
	for b.Loop() {
		query := filter.Query()
		for query.Next() {
			pos, vel := query.Get()
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}
