package ecs

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/mlange-42/ark/ecs/stats"
)

func TestWorldStats(t *testing.T) {
	w := NewWorld(128, 32)

	posVelMap := NewMap2[Position, Velocity](&w)
	posVelHeadMap := NewMap3[Position, Velocity, Heading](&w)
	posChildMap := NewMap3[Position, ChildOf, ChildOf2](&w)
	filter := NewFilter0(&w)

	p1 := w.NewEntity()
	p2 := w.NewEntity()
	p3 := w.NewEntity()

	posVelMap.NewBatchFn(100, nil)
	posChildMap.NewBatchFn(50, nil, RelIdx(1, p1), RelIdx(2, p2))
	posChildMap.NewBatchFn(50, nil, RelIdx(1, p3), RelIdx(2, p2))

	w.RemoveEntities(filter.Batch(), nil)
	_ = w.Stats()

	p1 = w.NewEntity()
	p2 = w.NewEntity()
	p3 = w.NewEntity()

	posVelMap.NewBatchFn(100, nil)
	posChildMap.NewBatchFn(50, nil, RelIdx(1, p1), RelIdx(2, p2))
	posChildMap.NewBatchFn(50, nil, RelIdx(1, p3), RelIdx(2, p2))

	posVelHeadMap.NewBatchFn(250, nil)
	posChildMap.NewBatchFn(50, nil, RelIdx(1, p2), RelIdx(2, p3))
	_ = w.Stats()

	s := w.Stats()
	fmt.Println(s.String())

	expectEqual(t, 5, len(s.ComponentTypes))
	expectEqual(t, 5, len(s.ComponentTypeNames))
	expectEqual(t, 4, len(s.Archetypes))
	expectEqual(t, 2, len(s.Archetypes[1].ComponentIDs))
	expectEqual(t, 2, len(s.Archetypes[1].ComponentTypes))
	expectEqual(t, 2, len(s.Archetypes[1].ComponentTypeNames))

	sReduced := w.Stats(stats.None)
	expectEqual(t, s.NumArchetypes, sReduced.NumArchetypes)

	sReduced = w.Stats(stats.Archetypes)
	expectEqual(t, s.Archetypes[0].NumTables, sReduced.Archetypes[0].NumTables)
	expectEqual(t, s.Archetypes[1].NumTables, sReduced.Archetypes[1].NumTables)
	expectEqual(t, s.Archetypes[2].NumTables, sReduced.Archetypes[2].NumTables)

	w.RemoveEntities(filter.Batch(), nil)
	s = w.Stats()
	fmt.Println(s.String())
}

func TestWorldStatsFlags(t *testing.T) {
	w := NewWorld(128, 32)

	posVelMap := NewMap2[Position, Velocity](&w)
	posVelHeadMap := NewMap3[Position, Velocity, Heading](&w)
	posChildMap := NewMap2[Position, ChildOf](&w)

	posVelMap.NewBatchFn(100, nil)
	posVelHeadMap.NewBatchFn(100, nil)

	parent := w.NewEntity()
	child := posChildMap.NewEntityFn(nil, Rel[ChildOf](parent))
	w.RemoveEntity(child)
	w.RemoveEntity(parent)

	s := w.Stats(stats.None)
	mem, memUsed := s.Memory, s.MemoryUsed

	expectEqual(t, 4, s.NumArchetypes)
	expectEqual(t, 0, len(s.Archetypes))

	_ = w.Stats(stats.Archetypes)
	s = w.Stats(stats.Archetypes)

	expectEqual(t, mem, s.Memory)
	expectEqual(t, memUsed, s.MemoryUsed)
	expectEqual(t, 4, s.NumArchetypes)
	expectEqual(t, 4, len(s.Archetypes))
	expectEqual(t, 1, s.Archetypes[0].NumTables)
	expectEqual(t, 0, len(s.Archetypes[0].Tables))

	s = w.Stats(stats.All)

	expectEqual(t, mem, s.Memory)
	expectEqual(t, memUsed, s.MemoryUsed)
	expectEqual(t, 4, s.NumArchetypes)
	expectEqual(t, 4, len(s.Archetypes))
	expectEqual(t, 1, s.Archetypes[0].NumTables)
	expectEqual(t, 1, len(s.Archetypes[0].Tables))
}

