package main

import (
	"testing"

	"github.com/mlange-42/ark/ecs"
)

func addRemove10(b *testing.B) {
	addRemove(b, 10)
}

func addRemove1000(b *testing.B) {
	addRemove(b, 1000)
}

func addRemove(b *testing.B, n int) {
	world := ecs.NewWorld()

	posMap := ecs.NewMap1[Position](world)
	velMap := ecs.NewMap1[Velocity](world)

	entities := make([]ecs.Entity, 0, n)

	for range n {
		e := posMap.NewEntityFn(nil)
		entities = append(entities, e)
	}

	for _, e := range entities {
		velMap.AddFn(e, nil)
	}
	for _, e := range entities {
		velMap.Remove(e)
	}

	for b.Loop() {
		for _, e := range entities {
			velMap.AddFn(e, nil)
		}
		for _, e := range entities {
			velMap.Remove(e)
		}
	}
}

func addRemoveNonTrivial10(b *testing.B) {
	addRemoveNonTrivial(b, 10)
}

func addRemoveNonTrivial1000(b *testing.B) {
	addRemoveNonTrivial(b, 1000)
}

func addRemoveNonTrivial(b *testing.B, n int) {
	world := ecs.NewWorld()

	map1 := ecs.NewMap1[SliceComp1](world)
	map2 := ecs.NewMap1[SliceComp2](world)

	entities := make([]ecs.Entity, 0, n)

	for range n {
		e := map1.NewEntityFn(nil)
		entities = append(entities, e)
	}

	for _, e := range entities {
		map2.AddFn(e, nil)
	}
	for _, e := range entities {
		map2.Remove(e)
	}

	for b.Loop() {
		for _, e := range entities {
			map2.AddFn(e, nil)
		}
		for _, e := range entities {
			map2.Remove(e)
		}
	}
}

func addRemoveBatch10(b *testing.B) {
	addRemoveBatch(b, 10)
}

func addRemoveBatch1000(b *testing.B) {
	addRemoveBatch(b, 1000)
}

func addRemoveBatch100000(b *testing.B) {
	addRemoveBatch(b, 100000)
}

func addRemoveBatch(b *testing.B, n int) {
	world := ecs.NewWorld()

	posMap := ecs.NewMap1[Position](world)
	velMap := ecs.NewMap1[Velocity](world)

	posFilter := ecs.NewFilter1[Position](world)
	velFilter := ecs.NewFilter1[Velocity](world)

	posMap.NewBatchFn(n, nil)

	velMap.AddBatchFn(posFilter.Batch(), nil)
	velMap.RemoveBatch(velFilter.Batch(), nil)

	for b.Loop() {
		velMap.AddBatchFn(posFilter.Batch(), nil)
		velMap.RemoveBatch(velFilter.Batch(), nil)
	}
}

func addRemoveBatchNonTrivial10(b *testing.B) {
	addRemoveBatchNonTrivial(b, 10)
}

func addRemoveBatchNonTrivial100(b *testing.B) {
	addRemoveBatchNonTrivial(b, 100)
}

func addRemoveBatchNonTrivial1000(b *testing.B) {
	addRemoveBatchNonTrivial(b, 1000)
}

func addRemoveBatchNonTrivial100000(b *testing.B) {
	addRemoveBatchNonTrivial(b, 100000)
}

func addRemoveBatchNonTrivial(b *testing.B, n int) {
	world := ecs.NewWorld()

	posMap := ecs.NewMap1[SliceComp1](world)
	velMap := ecs.NewMap1[SliceComp2](world)

	posFilter := ecs.NewFilter1[SliceComp1](world)
	velFilter := ecs.NewFilter1[SliceComp2](world)

	posMap.NewBatchFn(n, nil)

	velMap.AddBatchFn(posFilter.Batch(), nil)
	velMap.RemoveBatch(velFilter.Batch(), nil)

	for b.Loop() {
		velMap.AddBatchFn(posFilter.Batch(), nil)
		velMap.RemoveBatch(velFilter.Batch(), nil)
	}
}

func addRemoveBatchLarge10(b *testing.B) {
	addRemoveBatchLarge(b, 10)
}

func addRemoveBatchLarge1000(b *testing.B) {
	addRemoveBatchLarge(b, 1000)
}

func addRemoveBatchLarge100000(b *testing.B) {
	addRemoveBatchLarge(b, 100000)
}

func addRemoveBatchLarge(b *testing.B, n int) {
	world := ecs.NewWorld()

	posMap := ecs.NewMap11[Position, Comp1, Comp2, Comp3, Comp4, Comp5, Comp6, Comp7, Comp8, Comp9, Comp10](world)
	velMap := ecs.NewMap1[Velocity](world)

	posFilter := ecs.NewFilter1[Position](world)
	velFilter := ecs.NewFilter1[Velocity](world)

	posMap.NewBatchFn(n, nil)

	velMap.AddBatchFn(posFilter.Batch(), nil)
	velMap.RemoveBatch(velFilter.Batch(), nil)

	for b.Loop() {
		velMap.AddBatchFn(posFilter.Batch(), nil)
		velMap.RemoveBatch(velFilter.Batch(), nil)
	}
}
