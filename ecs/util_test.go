package ecs

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestCopyPtr(t *testing.T) {
	type itemType uint8

	var item itemType = 3
	typeOfItem := reflect.TypeOf(item)
	itemSize := sizeOf(typeOfItem)
	targetItemIndex := 6
	totalItems := 10
	data := reflect.New(reflect.ArrayOf(totalItems, typeOfItem)).Elem()
	dataPointer := data.Addr().UnsafePointer()

	getItem := func(index int) *itemType {
		return (*itemType)(unsafe.Add(dataPointer, uintptr(index)*itemSize))
	}

	for i := range totalItems {
		expectEqual(t, *getItem(i), itemType(0))
	}

	destination := unsafe.Add(dataPointer, uintptr(targetItemIndex)*itemSize)
	source := unsafe.Pointer(&item)
	copyPtr(source, destination, itemSize)

	for i := range totalItems {
		if i == targetItemIndex {
			expectEqual(t, *getItem(i), item)
		} else {
			expectEqual(t, *getItem(i), itemType(0))
		}
	}
}

func TestPagedSlice(t *testing.T) {
	a := pagedSlice[int32]{}

	for i := int32(0); i < 66; i++ {
		a.Add(i)
		expectEqual(t, *a.Get(i), i)
		expectEqual(t, a.Len(), i+1)
	}

	a.Set(3, 100)
	expectEqual(t, *a.Get(3), int32(100))
}

func TestIsTrivial(t *testing.T) {
	expectTrue(t, isTrivial(reflect.TypeFor[[5]int]()))
	expectTrue(t, isTrivial(reflect.TypeFor[struct{}]()))
	expectTrue(t, isTrivial(reflect.TypeFor[struct{ A int }]()))
	expectTrue(t, isTrivial(reflect.TypeFor[struct{ A struct{ A int } }]()))

	expectFalse(t, isTrivial(nil))
	expectFalse(t, isTrivial(reflect.TypeFor[[]int]()))
	expectFalse(t, isTrivial(reflect.TypeFor[[5]string]()))
	expectFalse(t, isTrivial(reflect.TypeFor[struct{ S []int }]()))
	expectFalse(t, isTrivial(reflect.TypeFor[struct{ S []string }]()))
	expectFalse(t, isTrivial(reflect.TypeFor[struct{ S [5]string }]()))
	expectFalse(t, isTrivial(reflect.TypeFor[struct{ A struct{ S string } }]()))
}

func BenchmarkSizeOf(b *testing.B) {
	tp := reflect.TypeFor[Position]()
	for b.Loop() {
		_ = sizeOf(tp)
	}
}
