package unsafe_test

import (
	"testing"

	"github.com/mlange-42/ark/ecs"
)

var world = ecs.NewWorld()
var posID = ecs.ComponentID[Position](world)
var velID = ecs.ComponentID[Velocity](world)

type Position struct {
	X float64
	Y float64
}

type Velocity struct {
	X float64
	Y float64
}

func TestUnsafeIDs(t *testing.T) {
	world := ecs.NewWorld()

	posID := ecs.ComponentID[Position](world)
	velID := ecs.ComponentID[Velocity](world)

	_, _ = posID, velID
}

func TestUnsafeNewEntity(t *testing.T) {
	entity := world.Unsafe().NewEntity(posID, velID)
	_ = entity
}

func TestUnsafeQuery(t *testing.T) {
	filter := ecs.NewUnsafeFilter(world, posID, velID)

	query := filter.Query()
	for query.Next() {
		pos := (*Position)(query.Get(posID))
		vel := (*Velocity)(query.Get(velID))
		pos.X += vel.X
		pos.Y += vel.Y
	}
}

func TestUnsafeGet(t *testing.T) {
	entity := world.Unsafe().NewEntity(posID, velID)

	pos := (*Position)(world.Unsafe().Get(entity, posID))
	_ = pos
}

func TestUnsafeComponents(t *testing.T) {
	entity := world.NewEntity()

	world.Unsafe().Add(entity, posID, velID)
	world.Unsafe().Remove(entity, posID, velID)
}
