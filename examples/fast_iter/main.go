// Demonstrates faster, table-based query iteration.
package main

import (
	"math/rand/v2"

	"github.com/mlange-42/ark/ecs"
)

// Position component
type Position struct {
	X, Y float64
}

// Velocity component
type Velocity struct {
	DX, DY float64
}

func main() {
	// Create a new World
	world := ecs.NewWorld()

	// Create a component mapper
	// Save mappers permanently and re-use them for best performance
	mapper := ecs.NewMap2[Position, Velocity](world)

	// Create entities with components
	for range 1000 {
		_ = mapper.NewEntity(
			&Position{X: rand.Float64() * 100, Y: rand.Float64() * 100},
			&Velocity{DX: rand.NormFloat64(), DY: rand.NormFloat64()},
		)
	}

	// Create a filter
	// Save filters permanently and re-use them for best performance
	filter := ecs.NewFilter2[Position, Velocity](world)

	// Time loop
	for range 5000 {
		// Get a fresh query and iterate it
		query := filter.Query()
		// Iterate over tables/archetypes instead of over entities
		for query.NextTable() {
			// Access component columns
			positions, velocities := query.GetColumns()
			// Iterate over individual components (i.e. per entity)
			for i := range positions {
				pos, vel := &positions[i], &velocities[i]
				pos.X += vel.DX
				pos.Y += vel.DY
			}
		}
	}
}
