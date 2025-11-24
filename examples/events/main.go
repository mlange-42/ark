// Demonstrates how to use the built-in event types.
package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/mlange-42/ark/ecs"
)

// Position component
type Position struct {
	X int
	Y int
}

// Velocity component
type Velocity struct {
	X int
	Y int
}

func main() {
	// Create a new World
	world := ecs.NewWorld()

	// Register an observer for entity creation, observing Position and Velocity
	ecs.Observe2[Position, Velocity](ecs.OnCreateEntity).
		Do(func(e ecs.Entity, pos *Position, vel *Velocity) {
			fmt.Printf("Create entity with %#v, %#v\n", *pos, *vel)
		}).
		Register(world)

	// Register an observer for entity removal, observing Position and Velocity
	ecs.Observe2[Position, Velocity](ecs.OnRemoveEntity).
		Do(func(e ecs.Entity, pos *Position, vel *Velocity) {
			fmt.Printf("Remove entity with %#v, %#v\n", *pos, *vel)
		}).
		Register(world)

	// A mapper for creating entities with Position and Velocity
	builder := ecs.NewMap2[Position, Velocity](world)

	// Create an entity with Position and Velocity
	entity := builder.NewEntity(
		&Position{X: rand.IntN(100), Y: rand.IntN(100)},
		&Velocity{},
	)

	// Remove the entity
	world.RemoveEntity(entity)
}
