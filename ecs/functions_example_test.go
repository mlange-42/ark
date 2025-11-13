package ecs_test

import "github.com/mlange-42/ark/ecs"

func ExampleAddResource() {
	world := ecs.NewWorld()

	gridResource := NewGrid(100, 100)
	ecs.AddResource(world, &gridResource)
	// Output:
}

func ExampleGetResource() {
	world := ecs.NewWorld()

	gridResource := NewGrid(100, 100)
	ecs.AddResource(world, &gridResource)

	grid := ecs.GetResource[Grid](world)
	entity := grid.Get(13, 42)
	_ = entity
	// Output:
}
