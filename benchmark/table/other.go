package main

import (
	"testing"

	"github.com/mlange-42/ark/benchmark"
	"github.com/mlange-42/ark/ecs"
)

func benchesOther() []benchmark.Benchmark {
	return []benchmark.Benchmark{
		{Name: "ecs.NewWorld", Desc: "", F: newWorld, N: 1, Factor: 0.001, Units: "μs"},
		{Name: "World.Reset", Desc: "empty world", F: resetWorld, N: 1},
		{Name: "ecs.ComponentID", Desc: "component already registered", F: componentID, N: 1},
	}
}

func newWorld(b *testing.B) {
	var w ecs.World
	for b.Loop() {
		w = ecs.NewWorld()
	}

	if w.IsLocked() {
		b.Errorf("World must be locked")
	}
}

func resetWorld(b *testing.B) {
	w := ecs.NewWorld()
	for b.Loop() {
		w.Reset()
	}

	if w.IsLocked() {
		b.Errorf("World must be locked")
	}
}

func componentID(b *testing.B) {
	w := ecs.NewWorld()
	origID := ecs.ComponentID[comp1](&w)

	var id ecs.ID

	for b.Loop() {
		id = ecs.ComponentID[comp1](&w)
	}

	if id != origID {
		b.Errorf("Expected %v, got %v", origID, id)
	}
}
