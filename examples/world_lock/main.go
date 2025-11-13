// Demonstrates how to manipulate entities despite the world being locked during query iteration.
package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/mlange-42/ark/ecs"
)

// Position component.
type Position struct {
	X float64
	Y float64
}

// Energy component.
type Energy struct {
	Value float64
}

func main() {
	// Create a new world.
	world := ecs.NewWorld()

	// Create a mapper to build entities with Position and Energy components.
	builder := ecs.NewMap2[Position, Energy](world)

	// Create 1000 entities with random positions and energy values.
	builder.NewBatchFn(1000, func(entity ecs.Entity, p *Position, e *Energy) {
		p.X = rand.Float64() * 10
		p.Y = rand.Float64() * 10
		e.Value = rand.Float64() * 100
	})

	// Run the simulation.
	run(world)
}

func run(world *ecs.World) {
	// Create a filter to query entities with an Energy component that also have an (unused) Position component.
	// Reuse the filter across multiple iterations for maximum performance.
	filter := ecs.NewFilter1[Energy](world).With(ecs.C[Position]())

	// List of entities to remove after the query iteration.
	// Reuse this to avoid allocations.
	toRemove := []ecs.Entity{}

	// Run 100 iterations of the simulation.
	for range 100 {
		// Create a fresh query from the filter.
		query := filter.Query()

		// Print the current number of entities.
		fmt.Printf("%4d entities\n", query.Count())

		// Iterate the query.
		for query.Next() {
			// Get the Energy component and decrease its value.
			en := query.Get()
			en.Value -= 1.0

			// If the energy is depleted, mark the entity for removal.
			// We can't remove it here, because the world is locked during query iteration.
			if en.Value <= 0 {
				toRemove = append(toRemove, query.Entity())
			}
		}

		// Remove all entities that were marked for removal.
		// This is safe to do here, because the world gets unlocked after query iteration.
		for _, e := range toRemove {
			world.RemoveEntity(e)
		}

		// Clear the toRemove list for the next iteration.
		toRemove = toRemove[:0]
	}
}
