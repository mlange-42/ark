package main

/*
import (
	"testing"

	"github.com/mlange-42/ark/benchmark"
	"github.com/mlange-42/ark/ecs"
)

func benchesComponentsBatch() []benchmark.Benchmark {
	return []benchmark.Benchmark{
		{Name: "Batch.Add 1 Comp", Desc: "1000, memory already allocated", F: componentsBatchAdd1_1000, N: 1000},
		{Name: "Batch.Add 5 Comps", Desc: "1000, memory already allocated", F: componentsBatchAdd5_1000, N: 1000},
		{Name: "Batch.Add 1 to 5 Comps", Desc: "1000, memory already allocated", F: componentsBatchAdd1to5_1000, N: 1000},

		{Name: "Batch.Remove 1 Comp", Desc: "1000, memory already allocated", F: componentsBatchRemove1_1000, N: 1000},
		{Name: "Batch.Remove 5 Comps", Desc: "1000, memory already allocated", F: componentsBatchRemove5_1000, N: 1000},
		{Name: "Batch.Remove 1 of 5 Comps", Desc: "1000, memory already allocated", F: componentsBatchRemove1of5_1000, N: 1000},

		{Name: "Batch.Exchange 1 Comp", Desc: "1000, memory already allocated", F: componentsBatchExchange1_1000, N: 1000},
		{Name: "Batch.Exchange 1 of 5 Comps", Desc: "1000, memory already allocated", F: componentsBatchExchange1of5_1000, N: 1000},
	}
}

func componentsBatchAdd1_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	filter := ecs.All(id1)

	w.Batch().New(1000)

	// Run once to allocate memory
	w.Batch().Add(ecs.All(), id1)
	w.Batch().Remove(filter, id1)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		w.Batch().Add(ecs.All(), id1)
		b.StopTimer()
		w.Batch().Remove(filter, id1)
	}
}

func componentsBatchAdd5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)
	ids := []ecs.ID{id1, id2, id3, id4, id5}
	filter := ecs.All(ids...)

	w.Batch().New(1000)

	// Run once to allocate memory
	w.Batch().Add(ecs.All(), ids...)
	w.Batch().Remove(filter, ids...)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		w.Batch().Add(ecs.All(), ids...)
		b.StopTimer()
		w.Batch().Remove(filter, ids...)
	}
}

func componentsBatchAdd1to5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)
	id6 := ecs.ComponentID[comp6](&w)
	filter := ecs.All(id1)

	w.Batch().New(1000, id2, id3, id4, id5, id6)

	// Run once to allocate memory
	w.Batch().Add(ecs.All(), id1)
	w.Batch().Remove(filter, id1)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		w.Batch().Add(ecs.All(), id1)
		b.StopTimer()
		w.Batch().Remove(filter, id1)
	}
}

func componentsBatchRemove1_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	filter := ecs.All(id1)

	w.Batch().New(1000, id1)

	// Run once to allocate memory
	w.Batch().Remove(filter, id1)
	w.Batch().Add(ecs.All(), id1)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		w.Batch().Remove(filter, id1)
		b.StopTimer()
		w.Batch().Add(ecs.All(), id1)
	}
}

func componentsBatchRemove5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)
	ids := []ecs.ID{id1, id2, id3, id4, id5}
	filter := ecs.All(ids...)

	w.Batch().New(1000, ids...)

	// Run once to allocate memory
	w.Batch().Remove(filter, ids...)
	w.Batch().Add(ecs.All(), ids...)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		w.Batch().Remove(filter, ids...)
		b.StopTimer()
		w.Batch().Add(ecs.All(), ids...)
	}
}

func componentsBatchRemove1of5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)
	filter := ecs.All(id1)

	w.Batch().New(1000, id1, id2, id3, id4, id5)

	// Run once to allocate memory
	w.Batch().Remove(filter, id1)
	w.Batch().Add(ecs.All(), id1)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		w.Batch().Remove(filter, id1)
		b.StopTimer()
		w.Batch().Add(ecs.All(), id1)
	}
}

func componentsBatchExchange1_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	filter1 := ecs.All(id1)
	filter2 := ecs.All(id2)
	ex1 := []ecs.ID{id1}
	ex2 := []ecs.ID{id2}

	w.Batch().New(1000, id1)

	// Run once to allocate memory
	w.Batch().Exchange(filter1, ex2, ex1)
	w.Batch().Exchange(filter2, ex1, ex2)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		w.Batch().Exchange(filter1, ex2, ex1)
		b.StopTimer()
		w.Batch().Exchange(filter2, ex1, ex2)
	}
}

func componentsBatchExchange1of5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)

	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)
	id6 := ecs.ComponentID[comp6](&w)

	filter1 := ecs.All(id1)
	filter2 := ecs.All(id2)
	ex1 := []ecs.ID{id1}
	ex2 := []ecs.ID{id2}

	w.Batch().New(1000, id1, id3, id4, id5, id6)

	// Run once to allocate memory
	w.Batch().Exchange(filter1, ex2, ex1)
	w.Batch().Exchange(filter2, ex1, ex2)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		w.Batch().Exchange(filter1, ex2, ex1)
		b.StopTimer()
		w.Batch().Exchange(filter2, ex1, ex2)
	}
}
*/
