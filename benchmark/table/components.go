package main

import (
	"testing"

	"github.com/mlange-42/ark/benchmark"
	"github.com/mlange-42/ark/ecs"
)

func benchesComponents() []benchmark.Benchmark {
	return []benchmark.Benchmark{
		{Name: "Map1.Add 1 Comp", Desc: "memory already alloc.", F: componentsAdd1_1000, N: 1000},
		{Name: "Map5.Add 5 Comps", Desc: "memory already alloc.", F: componentsAdd5_1000, N: 1000},
		{Name: "Map1.Add 1 to 5 Comps", Desc: "memory already alloc.", F: componentsAdd1to5_1000, N: 1000},

		{Name: "Map1.Remove 1 Comp", Desc: "memory already alloc.", F: componentsRemove1_1000, N: 1000},
		{Name: "Map5.Remove 5 Comps", Desc: "memory already alloc.", F: componentsRemove5_1000, N: 1000},
		{Name: "Map1.Remove 1 of 5 Comps", Desc: "memory already alloc.", F: componentsRemove1of5_1000, N: 1000},

		{Name: "Exchange1.Exchange 1 Comp", Desc: "memory already alloc.", F: componentsExchange1_1000, N: 1000},
		{Name: "Exchange1.Exchange 1 of 5 Comps", Desc: "memory already alloc.", F: componentsExchange1of5_1000, N: 1000},
	}
}

func componentsAdd1_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	mapper := ecs.NewMap1[comp1](&w)
	filter := ecs.NewFilter1[comp1](&w)

	entities := make([]ecs.Entity, 0, 1000)
	w.NewEntities(1000, func(entity ecs.Entity) {
		entities = append(entities, entity)
	})

	// Run once to allocate memory
	for _, e := range entities {
		mapper.AddFn(e, nil)
	}
	mapper.RemoveBatch(filter.Batch(), nil)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			mapper.AddFn(e, nil)
		}
		b.StopTimer()
		mapper.RemoveBatch(filter.Batch(), nil)
	}
}

func componentsAdd5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	mapper := ecs.NewMap5[comp1, comp2, comp3, comp4, comp5](&w)
	filter := ecs.NewFilter5[comp1, comp2, comp3, comp4, comp5](&w)

	entities := make([]ecs.Entity, 0, 1000)
	w.NewEntities(1000, func(entity ecs.Entity) {
		entities = append(entities, entity)
	})

	// Run once to allocate memory
	for _, e := range entities {
		mapper.AddFn(e, nil)
	}
	mapper.RemoveBatch(filter.Batch(), nil)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			mapper.AddFn(e, nil)
		}
		b.StopTimer()
		mapper.RemoveBatch(filter.Batch(), nil)
	}
}

func componentsAdd1to5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	builder := ecs.NewMap5[comp2, comp3, comp4, comp5, comp6](&w)
	mapper := ecs.NewMap1[comp1](&w)
	filter := ecs.NewFilter1[comp1](&w)

	entities := make([]ecs.Entity, 0, 1000)
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *comp2, b *comp3, c *comp4, d *comp5, e *comp6) {
		entities = append(entities, entity)
	})

	// Run once to allocate memory
	for _, e := range entities {
		mapper.AddFn(e, nil)
	}
	mapper.RemoveBatch(filter.Batch(), nil)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			mapper.AddFn(e, nil)
		}
		b.StopTimer()
		mapper.RemoveBatch(filter.Batch(), nil)
	}
}

func componentsRemove1_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	mapper := ecs.NewMap1[comp1](&w)
	filter := ecs.NewFilter0(&w)

	entities := make([]ecs.Entity, 0, 1000)
	mapper.NewBatchFn(1000, func(entity ecs.Entity, a *comp1) {
		entities = append(entities, entity)
	})

	// Run once to allocate memory
	for _, e := range entities {
		mapper.Remove(e)
	}
	mapper.AddBatchFn(filter.Batch(), nil)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			mapper.Remove(e)
		}
		b.StopTimer()
		mapper.AddBatchFn(filter.Batch(), nil)
	}
}

func componentsRemove5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	mapper := ecs.NewMap5[comp1, comp2, comp3, comp4, comp5](&w)
	filter := ecs.NewFilter0(&w)

	entities := make([]ecs.Entity, 0, 1000)
	mapper.NewBatchFn(1000, func(entity ecs.Entity, a *comp1, b *comp2, c *comp3, d *comp4, e *comp5) {
		entities = append(entities, entity)
	})

	// Run once to allocate memory
	for _, e := range entities {
		mapper.Remove(e)
	}
	mapper.AddBatchFn(filter.Batch(), nil)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			mapper.Remove(e)
		}
		b.StopTimer()
		mapper.AddBatchFn(filter.Batch(), nil)
	}
}

func componentsRemove1of5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	builder := ecs.NewMap5[comp1, comp2, comp3, comp4, comp5](&w)
	mapper := ecs.NewMap1[comp1](&w)
	filter := ecs.NewFilter0(&w)

	entities := make([]ecs.Entity, 0, 1000)
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *comp1, b *comp2, c *comp3, d *comp4, e *comp5) {
		entities = append(entities, entity)
	})

	// Run once to allocate memory
	for _, e := range entities {
		mapper.Remove(e)
	}
	mapper.AddBatchFn(filter.Batch(), nil)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			mapper.Remove(e)
		}
		b.StopTimer()
		mapper.AddBatchFn(filter.Batch(), nil)
	}
}

func componentsExchange1_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	ex1 := ecs.NewExchange1[comp1](&w).Removes(ecs.C[comp2]())
	ex2 := ecs.NewExchange1[comp2](&w).Removes(ecs.C[comp1]())
	builder := ecs.NewMap1[comp1](&w)
	filter := ecs.NewFilter1[comp2](&w)

	entities := make([]ecs.Entity, 0, 1000)
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *comp1) {
		entities = append(entities, entity)
	})

	// Run once to allocate memory
	for _, e := range entities {
		ex2.ExchangeFn(e, nil)
	}
	ex1.ExchangeBatchFn(filter.Batch(), nil)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			ex2.ExchangeFn(e, nil)
		}
		b.StopTimer()
		ex1.ExchangeBatchFn(filter.Batch(), nil)
	}
}

func componentsExchange1of5_1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	ex1 := ecs.NewExchange1[comp1](&w).Removes(ecs.C[comp2]())
	ex2 := ecs.NewExchange1[comp2](&w).Removes(ecs.C[comp1]())
	builder := ecs.NewMap5[comp1, comp3, comp4, comp5, comp6](&w)
	filter := ecs.NewFilter1[comp2](&w)

	entities := make([]ecs.Entity, 0, 1000)
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *comp1, b *comp3, c *comp4, d *comp5, e *comp6) {
		entities = append(entities, entity)
	})

	// Run once to allocate memory
	for _, e := range entities {
		ex2.ExchangeFn(e, nil)
	}
	ex1.ExchangeBatchFn(filter.Batch(), nil)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			ex2.ExchangeFn(e, nil)
		}
		b.StopTimer()
		ex1.ExchangeBatchFn(filter.Batch(), nil)
	}
}
