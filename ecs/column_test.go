package ecs

import (
	"reflect"
	"testing"
)

func TestColumnPointer(t *testing.T) {
	posType := reflect.TypeOf(Position{})
	column := newColumn(0, posType, sizeOf(posType), false, Entity{}, 8)

	if uintptr(column.pointer) != uintptr(column.data.Addr().UnsafePointer()) {
		t.Errorf("expected column.pointer to match column.data.Addr().UnsafePointer(), got %v and %v",
			column.pointer, column.data.Addr().UnsafePointer())
	}
}
