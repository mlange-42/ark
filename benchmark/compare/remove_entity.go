package main

import (
	"testing"

	"github.com/mlange-42/ark/ecs"
)

func removeEntity2Comp10(b *testing.B) {
	removeEntity2Comp(b, 10)
}

func removeEntity2Comp1000(b *testing.B) {
	removeEntity2Comp(b, 1000)
}

func removeEntity2Comp(b *testing.B, n int) {
	world := ecs.NewWorld()

	mapper := ecs.NewMap2[Position, Velocity](world)
	filter := ecs.NewFilter2[Position, Velocity](world)

	mapper.NewBatchFn(n, nil)
	world.RemoveEntities(filter.Batch(), nil)

	entities := make([]ecs.Entity, 0, n)

	for b.Loop() {
		b.StopTimer()
		entities = entities[:0]
		for range n {
			e := mapper.NewEntityFn(nil)
			// Just for fairness, because the others need to do that, too.
			entities = append(entities, e)
		}
		b.StartTimer()

		for i := len(entities) - 1; i >= 0; i-- {
			world.RemoveEntity(entities[i])
		}
	}
}

func removeEntity5Comp10(b *testing.B) {
	removeEntity5Comp(b, 10)
}

func removeEntity5Comp1000(b *testing.B) {
	removeEntity5Comp(b, 1000)
}

func removeEntity5Comp(b *testing.B, n int) {
	world := ecs.NewWorld()

	mapper := ecs.NewMap5[Comp1, Comp2, Comp3, Comp4, Comp5](world)
	filter := ecs.NewFilter5[Comp1, Comp2, Comp3, Comp4, Comp5](world)

	mapper.NewBatchFn(n, nil)
	world.RemoveEntities(filter.Batch(), nil)

	entities := make([]ecs.Entity, 0, n)

	for b.Loop() {
		b.StopTimer()
		entities = entities[:0]
		for range n {
			e := mapper.NewEntityFn(nil)
			// Just for fairness, because the others need to do that, too.
			entities = append(entities, e)
		}
		b.StartTimer()

		for i := len(entities) - 1; i >= 0; i-- {
			world.RemoveEntity(entities[i])
		}
	}
}
