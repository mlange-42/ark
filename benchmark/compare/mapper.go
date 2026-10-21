package main

import (
	"math/rand/v2"
	"runtime"
	"testing"

	"github.com/mlange-42/ark/ecs"
)

func mapperGet10(b *testing.B) {
	mapperGet(b, 10)
}

func mapperGet1000(b *testing.B) {
	mapperGet(b, 1000)
}

func mapperGet100000(b *testing.B) {
	mapperGet(b, 100000)
}

func mapperGet(b *testing.B, n int) {
	w := ecs.NewWorld()

	entities := make([]ecs.Entity, 0, n)
	builder := ecs.NewMap2[Position, Velocity](w)
	mapper := ecs.NewMap1[Position](w)

	builder.NewBatchFn(n, func(e ecs.Entity, p *Position, _ *Velocity) {
		p.X = 1
		entities = append(entities, e)
	})
	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var pos *Position

	// Wrapper to allow inlining, for more realistic results.
	loop := func() {
		for _, e := range entities {
			pos = mapper.Get(e)
		}
	}

	loop()

	for b.Loop() {
		loop()
	}

	runtime.KeepAlive(pos)
}
