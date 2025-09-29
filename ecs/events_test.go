package ecs

import (
	"fmt"
	"math"
	"testing"
)

var customEvent = NewEventType()

func TestNewEventType(t *testing.T) {
	e := NewEventType()
	expectEqual(t, eventsEnd+1, e)
	e = NewEventType()
	expectEqual(t, eventsEnd+2, e)

	for {
		e = NewEventType()
		if e == math.MaxUint8 {
			break
		}
	}

	expectPanicsWithValue(t, "reached maximum number of custom event types",
		func() { NewEventType() })

	nextUserEvent = eventsEnd
}

func TestCustomEvent(t *testing.T) {
	world := NewWorld()
	builder := NewMap2[Position, Velocity](&world)

	callCount := 0
	Observe(customEvent).
		For(C[Position]()).
		Do(func(e Entity) { callCount++ }).
		Register(&world)

	e := builder.NewEntity(&Position{1, 2}, &Velocity{3, 4})

	evt := world.Event(customEvent).For(C[Position]())

	evt.Emit(e)
	expectEqual(t, 1, callCount)
	world.Event(customEvent).For(C[Velocity]()).Emit(e)
	expectEqual(t, 1, callCount)
	evt.Emit(e)
	expectEqual(t, 2, callCount)
}

func TestCustomEventZero(t *testing.T) {
	world := NewWorld()

	callCount := 0
	Observe(customEvent).
		Do(func(e Entity) { callCount++ }).
		Register(&world)

	evt := world.Event(customEvent)
	evt.Emit(Entity{})
	expectEqual(t, 1, callCount)
}

func TestCustomEventGeneric(t *testing.T) {
	world := NewWorld()
	builder := NewMap2[Position, Velocity](&world)

	callCount := 0
	Observe1[Position](customEvent).
		Do(func(e Entity, pos *Position) {
			callCount++
			fmt.Printf("%#v", pos)
		}).
		Register(&world)

	e := builder.NewEntity(&Position{1, 2}, &Velocity{3, 4})

	evt := world.Event(customEvent).For(C[Position]())
	evt.Emit(e)
	expectEqual(t, 1, callCount)
}

func TestCustomEventEmpty(t *testing.T) {
	world := NewWorld()
	builder := NewMap2[Position, Velocity](&world)

	callCount := 0
	Observe(customEvent).
		Do(func(e Entity) {
			callCount++
		}).
		Register(&world)

	e := builder.NewEntity(&Position{1, 2}, &Velocity{3, 4})

	evt := world.Event(customEvent)
	evt.Emit(e)
	expectEqual(t, 1, callCount)
}

func TestCustomEventErrors(t *testing.T) {
	world := NewWorld()

	Observe1[Position](customEvent).
		Do(func(e Entity, p *Position) {}).
		Register(&world)

	e := world.NewEntity()

	expectPanicsWithValue(t, "only custom events can be emitted manually",
		func() {
			world.Event(OnCreateEntity)
		})

	expectPanicsWithValue(t, "entity does not have the required event components",
		func() {
			world.Event(customEvent).For(C[Position]()).Emit(e)
		})

	world.RemoveEntity(e)
	expectPanicsWithValue(t, "can't emit an event for a dead entity",
		func() {
			world.Event(customEvent).Emit(e)
		})
}

func BenchmarkEventEmit(b *testing.B) {
	w := NewWorld()
	builder := NewMap1[Position](&w)
	e := builder.NewEntity(&Position{})

	evt := w.Event(customEvent)

	for b.Loop() {
		evt.Emit(e)
	}
}

func BenchmarkEventCreateEmit(b *testing.B) {
	w := NewWorld()
	builder := NewMap1[Position](&w)
	e := builder.NewEntity(&Position{})

	for b.Loop() {
		w.Event(customEvent).Emit(e)
	}
}

func BenchmarkEventCreateForEmit(b *testing.B) {
	w := NewWorld()
	builder := NewMap1[Position](&w)
	e := builder.NewEntity(&Position{})

	for b.Loop() {
		w.Event(customEvent).For(C[Position]()).Emit(e)
	}
}
