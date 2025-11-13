package ecs_test

import (
	"math/rand/v2"

	"github.com/mlange-42/ark/ecs"
)

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
	simpleFilter := ecs.NewFilter2[Position, Velocity](world)

	// A simple filter with an additional (unused) and an excluded component.
	complexFilter := ecs.NewFilter2[Position, Velocity](world).
		With(ecs.C[Altitude]()).
		Without(ecs.C[ChildOf]())

	// A cached/registered filter, with an additional (unused) component.
	cachedFilter := ecs.NewFilter2[Position, Velocity](world).
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

func ExampleFilter2_New() {
	world := ecs.NewWorld()

	// Declare the filter, e.g. in your system struct.
	var filter *ecs.Filter2[Position, Velocity]

	// Construct the filter, avoiding repeated listing of generics.
	filter = filter.New(world)
	// Output:
}

func ExampleQuery2() {
	world := ecs.NewWorld()

	// A simple filter.
	filter := ecs.NewFilter2[Position, Velocity](world)

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

func ExampleQuery2_Count() {
	world := ecs.NewWorld()

	// Create a filter.
	filter := ecs.NewFilter2[Position, Velocity](world)

	// Query the filter
	query := filter.Query()

	// Get the number of entities matching the query
	count := query.Count()

	// Close the query
	// This is required if the query is not iterated
	query.Close()

	_ = count
	// Output:
}

func ExampleQuery2_EntityAt() {
	world := ecs.NewWorld()

	// Create entities.
	builder := ecs.NewMap2[Position, Velocity](world)
	builder.NewBatch(100, &Position{}, &Velocity{})

	// Create a filter.
	filter := ecs.NewFilter2[Position, Velocity](world)

	// Query the filter
	query := filter.Query()

	// Get the number of entities matching the query
	count := query.Count()

	// Get random entities from the query
	e1 := query.EntityAt(rand.IntN(count))
	e2 := query.EntityAt(rand.IntN(count))
	e3 := query.EntityAt(rand.IntN(count))

	// Close the query
	// This is required if the query is not iterated
	query.Close()

	_, _, _ = e1, e2, e3
	// Output:
}

func ExampleMap() {
	world := ecs.NewWorld()

	// Create a component mapper.
	mapper := ecs.NewMap[Position](world)

	// Create an entity.
	entity := mapper.NewEntity(&Position{X: 100, Y: 100})

	// Remove component from the entity.
	mapper.Remove(entity)
	// Add component to the entity.
	mapper.Add(entity, &Position{X: 100, Y: 100})
	// Output:
}

func ExampleMap_New() {
	world := ecs.NewWorld()

	// Declare the mapper, e.g. in your system struct.
	var mapper *ecs.Map[Position]

	// Construct the mapper, avoiding repeated generics.
	mapper = mapper.New(world)
	// Output:
}

func ExampleMap2() {
	world := ecs.NewWorld()

	// Create a component mapper.
	mapper := ecs.NewMap2[Position, Velocity](world)

	// Create an entity.
	entity := mapper.NewEntity(&Position{X: 100, Y: 100}, &Velocity{X: 1, Y: -1})

	// Remove components from the entity.
	mapper.Remove(entity)
	// Add components to the entity.
	mapper.Add(entity, &Position{X: 100, Y: 100}, &Velocity{X: 1, Y: -1})
	// Output:
}

func ExampleMap2_New() {
	world := ecs.NewWorld()

	// Declare the mapper, e.g. in your system struct.
	var mapper *ecs.Map2[Position, Velocity]

	// Construct the mapper, avoiding repeated listing of generics.
	mapper = mapper.New(world)
	// Output:
}

func ExampleExchange2() {
	world := ecs.NewWorld()

	// Create a component mapper.
	mapper := ecs.NewMap[Altitude](world)

	// Create an exchange helper.
	// Adds Position and Velocity, removes Altitude.
	exchange := ecs.NewExchange2[Position, Velocity](world).
		Removes(ecs.C[Altitude]())

	// Create an entity with an Altitude component.
	entity := mapper.NewEntity(&Altitude{Z: 10_000})

	// Remove Altitude and add Position and Velocity.
	exchange.Exchange(entity, &Position{X: 100, Y: 100}, &Velocity{X: 1, Y: -1})

	// Create another entity.
	entity = mapper.NewEntity(&Altitude{Z: 10_000})

	// Remove Altitude.
	exchange.Remove(entity)

	// Add Position and Velocity.
	exchange.Add(entity, &Position{X: 100, Y: 100}, &Velocity{X: 1, Y: -1})
	// Output:
}

func ExampleExchange2_New() {
	world := ecs.NewWorld()

	// Declare the exchange helper, e.g. in your system struct.
	var exchange *ecs.Exchange2[Position, Velocity]

	// Construct the exchange helper, avoiding repeated listing of generics.
	exchange = exchange.New(world).Removes(ecs.C[Altitude]())
	// Output:
}
