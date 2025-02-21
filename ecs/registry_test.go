package ecs

import "testing"

func BenchmarkRegistryGet(b *testing.B) {
	w := NewWorld(1024)

	_ = ComponentID[Position](&w)
	_ = ComponentID[Velocity](&w)

	for b.Loop() {
		_ = ComponentID[Velocity](&w)
	}
}

func BenchmarkTypeEquality(b *testing.B) {
	tp1 := typeOf[Position]()
	tp2 := typeOf[Velocity]()

	for b.Loop() {
		_ = tp1 == tp2
	}
}
