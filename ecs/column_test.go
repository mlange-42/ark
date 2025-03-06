package ecs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColumnPointer(t *testing.T) {
	posType := reflect.TypeOf(Position{})
	column := newColumn(posType, false, Entity{}, 8)

	assert.Equal(t, uintptr(column.pointer), uintptr(column.data.Addr().UnsafePointer()))
}

/*
func TestColumnAddRemove(t *testing.T) {
	posType := reflect.TypeOf(Position{})

	column := newColumn(posType, false, Entity{}, 8)

	assert.EqualValues(t, 8, column.cap)
	assert.EqualValues(t, 0, column.len)

	column.Add(unsafe.Pointer(&Position{1, 2}))
	column.Add(unsafe.Pointer(&Position{3, 4}))
	column.Add(unsafe.Pointer(&Position{0, 0}))
	column.Set(2, unsafe.Pointer(&Position{5, 6}))

	assert.EqualValues(t, 8, column.cap)
	assert.EqualValues(t, 3, column.len)

	pos := (*Position)(column.Get(2))
	assert.Equal(t, Position{5, 6}, *pos)

	pos = (*Position)(column.Get(0))
	assert.Equal(t, Position{1, 2}, *pos)

	for range 8 {
		column.Add(unsafe.Pointer(&Position{1, 2}))
	}
	assert.EqualValues(t, 16, column.cap)
	assert.EqualValues(t, 11, column.len)
}

func TestColumnAddRemoveLabel(t *testing.T) {
	labelType := reflect.TypeOf(Label{})

	column := newColumn(labelType, false, Entity{}, 8)

	assert.EqualValues(t, 8, column.cap)
	assert.EqualValues(t, 0, column.len)

	column.Add(unsafe.Pointer(&Label{}))
	column.Add(unsafe.Pointer(&Label{}))
	column.Add(unsafe.Pointer(&Label{}))
	column.Set(2, unsafe.Pointer(&Label{}))

	assert.EqualValues(t, 8, column.cap)
	assert.EqualValues(t, 3, column.len)

	pos := (*Label)(column.Get(2))
	assert.Equal(t, Label{}, *pos)

	pos = (*Label)(column.Get(0))
	assert.Equal(t, Label{}, *pos)

	for range 8 {
		column.Add(unsafe.Pointer(&Label{}))
	}
	assert.EqualValues(t, 16, column.cap)
	assert.EqualValues(t, 11, column.len)
}

func TestColumnReset(t *testing.T) {
	posType := reflect.TypeOf(Position{})
	posZeroValue := make([]byte, sizeOf(posType))
	posZeroPointer := unsafe.Pointer(&posZeroValue[0])

	labelType := reflect.TypeOf(Label{})
	var labelZeroPointer unsafe.Pointer

	labelColumn := newColumn(labelType, false, Entity{}, 8)
	posColumn := newColumn(posType, false, Entity{}, 8)

	labelColumn.Reset(labelZeroPointer)
	posColumn.Reset(posZeroPointer)

	for range 12 {
		labelColumn.Add(unsafe.Pointer(&Label{}))
		posColumn.Add(unsafe.Pointer(&Position{}))
	}

	labelColumn.Reset(labelZeroPointer)
	posColumn.Reset(posZeroPointer)

	assert.EqualValues(t, 0, labelColumn.len)
	assert.EqualValues(t, 0, posColumn.len)

	for range 123 {
		labelColumn.Add(unsafe.Pointer(&Label{}))
		posColumn.Add(unsafe.Pointer(&Position{}))
	}

	labelColumn.Reset(labelZeroPointer)
	posColumn.Reset(posZeroPointer)

	assert.EqualValues(t, 0, labelColumn.len)
	assert.EqualValues(t, 0, posColumn.len)
}
*/
