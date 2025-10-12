package ecs_test

import "github.com/mlange-42/ark/ecs"

// Create an event registry
var registry = ecs.EventRegistry{}

// Create event types
var OnCollisionDetected = registry.NewEventType()
var OnInputReceived = registry.NewEventType()

// If all custom events are defined in a single place, this is also possible:
const (
	OnLevelLoaded ecs.EventType = iota
	OnTimerElapsed
)

func ExampleEventRegistry() {
	_, _ = OnCollisionDetected, OnInputReceived
	// Output:
}

func ExampleEventRegistry_iota() {
	_, _ = OnLevelLoaded, OnTimerElapsed
	// Output:
}

func ExampleEventType_onCreateEntity() {
	world := ecs.NewWorld()

	// No filters. Matches any new entity.
	ecs.Observe(ecs.OnCreateEntity).
		Do(func(e ecs.Entity) { /* do something */ }).
		Register(&world)

	// Filter for new entities with certain components.
	ecs.Observe(ecs.OnCreateEntity).
		With(ecs.C[Position](), ecs.C[Velocity]()).
		Do(func(e ecs.Entity) { /* do something */ }).
		Register(&world)

	// Filter for new entities with certain components, but exclude other components.
	ecs.Observe(ecs.OnCreateEntity).
		With(ecs.C[Position](), ecs.C[Velocity]()).
		Without(ecs.C[Altitude]()).
		Do(func(e ecs.Entity) { /* do something */ }).
		Register(&world)
	// Output:
}

func ExampleEventType_onRemoveEntity() {
	world := ecs.NewWorld()

	// No filters. Matches any removed entity.
	ecs.Observe(ecs.OnRemoveEntity).
		Do(func(e ecs.Entity) { /* do something */ }).
		Register(&world)

	// Filter for removed entities with certain components.
	ecs.Observe(ecs.OnRemoveEntity).
		With(ecs.C[Position](), ecs.C[Velocity]()).
		Do(func(e ecs.Entity) { /* do something */ }).
		Register(&world)

	// Filter for removed entities with certain components, but exclude other components.
	ecs.Observe(ecs.OnRemoveEntity).
		With(ecs.C[Position](), ecs.C[Velocity]()).
		Without(ecs.C[Altitude]()).
		Do(func(e ecs.Entity) { /* do something */ }).
		Register(&world)
	// Output:
}

func ExampleEventType_onAddComponents() {
	world := ecs.NewWorld()

	// No filters. Matches any component additions.
	ecs.Observe(ecs.OnAddComponents).
		Do(func(e ecs.Entity) { /* do something */ }).
		Register(&world)

	// Filter for addition of a specific component.
	ecs.Observe(ecs.OnAddComponents).
		For(ecs.C[Position]()).
		Do(func(e ecs.Entity) { /* do something */ }).
		Register(&world)

	// Filter for addition of multiple components, all added at the same time.
	ecs.Observe(ecs.OnAddComponents).
		For(ecs.C[Position](), ecs.C[Velocity]()).
		Do(func(e ecs.Entity) { /* do something */ }).
		Register(&world)

	// Filter for addition of multiple components, equivalent to above example.
	ecs.Observe(ecs.OnAddComponents).
		For(ecs.C[Position]()).
		For(ecs.C[Velocity]()).
		Do(func(e ecs.Entity) { /* do something */ }).
		Register(&world)

	// Filter for addition of a specific component, to an entity with/without certain components.
	ecs.Observe(ecs.OnAddComponents).
		For(ecs.C[Altitude]()).
		With(ecs.C[Position]()).
		Without(ecs.C[Velocity]()).
		Do(func(e ecs.Entity) { /* do something */ }).
		Register(&world)
	// Output:
}

func ExampleEventType_onRemoveComponents() {
	world := ecs.NewWorld()

	// No filters. Matches any component removals.
	ecs.Observe(ecs.OnRemoveComponents).
		Do(func(e ecs.Entity) { /* do something */ }).
		Register(&world)

	// Filter for removals of a specific component.
	ecs.Observe(ecs.OnRemoveComponents).
		For(ecs.C[Position]()).
		Do(func(e ecs.Entity) { /* do something */ }).
		Register(&world)

	// Filter for removals of multiple components, all added at the same time.
	ecs.Observe(ecs.OnRemoveComponents).
		For(ecs.C[Position](), ecs.C[Velocity]()).
		Do(func(e ecs.Entity) { /* do something */ }).
		Register(&world)

	// Filter for removals of multiple components, equivalent to above example.
	ecs.Observe(ecs.OnRemoveComponents).
		For(ecs.C[Position]()).
		For(ecs.C[Velocity]()).
		Do(func(e ecs.Entity) { /* do something */ }).
		Register(&world)

	// Filter for removals of a specific component, from an entity with/without certain components.
	ecs.Observe(ecs.OnRemoveComponents).
		For(ecs.C[Altitude]()).
		With(ecs.C[Position]()).
		Without(ecs.C[Velocity]()).
		Do(func(e ecs.Entity) { /* do something */ }).
		Register(&world)
	// Output:
}

func ExampleEventType_onSetComponents() {
	world := ecs.NewWorld()

	// No filters. Matches setting any component.
	ecs.Observe(ecs.OnSetComponents).
		Do(func(e ecs.Entity) { /* do something */ }).
		Register(&world)

	// Filter for setting of a specific component.
	ecs.Observe(ecs.OnSetComponents).
		For(ecs.C[Position]()).
		Do(func(e ecs.Entity) { /* do something */ }).
		Register(&world)

	// Filter for setting of multiple components, all added at the same time.
	ecs.Observe(ecs.OnSetComponents).
		For(ecs.C[Position](), ecs.C[Velocity]()).
		Do(func(e ecs.Entity) { /* do something */ }).
		Register(&world)

	// Filter for setting of multiple components, equivalent to above example.
	ecs.Observe(ecs.OnSetComponents).
		For(ecs.C[Position]()).
		For(ecs.C[Velocity]()).
		Do(func(e ecs.Entity) { /* do something */ }).
		Register(&world)

	// Filter for setting of a specific component, on an entity with/without certain components.
	ecs.Observe(ecs.OnSetComponents).
		For(ecs.C[Altitude]()).
		With(ecs.C[Position]()).
		Without(ecs.C[Velocity]()).
		Do(func(e ecs.Entity) { /* do something */ }).
		Register(&world)
	// Output:
}
