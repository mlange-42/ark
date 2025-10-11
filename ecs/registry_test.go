package ecs

import (
	"reflect"
	"testing"
)

func TestComponentRegistry(t *testing.T) {
	reg := newComponentRegistry()

	posType := reflect.TypeOf((*Position)(nil)).Elem()
	rotType := reflect.TypeOf((*Velocity)(nil)).Elem()

	reg.registerComponent(posType, maskTotalBits)
	expectSlicesEqual(t, []uint8{uint8(0)}, reg.IDs)

	reg.registerComponent(rotType, maskTotalBits)
	reg.unregisterLastComponent()
	expectSlicesEqual(t, []uint8{uint8(0)}, reg.IDs)

	id0, _ := reg.ComponentID(posType)
	id1, _ := reg.ComponentID(rotType)
	expectEqual(t, uint8(0), id0)
	expectEqual(t, uint8(1), id1)

	expectSlicesEqual(t, []uint8{uint8(0), uint8(1)}, reg.IDs)

	t1, _ := reg.ComponentType(uint8(0))
	t2, _ := reg.ComponentType(uint8(1))

	expectEqual(t, posType, t1)
	expectEqual(t, rotType, t2)
}

func TestComponentRegistryOverflow(t *testing.T) {
	reg := newComponentRegistry()

	reg.registerComponent(reflect.TypeOf((*Position)(nil)).Elem(), 1)

	expectPanicsWithValue(t, "exceeded the maximum of 1 component types or resource types", func() {
		reg.registerComponent(reflect.TypeOf((*Velocity)(nil)).Elem(), 1)
	})
}

func TestAddArchetype(t *testing.T) {
	reg := newComponentRegistry()

	expectEqual(t, 1, reg.generation, "initial generation should be 1")

	id0 := reg.registerComponent(reflect.TypeOf((*Position)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*Velocity)(nil)).Elem(), maskTotalBits)
	id2 := reg.registerComponent(reflect.TypeOf((*Heading)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*ChildOf)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*ChildOf2)(nil)).Elem(), maskTotalBits)

	reg.addArchetype(id0)
	expectEqual(t, 2, reg.generation, "generation should be 2")
	reg.addArchetype(id2)
	expectEqual(t, 3, reg.generation, "generation should be 3")

	expectEqual(t, maskTotalBits, len(reg.Archetypes))
	expectSlicesEqual(t, []int{1, 0, 1, 0}, reg.Archetypes[:4])
}

func TestRareComponent(t *testing.T) {
	reg := newComponentRegistry()

	id0 := reg.registerComponent(reflect.TypeOf((*Position)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*Velocity)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*Heading)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*ChildOf)(nil)).Elem(), maskTotalBits)
	reg.registerComponent(reflect.TypeOf((*ChildOf2)(nil)).Elem(), maskTotalBits)

	reg.addArchetype(id0)

	expectEqual(t, ID{1}, reg.rareComponent([]ID{{0}, {1}}))
	expectEqual(t, ID{1}, reg.rareComponent([]ID{{1}, {0}}))
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
