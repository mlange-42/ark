package ecs

import (
	"reflect"
	"testing"
)

func TestColumnPointer(t *testing.T) {
	posType := reflect.TypeOf(Position{})
	column := newColumn(0, posType, posType.Size(), false, true, Entity{}, 8)

	expectEqual(t, uintptr(column.pointer), uintptr(column.data.Index(0).Addr().UnsafePointer()))
}

func TestColumnZeroRange(t *testing.T) {
	posType := reflect.TypeOf(Position{})
	column := newColumn(0, posType, posType.Size(), false, true, Entity{}, 8)

	data := column.data.Interface().([]Position)
	data[1] = Position{1, 2}

	expectEqual(t, Position{1, 2}, data[1])

	column.ZeroRange(0, 8)
	expectEqual(t, Position{0, 0}, data[1])
}
