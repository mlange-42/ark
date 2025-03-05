package main

/*
import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/mlange-42/ark/benchmark"
	"github.com/mlange-42/ark/ecs"
)

func benchesQuery() []benchmark.Benchmark {
	return []benchmark.Benchmark{
		{Name: "Query.Next", Desc: "", F: queryIter_100_000, N: 100_000},
		{Name: "Query.Next + 1x Query.Get", Desc: "", F: queryIterGet_1_100_000, N: 100_000},
		{Name: "Query.Next + 2x Query.Get", Desc: "", F: queryIterGet_2_100_000, N: 100_000},
		{Name: "Query.Next + 5x Query.Get", Desc: "", F: queryIterGet_5_100_000, N: 100_000},

		{Name: "Query.Next + Query.Entity", Desc: "", F: queryIterEntity_100_000, N: 100_000},

		{Name: "Query.Next + Query.Relation", Desc: "", F: queryRelation_100_000, N: 100_000},

		{Name: "Query.EntityAt, 1 arch", Desc: "", F: queryEntityAt_1Arch_1000, N: 1000},
		{Name: "Query.EntityAt, 1 arch", Desc: "registered filter", F: queryEntityAtRegistered_1Arch_1000, N: 1000},
		{Name: "Query.EntityAt, 5 arch", Desc: "", F: queryEntityAt_5Arch_1000, N: 1000},
		{Name: "Query.EntityAt, 5 arch", Desc: "registered filter", F: queryEntityAtRegistered_5Arch_1000, N: 1000},

		{Name: "World.Query", Desc: "", F: queryCreate, N: 1},
		{Name: "World.Query", Desc: "registered filter", F: queryCreateCached, N: 1},
	}
}

func queryIter_100_000(b *testing.B) {
	w := ecs.NewWorld()
	builder := ecs.NewBuilder(&w)
	builder.NewBatch(100_000)
	filter := ecs.All()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := w.Query(&filter)
		b.StartTimer()
		for query.Next() {
		}
	}
}

func queryIterGet_1_100_000(b *testing.B) {
	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)

	builder := ecs.NewBuilder(&w, id1)
	builder.NewBatch(100_000)
	filter := ecs.All(id1)

	var c1 *comp1

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := w.Query(&filter)
		b.StartTimer()
		for query.Next() {
			c1 = (*comp1)(query.Get(id1))
		}
	}
	b.StopTimer()
	sum := c1.V
	_ = sum
}

func queryIterGet_2_100_000(b *testing.B) {
	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)

	builder := ecs.NewBuilder(&w, id1, id2)
	builder.NewBatch(100_000)
	filter := ecs.All(id1, id2)

	var c1 *comp1
	var c2 *comp2

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := w.Query(&filter)
		b.StartTimer()
		for query.Next() {
			c1 = (*comp1)(query.Get(id1))
			c2 = (*comp2)(query.Get(id2))
		}
	}
	b.StopTimer()
	sum := c1.V + c2.V
	_ = sum
}

func queryIterGet_5_100_000(b *testing.B) {
	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)

	builder := ecs.NewBuilder(&w, id1, id2, id3, id4, id5)
	builder.NewBatch(100_000)
	filter := ecs.All(id1, id2, id3, id4, id5)

	var c1 *comp1
	var c2 *comp2
	var c3 *comp3
	var c4 *comp4
	var c5 *comp5

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := w.Query(&filter)
		b.StartTimer()
		for query.Next() {
			c1 = (*comp1)(query.Get(id1))
			c2 = (*comp2)(query.Get(id2))
			c3 = (*comp3)(query.Get(id3))
			c4 = (*comp4)(query.Get(id4))
			c5 = (*comp5)(query.Get(id5))
		}
	}
	b.StopTimer()
	sum := c1.V + c2.V + c3.V + c4.V + c5.V
	_ = sum
}

func queryIterEntity_100_000(b *testing.B) {
	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)

	builder := ecs.NewBuilder(&w, id1)
	builder.NewBatch(100_000)
	filter := ecs.All(id1)

	var e ecs.Entity
	var s string

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := w.Query(&filter)
		b.StartTimer()
		for query.Next() {
			e = query.Entity()
		}
		b.StopTimer()
		s = fmt.Sprint(e)
	}
	_ = s
}

func queryRelation_100_000(b *testing.B) {
	w := ecs.NewWorld()
	id1 := ecs.ComponentID[relComp1](&w)
	parent := w.NewEntity()

	builder := ecs.NewBuilder(&w, id1).WithRelation(id1)
	builder.NewBatch(100_000, parent)
	filter := ecs.All(id1)

	var par ecs.Entity

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query := w.Query(&filter)
		b.StartTimer()
		for query.Next() {
			par = query.Relation(id1)
		}
	}
	b.StopTimer()
	sum := par.IsZero()
	_ = sum
}

func queryEntityAt_1Arch_1000(b *testing.B) {
	b.StopTimer()
	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	builder := ecs.NewBuilder(&w, id1)
	builder.NewBatch(1000)

	indices := make([]int, 1000)
	for i := range indices {
		indices[i] = rand.Intn(1000)
	}

	filter := ecs.All(id1)
	query := w.Query(&filter)
	b.StartTimer()
	var e ecs.Entity
	for i := 0; i < b.N; i++ {
		for _, idx := range indices {
			e = query.EntityAt(idx)
		}
	}
	_ = e
}

func queryEntityAtRegistered_1Arch_1000(b *testing.B) {
	b.StopTimer()
	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	builder := ecs.NewBuilder(&w, id1)
	builder.NewBatch(1000)

	indices := make([]int, 1000)
	for i := range indices {
		indices[i] = i
	}
	rand.Shuffle(len(indices), func(i, j int) { indices[i], indices[j] = indices[j], indices[i] })

	f := ecs.All(id1)
	filter := w.Cache().Register(f)

	query := w.Query(&filter)
	b.StartTimer()
	var e ecs.Entity
	for i := 0; i < b.N; i++ {
		for _, idx := range indices {
			e = query.EntityAt(idx)
		}
	}
	_ = e
}

func queryEntityAt_5Arch_1000(b *testing.B) {
	b.StopTimer()
	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)

	b1 := ecs.NewBuilder(&w, id1)
	b2 := ecs.NewBuilder(&w, id1, id2)
	b3 := ecs.NewBuilder(&w, id1, id3)
	b4 := ecs.NewBuilder(&w, id1, id4)
	b5 := ecs.NewBuilder(&w, id1, id5)
	b1.NewBatch(200)
	b2.NewBatch(200)
	b3.NewBatch(200)
	b4.NewBatch(200)
	b5.NewBatch(200)

	indices := make([]int, 1000)
	for i := range indices {
		indices[i] = i
	}
	rand.Shuffle(len(indices), func(i, j int) { indices[i], indices[j] = indices[j], indices[i] })

	filter := ecs.All(id1)
	query := w.Query(&filter)
	b.StartTimer()
	var e ecs.Entity
	for i := 0; i < b.N; i++ {
		for _, idx := range indices {
			e = query.EntityAt(idx)
		}
	}
	_ = e
}

func queryEntityAtRegistered_5Arch_1000(b *testing.B) {
	b.StopTimer()
	w := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&w)
	id2 := ecs.ComponentID[comp2](&w)
	id3 := ecs.ComponentID[comp3](&w)
	id4 := ecs.ComponentID[comp4](&w)
	id5 := ecs.ComponentID[comp5](&w)

	b1 := ecs.NewBuilder(&w, id1)
	b2 := ecs.NewBuilder(&w, id1, id2)
	b3 := ecs.NewBuilder(&w, id1, id3)
	b4 := ecs.NewBuilder(&w, id1, id4)
	b5 := ecs.NewBuilder(&w, id1, id5)
	b1.NewBatch(200)
	b2.NewBatch(200)
	b3.NewBatch(200)
	b4.NewBatch(200)
	b5.NewBatch(200)

	indices := make([]int, 1000)
	for i := range indices {
		indices[i] = rand.Intn(1000)
	}

	f := ecs.All(id1)
	filter := w.Cache().Register(f)

	query := w.Query(&filter)
	b.StartTimer()
	var e ecs.Entity
	for i := 0; i < b.N; i++ {
		for _, idx := range indices {
			e = query.EntityAt(idx)
		}
	}
	_ = e
}

func queryCreate(b *testing.B) {
	b.StopTimer()

	world := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&world)
	builder := ecs.NewBuilder(&world, id1)
	builder.NewBatch(100)

	filter := ecs.All(id1)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(&filter)
		query.Close()
	}
}

func queryCreateCached(b *testing.B) {
	b.StopTimer()

	world := ecs.NewWorld()
	id1 := ecs.ComponentID[comp1](&world)
	builder := ecs.NewBuilder(&world, id1)
	builder.NewBatch(100)

	f := ecs.All(id1)
	filter := world.Cache().Register(f)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(&filter)
		query.Close()
	}
}
*/
