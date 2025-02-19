package ecs

import (
	"math/rand/v2"
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type columnInterface interface {
	Add(any)
	Get(uint32) unsafe.Pointer
}

type sliceColumn[T any] struct {
	data     []T
	pointer  unsafe.Pointer
	itemSize uintptr
}

func newSliceColumn[T any]() *sliceColumn[T] {
	return &sliceColumn[T]{
		itemSize: sizeOf(typeOf[T]()),
	}
}

func (c *sliceColumn[T]) Get(index uint32) unsafe.Pointer {
	return unsafe.Add(c.pointer, int(index)*int(c.itemSize))
}

func (c *sliceColumn[T]) Add(value any) {
	c.data = append(c.data, value.(T))
	c.pointer = unsafe.Pointer(&c.data[0])
}

func BenchmarkColumnGet(b *testing.B) {
	n := 1000
	posType := reflect.TypeOf(Position{})
	column := newColumn(posType, 1000)

	indices := make([]uint32, n)
	for i := 0; i < n; i++ {
		column.Add(unsafe.Pointer(&Position{1, 2}))
		indices[i] = uint32(i)
	}
	rand.Shuffle(n, func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})

	var ptr unsafe.Pointer
	for b.Loop() {
		ptr = column.Get(500)
	}

	assert.NotNil(b, ptr)
}

func BenchmarkSliceColumnGet(b *testing.B) {
	n := 1000
	var column columnInterface = &sliceColumn[Position]{}

	indices := make([]uint32, n)
	for i := 0; i < n; i++ {
		column.Add(Position{1, 2})
		indices[i] = uint32(i)
	}
	rand.Shuffle(n, func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})

	var ptr unsafe.Pointer
	for b.Loop() {
		ptr = column.Get(500)
	}

	assert.NotNil(b, ptr)
}
