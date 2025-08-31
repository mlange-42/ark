package ecs

import (
	"reflect"
	"testing"
)

func TestComponentRegistry(t *testing.T) {
	reg := newComponentRegistry()
	posType := reflect.TypeOf((*Position)(nil)).Elem()
	velType := reflect.TypeOf((*Velocity)(nil)).Elem()
	reg.registerComponent(posType, maskTotalBits)
	if len(reg.IDs) != 1 || reg.IDs[0] != 0 {
		t.Errorf("expected IDs to be [0], got %v", reg.IDs)
	}
	reg.registerComponent(velType, maskTotalBits)
	reg.unregisterLastComponent()
	if len(reg.IDs) != 1 || reg.IDs[0] != 0 {
		t.Errorf("expected IDs to be [0], got %v", reg.IDs)
	}
	id0, _ := reg.ComponentID(posType)
	id1, _ := reg.ComponentID(velType)
	if id0 != 0 {
		t.Errorf("expected ID of Position to be 0, got %d", id0)
	}
	if id1 != 1 {
		t.Errorf("expected ID of Velocity to be 1, got %d", id1)
	}
	if len(reg.IDs) != 2 || reg.IDs[0] != 0 || reg.IDs[1] != 1 {
		t.Errorf("expected IDs to be [0, 1], got %v", reg.IDs)
	}
	t1, _ := reg.ComponentType(uint8(0))
	t2, _ := reg.ComponentType(uint8(1))
	if t1 != posType {
		t.Errorf("expected type of ID 0 to be Position, got %v", t1)
	}
	if t2 != velType {
		t.Errorf("expected type of ID 1 to be Velocity, got %v", t2)
	}
}

func TestComponentRegistryOverflow(t *testing.T) {
	reg := newComponentRegistry()
	reg.registerComponent(reflect.TypeOf((*Position)(nil)).Elem(), 1)
	expectPanicWithValue(t, "exceeded the maximum of 1 component types or resource types", func() {
		reg.registerComponent(reflect.TypeOf((*Velocity)(nil)).Elem(), 1)
	})
}

func TestAddArchetype(t *testing.T) {
	reg := newComponentRegistry()
	id0 := reg.registerComponent(reflect.TypeOf((*Position)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*Velocity)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*Heading)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*ChildOf)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*ChildOf2)(nil)).Elem(), maskTotalBits)
	reg.addArchetype(id0)
	if len(reg.Archetypes) != maskTotalBits {
		t.Errorf("expected length of Archetypes to be %d, got %d", maskTotalBits, len(reg.Archetypes))
	}
	if reg.Archetypes[0] != 1 || reg.Archetypes[1] != 0 {
		t.Errorf("expected first two elements of Archetypes to be [1, 0], got [%d, %d]", reg.Archetypes[0], reg.Archetypes[1])
	}
}

func TestRareComponent(t *testing.T) {
	reg := newComponentRegistry()
	id0 := reg.registerComponent(reflect.TypeOf((*Position)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*Velocity)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*Heading)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*ChildOf)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*ChildOf2)(nil)).Elem(), maskTotalBits)
	reg.addArchetype(id0)
	if reg.rareComponent([]ID{{0}, {1}}) != (ID{1}) {
		t.Errorf("expected rareComponent to return ID{1}, got %v", reg.rareComponent([]ID{{0}, {1}}))
	}
	if reg.rareComponent([]ID{{1}, {0}}) != (ID{1}) {
		t.Errorf("expected rareComponent to return ID{1}, got %v", reg.rareComponent([]ID{{1}, {0}}))
	}
}

func BenchmarkRegistryGet(b *testing.B) {
	w := NewWorld(1024)
	_ = ComponentID[Position](&w)
	_ = ComponentID[Velocity](&w)
	for b.Loop() {
		_ = ComponentID[Velocity](&w)
	}
}

func BenchmarkTypeEquality(b *testing.B) {
	tp1 := reflect.TypeFor[Position]()
	tp2 := reflect.TypeFor[Velocity]()
	for b.Loop() {
		_ = tp1 == tp2
	}
}

func BenchmarkRareComponent2(b *testing.B) {
	reg := newComponentRegistry()
	id0 := reg.registerComponent(reflect.TypeOf((*Position)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*Velocity)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*Heading)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*ChildOf)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*ChildOf2)(nil)).Elem(), maskTotalBits)
	reg.addArchetype(id0)
	ids := []ID{{0}, {1}}
	for b.Loop() {
		reg.rareComponent(ids)
	}
}

func BenchmarkRareComponent5(b *testing.B) {
	reg := newComponentRegistry()
	id0 := reg.registerComponent(reflect.TypeOf((*Position)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*Velocity)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*Heading)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*ChildOf)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*ChildOf2)(nil)).Elem(), maskTotalBits)
	reg.addArchetype(id0)
	ids := []ID{{0}, {1}, {2}, {3}, {4}}
	for b.Loop() {
		reg.rareComponent(ids)
	}
}
