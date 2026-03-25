package main

import (
	"testing"

	"github.com/mlange-42/ark/ecs"
)

func newBatch2Comp10(b *testing.B) {
	newBatch2Comp(b, 10)
}

func newBatch2Comp1000(b *testing.B) {
	newBatch2Comp(b, 1000)
}

func newBatch2Comp100000(b *testing.B) {
	newBatch2Comp(b, 100000)
}

func newBatch2Comp(b *testing.B, n int) {
	world := ecs.NewWorld()

	mapper := ecs.NewMap2[Position, Velocity](world)
	mapper.NewBatchFn(n, nil)

	filter := ecs.NewFilter2[Position, Velocity](world)
	world.RemoveEntities(filter.Batch(), nil)

	for b.Loop() {
		mapper.NewBatchFn(n, nil)
		b.StopTimer()
		world.RemoveEntities(filter.Batch(), nil)
		b.StartTimer()
	}
}

func newBatch5Comp10(b *testing.B) {
	newBatch5Comp(b, 10)
}

func newBatch5Comp1000(b *testing.B) {
	newBatch5Comp(b, 1000)
}

func newBatch5Comp100000(b *testing.B) {
	newBatch5Comp(b, 100000)
}

func newBatch5Comp(b *testing.B, n int) {
	world := ecs.NewWorld()

	mapper := ecs.NewMap5[Comp1, Comp2, Comp3, Comp4, Comp5](world)
	mapper.NewBatchFn(n, nil)

	filter := ecs.NewFilter5[Comp1, Comp2, Comp3, Comp4, Comp5](world)
	world.RemoveEntities(filter.Batch(), nil)

	for b.Loop() {
		mapper.NewBatchFn(n, nil)
		b.StopTimer()
		world.RemoveEntities(filter.Batch(), nil)
		b.StartTimer()
	}
}
