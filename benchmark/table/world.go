package main

import (
	"math/rand"
	"runtime"
	"testing"

	"github.com/mlange-42/ark/benchmark"
	"github.com/mlange-42/ark/ecs"
)

func benchesWorld() []benchmark.Benchmark {
	return []benchmark.Benchmark{
		{Name: "World.Alive", Desc: "random, 1000 entities", F: worldAlive1000, N: 1000},
		{Name: "Map.Get", Desc: "random, 1000 entities", F: worldGet1000, N: 1000},
		{Name: "Map.GetUnchecked", Desc: "random, 1000 entities", F: worldGetUnchecked1000, N: 1000},
		{Name: "Map.Has", Desc: "random, 1000 entities", F: worldHas1000, N: 1000},
		{Name: "Map.HasUnchecked", Desc: "random, 1000 entities", F: worldHasUnchecked1000, N: 1000},
		{Name: "Map5.Get 5", Desc: "random, 1000 entities", F: worldGet5_1000, N: 1000},
		{Name: "Map5.HasAll 5", Desc: "random, 1000 entities", F: worldHasAll5_1000, N: 1000},
		{Name: "Map.GetRelation", Desc: "random, 1000 entities", F: worldRelation1000, N: 1000},
		{Name: "Map.GetRelationUnchecked", Desc: "random, 1000 entities", F: worldRelationUnchecked1000, N: 1000},
	}
}

func worldGet1000(b *testing.B) {
	w := ecs.NewWorld()

	mapper := ecs.NewMap[comp1](w)
	entities := make([]ecs.Entity, 0, 1000)
	builder := ecs.NewMap1[comp1](w)
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *comp1) {
		entities = append(entities, entity)
	})

	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var comp *comp1

	// Wrapper to allow inlining, for more realistic results.
	loop := func() {
		for _, e := range entities {
			comp = mapper.Get(e)
		}
	}
	for b.Loop() {
		loop()
	}

	runtime.KeepAlive(comp)
}

func worldGetUnchecked1000(b *testing.B) {
	w := ecs.NewWorld()

	mapper := ecs.NewMap[comp1](w)
	entities := make([]ecs.Entity, 0, 1000)
	builder := ecs.NewMap1[comp1](w)
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *comp1) {
		entities = append(entities, entity)
	})

	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var comp *comp1

	// Wrapper to allow inlining, for more realistic results.
	loop := func() {
		for _, e := range entities {
			comp = mapper.GetUnchecked(e)
		}
	}
	for b.Loop() {
		loop()
	}
	runtime.KeepAlive(comp)
}

func worldHas1000(b *testing.B) {
	w := ecs.NewWorld()

	mapper := ecs.NewMap[comp1](w)
	entities := make([]ecs.Entity, 0, 1000)
	builder := ecs.NewMap1[comp1](w)
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *comp1) {
		entities = append(entities, entity)
	})

	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var has bool

	// Wrapper to allow inlining, for more realistic results.
	loop := func() {
		for _, e := range entities {
			has = mapper.Has(e)
		}
	}
	for b.Loop() {
		loop()
	}
	runtime.KeepAlive(has)
}

func worldHasUnchecked1000(b *testing.B) {
	w := ecs.NewWorld()

	mapper := ecs.NewMap[comp1](w)
	entities := make([]ecs.Entity, 0, 1000)
	builder := ecs.NewMap1[comp1](w)
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *comp1) {
		entities = append(entities, entity)
	})

	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var has bool

	// Wrapper to allow inlining, for more realistic results.
	loop := func() {
		for _, e := range entities {
			has = mapper.HasUnchecked(e)
		}
	}
	for b.Loop() {
		loop()
	}
	runtime.KeepAlive(has)
}

func worldGet5_1000(b *testing.B) {
	w := ecs.NewWorld()

	entities := make([]ecs.Entity, 0, 1000)
	mapper := ecs.NewMap5[comp1, comp2, comp3, comp4, comp5](w)
	mapper.NewBatchFn(1000, func(entity ecs.Entity, a *comp1, b *comp2, c *comp3, d *comp4, e *comp5) {
		entities = append(entities, entity)
	})

	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var comp1 *comp1
	var comp2 *comp2
	var comp3 *comp3
	var comp4 *comp4
	var comp5 *comp5

	// Wrapper to allow inlining, for more realistic results.
	loop := func() {
		for _, e := range entities {
			comp1, comp2, comp3, comp4, comp5 = mapper.Get(e)
		}
	}
	for b.Loop() {
		loop()
	}
	runtime.KeepAlive(comp1)
	runtime.KeepAlive(comp2)
	runtime.KeepAlive(comp3)
	runtime.KeepAlive(comp4)
	runtime.KeepAlive(comp5)
}

func worldHasAll5_1000(b *testing.B) {
	w := ecs.NewWorld()

	entities := make([]ecs.Entity, 0, 1000)
	mapper := ecs.NewMap5[comp1, comp2, comp3, comp4, comp5](w)
	mapper.NewBatchFn(1000, func(entity ecs.Entity, a *comp1, b *comp2, c *comp3, d *comp4, e *comp5) {
		entities = append(entities, entity)
	})

	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var has bool

	// Wrapper to allow inlining, for more realistic results.
	loop := func() {
		for _, e := range entities {
			has = mapper.HasAll(e)
		}
	}
	for b.Loop() {
		loop()
	}
	runtime.KeepAlive(has)
}

func worldAlive1000(b *testing.B) {
	w := ecs.NewWorld()

	entities := make([]ecs.Entity, 0, 1000)
	builder := ecs.NewMap1[comp1](w)
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *comp1) {
		entities = append(entities, entity)
	})

	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var has bool

	// Wrapper to allow inlining, for more realistic results.
	loop := func() {
		for _, e := range entities {
			has = w.Alive(e)
		}
	}
	for b.Loop() {
		loop()
	}
	runtime.KeepAlive(has)
}

func worldRelation1000(b *testing.B) {
	w := ecs.NewWorld()
	parent := w.NewEntity()

	mapper := ecs.NewMap[relComp1](w)
	entities := make([]ecs.Entity, 0, 1000)
	builder := ecs.NewMap1[relComp1](w)
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *relComp1) {
		entities = append(entities, entity)
	}, ecs.Rel[relComp1](parent))

	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var par ecs.Entity

	// Wrapper to allow inlining, for more realistic results.
	loop := func() {
		for _, e := range entities {
			par = mapper.GetRelation(e)
		}
	}
	for b.Loop() {
		loop()
	}
	runtime.KeepAlive(par)
}

func worldRelationUnchecked1000(b *testing.B) {
	w := ecs.NewWorld()
	parent := w.NewEntity()

	mapper := ecs.NewMap[relComp1](w)
	entities := make([]ecs.Entity, 0, 1000)
	builder := ecs.NewMap1[relComp1](w)
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *relComp1) {
		entities = append(entities, entity)
	}, ecs.Rel[relComp1](parent))

	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var par ecs.Entity

	// Wrapper to allow inlining, for more realistic results.
	loop := func() {
		for _, e := range entities {
			par = mapper.GetRelationUnchecked(e)
		}
	}
	for b.Loop() {
		loop()
	}
	runtime.KeepAlive(par)
}
