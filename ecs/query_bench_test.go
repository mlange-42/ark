package ecs

import "testing"

func BenchmarkPosVelQuery_1000(b *testing.B) {
	n := 1000
	world := NewWorld(128)

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
	world := NewWorld(128)

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
	world := NewWorld(128)

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
