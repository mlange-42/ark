package main

import (
	"testing"

	"github.com/mlange-42/ark/ecs"
)

func newEntity2Comp10(b *testing.B) {
	newEntity2Comp(b, 10)
}

func newEntity2Comp1000(b *testing.B) {
	newEntity2Comp(b, 1000)
}

func newEntity2Comp100000(b *testing.B) {
	newEntity2Comp(b, 100000)
}

func newEntity2Comp(b *testing.B, n int) {
	world := ecs.NewWorld()

	mapper := ecs.NewMap2[Position, Velocity](world)
	mapper.NewBatchFn(n, nil)

	filter := ecs.NewFilter2[Position, Velocity](world)
	world.RemoveEntities(filter.Batch(), nil)

	entities := make([]ecs.Entity, 0, n)

	for b.Loop() {
		for range n {
			e := mapper.NewEntityFn(nil)
			// Just for fairness, because the others need to do that, too.
			entities = append(entities, e)
		}
		b.StopTimer()

		if n < 64 {
			// Speed up cleanup for low entity counts
			for i := len(entities) - 1; i >= 0; i-- {
				world.RemoveEntity(entities[i])
			}
		} else {
			world.RemoveEntities(filter.Batch(), nil)
		}

		entities = entities[:0]
		b.StartTimer()
	}
}
