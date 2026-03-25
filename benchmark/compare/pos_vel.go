package main

import (
	"math/rand/v2"
	"testing"

	"github.com/mlange-42/ark/ecs"
)

func posVelQuery10(b *testing.B) {
	posVelQuery(b, 10)
}

func posVelQuery1000(b *testing.B) {
	posVelQuery(b, 1000)
}

func posVelQuery100000(b *testing.B) {
	posVelQuery(b, 100000)
}

func posVelQuery(b *testing.B, n int) {
	w := ecs.NewWorld()

	builder := ecs.NewMap2[Position, Velocity](w)
	builder.NewBatch(n, &Position{}, &Velocity{1, 1})

	filter := ecs.NewFilter2[Position, Velocity](w)

	// Wrapper to allow inlining, for more realistic results.
	loop := func() {
		query := filter.Query()
		for query.Next() {
			pos, vel := query.Get()
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}

	for b.Loop() {
		loop()
	}
}

func posVelMap10(b *testing.B) {
	posVelMap(b, 10)
}

func posVelMap1000(b *testing.B) {
	posVelMap(b, 1000)
}

func posVelMap100000(b *testing.B) {
	posVelMap(b, 100000)
}

func posVelMap(b *testing.B, n int) {
	w := ecs.NewWorld()

	builder := ecs.NewMap2[Position, Velocity](w)
	builder.NewBatch(n, &Position{}, &Velocity{1, 1})

	entities := make([]ecs.Entity, 0, n)

	filter := ecs.NewFilter2[Position, Velocity](w)
	query := filter.Query()
	for query.Next() {
		entities = append(entities, query.Entity())
	}
	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	// Wrapper to allow inlining, for more realistic results.
	loop := func() {
		for _, e := range entities {
			pos, vel := builder.Get(e)
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}

	for b.Loop() {
		loop()
	}
}
