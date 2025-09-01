package ecs

import (
	"reflect"
	"testing"
)

func TestColumnPointer(t *testing.T) {
	posType := reflect.TypeOf(Position{})
	column := newColumn(0, posType, sizeOf(posType), false, Entity{}, 8)

	expectEqual(t, uintptr(column.pointer), uintptr(column.data.Addr().UnsafePointer()))
}
