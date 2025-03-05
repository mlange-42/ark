package main

import (
	"testing"

	"github.com/mlange-42/ark/benchmark"
	"github.com/mlange-42/ark/ecs"
)

func benchesComponentsBatch() []benchmark.Benchmark {
	return []benchmark.Benchmark{
		{Name: "Map1.AddBatch 1 Comp", Desc: "1000, memory already allocated", F: componentsBatchAdd1_1000, N: 1000},
		{Name: "Map5.AddBatch 5 Comps", Desc: "1000, memory already allocated", F: componentsBatchAdd5_1000, N: 1000},
		{Name: "Map1.AddBatch 1 to 5 Comps", Desc: "1000, memory already allocated", F: componentsBatchAdd1to5_1000, N: 1000},

		{Name: "Map1.RemoveBatch 1 Comp", Desc: "1000, memory already allocated", F: componentsBatchRemove1_1000, N: 1000},
		{Name: "Map5.RemoveBatch 5 Comps", Desc: "1000, memory already allocated", F: componentsBatchRemove5_1000, N: 1000},
		{Name: "Map1.RemoveBatch 1 of 5 Comps", Desc: "1000, memory already allocated", F: componentsBatchRemove1of5_1000, N: 1000},

		{Name: "Exchange1.ExchangeBatch 1 Comp", Desc: "1000, memory already allocated", F: componentsBatchExchange1_1000, N: 1000},
		{Name: "Exchange1.ExchangeBatch 1 of 5 Comps", Desc: "1000, memory already allocated", F: componentsBatchExchange1of5_1000, N: 1000},
	}
}

func componentsBatchAdd1_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	mapper := ecs.NewMap1[comp1](&w)
	filter1 := ecs.NewFilter0(&w)
	filter2 := ecs.NewFilter1[comp1](&w)

	entities := make([]ecs.Entity, 0, 1000)
	w.NewEntities(1000, func(entity ecs.Entity) {
		entities = append(entities, entity)
	})

	// Run once to allocate memory
	mapper.AddBatchFn(filter1.Batch(), nil)
	mapper.RemoveBatch(filter2.Batch(), nil)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		mapper.AddBatchFn(filter1.Batch(), nil)
		b.StopTimer()
		mapper.RemoveBatch(filter2.Batch(), nil)
	}
}

func componentsBatchAdd5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	mapper := ecs.NewMap5[comp1, comp2, comp3, comp4, comp5](&w)
	filter1 := ecs.NewFilter0(&w)
	filter2 := ecs.NewFilter5[comp1, comp2, comp3, comp4, comp5](&w)

	entities := make([]ecs.Entity, 0, 1000)
	w.NewEntities(1000, func(entity ecs.Entity) {
		entities = append(entities, entity)
	})

	// Run once to allocate memory
	mapper.AddBatchFn(filter1.Batch(), nil)
	mapper.RemoveBatch(filter2.Batch(), nil)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		mapper.AddBatchFn(filter1.Batch(), nil)
		b.StopTimer()
		mapper.RemoveBatch(filter2.Batch(), nil)
	}
}

func componentsBatchAdd1to5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	builder := ecs.NewMap5[comp2, comp3, comp4, comp5, comp6](&w)
	mapper := ecs.NewMap1[comp1](&w)
	filter1 := ecs.NewFilter0(&w)
	filter2 := ecs.NewFilter1[comp1](&w)

	entities := make([]ecs.Entity, 0, 1000)
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *comp2, b *comp3, c *comp4, d *comp5, e *comp6) {
		entities = append(entities, entity)
	})

	// Run once to allocate memory
	mapper.AddBatchFn(filter1.Batch(), nil)
	mapper.RemoveBatch(filter2.Batch(), nil)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		mapper.AddBatchFn(filter1.Batch(), nil)
		b.StopTimer()
		mapper.RemoveBatch(filter2.Batch(), nil)
	}
}

