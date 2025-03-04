package operations

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

type Altitude struct {
	Z float64
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

func TestExchange(t *testing.T) {
	// Create an exchange helper
	exchange := ecs.NewExchange1[Altitude](&world).
		Removes(ecs.C[Position](), ecs.C[Velocity]())

	_ = exchange
}
