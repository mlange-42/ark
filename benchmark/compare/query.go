package main

import (
	"testing"

	"github.com/mlange-42/ark/ecs"
)

func queryCreateClose(b *testing.B) {
	w := ecs.NewWorld()

	builder := ecs.NewMap2[Position, Velocity](w)
	builder.NewBatch(1000, &Position{}, &Velocity{1, 1})

	filter := ecs.NewFilter2[Position, Velocity](w)
	query := filter.Query()
	query.Close()

	// Wrapper to allow inlining, for more realistic results.
	loop := func() {
		query := filter.Query()
		query.Close()
	}

	for b.Loop() {
		loop()
	}
}

func queryCreateCloseRegistered(b *testing.B) {
	w := ecs.NewWorld()

	builder := ecs.NewMap2[Position, Velocity](w)
	builder.NewBatch(1000, &Position{}, &Velocity{1, 1})

	filter := ecs.NewFilter2[Position, Velocity](w).Register()
	query := filter.Query()
	query.Close()

	// Wrapper to allow inlining, for more realistic results.
	loop := func() {
		query := filter.Query()
		query.Close()
	}

	for b.Loop() {
		loop()
	}
}
