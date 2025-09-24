package ecs

import (
	"testing"
)

func TestNewObserver(t *testing.T) {
	w := NewWorld()
	obs := NewObserver(OnCreateEntity, func(e Entity) {}).Register(&w)

	expectPanicsWithValue(t, "can't modify a registered observer",
		func() {
			obs.With(C[Position]())
		})
	expectPanicsWithValue(t, "can't modify a registered observer",
		func() {
			obs.Without(C[Position]())
		})

	obs = NewObserver(OnAdd, func(e Entity) {})
	expectPanicsWithValue(t, "can use Observer.Without only for OnCreateEntity and OnRemoveEntity events",
		func() {
			obs.Without(C[Position]())
		})
}

func TestObserverRegister(t *testing.T) {
	w := NewWorld()

	obs := NewObserver(OnCreateEntity, func(e Entity) {}).
		With(C[Position]()).
		With(C[Velocity]()).
		Without(C[Heading]()).
		Register(&w)
	expectTrue(t, w.storage.observers.HasObservers(OnCreateEntity))

	expectPanicsWithValue(t, "observer is already registered",
		func() {
			obs.Register(&w)
		})

	obs.Unregister(&w)
	expectFalse(t, w.storage.observers.HasObservers(OnCreateEntity))

	expectPanicsWithValue(t, "observer is not registered",
		func() {
			obs.Unregister(&w)
		})
}

func TestObserverManager(t *testing.T) {
	m := newObserverManager()
	expectPanicsWithValue(t, "can't unregister observer, not found",
		func() {
			m.RemoveObserver(&Observer{id: 13})
		})
}

func TestObserverOnCreateEntity(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)
	builder1 := NewMap1[Position](&w)
	builder2 := NewMap1[Velocity](&w)

	callCount := 0

	NewObserver(OnCreateEntity,
		func(e Entity) {
			callCount++
		}).
		With(C[Position]()).
		Register(&w)

	builder1.NewEntity(&Position{})
	expectEqual(t, 1, callCount)
	builder2.NewEntity(&Velocity{})
	expectEqual(t, 1, callCount)

	w.Unsafe().NewEntity(posID)
	expectEqual(t, 2, callCount)
	w.Unsafe().NewEntity(velID)
	expectEqual(t, 2, callCount)

	builder1.NewBatch(10, &Position{})
	expectEqual(t, 12, callCount)
	builder2.NewBatch(10, &Velocity{})
	expectEqual(t, 12, callCount)
}

func TestObserverOnRemoveEntity(t *testing.T) {
	w := NewWorld()

	builder1 := NewMap1[Position](&w)
	builder2 := NewMap1[Velocity](&w)
	filter1 := NewFilter1[Position](&w)
	filter2 := NewFilter1[Velocity](&w)

	callCount := 0

	NewObserver(OnRemoveEntity,
		func(e Entity) {
			callCount++
		}).
		With(C[Position]()).
		Register(&w)

	e := builder1.NewEntity(&Position{})
	w.RemoveEntity(e)
	expectEqual(t, 1, callCount)
	e = builder2.NewEntity(&Velocity{})
	w.RemoveEntity(e)
	expectEqual(t, 1, callCount)

	builder1.NewBatch(10, &Position{})
	w.RemoveEntities(filter1.Batch(), nil)
	expectEqual(t, 11, callCount)
	builder2.NewBatch(10, &Velocity{})
	w.RemoveEntities(filter2.Batch(), nil)
	expectEqual(t, 11, callCount)
}

func TestObserverOnAddRemoveEntity(t *testing.T) {
	w := NewWorld()

	builder1 := NewMap1[Position](&w)
	builder2 := NewMap1[Velocity](&w)
	filter1 := NewFilter1[Position](&w)
	filter2 := NewFilter1[Velocity](&w)

	callAdd := 0
	callRemove := 0

	NewObserver(OnAdd,
		func(e Entity) {
			callAdd++
		}).
		With(C[Position]()).
		Register(&w)

	NewObserver(OnRemove,
		func(e Entity) {
			callRemove++
		}).
		With(C[Position]()).
		Register(&w)

	e := w.NewEntity()
	builder1.Add(e, &Position{})
	expectEqual(t, 1, callAdd)
	builder2.Add(e, &Velocity{})
	expectEqual(t, 1, callAdd)

	builder2.Remove(e)
	expectEqual(t, 0, callRemove)
	builder1.Remove(e)
	expectEqual(t, 1, callRemove)

	builder2.NewBatch(10, &Velocity{})
	builder1.AddBatch(filter2.Batch(), &Position{})
	expectEqual(t, 11, callAdd)

	builder1.RemoveBatch(filter1.Batch(), nil)
	expectEqual(t, 11, callRemove)
}
