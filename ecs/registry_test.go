package ecs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComponentRegistry(t *testing.T) {
	reg := newComponentRegistry()

	posType := reflect.TypeOf((*Position)(nil)).Elem()
	rotType := reflect.TypeOf((*Velocity)(nil)).Elem()

	reg.registerComponent(posType, maskTotalBits)
	assert.Equal(t, []uint8{uint8(0)}, reg.IDs)

	reg.registerComponent(rotType, maskTotalBits)
	reg.unregisterLastComponent()
	assert.Equal(t, []uint8{uint8(0)}, reg.IDs)

	id0, _ := reg.ComponentID(posType)
	id1, _ := reg.ComponentID(rotType)
	assert.Equal(t, uint8(0), id0)
	assert.Equal(t, uint8(1), id1)

	assert.Equal(t, []uint8{uint8(0), uint8(1)}, reg.IDs)

	t1, _ := reg.ComponentType(uint8(0))
	t2, _ := reg.ComponentType(uint8(1))

	assert.Equal(t, posType, t1)
	assert.Equal(t, rotType, t2)
}

func TestComponentRegistryOverflow(t *testing.T) {
	reg := newComponentRegistry()

	reg.registerComponent(reflect.TypeOf((*Position)(nil)).Elem(), 1)

	assert.PanicsWithValue(t, "exceeded the maximum of 1 component types or resource types", func() {
		reg.registerComponent(reflect.TypeOf((*Velocity)(nil)).Elem(), 1)
	})
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
	tp1 := typeOf[Position]()
	tp2 := typeOf[Velocity]()

	for b.Loop() {
		_ = tp1 == tp2
	}
}
