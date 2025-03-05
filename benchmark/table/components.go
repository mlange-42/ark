package main

/*
import (
	"testing"

	"github.com/mlange-42/ark/benchmark"
	"github.com/mlange-42/ark/ecs"
	"github.com/mlange-42/ark/generic"
)

func benchesComponents() []benchmark.Benchmark {
	return []benchmark.Benchmark{
		{Name: "World.Add 1 Comp", Desc: "memory already alloc.", F: componentsAdd1_1000, N: 1000},
		{Name: "World.Add 5 Comps", Desc: "memory already alloc.", F: componentsAdd5_1000, N: 1000},
		{Name: "World.Add 1 to 5 Comps", Desc: "memory already alloc.", F: componentsAdd1to5_1000, N: 1000},

		{Name: "World.Remove 1 Comp", Desc: "memory already alloc.", F: componentsRemove1_1000, N: 1000},
		{Name: "World.Remove 5 Comps", Desc: "memory already alloc.", F: componentsRemove5_1000, N: 1000},
		{Name: "World.Remove 1 of 5 Comps", Desc: "memory already alloc.", F: componentsRemove1of5_1000, N: 1000},

		{Name: "World.Exchange 1 Comp", Desc: "memory already alloc.", F: componentsExchange1_1000, N: 1000},
		{Name: "World.Exchange 1 of 5 Comps", Desc: "memory already alloc.", F: componentsExchange1of5_1000, N: 1000},

		{Name: "Map1.Assign 1 Comps", Desc: "memory already alloc.", F: componentsAssignGeneric1_1000, N: 1000},
		{Name: "Map5.Assign 5 Comps", Desc: "memory already alloc.", F: componentsAssignGeneric5_1000, N: 1000},
	}
}

func componentsAdd1_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	ids := []ecs.ID{id1}
	filter := ecs.All(id1)

	query := w.Batch().NewQ(1000)

	entities := make([]ecs.Entity, 0, 1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	// Run once to allocate memory
	for _, e := range entities {
		w.Add(e, ids...)
	}
	w.Batch().Remove(filter, ids...)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			w.Add(e, ids...)
		}
		b.StopTimer()
		w.Batch().Remove(filter, ids...)
	}
}

func componentsAdd5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)

	ids := []ecs.ID{id1, id2, id3, id4, id5}
	filter := ecs.All(ids...)

	query := w.Batch().NewQ(1000)

	entities := make([]ecs.Entity, 0, 1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	// Run once to allocate memory
	for _, e := range entities {
		w.Add(e, ids...)
	}
	w.Batch().Remove(filter, ids...)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			w.Add(e, ids...)
		}
		b.StopTimer()
		w.Batch().Remove(filter, ids...)
	}
}

func componentsAdd1to5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)
	id6 := ecs.ComponentID[comp6](&w)
	ids := []ecs.ID{id1}

	filter := ecs.All(id1)

	query := w.Batch().NewQ(1000, id2, id3, id4, id5, id6)

	entities := make([]ecs.Entity, 0, 1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	// Run once to allocate memory
	for _, e := range entities {
		w.Add(e, ids...)
	}
	w.Batch().Remove(filter, ids...)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			w.Add(e, ids...)
		}
		b.StopTimer()
		w.Batch().Remove(filter, ids...)
	}
}

func componentsRemove1_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	ids := []ecs.ID{id1}
	filter := ecs.All()

	query := w.Batch().NewQ(1000, id1)

	entities := make([]ecs.Entity, 0, 1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	// Run once to allocate memory
	for _, e := range entities {
		w.Remove(e, ids...)
	}
	w.Batch().Add(filter, ids...)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			w.Remove(e, ids...)
		}
		b.StopTimer()
		w.Batch().Add(filter, ids...)
	}
}

func componentsRemove5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)
	ids := []ecs.ID{id1, id2, id3, id4, id5}
	filter := ecs.All()

	query := w.Batch().NewQ(1000, ids...)

	entities := make([]ecs.Entity, 0, 1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	// Run once to allocate memory
	for _, e := range entities {
		w.Remove(e, ids...)
	}
	w.Batch().Add(filter, ids...)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			w.Remove(e, ids...)
		}
		b.StopTimer()
		w.Batch().Add(filter, ids...)
	}
}

func componentsRemove1of5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)
	filter := ecs.All()
	ids := []ecs.ID{id1}

	query := w.Batch().NewQ(1000, id1, id2, id3, id4, id5)

	entities := make([]ecs.Entity, 0, 1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	// Run once to allocate memory
	for _, e := range entities {
		w.Remove(e, ids...)
	}
	w.Batch().Add(filter, id1)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			w.Remove(e, ids...)
		}
		b.StopTimer()
		w.Batch().Add(filter, ids...)
	}
}

func componentsExchange1_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	filter := ecs.All(id2)
	ex1 := []ecs.ID{id1}
	ex2 := []ecs.ID{id2}

	query := w.Batch().NewQ(1000, id1)

	entities := make([]ecs.Entity, 0, 1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	// Run once to allocate memory
	for _, e := range entities {
		w.Exchange(e, ex2, ex1)
	}
	w.Batch().Exchange(filter, ex1, ex2)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			w.Exchange(e, ex2, ex1)
		}
		b.StopTimer()
		w.Batch().Exchange(filter, ex1, ex2)
	}
}

func componentsExchange1of5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)

	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)
	id6 := ecs.ComponentID[comp6](&w)

	filter := ecs.All(id2)
	ex1 := []ecs.ID{id1}
	ex2 := []ecs.ID{id2}

	query := w.Batch().NewQ(1000, id1, id3, id4, id5, id6)

	entities := make([]ecs.Entity, 0, 1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	// Run once to allocate memory
	for _, e := range entities {
		w.Exchange(e, ex2, ex1)
	}
	w.Batch().Exchange(filter, ex1, ex2)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			w.Exchange(e, ex2, ex1)
		}
		b.StopTimer()
		w.Batch().Exchange(filter, ex1, ex2)
	}
}

func componentsAssignGeneric1_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)

	filter := ecs.All(id1)

	query := w.Batch().NewQ(1000)

	entities := make([]ecs.Entity, 0, 1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	c1 := comp1{}

	mapper := generic.NewMap1[comp1](&w)

	// Run once to allocate memory
	for _, e := range entities {
		mapper.Assign(e, &c1)
	}
	w.Batch().Remove(filter, id1)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			mapper.Assign(e, &c1)
		}
		b.StopTimer()
		w.Batch().Remove(filter, id1)
	}
}

func componentsAssignGeneric5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)

	ids := []ecs.ID{id1, id2, id3, id4, id5}
	filter := ecs.All(ids...)

	query := w.Batch().NewQ(1000)

	entities := make([]ecs.Entity, 0, 1000)
	for query.Next() {
		entities = append(entities, query.Entity())
	}

	c1 := comp1{}
	c2 := comp2{}
	c3 := comp3{}
	c4 := comp4{}
	c5 := comp5{}

	mapper := generic.NewMap5[comp1, comp2, comp3, comp4, comp5](&w)

	// Run once to allocate memory
	for _, e := range entities {
		mapper.Assign(e, &c1, &c2, &c3, &c4, &c5)
	}
	w.Batch().Remove(filter, ids...)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			mapper.Assign(e, &c1, &c2, &c3, &c4, &c5)
		}
		b.StopTimer()
		w.Batch().Remove(filter, ids...)
	}
}
*/
