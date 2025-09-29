package ecs_test

import "github.com/mlange-42/ark/ecs"

func ExampleObserver() {
	world := ecs.NewWorld()

	// Create an Observer.
	obs := ecs.Observe(ecs.OnCreateEntity).
		With(ecs.C[Position]()).
		Do(func(e ecs.Entity) {
			// do something
		}).
		Register(&world)

	// Unregister later if required.
	obs.Unregister(&world)
	// Output:
}

func ExampleObserver_For() {
	world := ecs.NewWorld()

	_ = ecs.Observe(ecs.OnAddComponents).
		For(ecs.C[Position]()).
		For(ecs.C[Velocity]()).
		Do(func(e ecs.Entity) {
			// do something
		}).
		Register(&world)
	// Output:
}

func ExampleObserver_With() {
	world := ecs.NewWorld()

	_ = ecs.Observe(ecs.OnCreateEntity).
		With(ecs.C[Position]()).
		With(ecs.C[Velocity]()).
		Do(func(e ecs.Entity) {
			// do something
		}).
		Register(&world)
	// Output:
}

func ExampleObserver_Without() {
	world := ecs.NewWorld()

	_ = ecs.Observe(ecs.OnCreateEntity).
		Without(ecs.C[Position]()).
		Do(func(e ecs.Entity) {
			// do something
		}).
		Register(&world)
	// Output:
}

func ExampleObserver_Exclusive() {
	world := ecs.NewWorld()

	_ = ecs.Observe(ecs.OnCreateEntity).
		With(ecs.C[Position]()).
		With(ecs.C[Velocity]()).
		Exclusive().
		Do(func(e ecs.Entity) {
			// do something
		}).
		Register(&world)
	// Output:
}

func ExampleObserver1() {
	world := ecs.NewWorld()

	obs := ecs.Observe1[Position](ecs.OnAddComponents).
		Do(func(e ecs.Entity, pos *Position) {
			// do something
		}).
		Register(&world)

	// Unregister later if required.
	obs.Unregister(&world)
	// Output:
}

func ExampleObserver1_For() {
	world := ecs.NewWorld()

	_ = ecs.Observe1[Position](ecs.OnAddComponents).
		For(ecs.C[Velocity]()).
		Do(func(e ecs.Entity, pos *Position) {
			// do something
		}).
		Register(&world)
	// Output:
}

func ExampleObserver1_With() {
	world := ecs.NewWorld()

	_ = ecs.Observe1[Position](ecs.OnAddComponents).
		With(ecs.C[Velocity]()).
		With(ecs.C[Altitude]()).
		Do(func(e ecs.Entity, pos *Position) {
			// do something
		}).
		Register(&world)
	// Output:
}

func ExampleObserver1_Without() {
	world := ecs.NewWorld()

	_ = ecs.Observe1[Position](ecs.OnAddComponents).
		Without(ecs.C[Velocity]()).
		Do(func(e ecs.Entity, pos *Position) {
			// do something
		}).
		Register(&world)
	// Output:
}

func ExampleObserver1_Exclusive() {
	world := ecs.NewWorld()

	_ = ecs.Observe1[Position](ecs.OnAddComponents).
		With(ecs.C[Velocity]()).
		With(ecs.C[Altitude]()).
		Exclusive().
		Do(func(e ecs.Entity, pos *Position) {
			// do something
		}).
		Register(&world)
	// Output:
}

func ExampleEvent() {
	// Define a custom event type.
	var registry = ecs.NewEventRegistry()
	var OnTeleport = registry.NewEventType()

	world := ecs.NewWorld()

	// Create an event.
	event := world.Event(OnTeleport).
		For(ecs.C[Position]())

	// Create an entity.
	builder := ecs.NewMap1[Position](&world)
	entity := builder.NewEntity(&Position{1, 2})

	// Emit the event.
	event.Emit(entity)
	// Output:
}
