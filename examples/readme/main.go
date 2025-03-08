// The minimal example from the README,
// for automatic testing by the GitHub CI.
package main

import (
	"math/rand/v2"

	"github.com/mlange-42/ark/ecs"
)

// Position component
type Position struct {
	X float64
	Y float64
}

// Velocity component
type Velocity struct {
	X float64
	Y float64
}

func main() {
	// Create a new World.
	world := ecs.NewWorld()

	// Create a component mapper.
	mapper := ecs.NewMap2[Position, Velocity](&world)

	// Create entities.
	for range 1000 {
		// Create a new Entity with components.
		_ = mapper.NewEntity(
			&Position{X: rand.Float64() * 100, Y: rand.Float64() * 100},
			&Velocity{X: rand.NormFloat64(), Y: rand.NormFloat64()},
		)
	}

	// Create a filter.
	filter := ecs.NewFilter2[Position, Velocity](&world)

	// Time loop.
	for range 5000 {
		// Get a fresh query.
		query := filter.Query()
		// Iterate it
		for query.Next() {
			// Component access through the Query.
			pos, vel := query.Get()
			// Update component fields.
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}
