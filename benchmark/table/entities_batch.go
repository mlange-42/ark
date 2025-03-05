package main

import (
	"testing"

	"github.com/mlange-42/ark/benchmark"
	"github.com/mlange-42/ark/ecs"
)

func benchesEntitiesBatch() []benchmark.Benchmark {
	return []benchmark.Benchmark{
		{Name: "World.NewEntities", Desc: "1000, memory already allocated", F: entitiesBatchCreate1000, N: 1000},
		{Name: "Map1.NewEntity w/ 1 Comp", Desc: "1000, memory already allocated", F: entitiesBatchCreate1Comp1000, N: 1000},
		{Name: "Map5.NewEntity w/ 5 Comps", Desc: "1000, memory already allocated", F: entitiesBatchCreate5Comp1000, N: 1000},

		{Name: "World.RemoveEntities", Desc: "1000", F: entitiesBatchRemove1000, N: 1000},
		{Name: "World.RemoveEntities w/ 1 Comp", Desc: "1000", F: entitiesBatchRemove1Comp1000, N: 1000},
		{Name: "World.RemoveEntities w/ 5 Comps", Desc: "1000", F: entitiesBatchRemove5Comp1000, N: 1000},
	}
}

func entitiesBatchCreate1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	filter := ecs.NewFilter0(&w)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		w.NewEntities(1000, nil)
		b.StopTimer()
		w.RemoveEntities(filter.Batch(), nil)
	}
}

func entitiesBatchCreate1Comp1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	builder := ecs.NewMap1[comp1](&w)
	filter := ecs.NewFilter0(&w)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		builder.NewBatchFn(1000, nil)
		b.StopTimer()
		w.RemoveEntities(filter.Batch(), nil)
	}
}

func entitiesBatchCreate5Comp1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	builder := ecs.NewMap5[comp1, comp2, comp3, comp4, comp5](&w)
	filter := ecs.NewFilter0(&w)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		builder.NewBatchFn(1000, nil)
		b.StopTimer()
		w.RemoveEntities(filter.Batch(), nil)
	}
}

func entitiesBatchRemove1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	filter := ecs.NewFilter0(&w)

	for i := 0; i < b.N; i++ {
		w.NewEntities(1000, nil)
		b.StartTimer()
		w.RemoveEntities(filter.Batch(), nil)
		b.StopTimer()
	}
}

func entitiesBatchRemove1Comp1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	builder := ecs.NewMap1[comp1](&w)
	filter := ecs.NewFilter0(&w)

	for i := 0; i < b.N; i++ {
		builder.NewBatchFn(1000, nil)
		b.StartTimer()
		w.RemoveEntities(filter.Batch(), nil)
		b.StopTimer()
	}
}

func entitiesBatchRemove5Comp1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	builder := ecs.NewMap5[comp1, comp2, comp3, comp4, comp5](&w)
	filter := ecs.NewFilter0(&w)

	for i := 0; i < b.N; i++ {
		builder.NewBatchFn(1000, nil)
		b.StartTimer()
		w.RemoveEntities(filter.Batch(), nil)
		b.StopTimer()
	}
}
