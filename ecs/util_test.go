package ecs

import (
	"reflect"
	"runtime"
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
	data := reflect.New(reflect.ArrayOf(totalItems, typeOfItem)).Elem()
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

func BenchmarkReflectArrayToSlice(b *testing.B) {
	capacity := 100
	typ := reflect.TypeFor[Position]()
	data := reflect.New(reflect.ArrayOf(capacity, typ)).Elem()

	for i := range data.Len() {
		data.Index(i).Set(reflect.ValueOf(Position{X: 1, Y: 1}))
	}

	var slice []Position
	loop := func() {
		slice = data.Slice(0, data.Len()).Interface().([]Position)
	}

	for b.Loop() {
		loop()
	}

	runtime.KeepAlive(slice)
}

func BenchmarkReflectSliceToSlice(b *testing.B) {
	capacity := 100
	typ := reflect.TypeFor[Position]()
	data := reflect.MakeSlice(reflect.SliceOf(typ), capacity, capacity)

	for i := range data.Len() {
		data.Index(i).Set(reflect.ValueOf(Position{X: 1, Y: 1}))
	}

	var slice []Position
	loop := func() {
		slice = data.Interface().([]Position)
	}

	for b.Loop() {
		loop()
	}

	runtime.KeepAlive(slice)
}

func BenchmarkReflectAnyToSlice(b *testing.B) {
	capacity := 100
	typ := reflect.TypeFor[Position]()
	data := reflect.MakeSlice(reflect.SliceOf(typ), capacity, capacity)

	for i := range data.Len() {
		data.Index(i).Set(reflect.ValueOf(Position{X: 1, Y: 1}))
	}

	var anySlice any = data.Interface().([]Position)

	var slice []Position
	loop := func() {
		slice = anySlice.([]Position)
	}

	for b.Loop() {
		loop()
	}

	runtime.KeepAlive(slice)
}
