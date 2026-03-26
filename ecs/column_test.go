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

func TestColumnZeroRangeNonTrivial(t *testing.T) {
	pos := Position{1, 2}
	ptr := PointerType{&pos}

	ptrType := reflect.TypeOf(PointerComp{})
	column := newColumn(0, ptrType, ptrType.Size(), false, false, Entity{}, 8)

	data := column.data.Interface().([]PointerComp)
	data[1] = PointerComp{Ptr: &ptr, Value: 1}

	expectEqual(t, PointerComp{Ptr: &ptr, Value: 1}, data[1])

	column.ZeroRange(0, 8)
	expectEqual(t, PointerComp{Ptr: nil, Value: 0}, data[1])
}
