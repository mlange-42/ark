package main

import (
	"testing"

	"github.com/mlange-42/ark/benchmark"
	"github.com/mlange-42/ark/ecs"
)

func benchesEntitiesBatch() []benchmark.Benchmark {
	return []benchmark.Benchmark{
		{Name: "World.NewEntities", Desc: "1000, memory already alloc.", F: entitiesBatchCreate1000, N: 1000},
		{Name: "Map1.NewBatchFn w/ 1 Comp", Desc: "1000, memory already alloc.", F: entitiesBatchCreate1Comp1000, N: 1000},
		{Name: "Map5.NewBatchFn w/ 5 Comps", Desc: "1000, memory already alloc.", F: entitiesBatchCreate5Comp1000, N: 1000},

		{Name: "World.RemoveEntities", Desc: "1000", F: entitiesBatchRemove1000, N: 1000},
		{Name: "World.RemoveEntities w/ 1 Comp", Desc: "1000", F: entitiesBatchRemove1Comp1000, N: 1000},
		{Name: "World.RemoveEntities w/ 5 Comps", Desc: "1000", F: entitiesBatchRemove5Comp1000, N: 1000},
	}
}

func entitiesBatchCreate1000(b *testing.B) {
	w := ecs.NewWorld()
	filter := ecs.NewFilter0(w)

	for b.Loop() {
		w.NewEntities(1000, nil)
		b.StopTimer()
		w.RemoveEntities(filter.Batch(), nil)
		b.StartTimer()
	}
}

func entitiesBatchCreate1Comp1000(b *testing.B) {
	w := ecs.NewWorld()
	builder := ecs.NewMap1[comp1](w)
	filter := ecs.NewFilter0(w)

	for b.Loop() {
		builder.NewBatchFn(1000, nil)
		b.StopTimer()
		w.RemoveEntities(filter.Batch(), nil)
		b.StartTimer()
	}
}

func entitiesBatchCreate5Comp1000(b *testing.B) {
	w := ecs.NewWorld()
	builder := ecs.NewMap5[comp1, comp2, comp3, comp4, comp5](w)
	filter := ecs.NewFilter0(w)

	for b.Loop() {
		builder.NewBatchFn(1000, nil)
		b.StopTimer()
		w.RemoveEntities(filter.Batch(), nil)
		b.StartTimer()
	}
}

func entitiesBatchRemove1000(b *testing.B) {
	w := ecs.NewWorld()
	filter := ecs.NewFilter0(w)

	for b.Loop() {
		b.StopTimer()
		w.NewEntities(1000, nil)
		b.StartTimer()
		w.RemoveEntities(filter.Batch(), nil)
	}
}

func entitiesBatchRemove1Comp1000(b *testing.B) {
	w := ecs.NewWorld()
	builder := ecs.NewMap1[comp1](w)
	filter := ecs.NewFilter0(w)

	for b.Loop() {
		b.StopTimer()
		builder.NewBatchFn(1000, nil)
		b.StartTimer()
		w.RemoveEntities(filter.Batch(), nil)
	}
}

func entitiesBatchRemove5Comp1000(b *testing.B) {
	w := ecs.NewWorld()
	builder := ecs.NewMap5[comp1, comp2, comp3, comp4, comp5](w)
	filter := ecs.NewFilter0(w)

	for b.Loop() {
		b.StopTimer()
		builder.NewBatchFn(1000, nil)
		b.StartTimer()
		w.RemoveEntities(filter.Batch(), nil)
	}
}
