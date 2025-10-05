package ecs

import (
	"sync"
	"testing"
)

func BenchmarkPosVelQueryInlined_1000(b *testing.B) {
	n := 1000
	world := NewWorld(1024)

	mapper := NewMap2[Position, Velocity](&world)
	mapper.NewBatch(n, &Position{}, &Velocity{X: 1, Y: 0})

	filter := NewFilter2[Position, Velocity](&world)
	loop := func() {
		query := filter.Query()
		for query.Next() {
			pos, vel := query.Get()
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
	for b.Loop() {
		loop()
	}
}

func BenchmarkPosVelQuery_1000(b *testing.B) {
	n := 1000
	world := NewWorld(1024)

	mapper := NewMap2[Position, Velocity](&world)
	mapper.NewBatch(n, &Position{}, &Velocity{X: 1, Y: 0})

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

func BenchmarkPosVelQueryCached_1000(b *testing.B) {
	n := 1000
	world := NewWorld(1024)

	mapper := NewMap2[Position, Velocity](&world)
	mapper.NewBatch(n, &Position{}, &Velocity{X: 1, Y: 0})

	filter := NewFilter2[Position, Velocity](&world).Register()
	for b.Loop() {
		query := filter.Query()
		for query.Next() {
			pos, vel := query.Get()
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}

func BenchmarkPosVelQueryUnsafe_1000(b *testing.B) {
	n := 1000
	world := NewWorld(1024)

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	mapper := NewMap2[Position, Velocity](&world)
	mapper.NewBatch(n, &Position{}, &Velocity{X: 1, Y: 0})

	filter := NewUnsafeFilter(&world, posID, velID)
	for b.Loop() {
		query := filter.Query()
		for query.Next() {
			pos := (*Position)(query.Get(posID))
			vel := (*Velocity)(query.Get(velID))
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}

func BenchmarkPosVelQuerySerial_100k(b *testing.B) {
	n := 100_000
	world := NewWorld(1024)

	mapper := NewMap2[Position, Velocity](&world)
	mapper.NewBatch(n, &Position{}, &Velocity{X: 1, Y: 0})

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

func BenchmarkPosVelQueryParallel4_100k(b *testing.B) {
	n := 100_000
	threads := 4
	world := NewWorld(1024)

	parents := make([]Entity, 0, threads)
	for range threads {
		parent := world.NewEntity()
		parents = append(parents, parent)
	}

	mapper := NewMap3[Position, Velocity, ChildOf](&world)
	for _, p := range parents {
		mapper.NewBatch(n/threads, &Position{}, &Velocity{X: 1, Y: 0}, &ChildOf{}, RelIdx(2, p))
	}

	filter := NewFilter2[Position, Velocity](&world).
		With(C[ChildOf]())

	task := func(t Entity, wg *sync.WaitGroup) {
		defer wg.Done()
		query := filter.Query(RelIdx(2, t))
		for query.Next() {
			pos, vel := query.Get()
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}

	for b.Loop() {
		var wg sync.WaitGroup
		wg.Add(threads)
		for _, t := range parents {
			go task(t, &wg)
		}
		wg.Wait()
	}
}

func BenchmarkPosVelMap_1000(b *testing.B) {
	n := 1000
	world := NewWorld(1024)

	mapper := NewMap2[Position, Velocity](&world)

	entities := make([]Entity, 0, n)
	for range n {
		e := world.NewEntity()
		mapper.Add(e, &Position{}, &Velocity{X: 1, Y: 0})
		entities = append(entities, e)
	}

	for b.Loop() {
		for _, e := range entities {
			pos, vel := mapper.Get(e)
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}

func BenchmarkPosVelMap_1000_Unchecked(b *testing.B) {
	n := 1000
	world := NewWorld(1024)

	mapper := NewMap2[Position, Velocity](&world)

	entities := make([]Entity, 0, n)
	for range n {
		e := mapper.NewEntity(&Position{}, &Velocity{X: 1, Y: 0})
		entities = append(entities, e)
	}

	for b.Loop() {
		for _, e := range entities {
			pos, vel := mapper.GetUnchecked(e)
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}

func BenchmarkPosVelMapUnsafe_1000(b *testing.B) {
	n := 1000
	world := NewWorld(1024)
	u := world.Unsafe()

	posID := ComponentID[Position](&world)
	velID := ComponentID[Velocity](&world)

	mapper := NewMap2[Position, Velocity](&world)

	entities := make([]Entity, 0, n)
	for range n {
		e := mapper.NewEntity(&Position{}, &Velocity{X: 1, Y: 0})
		entities = append(entities, e)
	}

	for b.Loop() {
		for _, e := range entities {
			pos := (*Position)(u.Get(e, posID))
			vel := (*Velocity)(u.Get(e, velID))
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}

func BenchmarkCreateQuery2(b *testing.B) {
	w := NewWorld()

	builder := NewMap2[CompA, CompB](&w)
	builder.NewBatchFn(100, nil)
	filter := NewFilter2[CompA, CompB](&w)

	for b.Loop() {
		query := filter.Query()
		query.Close()
	}
}

func BenchmarkCreateQuery2Cached(b *testing.B) {
	w := NewWorld()

	builder := NewMap2[CompA, CompB](&w)
	builder.NewBatchFn(100, nil)
	filter := NewFilter2[CompA, CompB](&w).Register()

	for b.Loop() {
		query := filter.Query()
		query.Close()
	}
}
