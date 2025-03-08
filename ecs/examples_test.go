package ecs_test

import "github.com/mlange-42/ark/ecs"

func NewGrid(width, height int) Grid {
	return Grid{
		Width:  width,
		Height: height,
	}
}

func (g *Grid) Get(x, y int) ecs.Entity {
	return ecs.Entity{}
}

func ExampleFilter2() {
	world := ecs.NewWorld()

	// A simple filter for components (non-exclusive).
	simpleFilter := ecs.NewFilter2[Position, Velocity](&world)

	// A simple filter with an additional (unused) and an excluded component.
	complexFilter := ecs.NewFilter2[Position, Velocity](&world).
		With(ecs.C[Altitude]()).
		Without(ecs.C[ChildOf]())

	// A cached/registered filter, with an additional (unused) component.
	cachedFilter := ecs.NewFilter2[Position, Velocity](&world).
		With(ecs.C[Altitude]()).
		Register()

	// Create a query from a filter, and iterate it.
	query := simpleFilter.Query()
	for query.Next() {
		// ...
	}

	_, _ = complexFilter, cachedFilter
	// Output:
}

func ExampleQuery2() {
	world := ecs.NewWorld()

	// A simple filter.
	filter := ecs.NewFilter2[Position, Velocity](&world)

	// Create a fresh query before iterating.
	query := filter.Query()
	for query.Next() {
		// Access components of the current entity.
		pos, vel := query.Get()
		// Access the current entity itself.
		entity := query.Entity()
		// ...
		_, _, _ = pos, vel, entity
	}
	// Output:
}
