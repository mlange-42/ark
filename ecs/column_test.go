package ecs

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestColumnPointer(t *testing.T) {
	posType := reflect.TypeOf(Position{})
	column := newColumn(posType, false, Entity{}, 8)

	assert.Equal(t, uintptr(column.pointer), uintptr(column.data.Addr().UnsafePointer()))
}

func TestColumnAddRemove(t *testing.T) {
	posType := reflect.TypeOf(Position{})
	zeroValue := make([]byte, sizeOf(posType))
	zeroPointer := unsafe.Pointer(&zeroValue[0])

	column := newColumn(posType, false, Entity{}, 8)

	assert.Equal(t, 8, column.Cap())
	assert.Equal(t, 0, column.Len())

	column.Add(unsafe.Pointer(&Position{1, 2}))
	column.Add(unsafe.Pointer(&Position{3, 4}))
	column.Add(unsafe.Pointer(&Position{0, 0}))
	column.Set(2, unsafe.Pointer(&Position{5, 6}))

	assert.Equal(t, 8, column.Cap())
	assert.Equal(t, 3, column.Len())

	pos := (*Position)(column.Get(2))
	assert.Equal(t, Position{5, 6}, *pos)

	swapped := column.Remove(0, zeroPointer)
	assert.True(t, swapped)
	assert.Equal(t, 2, column.Len())

	pos = (*Position)(column.Get(0))
	assert.Equal(t, Position{5, 6}, *pos)

	swapped = column.Remove(1, zeroPointer)
	assert.False(t, swapped)
	assert.Equal(t, 1, column.Len())

	for range 8 {
		column.Add(unsafe.Pointer(&Position{1, 2}))
	}
	assert.Equal(t, 16, column.Cap())
	assert.Equal(t, 9, column.Len())
}

func TestColumnAddRemoveLabel(t *testing.T) {
	labelType := reflect.TypeOf(Label{})
	var zeroPointer unsafe.Pointer

	column := newColumn(labelType, false, Entity{}, 8)

	assert.Equal(t, 8, column.Cap())
	assert.Equal(t, 0, column.Len())

	column.Add(unsafe.Pointer(&Label{}))
	column.Add(unsafe.Pointer(&Label{}))
	column.Add(unsafe.Pointer(&Label{}))
	column.Set(2, unsafe.Pointer(&Label{}))

	assert.Equal(t, 8, column.Cap())
	assert.Equal(t, 3, column.Len())

	pos := (*Label)(column.Get(2))
	assert.Equal(t, Label{}, *pos)

	swapped := column.Remove(0, zeroPointer)
	assert.True(t, swapped)
	assert.Equal(t, 2, column.Len())

	pos = (*Label)(column.Get(0))
	assert.Equal(t, Label{}, *pos)

	swapped = column.Remove(1, zeroPointer)
	assert.False(t, swapped)
	assert.Equal(t, 1, column.Len())

	for range 8 {
		column.Add(unsafe.Pointer(&Label{}))
	}
	assert.Equal(t, 16, column.Cap())
	assert.Equal(t, 9, column.Len())
}
