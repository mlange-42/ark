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
