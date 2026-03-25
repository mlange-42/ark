package main

import (
	"testing"

	"github.com/mlange-42/ark/ecs"
)

func removeEntityBatch2Comp10(b *testing.B) {
	removeEntityBatch2Comp(b, 10)
}

func removeEntityBatch2Comp1000(b *testing.B) {
	removeEntityBatch2Comp(b, 1000)
}

func removeEntityBatch2Comp(b *testing.B, n int) {
	world := ecs.NewWorld()

	mapper := ecs.NewMap2[Position, Velocity](world)
	filter := ecs.NewFilter2[Position, Velocity](world)

	mapper.NewBatchFn(n, nil)
	world.RemoveEntities(filter.Batch(), nil)

	for b.Loop() {
		b.StopTimer()
		mapper.NewBatchFn(n, nil)
		b.StartTimer()
		world.RemoveEntities(filter.Batch(), nil)
	}
}

func removeEntityBatch5Comp10(b *testing.B) {
	removeEntityBatch5Comp(b, 10)
}

func removeEntityBatch5Comp1000(b *testing.B) {
	removeEntityBatch5Comp(b, 1000)
}

func removeEntityBatch5Comp(b *testing.B, n int) {
	world := ecs.NewWorld()

	mapper := ecs.NewMap5[Comp1, Comp2, Comp3, Comp4, Comp5](world)
	filter := ecs.NewFilter5[Comp1, Comp2, Comp3, Comp4, Comp5](world)

	mapper.NewBatchFn(n, nil)
	world.RemoveEntities(filter.Batch(), nil)

	for b.Loop() {
		b.StopTimer()
		mapper.NewBatchFn(n, nil)
		b.StartTimer()
		world.RemoveEntities(filter.Batch(), nil)
	}
}
