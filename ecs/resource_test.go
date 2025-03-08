package ecs_test

import (
	"testing"

	"github.com/mlange-42/ark/ecs"
	"github.com/stretchr/testify/assert"
)

func TestResource(t *testing.T) {
	w := ecs.NewWorld(1024)
	get := ecs.NewResource[Grid](&w)

	assert.False(t, get.Has())
	gridResource := NewGrid(100, 200)
	get.Add(&gridResource)

	assert.True(t, get.Has())
	grid := get.Get()

	assert.Equal(t, Grid{100, 200}, *grid)

	get.Remove()
	assert.False(t, get.Has())
}

func ExampleResource() {
	// Create a world.
	world := ecs.NewWorld(1024)

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
