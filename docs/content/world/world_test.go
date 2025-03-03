package world

import (
	"testing"

	"github.com/mlange-42/ark/ecs"
)

func TestWorldSimple(t *testing.T) {
	world := ecs.NewWorld()
	_ = world
}

func TestWorldConfig(t *testing.T) {
	world := ecs.NewWorld(1024)
	_ = world
}

func TestWorldConfigRelations(t *testing.T) {
	world := ecs.NewWorld(1028, 128)
	_ = world
}

func TestWorldReset(t *testing.T) {
	world := ecs.NewWorld()
	// ... do something with the world

	world.Reset()
	// ... start over again
}
