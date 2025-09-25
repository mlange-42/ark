package ecs

import (
	"testing"
)

func TestNewObserver(t *testing.T) {
	w := NewWorld()
	obs := NewObserver(OnCreateEntity).Do(func(e Entity) {}).Register(&w)

	expectPanicsWithValue(t, "can't modify a registered observer",
		func() {
			obs.With(C[Position]())
		})
	expectPanicsWithValue(t, "can't modify a registered observer",
		func() {
			obs.Without(C[Position]())
		})

	expectPanicsWithValue(t, "can use Observer.For only for OnAddComponents, OnRemoveComponents and OnSetComponents events",
		func() {
			NewObserver(OnCreateEntity).For(C[Position]())
		})
	expectPanicsWithValue(t, "can use Observer.For only for OnAddComponents, OnRemoveComponents and OnSetComponents events",
		func() {
			NewObserver(OnRemoveEntity).For(C[Position]())
		})

	expectPanicsWithValue(t, "observer callback must be set via Do before registering",
		func() {
			NewObserver(OnCreateEntity).Register(&w)
		})
	expectPanicsWithValue(t, "observer already has a callback",
		func() {
			NewObserver(OnCreateEntity).Do(func(e Entity) {}).Do(func(e Entity) {})
		})

	obs = NewObserver(OnAddComponents).Do(func(e Entity) {})

	obs = obs.For(C[Position]())
	expectEqual(t, 1, len(obs.comps))

	obs = obs.With(C[Position]())
	expectEqual(t, 1, len(obs.with))

	obs = obs.Without()
	expectEqual(t, 0, len(obs.without))
	expectFalse(t, obs.hasWithout)

	obs = obs.Without(C[Position]())
	expectEqual(t, 1, len(obs.without))
	expectTrue(t, obs.hasWithout)
}

func TestObserverRegister(t *testing.T) {
	w := NewWorld()

	obs1 := NewObserver(OnCreateEntity).
		With(C[Position]()).
		With(C[Velocity]()).
		Without(C[Heading]()).
		Do(func(e Entity) {}).
		Register(&w)
	expectTrue(t, w.storage.observers.HasObservers(OnCreateEntity))

	obs2 := NewObserver(OnCreateEntity).
		With(C[Position]()).
		Do(func(e Entity) {}).
		Register(&w)

	expectPanicsWithValue(t, "observer is already registered",
		func() {
			obs1.Register(&w)
		})

	obs1.Unregister(&w)
	expectTrue(t, w.storage.observers.HasObservers(OnCreateEntity))
	obs2.Unregister(&w)
	expectFalse(t, w.storage.observers.HasObservers(OnCreateEntity))

	expectPanicsWithValue(t, "observer is not registered",
		func() {
			obs1.Unregister(&w)
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

	NewObserver(OnCreateEntity).
		With(C[Position]()).
		Do(func(e Entity) {
			callCount++
		}).
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

func TestObserverOnCreateEntities(t *testing.T) {
	w := NewWorld()

	callCount := 0

	NewObserver(OnCreateEntity).
		Do(func(e Entity) {
			callCount++
		}).
		Register(&w)

	w.NewEntities(10, nil)
	expectEqual(t, 10, callCount)
}

func TestObserverOnRemoveEntity(t *testing.T) {
	w := NewWorld()

	builder1 := NewMap1[Position](&w)
	builder2 := NewMap1[Velocity](&w)
	filter1 := NewFilter1[Position](&w)
	filter2 := NewFilter1[Velocity](&w)

	callCount := 0

	NewObserver(OnRemoveEntity).
		With(C[Position]()).
		Do(func(e Entity) {
			callCount++
		}).
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

func TestObserverOnAddRemove(t *testing.T) {
	w := NewWorld()

	builder1 := NewMap1[Position](&w)
	builder2 := NewMap1[Velocity](&w)
	filter1 := NewFilter1[Position](&w)
	filter2 := NewFilter1[Velocity](&w)

	callAdd := 0
	callRemove := 0

	NewObserver(OnAddComponents).
		For(C[Position]()).
		Do(func(e Entity) {
			callAdd++
		}).
		Register(&w)

	NewObserver(OnRemoveComponents).
		For(C[Position]()).
		Do(func(e Entity) {
			callRemove++
		}).
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

func TestObserverOnSet(t *testing.T) {
	w := NewWorld()

	builder1 := NewMap1[Position](&w)
	builder2 := NewMap1[Velocity](&w)

	callCount := 0

	NewObserver(OnSetComponents).
		With(C[Position]()).
		Do(func(e Entity) {
			callCount++
		}).
		Register(&w)

	e := builder1.NewEntity(&Position{})
	builder2.Add(e, &Velocity{})
	expectEqual(t, 0, callCount)

	builder1.Set(e, &Position{})
	expectEqual(t, 1, callCount)
	builder2.Set(e, &Velocity{})
	expectEqual(t, 1, callCount)
}

func TestObserverWith(t *testing.T) {
	w := NewWorld()

	builder1 := NewMap1[Position](&w)
	builder2 := NewMap1[Velocity](&w)

	callAdd := 0
	callRemove := 0
	callSet := 0

	NewObserver(OnAddComponents).
		For(C[Position]()).
		With(C[Velocity]()).
		Do(func(e Entity) {
			callAdd++
		}).
		Register(&w)

	NewObserver(OnRemoveComponents).
		For(C[Position]()).
		With(C[Velocity]()).
		Do(func(e Entity) {
			callRemove++
		}).
		Register(&w)

	NewObserver(OnSetComponents).
		For(C[Position]()).
		With(C[Velocity]()).
		Do(func(e Entity) {
			callSet++
		}).
		Register(&w)

	e := builder2.NewEntity(&Velocity{})

	builder1.Add(e, &Position{})
	expectEqual(t, 1, callAdd)

	builder1.Set(e, &Position{})
	expectEqual(t, 1, callSet)

	builder1.Remove(e)
	expectEqual(t, 1, callRemove)

	e = w.NewEntity()

	builder1.Add(e, &Position{})
	expectEqual(t, 1, callAdd)

	builder1.Set(e, &Position{})
	expectEqual(t, 1, callSet)

	builder1.Remove(e)
	expectEqual(t, 1, callRemove)
}

func TestObserverWithout(t *testing.T) {
	w := NewWorld()

	builder1 := NewMap1[Position](&w)
	builder2 := NewMap1[Velocity](&w)

	callAdd := 0
	callRemove := 0
	callSet := 0

	NewObserver(OnAddComponents).
		For(C[Position]()).
		Without(C[Velocity]()).
		Do(func(e Entity) {
			callAdd++
		}).
		Register(&w)

	NewObserver(OnRemoveComponents).
		For(C[Position]()).
		Without(C[Velocity]()).
		Do(func(e Entity) {
			callRemove++
		}).
		Register(&w)

	NewObserver(OnSetComponents).
		For(C[Position]()).
		Without(C[Velocity]()).
		Do(func(e Entity) {
			callSet++
		}).
		Register(&w)

	e := w.NewEntity()

	builder1.Add(e, &Position{})
	expectEqual(t, 1, callAdd)

	builder1.Set(e, &Position{})
	expectEqual(t, 1, callSet)

	builder1.Remove(e)
	expectEqual(t, 1, callRemove)

	e = builder2.NewEntity(&Velocity{})

	builder1.Add(e, &Position{})
	expectEqual(t, 1, callAdd)

	builder1.Set(e, &Position{})
	expectEqual(t, 1, callSet)

	builder1.Remove(e)
	expectEqual(t, 1, callRemove)
}
