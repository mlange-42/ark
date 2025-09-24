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

	NewObserver(OnCreateEntity,
		func(e Entity) {
			fmt.Println(e)
		}).
		With(C[Position](), C[Heading]()).
		Register(&w)

	w.storage.observers.FireCreateEntity(Entity{id: 1}, &posVelMask)
}
