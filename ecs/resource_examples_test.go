package ecs_test

import "github.com/mlange-42/ark/ecs"

func ExampleResource() {
	// Create a world.
	world := ecs.NewWorld()

	// Create a resource.
	gridResource := NewGrid(100, 100)
	// Add it to the world.
	ecs.AddResource(&world, &gridResource)

	// Resource access in systems.
	// Create and store a resource accessor.
	gridAccess := ecs.NewResource[Grid](&world)

	// Use the resource.
	grid := gridAccess.Get()
	entity := grid.Get(13, 42)
	_ = entity
	// Output:
}
