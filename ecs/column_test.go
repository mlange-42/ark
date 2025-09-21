package ecs

import (
	"reflect"
	"testing"
)

func TestColumnPointer(t *testing.T) {
	posType := reflect.TypeOf(Position{})
	column := newColumn(0, posType, posType.Size(), false, true, Entity{}, 8)

	expectEqual(t, uintptr(column.pointer), uintptr(column.data.Addr().UnsafePointer()))
}
