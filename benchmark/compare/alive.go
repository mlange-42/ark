package main

import (
	"math/rand/v2"
	"runtime"
	"testing"

	"github.com/mlange-42/ark/ecs"
)

func entityAlive10(b *testing.B) {
	entityAlive(b, 10)
}

func entityAlive1000(b *testing.B) {
	entityAlive(b, 1000)
}

func entityAlive(b *testing.B, n int) {
	w := ecs.NewWorld()

	entities := make([]ecs.Entity, 0, n)
	filter := ecs.NewFilter2[Position, Velocity](w)
	builder := ecs.NewMap2[Position, Velocity](w)

	builder.NewBatchFn(n/2, func(e ecs.Entity, _ *Position, _ *Velocity) {
		entities = append(entities, e)
	})
	w.RemoveEntities(filter.Batch(), nil)
	builder.NewBatchFn(n/2, func(e ecs.Entity, _ *Position, _ *Velocity) {
		entities = append(entities, e)
	})

	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	test := false

	// Wrapper to allow inlining, for more realistic results.
	loop := func() {
		for _, e := range entities {
			test = test || w.Alive(e)
		}
	}

	loop()

	for b.Loop() {
		loop()
	}

	runtime.KeepAlive(test)
}
