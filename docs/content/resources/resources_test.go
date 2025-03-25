package resources

import (
	"testing"

	"github.com/mlange-42/ark/ecs"
)

var world = ecs.NewWorld()

type Grid struct{}

func NewGrid(sx, sy int) Grid {
	return Grid{}
}

func TestAddResource(t *testing.T) {
	// Create a resource.
	var worldGrid Grid = NewGrid(100, 100)
	// Add it to the world.
	ecs.AddResource(&world, &worldGrid)
}

func TestResourceWorld(t *testing.T) {
	// Get a resource from the world.
	grid := ecs.GetResource[Grid](&world)
	_ = grid
}

func TestResourceMapper(t *testing.T) {
	// In your system, create a resource mapper.
	// Store it permanently and reuse it for best performance.
	gridRes := ecs.NewResource[Grid](&world)

	// Access the resource.
	grid := gridRes.Get()
	_ = grid
}

func TestResourceMapperAddRemove(t *testing.T) {
	// In your system, create a resource mapper.
	// Store it permanently and reuse it for best performance.
	gridRes := ecs.NewResource[Grid](&world)

	// Check for existence of the resource.
	if gridRes.Has() {
		// Remove the resource if it exists.
		gridRes.Remove()
	} else {
		// Add a new one otherwise.
		grid := NewGrid(100, 100)
		gridRes.Add(&grid)
	}
}
