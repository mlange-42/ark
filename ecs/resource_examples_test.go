package ecs_test

import "github.com/mlange-42/ark/ecs"

func ExampleResource() {
	// Create a world.
	world := ecs.NewWorld()

	// Create a resource.
	gridResource := NewGrid(100, 100)
	// Add it to the world.
	ecs.AddResource(world, &gridResource)

	// Resource access in systems.
	// Create and store a resource accessor.
	gridAccess := ecs.NewResource[Grid](world)

	// Use the resource.
	grid := gridAccess.Get()
	entity := grid.Get(13, 42)
	_ = entity
	// Output:
}

func ExampleResource_New() {
	// Create a world.
	world := ecs.NewWorld()

	// Declare resource accessor, e.g. in struct definition.
	var gridAccess ecs.Resource[Grid]

	// Construct the accessor elsewhere, e.g. in the constructor.
	gridAccess = gridAccess.New(world)
	// Output:
}
