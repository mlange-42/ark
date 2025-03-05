package main

import (
	"testing"

	"github.com/mlange-42/ark/benchmark"
	"github.com/mlange-42/ark/ecs"
)

func benchesEntities() []benchmark.Benchmark {
	return []benchmark.Benchmark{
		{Name: "Entity.IsZero", Desc: "", F: entitiesIsZero2, N: 2},

		{Name: "World.NewEntity", Desc: "memory already alloc.", F: entitiesCreate1000, N: 1000},
		{Name: "Map1.NewEntity w/ 1 Comp", Desc: "memory already alloc.", F: entitiesCreate1Comp1000, N: 1000},
		{Name: "Map5.NewEntity w/ 5 Comps", Desc: "memory already alloc.", F: entitiesCreate5Comp1000, N: 1000},

		{Name: "World.RemoveEntity", Desc: "", F: entitiesRemove1000, N: 1000},
		{Name: "World.RemoveEntity w/ 1 Comp", Desc: "", F: entitiesRemove1Comp1000, N: 1000},
		{Name: "World.RemoveEntity w/ 5 Comps", Desc: "", F: entitiesRemove5Comp1000, N: 1000},
	}
}

func entitiesIsZero2(b *testing.B) {
	w := ecs.NewWorld()
	e := w.NewEntity()
	z := ecs.Entity{}
	var zero1 bool
	var zero2 bool

	for b.Loop() {
		zero1 = e.IsZero()
		zero2 = z.IsZero()
	}
	s := zero1 || zero2
	_ = s
}

func entitiesCreate1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	filter := ecs.NewFilter0(&w)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for j := 0; j < 1000; j++ {
			_ = w.NewEntity()
		}
		b.StopTimer()
		w.RemoveEntities(filter.Batch(), nil)
	}
}

func entitiesCreate1Comp1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	builder := ecs.NewMap1[comp1](&w)
	filter := ecs.NewFilter0(&w)

	c1 := comp1{}
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for j := 0; j < 1000; j++ {
			_ = builder.NewEntity(&c1)
		}
		b.StopTimer()
		w.RemoveEntities(filter.Batch(), nil)
	}
}

func entitiesCreate5Comp1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	builder := ecs.NewMap5[comp1, comp2, comp3, comp4, comp5](&w)
	filter := ecs.NewFilter0(&w)

	c1 := comp1{}
	c2 := comp2{}
	c3 := comp3{}
	c4 := comp4{}
	c5 := comp5{}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for j := 0; j < 1000; j++ {
			_ = builder.NewEntity(&c1, &c2, &c3, &c4, &c5)
		}
		b.StopTimer()
		w.RemoveEntities(filter.Batch(), nil)
	}
}

func entitiesRemove1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()

	entities := make([]ecs.Entity, 0, 1000)

	for i := 0; i < b.N; i++ {
		w.NewEntities(1000, func(entity ecs.Entity) {
			entities = append(entities, entity)
		})
		b.StartTimer()
		for _, e := range entities {
			w.RemoveEntity(e)
		}
		b.StopTimer()
		entities = entities[:0]
	}
}

func entitiesRemove1Comp1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	builder := ecs.NewMap1[comp1](&w)

	entities := make([]ecs.Entity, 0, 1000)

	for i := 0; i < b.N; i++ {
		builder.NewBatchFn(1000, func(entity ecs.Entity, a *comp1) {
			entities = append(entities, entity)
		})
		b.StartTimer()
		for _, e := range entities {
			w.RemoveEntity(e)
		}
		b.StopTimer()
		entities = entities[:0]
	}
}

func entitiesRemove5Comp1000(b *testing.B) {
	b.StopTimer()

	w := ecs.NewWorld()
	builder := ecs.NewMap5[comp1, comp2, comp3, comp4, comp5](&w)

	entities := make([]ecs.Entity, 0, 1000)

	for i := 0; i < b.N; i++ {
		builder.NewBatchFn(1000, func(entity ecs.Entity, a *comp1, b *comp2, c *comp3, d *comp4, e *comp5) {
			entities = append(entities, entity)
		})
		b.StartTimer()
		for _, e := range entities {
			w.RemoveEntity(e)
		}
		b.StopTimer()
		entities = entities[:0]
	}
}
