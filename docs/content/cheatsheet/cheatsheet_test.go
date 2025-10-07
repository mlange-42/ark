package cheatsheet

import (
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

func TestCreateEntities(t *testing.T) {
	// Create a component mapper
	mapper := ecs.NewMap2[Position, Velocity](&world)

	// Create an entity, giving components
	mapper.NewEntity(&Position{}, &Velocity{X: 1, Y: -1})
	// Create an entity, using a callback
	mapper.NewEntityFn(func(pos *Position, vel *Velocity) {
		pos.X, pos.Y = 0, 0
		vel.X, vel.Y = 1, -1
	})

	// Create many entities, all with the same component values
	mapper.NewBatch(100, &Position{}, &Velocity{X: 1, Y: -1})
	// Create many entities, using a callback (set values individually)
	mapper.NewBatchFn(100, func(e ecs.Entity, pos *Position, vel *Velocity) {
		pos.X, pos.Y = rand.Float64()*100, rand.Float64()*100
		vel.X, vel.Y = rand.NormFloat64(), rand.NormFloat64()
	})
}

func TestRemoveEntities(t *testing.T) {

}
