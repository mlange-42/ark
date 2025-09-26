package events

import (
	"fmt"
	"testing"

	"github.com/mlange-42/ark/ecs"
)

var world = ecs.NewWorld()

type Position struct {
	X float64
	Y float64
}

func TestEventsBasic(t *testing.T) {
	// Create an observer.
	ecs.Observe(ecs.OnCreateEntity).
		For(ecs.C[Position]()).
		Do(func(e ecs.Entity) {
			fmt.Printf("%#v\n", e)
		}).
		Register(&world)

	// Create an entity that triggers the observer's callback.
	builder := ecs.NewMap1[Position](&world)
	builder.NewEntity(&Position{X: 10, Y: 11})
}
