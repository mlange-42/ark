package ecs

import (
	"math"
	"testing"
)

func BenchmarkCreateEntity1Comp_1000(b *testing.B) {
	w := NewWorld()
	builder := NewMap1[Position](&w)
	filter := NewFilter0(&w)

	builder.NewBatchFn(1000, nil)
	w.RemoveEntities(filter.Batch(), nil)

	for b.Loop() {
		for range 1000 {
			_ = builder.NewEntityFn(nil)
		}
		b.StopTimer()
		w.RemoveEntities(filter.Batch(), nil)
		b.StartTimer()
	}
}

func BenchmarkWorldStats4Arch(b *testing.B) {
	w := NewWorld()

	idA := ComponentID[CompA](&w)
	idB := ComponentID[CompB](&w)

	u := w.Unsafe()
	u.NewEntity()
	u.NewEntity(idA)
	u.NewEntity(idB)
	u.NewEntity(idA, idB)

	stats := w.Stats()

	for b.Loop() {
		stats = w.Stats()
	}
	expectEqual(b, 4, len(stats.Archetypes))
}

func BenchmarkWorldStats16Arch(b *testing.B) {
	w := NewWorld()

	allIDs := []ID{
		ComponentID[CompA](&w),
		ComponentID[CompB](&w),
		ComponentID[CompC](&w),
		ComponentID[CompD](&w),
	}

	ids := []ID{}
	for i := range int(math.Pow(2, float64(len(allIDs)))) {
		for j, id := range allIDs {
			m := 1 << j
			if i&m == m {
				ids = append(ids, id)
			}
		}
		w.Unsafe().NewEntity(ids...)
		ids = ids[:0]
	}

	stats := w.Stats()

	for b.Loop() {
		stats = w.Stats()
	}
	expectEqual(b, 16, len(stats.Archetypes))
}

func BenchmarkWorldStats64Arch(b *testing.B) {
	w := NewWorld()

	allIDs := []ID{
		ComponentID[CompA](&w),
		ComponentID[CompB](&w),
		ComponentID[CompC](&w),
		ComponentID[CompD](&w),
		ComponentID[CompE](&w),
		ComponentID[CompF](&w),
	}

	ids := []ID{}
	for i := range int(math.Pow(2, float64(len(allIDs)))) {
		for j, id := range allIDs {
			m := 1 << j
			if i&m == m {
				ids = append(ids, id)
			}
		}
		w.Unsafe().NewEntity(ids...)
		ids = ids[:0]
	}

	stats := w.Stats()

	for b.Loop() {
		stats = w.Stats()
	}
	expectEqual(b, 64, len(stats.Archetypes))
}

func BenchmarkWorldStats1024Arch(b *testing.B) {
	w := NewWorld()

	allIDs := []ID{
		ComponentID[CompA](&w),
		ComponentID[CompB](&w),
		ComponentID[CompC](&w),
		ComponentID[CompD](&w),
		ComponentID[CompE](&w),
		ComponentID[CompF](&w),
		ComponentID[CompG](&w),
		ComponentID[CompH](&w),
		ComponentID[CompI](&w),
		ComponentID[CompJ](&w),
	}

	ids := []ID{}
	for i := range int(math.Pow(2, float64(len(allIDs)))) {
		for j, id := range allIDs {
			m := 1 << j
			if i&m == m {
				ids = append(ids, id)
			}
		}
		w.Unsafe().NewEntity(ids...)
		ids = ids[:0]
	}

	stats := w.Stats()

	for b.Loop() {
		stats = w.Stats()
	}
	expectEqual(b, 1024, len(stats.Archetypes))
}
