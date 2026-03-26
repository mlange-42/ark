package main

import (
	"testing"

	"github.com/mlange-42/ark/ecs"
)

func ark16Byte(b *testing.B, n int) {
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

func ark32Byte(b *testing.B, n int) {
	world := ecs.NewWorld()

	builder := ecs.NewMap4[Position, Velocity, Payload1, Payload2](world)
	builder.NewBatch(n, &Position{0, 0}, &Velocity{1, 1}, &Payload1{}, &Payload2{})

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

	builder := ecs.NewMap8[Position, Velocity, Payload1, Payload2, Payload3, Payload4, Payload5, Payload6](world)
	builder.NewBatch(n, &Position{0, 0}, &Velocity{1, 1}, &Payload1{}, &Payload2{}, &Payload3{}, &Payload4{}, &Payload5{}, &Payload6{})

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
