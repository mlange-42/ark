package ecs

import (
	"fmt"
	"testing"
)

func TestCustomEvent(t *testing.T) {
	world := NewWorld()
	customEvent := NewEventType()
	builder := NewMap2[Position, Velocity](&world)

	callCount := 0
	Observe(customEvent).
		For(C[Position]()).
		Do(func(e Entity) { callCount++ }).
		Register(&world)

	e := builder.NewEntity(&Position{1, 2}, &Velocity{3, 4})

	evt := NewEvent(customEvent, &world).For(C[Position]())

	evt.Emit(e)
	expectEqual(t, 1, callCount)
	NewEvent(customEvent, &world).For(C[Velocity]()).Emit(e)
	expectEqual(t, 1, callCount)
	evt.Emit(e)
	expectEqual(t, 2, callCount)
}

func TestCustomEventGeneric(t *testing.T) {
	world := NewWorld()
	customEvent := NewEventType()
	builder := NewMap2[Position, Velocity](&world)

	callCount := 0
	Observe1[Position](customEvent).
		Do(func(e Entity, pos *Position) {
			callCount++
			fmt.Printf("%#v", pos)
		}).
		Register(&world)

	e := builder.NewEntity(&Position{1, 2}, &Velocity{3, 4})

	evt := NewEvent(customEvent, &world).For(C[Position]())
	evt.Emit(e)
	expectEqual(t, 1, callCount)
}

func TestCustomEventEmpty(t *testing.T) {
	world := NewWorld()
	customEvent := NewEventType()
	builder := NewMap2[Position, Velocity](&world)

	callCount := 0
	Observe(customEvent).
		Do(func(e Entity) {
			callCount++
		}).
		Register(&world)

	e := builder.NewEntity(&Position{1, 2}, &Velocity{3, 4})

	evt := NewEvent(customEvent, &world)
	evt.Emit(e)
	expectEqual(t, 1, callCount)
}
