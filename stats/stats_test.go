package stats

import (
	"fmt"
	"testing"

	"github.com/mlange-42/ark/ecs"
)

// Position component
type Position struct {
	X float64
	Y float64
}

// Heading component
type Heading struct {
	Angle float64
}

func TestWorldStats(t *testing.T) {
	world := ecs.NewWorld()

	builder := ecs.NewMap2[Position, Heading](&world)
	builder.NewBatchFn(100, nil)

	stats := world.Stats()
	fmt.Println(stats)
}
