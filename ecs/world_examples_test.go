package ecs_test

import (
	"time"

	"github.com/mlange-42/ark/ecs"
)

func ExampleNewWorld() {
	// Create a world with default setting
	world := ecs.NewWorld()
	_ = &world
	// Output:
}

func ExampleNewWorld_capacity() {
	// Create a world with an initial capacity
	world := ecs.NewWorld(2048)
	_ = &world
	// Output:
}

func ExampleNewWorld_relationCapacity() {
	// Create a world with an initial capacity
	// and a lower capacity for tables with relations
	world := ecs.NewWorld(2048, 256)
	_ = &world
	// Output:
}

func ExampleWorld_Shrink() {
	// Create a world
	world := ecs.NewWorld()

	// ... do a lot of work with the world

	// Shrink to free unoccupied memory
	world.Shrink()
	// Output:
}

func ExampleWorld_Shrink_limit() {
	// Create a world
	world := ecs.NewWorld()

	// ... do a lot of work with the world

	// Shrink with a time limit,
	// distribute repeated calls over frames/updates
	for world.Shrink(2 * time.Millisecond) {
	}
	// Output:
}