func TestWorldStatsSizeJSON(t *testing.T) {
	names := []string{"none", "arch", "tabl", "all "}
	for _, comps := range [...]int{2, 4, 6, 10} {
		for _, flags := range [...]stats.Option{stats.All, stats.Archetypes, stats.None} {
			w := createWorld(comps)
			stats := w.Stats(stats.Option(flags))
			js, err := json.Marshal(&stats)
			expectNil(t, err)
			fmt.Printf("JSON size for stats (%s), %4d archetypes: %7.3f kB\n",
				names[flags],
				int(math.Pow(2, float64(comps))),
				float64(len(js))/float64(1024))
		}
	}
}

func benchmarkWorldStats(b *testing.B, comps int, flags stats.Option) {
	w := createWorld(comps)

	stat := w.Stats(flags)

	for b.Loop() {
		stat = w.Stats(flags)
	}
	expectEqual(b, int(math.Pow(2, float64(comps))), stat.NumArchetypes)
	if flags&stats.Archetypes != 0 {
		expectEqual(b, int(math.Pow(2, float64(comps))), len(stat.Archetypes))
		expectEqual(b, 1, stat.Archetypes[0].NumTables)
		if flags&stats.Tables != 0 {
			expectEqual(b, 1, len(stat.Archetypes[0].Tables))
		} else {
			expectEqual(b, 0, len(stat.Archetypes[0].Tables))
		}
	} else {
		expectEqual(b, 0, len(stat.Archetypes))
	}
}

func BenchmarkWorldStats_4Arch_All(b *testing.B) {
	benchmarkWorldStats(b, 2, stats.All)
}

func BenchmarkWorldStats_4Arch_Arches(b *testing.B) {
	benchmarkWorldStats(b, 2, stats.Archetypes)
}

func BenchmarkWorldStats_4Arch_None(b *testing.B) {
	benchmarkWorldStats(b, 2, stats.None)
}

func BenchmarkWorldStats_16Arch_All(b *testing.B) {
	benchmarkWorldStats(b, 4, stats.All)
}

func BenchmarkWorldStats_16Arch_Arches(b *testing.B) {
	benchmarkWorldStats(b, 4, stats.Archetypes)
}

func BenchmarkWorldStats_16Arch_None(b *testing.B) {
	benchmarkWorldStats(b, 4, stats.None)
}

func BenchmarkWorldStats_64Arch_All(b *testing.B) {
	benchmarkWorldStats(b, 6, stats.All)
}

func BenchmarkWorldStats_64Arch_Arches(b *testing.B) {
	benchmarkWorldStats(b, 6, stats.Archetypes)
}

func BenchmarkWorldStats_64Arch_None(b *testing.B) {
	benchmarkWorldStats(b, 6, stats.None)
}

func BenchmarkWorldStats_1024Arch_All(b *testing.B) {
	benchmarkWorldStats(b, 10, stats.All)
}

func BenchmarkWorldStats_1024Arch_Arches(b *testing.B) {
	benchmarkWorldStats(b, 10, stats.Archetypes)
}

func BenchmarkWorldStats_1024Arch_None(b *testing.B) {
	benchmarkWorldStats(b, 10, stats.None)
}

func createWorld(comps int) *World {
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
	usedIDs := allIDs[:comps]

	ids := []ID{}
	for i := range int(math.Pow(2, float64(len(usedIDs)))) {
		for j, id := range usedIDs {
			m := 1 << j
			if i&m == m {
				ids = append(ids, id)
			}
		}
		w.Unsafe().NewEntity(ids...)
		ids = ids[:0]
	}
	return &w
}
