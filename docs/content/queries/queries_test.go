package queries

import (
	"testing"

	"github.com/mlange-42/ark/ecs"
)

var world = ecs.NewWorld()

type Position struct {
	X float64
	Y float64
}

type Velocity struct {
	X float64
	Y float64
}

type Altitude struct {
	Z float64
}

func TestQueriesBasic(t *testing.T) {
	// Create a filter.
	filter := ecs.NewFilter2[Position, Velocity](world)
	// Obtain a query.
	query := filter.Query()
	// Iterate the query.
	for query.Next() {
		pos, vel := query.Get()
		pos.X += vel.X
		pos.Y += vel.Y
	}
}

func TestQueriesLock(t *testing.T) {
	// Create a filter.
	filter := ecs.NewFilter1[Altitude](world)
	// Create a slice to collect entities.
	// Ideally, store this permanently for re-use.
	toRemove := []ecs.Entity{}

	query := filter.Query()
	for query.Next() {
		alt := query.Get()
		alt.Z--
		if alt.Z < 0 {
			// Collect entities to remove.
			toRemove = append(toRemove, query.Entity())
		}
	}

	// Do the removal.
	for _, e := range toRemove {
		world.RemoveEntity(e)
	}
	// Reset the slice for re-use.
	toRemove = toRemove[:0]
}

func TestQueriesWith(t *testing.T) {
	// Create a filter.
	filter := ecs.NewFilter1[Position](world).
		With(ecs.C[Velocity](), ecs.C[Altitude]())

	// Obtain a query.
	_ = filter.Query()
	// ...
}

func TestQueriesWith2(t *testing.T) {
	// Create a filter.
	filter := ecs.NewFilter1[Position](world).
		With(ecs.C[Velocity]()).
		With(ecs.C[Altitude]())

	// Obtain a query.
	_ = filter.Query()
	// ...
}

func TestQueriesWithout(t *testing.T) {
	// Create a filter.
	filter := ecs.NewFilter1[Position](world).
		Without(ecs.C[Velocity](), ecs.C[Altitude]())

	// Obtain a query.
	_ = filter.Query()
	// ...
}

func TestQueriesWithout2(t *testing.T) {
	// Create a filter.
	filter := ecs.NewFilter1[Position](world).
		Without(ecs.C[Velocity]()).
		Without(ecs.C[Altitude]())

	// Obtain a query.
	_ = filter.Query()
	// ...
}

func TestQueriesExclusive(t *testing.T) {
	// Create a filter.
	filter := ecs.NewFilter1[Position](world).
		Exclusive()

	// Obtain a query.
	_ = filter.Query()
	// ...
}

func TestQueriesOptional(t *testing.T) {
	// Create a filter.
	filter := ecs.NewFilter2[Position, Velocity](world)
	// Create a component mapper.
	altMap := ecs.NewMap[Altitude](world)

	// Obtain a query.
	query := filter.Query()
	for query.Next() {
		// Get the current entity.
		entity := query.Entity()
		// Check whether the current entity has an Altitude component.
		if altMap.Has(entity) {
			alt := altMap.Get(entity)
			alt.Z += 1.0
		}
		// Do other stuff...
	}
}

func TestQueriesCached(t *testing.T) {
	// Create a filter.
	filter := ecs.NewFilter2[Position, Velocity](world).
		Register() // Register it to the cache.

	// Obtain a query.
	_ = filter.Query()
	// ...
}