func componentsBatchRemove1_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	mapper := ecs.NewMap1[comp1](&w)
	filter1 := ecs.NewFilter1[comp1](&w)
	filter2 := ecs.NewFilter0(&w)

	entities := make([]ecs.Entity, 0, 1000)
	mapper.NewBatchFn(1000, func(entity ecs.Entity, a *comp1) {
		entities = append(entities, entity)
	})

	// Run once to allocate memory
	mapper.RemoveBatch(filter1.Batch(), nil)
	mapper.AddBatchFn(filter2.Batch(), nil)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		mapper.RemoveBatch(filter1.Batch(), nil)
		b.StopTimer()
		mapper.AddBatchFn(filter2.Batch(), nil)
	}
}

func componentsBatchRemove5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	mapper := ecs.NewMap5[comp1, comp2, comp3, comp4, comp5](&w)
	filter1 := ecs.NewFilter5[comp1, comp2, comp3, comp4, comp5](&w)
	filter2 := ecs.NewFilter0(&w)

	entities := make([]ecs.Entity, 0, 1000)
	mapper.NewBatchFn(1000, func(entity ecs.Entity, a *comp1, b *comp2, c *comp3, d *comp4, e *comp5) {
		entities = append(entities, entity)
	})

	// Run once to allocate memory
	mapper.RemoveBatch(filter1.Batch(), nil)
	mapper.AddBatchFn(filter2.Batch(), nil)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		mapper.RemoveBatch(filter1.Batch(), nil)
		b.StopTimer()
		mapper.AddBatchFn(filter2.Batch(), nil)
	}
}

func componentsBatchRemove1of5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	builder := ecs.NewMap5[comp1, comp2, comp3, comp4, comp5](&w)
	mapper := ecs.NewMap1[comp1](&w)
	filter1 := ecs.NewFilter1[comp1](&w)
	filter2 := ecs.NewFilter0(&w)

	entities := make([]ecs.Entity, 0, 1000)
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *comp1, b *comp2, c *comp3, d *comp4, e *comp5) {
		entities = append(entities, entity)
	})

	// Run once to allocate memory
	mapper.RemoveBatch(filter1.Batch(), nil)
	mapper.AddBatchFn(filter2.Batch(), nil)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		mapper.RemoveBatch(filter1.Batch(), nil)
		b.StopTimer()
		mapper.AddBatchFn(filter2.Batch(), nil)
	}
}

func componentsBatchExchange1_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	ex1 := ecs.NewExchange1[comp1](&w).Removes(ecs.C[comp2]())
	ex2 := ecs.NewExchange1[comp2](&w).Removes(ecs.C[comp1]())
	builder := ecs.NewMap1[comp1](&w)
	filter1 := ecs.NewFilter1[comp1](&w)
	filter2 := ecs.NewFilter1[comp2](&w)

	entities := make([]ecs.Entity, 0, 1000)
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *comp1) {
		entities = append(entities, entity)
	})

	// Run once to allocate memory
	ex2.ExchangeBatchFn(filter1.Batch(), nil)
	ex1.ExchangeBatchFn(filter2.Batch(), nil)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		ex2.ExchangeBatchFn(filter1.Batch(), nil)
		b.StopTimer()
		ex1.ExchangeBatchFn(filter2.Batch(), nil)
	}
}

func componentsBatchExchange1of5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	ex1 := ecs.NewExchange1[comp1](&w).Removes(ecs.C[comp2]())
	ex2 := ecs.NewExchange1[comp2](&w).Removes(ecs.C[comp1]())
	builder := ecs.NewMap5[comp1, comp3, comp4, comp5, comp6](&w)
	filter1 := ecs.NewFilter1[comp1](&w)
	filter2 := ecs.NewFilter1[comp2](&w)

	entities := make([]ecs.Entity, 0, 1000)
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *comp1, b *comp3, c *comp4, d *comp5, e *comp6) {
		entities = append(entities, entity)
	})

	// Run once to allocate memory
	ex2.ExchangeBatchFn(filter1.Batch(), nil)
	ex1.ExchangeBatchFn(filter2.Batch(), nil)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		ex2.ExchangeBatchFn(filter1.Batch(), nil)
		b.StopTimer()
		ex1.ExchangeBatchFn(filter2.Batch(), nil)
	}
}
