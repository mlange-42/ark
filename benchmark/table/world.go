package main

import (
	"math/rand"
	"testing"

	"github.com/mlange-42/ark/benchmark"
	"github.com/mlange-42/ark/ecs"
)

func benchesWorld() []benchmark.Benchmark {
	return []benchmark.Benchmark{
		{Name: "Map.Get", Desc: "random, 1000 entities", F: worldGet1000, N: 1000},
		{Name: "Map.GetUnchecked", Desc: "random, 1000 entities", F: worldGetUnchecked1000, N: 1000},
		{Name: "Map.Has", Desc: "random, 1000 entities", F: worldHas1000, N: 1000},
		{Name: "Map.HasUnchecked", Desc: "random, 1000 entities", F: worldHasUnchecked1000, N: 1000},
		{Name: "World.Alive", Desc: "random, 1000 entities", F: worldAlive1000, N: 1000},
		{Name: "Map.GetRelation", Desc: "random, 1000 entities", F: worldRelation1000, N: 1000},
		{Name: "Map.GetRelationUnchecked", Desc: "random, 1000 entities", F: worldRelationUnchecked1000, N: 1000},
	}
}

func worldGet1000(b *testing.B) {
	w := ecs.NewWorld()

	mapper := ecs.NewMap[comp1](&w)
	entities := make([]ecs.Entity, 0, 1000)
	builder := ecs.NewMap1[comp1](&w)
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *comp1) {
		entities = append(entities, entity)
	})

	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var comp *comp1
	for b.Loop() {
		for _, e := range entities {
			comp = mapper.Get(e)
		}
	}
	_ = comp
}

func worldGetUnchecked1000(b *testing.B) {
	w := ecs.NewWorld()

	mapper := ecs.NewMap[comp1](&w)
	entities := make([]ecs.Entity, 0, 1000)
	builder := ecs.NewMap1[comp1](&w)
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *comp1) {
		entities = append(entities, entity)
	})

	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var comp *comp1
	for b.Loop() {
		for _, e := range entities {
			comp = mapper.GetUnchecked(e)
		}
	}
	_ = comp
}

func worldHas1000(b *testing.B) {
	w := ecs.NewWorld()

	mapper := ecs.NewMap[comp1](&w)
	entities := make([]ecs.Entity, 0, 1000)
	builder := ecs.NewMap1[comp1](&w)
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *comp1) {
		entities = append(entities, entity)
	})

	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var has bool
	for b.Loop() {
		for _, e := range entities {
			has = mapper.Has(e)
		}
	}
	_ = has
}

func worldHasUnchecked1000(b *testing.B) {
	w := ecs.NewWorld()

	mapper := ecs.NewMap[comp1](&w)
	entities := make([]ecs.Entity, 0, 1000)
	builder := ecs.NewMap1[comp1](&w)
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *comp1) {
		entities = append(entities, entity)
	})

	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var has bool
	for b.Loop() {
		for _, e := range entities {
			has = mapper.HasUnchecked(e)
		}
	}
	_ = has
}

func worldAlive1000(b *testing.B) {
	w := ecs.NewWorld()

	entities := make([]ecs.Entity, 0, 1000)
	builder := ecs.NewMap1[comp1](&w)
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *comp1) {
		entities = append(entities, entity)
	})

	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var has bool
	for b.Loop() {
		for _, e := range entities {
			has = w.Alive(e)
		}
	}
	_ = has
}

func worldRelation1000(b *testing.B) {
	w := ecs.NewWorld()
	parent := w.NewEntity()

	mapper := ecs.NewMap[relComp1](&w)
	entities := make([]ecs.Entity, 0, 1000)
	builder := ecs.NewMap1[relComp1](&w)
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *relComp1) {
		entities = append(entities, entity)
	}, ecs.Rel[relComp1](parent))

	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var par ecs.Entity
	for b.Loop() {
		for _, e := range entities {
			par = mapper.GetRelation(e)
		}
	}
	_ = par
}

func worldRelationUnchecked1000(b *testing.B) {
	w := ecs.NewWorld()
	parent := w.NewEntity()

	mapper := ecs.NewMap[relComp1](&w)
	entities := make([]ecs.Entity, 0, 1000)
	builder := ecs.NewMap1[relComp1](&w)
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *relComp1) {
		entities = append(entities, entity)
	}, ecs.Rel[relComp1](parent))

	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	var par ecs.Entity
	for b.Loop() {
		for _, e := range entities {
			par = mapper.GetRelationUnchecked(e)
		}
	}
	_ = par
}
