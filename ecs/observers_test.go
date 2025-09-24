package ecs

import (
	"fmt"
	"testing"
)

func TestObserverManager(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)
	posVelMask := newMask(posID, velID)

	called := false
	obs := NewObserver(OnCreateEntity,
		func(e Entity) {
			called = true
			fmt.Println(e)
		}).
		With(C[Position]()).
		With(C[Velocity]()).
		Register(&w)

	w.storage.observers.FireCreateEntity(Entity{id: 1}, &posVelMask)
	expectTrue(t, called)

	called = false
	obs.Unregister(&w)
	w.storage.observers.FireCreateEntity(Entity{id: 1}, &posVelMask)
	expectFalse(t, called)
}
