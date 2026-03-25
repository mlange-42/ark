package main

import (
	"testing"

	"github.com/mlange-42/ark/ecs"
)

func addRemove10(b *testing.B) {
	addRemove(b, 10)
}

func addRemove1000(b *testing.B) {
	addRemove(b, 1000)
}

func addRemove(b *testing.B, n int) {
	world := ecs.NewWorld()

	posMap := ecs.NewMap1[Position](world)
	velMap := ecs.NewMap1[Velocity](world)

	entities := make([]ecs.Entity, 0, n)

	for range n {
		e := posMap.NewEntityFn(nil)
		entities = append(entities, e)
	}

	for _, e := range entities {
		velMap.AddFn(e, nil)
	}
	for _, e := range entities {
		velMap.Remove(e)
	}

	for b.Loop() {
		for _, e := range entities {
			velMap.AddFn(e, nil)
		}
		for _, e := range entities {
			velMap.Remove(e)
		}
	}
}
