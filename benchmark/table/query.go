package main

import (
	"testing"

	"github.com/mlange-42/ark/benchmark"
	"github.com/mlange-42/ark/ecs"
)

func benchesQuery() []benchmark.Benchmark {
	return []benchmark.Benchmark{
		{Name: "Query.Next", Desc: "", F: queryIter100k, N: 100_000},
		{Name: "Query.Next + Query.Get 1", Desc: "", F: queryIterGet1Comp100k, N: 100_000},
		{Name: "Query.Next + Query.Get 2", Desc: "", F: queryIterGet2Comp100k, N: 100_000},
		{Name: "Query.Next + Query.Get 5", Desc: "", F: queryIterGet5Comp100k, N: 100_000},

		{Name: "Query.Next + Query.Entity", Desc: "", F: queryIterEntity100k, N: 100_000},

		{Name: "Query.Next + Query.Relation", Desc: "", F: queryRelation100k, N: 100_000},

		{Name: "Filter1.Query + Query1.Close", Desc: "", F: queryCreate, N: 1},
		{Name: "Filter1.Query + Query1.Close", Desc: "registered filter", F: queryCreateCached, N: 1},
	}
}

func queryIter100k(b *testing.B) {
	w := ecs.NewWorld()

	w.NewEntities(100_000, nil)
	filter := ecs.NewFilter0(&w)

	for b.Loop() {
		query := filter.Query()
		for query.Next() {
		}
	}
}

func queryIterGet1Comp100k(b *testing.B) {
	w := ecs.NewWorld()

	builder := ecs.NewMap1[comp1](&w)
	builder.NewBatchFn(100_000, nil)

	filter := ecs.NewFilter1[comp1](&w)

	var c1 *comp1

	for b.Loop() {
		query := filter.Query()
		for query.Next() {
			c1 = query.Get()
		}
	}
	_ = c1
}

func queryIterGet2Comp100k(b *testing.B) {
	w := ecs.NewWorld()

	builder := ecs.NewMap2[comp1, comp2](&w)
	builder.NewBatchFn(100_000, nil)

	filter := ecs.NewFilter2[comp1, comp2](&w)

	var c1 *comp1
	var c2 *comp2

	for b.Loop() {
		query := filter.Query()
		for query.Next() {
			c1, c2 = query.Get()
		}
	}
	_, _ = c1, c2
}

func queryIterGet5Comp100k(b *testing.B) {
	w := ecs.NewWorld()

	builder := ecs.NewMap5[comp1, comp2, comp3, comp4, comp5](&w)
	builder.NewBatchFn(100_000, nil)

	filter := ecs.NewFilter5[comp1, comp2, comp3, comp4, comp5](&w)

	var c1 *comp1
	var c2 *comp2
	var c3 *comp3
	var c4 *comp4
	var c5 *comp5

	for b.Loop() {
		query := filter.Query()
		for query.Next() {
			c1, c2, c3, c4, c5 = query.Get()
		}
	}
	sum := c1.V + c2.V + c3.V + c4.V + c5.V
	_ = sum
}

func queryIterEntity100k(b *testing.B) {
	w := ecs.NewWorld()

	builder := ecs.NewMap1[comp1](&w)
	builder.NewBatchFn(100_000, nil)
	filter := ecs.NewFilter1[comp1](&w)

	var e ecs.Entity

	b.ResetTimer()
	for b.Loop() {
		query := filter.Query()
		for query.Next() {
			e = query.Entity()
		}
	}
	_ = e
}

func queryRelation100k(b *testing.B) {
	w := ecs.NewWorld()
	parent := w.NewEntity()

	builder := ecs.NewMap1[relComp1](&w)
	builder.NewBatchFn(100_000, nil, ecs.Rel[relComp1](parent))
	filter := ecs.NewFilter1[relComp1](&w)

	var par ecs.Entity
	for b.Loop() {
		query := filter.Query()
		for query.Next() {
			par = query.GetRelation(0)
		}
	}
	_ = par
}

func queryCreate(b *testing.B) {
	w := ecs.NewWorld()

	builder := ecs.NewMap1[comp1](&w)
	builder.NewBatchFn(100, nil)
	filter := ecs.NewFilter1[comp1](&w)

	for b.Loop() {
		query := filter.Query()
		query.Close()
	}
}

func queryCreateCached(b *testing.B) {
	w := ecs.NewWorld()

	builder := ecs.NewMap1[comp1](&w)
	builder.NewBatchFn(100, nil)
	filter := ecs.NewFilter1[comp1](&w).Register()

	for b.Loop() {
		query := filter.Query()
		query.Close()
	}
}
