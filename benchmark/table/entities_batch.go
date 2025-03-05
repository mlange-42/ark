package main

import (
	"testing"

	"github.com/mlange-42/ark/benchmark"
	"github.com/mlange-42/ark/ecs"
)

func benchesEntitiesBatch() []benchmark.Benchmark {
	return []benchmark.Benchmark{
		{Name: "Batch.New", Desc: "1000, memory already allocated", F: entitiesBatchCreate_1000, N: 1000},
		{Name: "Batch.New w/ 1 Comp", Desc: "1000, memory already allocated", F: entitiesBatchCreate_1Comp_1000, N: 1000},
		{Name: "Batch.New w/ 5 Comps", Desc: "1000, memory already allocated", F: entitiesBatchCreate_5Comp_1000, N: 1000},

		{Name: "Batch.RemoveEntities", Desc: "1000", F: entitiesBatchRemove_1000, N: 1000},
		{Name: "Batch.RemoveEntities w/ 1 Comp", Desc: "1000", F: entitiesBatchRemove_1Comp_1000, N: 1000},
		{Name: "Batch.RemoveEntities w/ 5 Comps", Desc: "1000", F: entitiesBatchRemove_5Comp_1000, N: 1000},
	}
}

func entitiesBatchCreate_1000(b *testing.B) {
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

func entitiesBatchCreate_1Comp_1000(b *testing.B) {
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

func entitiesBatchCreate_5Comp_1000(b *testing.B) {
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

func entitiesBatchRemove_1000(b *testing.B) {
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

func entitiesBatchRemove_1Comp_1000(b *testing.B) {
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

func entitiesBatchRemove_5Comp_1000(b *testing.B) {
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
