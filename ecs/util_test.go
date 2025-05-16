package ecs

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestCopyPtr(t *testing.T) {
	assert := assert.New(t)

	type itemType uint8 // can be any type, result stays the same

	// setup
	var item itemType = 3
	typeOfItem := reflect.TypeOf(item)
	itemSize := sizeOf(typeOfItem)
	targetItemIndex := 6
	totalItems := 10
	data := reflect.New(reflect.ArrayOf(int(totalItems), typeOfItem)).Elem()
	dataPointer := data.Addr().UnsafePointer() // points to the start of data

	getItem := func(index int) *itemType {
		return (*itemType)(unsafe.Add(dataPointer, uintptr(index)*itemSize))
	}

	// check that the expected item is not there yet
	for i := range totalItems {
		assert.Equal(itemType(0), *getItem(i))
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
			assert.Equal(item, *getItem(i))
		} else {
			assert.Equal(itemType(0), *getItem(i))
		}
	}
}

func TestPagedSlice(t *testing.T) {
	a := pagedSlice[int32]{}

	var i int32
	for i = range 66 {
		a.Add(i)
		assert.Equal(t, i, *a.Get(i))
		assert.Equal(t, i+1, a.Len())
	}

	a.Set(3, 100)
	assert.Equal(t, int32(100), *a.Get(3))
}

func TestIsTrivial(t *testing.T) {
	assert.True(t, isTrivial(reflect.TypeFor[[5]int]()))
	assert.True(t, isTrivial(reflect.TypeFor[struct{}]()))
	assert.True(t, isTrivial(reflect.TypeFor[struct {
		A int
	}]()))
	assert.True(t, isTrivial(reflect.TypeFor[struct {
		A struct{ A int }
	}]()))

	assert.False(t, isTrivial(nil))
	assert.False(t, isTrivial(reflect.TypeFor[[]int]()))
	assert.False(t, isTrivial(reflect.TypeFor[[5]string]()))

	assert.False(t, isTrivial(reflect.TypeFor[struct {
		S []int
	}]()))
	assert.False(t, isTrivial(reflect.TypeFor[struct {
		S []string
	}]()))
	assert.False(t, isTrivial(reflect.TypeFor[struct {
		S [5]string
	}]()))
	assert.False(t, isTrivial(reflect.TypeFor[struct {
		A struct{ S string }
	}]()))
}

func BenchmarkSizeOf(b *testing.B) {
	tp := reflect.TypeFor[Position]()
	for b.Loop() {
		_ = sizeOf(tp)
	}
}
