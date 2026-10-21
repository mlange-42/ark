package ecs

import (
	"math/rand/v2"
	"runtime"
	"testing"
)

// Benchmarks with millions of entities. Tables of this size do not fit into
// the CPU caches, so these benchmarks measure the cost of memory access and
// table growth, which the benchmarks with 1000 entities do not show.

const (
	large1M = 1 << 20
	large8M = 1 << 23
)

// newShuffledWorld creates n entities with Position and Velocity. With shuffle,
// all entities are removed in random order and created again. The recycled
// entity IDs then follow a random order relative to the table rows, as in a
// long simulation where entities die and are born at random.
func newShuffledWorld(n int, shuffle bool) (*World, []Entity) {
	w := NewWorld(n)
	mapper := NewMap2[Position, Velocity](w)

	entities := make([]Entity, 0, n)
	collect := func(e Entity, _ *Position, _ *Velocity) {
		entities = append(entities, e)
	}
	mapper.NewBatchFn(n, collect)
	if !shuffle {
		return w, entities
	}

	rng := rand.New(rand.NewPCG(1, 2))
	rng.Shuffle(len(entities), func(i, j int) {
		entities[i], entities[j] = entities[j], entities[i]
	})
	for _, e := range entities {
		w.RemoveEntity(e)
	}
	entities = entities[:0]
	mapper.NewBatchFn(n, collect)
	return w, entities
}

// Query iteration over both components, the reference for the benchmarks below.
func benchmarkQuery2(b *testing.B, n int, shuffle bool) {
	w, _ := newShuffledWorld(n, shuffle)
	filter := NewFilter2[Position, Velocity](w)

	sum := 0.0
	for b.Loop() {
		query := filter.Query()
		for query.Next() {
			pos, vel := query.Get()
			sum += pos.X + vel.X
		}
	}
	runtime.KeepAlive(sum)
}

// Query iteration over one component, with the second component read through
// a Map for the current entity. Map.Get checks the entity's generation and
// looks up its table and row, both in the entity index.
func benchmarkQueryMap(b *testing.B, n int, shuffle bool) {
	w, _ := newShuffledWorld(n, shuffle)
	filter := NewFilter1[Position](w)
	mapper := NewMap1[Velocity](w)

	sum := 0.0
	for b.Loop() {
		query := filter.Query()
		for query.Next() {
			pos := query.Get()
			vel := mapper.Get(query.Entity())
			sum += pos.X + vel.X
		}
	}
	runtime.KeepAlive(sum)
}

func BenchmarkLargeQuery2_1M(b *testing.B)           { benchmarkQuery2(b, large1M, false) }
func BenchmarkLargeQuery2Shuffled_1M(b *testing.B)   { benchmarkQuery2(b, large1M, true) }
func BenchmarkLargeQueryMap_1M(b *testing.B)         { benchmarkQueryMap(b, large1M, false) }
func BenchmarkLargeQueryMapShuffled_1M(b *testing.B) { benchmarkQueryMap(b, large1M, true) }
func BenchmarkLargeQuery2Shuffled_8M(b *testing.B)   { benchmarkQuery2(b, large8M, true) }
func BenchmarkLargeQueryMapShuffled_8M(b *testing.B) { benchmarkQueryMap(b, large8M, true) }

// Removal of random entities from a large table. Each removal moves the last
// row into the gap (swap-remove) and updates the entity index of the moved
// entity.
func benchmarkRemoveRandom(b *testing.B, n int) {
	w, entities := newShuffledWorld(n, false)
	mapper := NewMap2[Position, Velocity](w)
	rng := rand.New(rand.NewPCG(1, 2))

	collect := func(e Entity, _ *Position, _ *Velocity) {
		entities = append(entities, e)
	}

	for b.Loop() {
		for range 1024 {
			i := rng.IntN(len(entities))
			w.RemoveEntity(entities[i])
			last := len(entities) - 1
			entities[i] = entities[last]
			entities = entities[:last]
		}
		b.StopTimer()
		mapper.NewBatchFn(1024, collect)
		b.StartTimer()
	}
}

func BenchmarkRemoveRandom_1M(b *testing.B) { benchmarkRemoveRandom(b, large1M) }
func BenchmarkRemoveRandom_8M(b *testing.B) { benchmarkRemoveRandom(b, large8M) }
