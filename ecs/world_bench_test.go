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

func BenchmarkCreateEntitiesAlloc(b *testing.B) {
	w := NewWorld()
	builder := NewMap1[Position](&w)
	filter := NewFilter0(&w)

	builder.NewBatchFn(1000, nil)
	w.RemoveEntities(filter.Batch(), nil)

	for b.Loop() {
		b.StopTimer()
		w := NewWorld(8)
		builder := NewMap2[Position, Velocity](&w)
		b.StartTimer()
		for range 1000 {
			builder.NewEntityFn(nil)
		}
	}
}

func BenchmarkWorldLockUnlock(b *testing.B) {
	w := NewWorld()

	for b.Loop() {
		l := w.lock()
		w.unlock(l)
	}
}

func BenchmarkWorldLockUnlockSafe(b *testing.B) {
	w := NewWorld()

	for b.Loop() {
		l := w.lockSafe()
		w.unlockSafe(l)
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

func BenchmarkRemoveTrivial_1000(b *testing.B) {
	w := NewWorld()
	builder := NewMap1[Heading](&w)

	toRemove := []Entity{}
	builder.NewBatchFn(1000, func(entity Entity, a *Heading) {
		toRemove = append(toRemove, entity)
	})
	for _, e := range toRemove {
		w.RemoveEntity(e)
	}

	for b.Loop() {
		b.StopTimer()
		toRemove = toRemove[:0]
		builder.NewBatchFn(1000, func(entity Entity, a *Heading) {
			toRemove = append(toRemove, entity)
		})
		b.StartTimer()
		for _, e := range toRemove {
			w.RemoveEntity(e)
		}
	}
}

func BenchmarkRemoveNonTrivial_1000(b *testing.B) {
	w := NewWorld()
	builder := NewMap1[PointerType](&w)

	toRemove := []Entity{}
	builder.NewBatchFn(1000, func(entity Entity, a *PointerType) {
		toRemove = append(toRemove, entity)
	})
	for _, e := range toRemove {
		w.RemoveEntity(e)
	}

	for b.Loop() {
		b.StopTimer()
		toRemove = toRemove[:0]
		builder.NewBatchFn(1000, func(entity Entity, a *PointerType) {
			toRemove = append(toRemove, entity)
		})
		b.StartTimer()
		for _, e := range toRemove {
			w.RemoveEntity(e)
		}
	}
}

func benchmarkQueryNumArches(b *testing.B, arches int, n int) {
	world := NewWorld(1024)

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)
	allIDs := []ID{
		ComponentID[CompA](&world),
		ComponentID[CompB](&world),
		ComponentID[CompC](&world),
		ComponentID[CompD](&world),
		ComponentID[CompE](&world),
		ComponentID[CompF](&world),
		ComponentID[CompG](&world),
		ComponentID[CompH](&world),
		ComponentID[CompI](&world),
		ComponentID[CompJ](&world),
	}

	extraIDs := allIDs[:int(math.Log2(float64(arches)))]

	ids := []ID{}
	for i := range n {
		ids = append(ids, posID, velID)
		for j, id := range extraIDs {
			m := 1 << j
			if i&m == m {
				ids = append(ids, id)
			}
		}
		world.Unsafe().NewEntity(ids...)

		ids = ids[:0]
	}

	filter := NewFilter2[Position, Velocity](&world)

	for b.Loop() {
		query := filter.Query()
		for query.Next() {
			pos, vel := query.Get()
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}

func BenchmarkQuery1Arch_1024(b *testing.B) {
	benchmarkQueryNumArches(b, 1, 1024)
}

func BenchmarkQuery32Arch_1024(b *testing.B) {
	benchmarkQueryNumArches(b, 32, 1024)
}

func BenchmarkQuery128Arch_1024(b *testing.B) {
	benchmarkQueryNumArches(b, 128, 1024)
}

func BenchmarkQuery1024Arch_1024(b *testing.B) {
	benchmarkQueryNumArches(b, 1024, 1024)
}
