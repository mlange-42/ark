package main

/*
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

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		w.Batch().New(1000)
		b.StopTimer()
		w.Batch().RemoveEntities(ecs.All())
	}
}

func entitiesBatchCreate_1Comp_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	ids := []ecs.ID{
		ecs.ComponentID[comp1](&w),
	}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		w.Batch().New(1000, ids...)
		b.StopTimer()
		w.Batch().RemoveEntities(ecs.All())
	}
}

func entitiesBatchCreate_5Comp_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)
	ids := []ecs.ID{id1, id2, id3, id4, id5}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		w.Batch().New(1000, ids...)
		b.StopTimer()
		w.Batch().RemoveEntities(ecs.All())
	}
}

func entitiesBatchRemove_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	builder := ecs.NewBuilder(&w)

	entities := make([]ecs.Entity, 0, 1000)

	for i := 0; i < b.N; i++ {
		builder.NewBatch(1000)
		b.StartTimer()
		w.Batch().RemoveEntities(ecs.All())
		b.StopTimer()
		entities = entities[:0]
	}
}

func entitiesBatchRemove_1Comp_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	builder := ecs.NewBuilder(&w, id1)

	entities := make([]ecs.Entity, 0, 1000)

	for i := 0; i < b.N; i++ {
		builder.NewBatch(1000)
		b.StartTimer()
		w.Batch().RemoveEntities(ecs.All())
		b.StopTimer()
		entities = entities[:0]
	}
}

func entitiesBatchRemove_5Comp_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)
	ids := []ecs.ID{id1, id2, id3, id4, id5}
	builder := ecs.NewBuilder(&w, ids...)

	entities := make([]ecs.Entity, 0, 1000)

	for i := 0; i < b.N; i++ {
		builder.NewBatch(1000)
		b.StartTimer()
		w.Batch().RemoveEntities(ecs.All())
		b.StopTimer()
		entities = entities[:0]
	}
}
*/
