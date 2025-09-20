package ecs

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestCapPow2(t *testing.T) {
	expectEqual(t, 1, capPow2(0))
	expectEqual(t, 64, capPow2(64))
	expectEqual(t, 128, capPow2(65))
	expectEqual(t, 1024, capPow2(1000))
	expectEqual(t, 1048576, capPow2(1_000_000))
}

func TestCopyPtr(t *testing.T) {
	type itemType uint8 // can be any type, result stays the same

	// setup
	var item itemType = 3
	typeOfItem := reflect.TypeOf(item)
	itemSize := typeOfItem.Size()
	targetItemIndex := 6
	totalItems := 10
	data := reflect.New(reflect.ArrayOf(int(totalItems), typeOfItem)).Elem()
	dataPointer := data.Addr().UnsafePointer() // points to the start of data

	getItem := func(index int) *itemType {
		return (*itemType)(unsafe.Add(dataPointer, uintptr(index)*itemSize))
	}

	// check that the expected item is not there yet
	for i := range totalItems {
		expectEqual(t, itemType(0), *getItem(i))
	}

	// copy the item to the right place
	destination := unsafe.Add(
		dataPointer,
		uintptr(targetItemIndex)*itemSize,
	)

	source := unsafe.Pointer(&item)
	copyPtr(source, destination, itemSize)

	// check that only the expected item is now set
	for i := range totalItems {
		if i == targetItemIndex {
			expectEqual(t, item, *getItem(i))
		} else {
			expectEqual(t, itemType(0), *getItem(i))
		}
	}
}

func TestPagedSlice(t *testing.T) {
	a := pagedSlice[int32]{}

	var i int32
	for i = range 66 {
		a.Add(i)
		expectEqual(t, i, *a.Get(i))
		expectEqual(t, i+1, a.Len())
	}

	a.Set(3, 100)
	expectEqual(t, int32(100), *a.Get(3))
}

func TestIsTrivial(t *testing.T) {
	expectTrue(t, isTrivial(reflect.TypeFor[[5]int]()))
	expectTrue(t, isTrivial(reflect.TypeFor[struct{}]()))
	expectTrue(t, isTrivial(reflect.TypeFor[struct {
		A int
	}]()))
	expectTrue(t, isTrivial(reflect.TypeFor[struct {
		A struct{ A int }
	}]()))

	expectFalse(t, isTrivial(nil))
	expectFalse(t, isTrivial(reflect.TypeFor[[]int]()))
	expectFalse(t, isTrivial(reflect.TypeFor[[5]string]()))

	expectFalse(t, isTrivial(reflect.TypeFor[struct {
		S []int
	}]()))
	expectFalse(t, isTrivial(reflect.TypeFor[struct {
		S []string
	}]()))
	expectFalse(t, isTrivial(reflect.TypeFor[struct {
		S [5]string
	}]()))
	expectFalse(t, isTrivial(reflect.TypeFor[struct {
		A struct{ S string }
	}]()))
}

func BenchmarkSizeOf(b *testing.B) {
	tp := reflect.TypeFor[Position]()
	for b.Loop() {
		_ = tp.Size()
	}
}

func BenchmarkCapPow2(b *testing.B) {
	for b.Loop() {
		_ = capPow2(513)
	}
}
