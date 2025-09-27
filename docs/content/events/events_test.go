package events

import (
	"fmt"
	"testing"

	"github.com/mlange-42/ark/ecs"
)

var world = ecs.NewWorld()
var builder = ecs.NewMap1[Position](&world)

var entity = builder.NewEntity(&Position{})
var uiElement = world.NewEntity()

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

func TestEventsBasic(t *testing.T) {
	// Create an observer.
	ecs.Observe1[Position](ecs.OnCreateEntity).
		Do(func(e ecs.Entity, pos *Position) {
			fmt.Printf("%#v\n", pos)
		}).
		Register(&world)

	// Create an entity that triggers the observer's callback.
	builder := ecs.NewMap1[Position](&world)
	builder.NewEntity(&Position{X: 10, Y: 11})
}

func TestCombineObservers(t *testing.T) {
	// Common callback.
	fn := func(evt ecs.EventType, entity ecs.Entity, pos *Position) {
		if evt == ecs.OnAddComponents {
			// do something
		}
		if evt == ecs.OnRemoveComponents {
			// do something
		}
	}

	// Observer for adding components.
	ecs.Observe1[Position](ecs.OnAddComponents).
		Do(func(e ecs.Entity, pos *Position) { fn(ecs.OnAddComponents, e, pos) }).
		Register(&world)

	// Observer for removing components.
	ecs.Observe1[Position](ecs.OnRemoveComponents).
		Do(func(e ecs.Entity, pos *Position) { fn(ecs.OnRemoveComponents, e, pos) }).
		Register(&world)
}

func TestObserveCreate(t *testing.T) {
	ecs.Observe1[Position](ecs.OnCreateEntity).
		Do(func(e ecs.Entity, p *Position) { /* ... */ })

	ecs.Observe(ecs.OnCreateEntity).
		With(ecs.C[Position]()).
		Do(func(e ecs.Entity) { /* ... */ })
}

func TestObserve2Create(t *testing.T) {
	ecs.Observe2[Position, Velocity](ecs.OnCreateEntity).
		Do(func(e ecs.Entity, p *Position, v *Velocity) { /* ... */ })

	ecs.Observe1[Position](ecs.OnCreateEntity).
		With(ecs.C[Velocity]()).
		Do(func(e ecs.Entity, p *Position) { /* ... */ })
}

func TestObserveCreateEmpty(t *testing.T) {
	ecs.Observe(ecs.OnCreateEntity).
		Do(func(e ecs.Entity) { /* ... */ })
}

func TestObserveAdd(t *testing.T) {
	ecs.Observe1[Position](ecs.OnAddComponents).
		Do(func(e ecs.Entity, p *Position) { /* ... */ })
}

func TestObserveAddWith(t *testing.T) {
	ecs.Observe1[Position](ecs.OnAddComponents).
		With(ecs.C[Velocity]()).
		Without(ecs.C[Altitude]()).
		Do(func(e ecs.Entity, p *Position) { /* ... */ })
}

func TestNewEventType(t *testing.T) {
	var OnSomethingHappened = ecs.NewEventType()

	_ = OnSomethingHappened
}

func TestEventEmit(t *testing.T) {
	// Define the event type.
	var OnSynchronized = ecs.NewEventType()

	// Add an observer for the event type.
	ecs.Observe1[Position](OnSynchronized).
		Do(func(e ecs.Entity, p *Position) { /*...*/ }).
		Register(&world)

	// Define the event.
	event := world.Event(OnSynchronized).
		For(ecs.C[Position]())

	// Emit the event for an entity.
	event.Emit(entity)
}

func TestEventClick(t *testing.T) {
	// Define the event type.
	var OnClick = ecs.NewEventType()

	// Emit a click event.
	world.Event(OnClick).Emit(uiElement)
}

func TestEventZeroEntity(t *testing.T) {
	// Define the event type.
	var OnGameOver = ecs.NewEventType()

	// Emit a game over event.
	world.Event(OnGameOver).Emit(ecs.Entity{})
}
