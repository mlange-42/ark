package batch

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/mlange-42/ark/ecs"
)

var world = ecs.NewWorld()

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

func TestNewBatch(t *testing.T) {
	// Create a component mapper.
	mapper := ecs.NewMap2[Position, Velocity](&world)
	// Create a batch of 100 entities.
	mapper.NewBatch(100, &Position{}, &Velocity{X: 1, Y: -1})
}

func TestNewBatchFn(t *testing.T) {
	// Create a component mapper.
	mapper := ecs.NewMap2[Position, Velocity](&world)
	// Create a batch of 100 entities using a callback.
	mapper.NewBatchFn(100, func(entity ecs.Entity, pos *Position, vel *Velocity) {
		pos.X = rand.Float64() * 100
		pos.Y = rand.Float64() * 100

		vel.X = rand.NormFloat64()
		vel.Y = rand.NormFloat64()
	})
}

func TestBatchComponents(t *testing.T) {
	// Create a component mapper.
	mapper := ecs.NewMap2[Position, Velocity](&world)
	// Create some entities.
	for range 100 {
		world.NewEntity()
	}

	// Create a filter.
	filter := ecs.NewFilter0(&world)
	// Batch-add components.
	mapper.AddBatch(filter.Batch(), &Position{}, &Velocity{X: 1, Y: -1})

	// Batch-Remove components. The optional callback is not used here.
	mapper.RemoveBatch(filter.Batch(), nil)

	// Alternatively, add and initialize components with a callback.
	mapper.AddBatchFn(filter.Batch(), func(entity ecs.Entity, pos *Position, vel *Velocity) {
		// ...
	})
}

func TestRemoveEntities(t *testing.T) {
	// Create a component mapper.
	mapper := ecs.NewMap2[Position, Velocity](&world)
	// Create some entities.
	mapper.NewBatch(10, &Position{}, &Velocity{X: 1, Y: -1})

	// Create a filter.
	filter := ecs.NewFilter2[Position, Velocity](&world)
	// Remove all matching entities. The callback can also be nil.
	world.RemoveEntities(filter.Batch(), func(entity ecs.Entity) {
		fmt.Println("Removing", entity)
	})
}
