// Demonstrates how to use the custom event types.
package main

import (
	"fmt"

	"github.com/mlange-42/ark/ecs"
)

// Position component
type Position struct {
	X int
	Y int
}

// OnTeleport is a custom event type.
var OnTeleport = ecs.NewEventType()

func main() {
	// Create a new World
	world := ecs.NewWorld()

	// Register an observer for the custom event, observing Position
	ecs.Observe1[Position](OnTeleport).
		Do(func(e ecs.Entity, pos *Position) {
			fmt.Printf("Entity teleported to: %#v\n", *pos)
		}).
		Register(&world)

	// Create an entity with Position
	builder := ecs.NewMap1[Position](&world)
	entity := builder.NewEntity(&Position{X: 0, Y: 0})

	// Create a custom event for teleports
	event := world.Event(OnTeleport).For(ecs.C[Position]())
	// Emit the event for the entity
	event.Emit(entity)
}
