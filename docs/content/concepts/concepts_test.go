package world

import (
	"testing"

	"github.com/mlange-42/ark/ecs"
)

var world = ecs.NewWorld()
var entity = world.NewEntity()

type Position struct {
	X float64
	Y float64
}

type Velocity struct {
	X float64
	Y float64
}

type Grid struct{}

func NewGrid(sx, sy int) Grid {
	return Grid{}
}

func TestWorldSimple(t *testing.T) {
	world := ecs.NewWorld()
	_ = world
}

func TestWorldConfig(t *testing.T) {
	world := ecs.NewWorld(1024)
	_ = world
}

func TestWorldReset(t *testing.T) {
	world := ecs.NewWorld()
	// ... do something with the world.

	world.Reset()
	// ... start over again.
}

func TestCreateEntitySimple(t *testing.T) {
	entity := world.NewEntity()
	_ = entity
}

func TestRemoveEntity(t *testing.T) {
	world.RemoveEntity(entity)
}

func TestEntityAlive(t *testing.T) {
	alive := world.Alive(entity)
	if alive {
		// ...
	}
}

func TestComponentMapper(t *testing.T) {
	// Create a component mapper.
	mapper := ecs.NewMap2[Position, Velocity](&world)

	// Create an entity with components.
	entity1 := mapper.NewEntity(
		&Position{X: 0, Y: 0},
		&Velocity{X: 1, Y: -1},
	)

	// Create an entity without components.
	entity2 := world.NewEntity()
	// Add components to it.
	mapper.Add(
		entity2,
		&Position{X: 0, Y: 0},
		&Velocity{X: 1, Y: -1},
	)
	// Remove components.
	mapper.Remove(entity2)

	// Remove the entities.
	world.RemoveEntity(entity1)
	world.RemoveEntity(entity2)
}

func TestComponentMapperGet(t *testing.T) {
	// Create a component mapper.
	mapper := ecs.NewMap2[Position, Velocity](&world)

	// Create an entity with components.
	entity1 := mapper.NewEntity(
		&Position{X: 0, Y: 0},
		&Velocity{X: 1, Y: -1},
	)

	// Get mapped components for an entity.
	pos, vel := mapper.Get(entity1)

	_, _ = pos, vel
}

func TestQuery(t *testing.T) {
	// Create a filter.
	filter := ecs.NewFilter2[Position, Velocity](&world)
	// Obtain a query.
	query := filter.Query()
	// Iterate the query.
	for query.Next() {
		pos, vel := query.Get()
		pos.X += vel.X
		pos.Y += vel.Y
	}
}

func TestResource(t *testing.T) {
	// Create a resource.
	var worldGrid Grid = NewGrid(100, 100)
	// Add it to the world.
	ecs.AddResource(&world, &worldGrid)

	// In your system, create a resource mapper.
	gridRes := ecs.NewResource[Grid](&world)
	// Access the resource.
	grid := gridRes.Get()
	_ = grid
}
