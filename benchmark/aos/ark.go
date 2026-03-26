package main

import (
	"testing"

	"github.com/mlange-42/ark/ecs"
)

func ark32Byte(b *testing.B, n int) {
	world := ecs.NewWorld()

	builder := ecs.NewMap2[Position, Velocity](world)
	builder.NewBatch(n, &Position{0, 0}, &Velocity{1, 1})

	filter := ecs.NewFilter2[Position, Velocity](world)

	loop := func() {
		query := filter.Query()
		for query.NextTable() {
			positions, velocities := query.GetColumns()
			for i := range positions {
				pos, vel := &positions[i], &velocities[i]
				pos.X += vel.X
				pos.Y += vel.Y
			}
		}
	}

	for b.Loop() {
		loop()
	}
}

func ark64Byte(b *testing.B, n int) {
	world := ecs.NewWorld()

	builder := ecs.NewMap3[Position, Velocity, Payload32B](world)
	builder.NewBatch(n, &Position{0, 0}, &Velocity{1, 1}, &Payload32B{})

	filter := ecs.NewFilter2[Position, Velocity](world)

	loop := func() {
		query := filter.Query()
		for query.NextTable() {
			positions, velocities := query.GetColumns()
			for i := range positions {
				pos, vel := &positions[i], &velocities[i]
				pos.X += vel.X
				pos.Y += vel.Y
			}
		}
	}

	for b.Loop() {
		loop()
	}
}

func ark128Byte(b *testing.B, n int) {
	world := ecs.NewWorld()

	builder := ecs.NewMap4[Position, Velocity, Payload32B, Payload64B](world)
	builder.NewBatch(n, &Position{0, 0}, &Velocity{1, 1}, &Payload32B{}, &Payload64B{})

	filter := ecs.NewFilter2[Position, Velocity](world)

	loop := func() {
		query := filter.Query()
		for query.NextTable() {
			positions, velocities := query.GetColumns()
			for i := range positions {
				pos, vel := &positions[i], &velocities[i]
				pos.X += vel.X
				pos.Y += vel.Y
			}
		}
	}

	for b.Loop() {
		loop()
	}
}

func ark256Byte(b *testing.B, n int) {
	world := ecs.NewWorld()

	builder := ecs.NewMap5[Position, Velocity, Payload32B, Payload64B, Payload128B](world)
	builder.NewBatch(n, &Position{0, 0}, &Velocity{1, 1}, &Payload32B{}, &Payload64B{}, &Payload128B{})

	filter := ecs.NewFilter2[Position, Velocity](world)

	loop := func() {
		query := filter.Query()
		for query.NextTable() {
			positions, velocities := query.GetColumns()
			for i := range positions {
				pos, vel := &positions[i], &velocities[i]
				pos.X += vel.X
				pos.Y += vel.Y
			}
		}
	}

	for b.Loop() {
		loop()
	}
}
