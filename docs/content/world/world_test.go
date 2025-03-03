package world

import (
	"testing"

	"github.com/mlange-42/ark/ecs"
)

func TestWorldSimple(t *testing.T) {
	world := ecs.NewWorld(1024)
	_ = world
}

func TestWorldReset(t *testing.T) {
	world := ecs.NewWorld(1024)
	// ... do something with the world

	world.Reset()
	// ... start over again
}
