package ecs

import (
	"math"
	"reflect"
	"unsafe"
)

// Page size of pagedSlice type
const pageSize = 32

// copyPtr copies from one pointer to another.
func copyPtr(src, dst unsafe.Pointer, itemSize uintptr) {
	dstSlice := (*[math.MaxInt32]byte)(dst)[:itemSize:itemSize]
	srcSlice := (*[math.MaxInt32]byte)(src)[:itemSize:itemSize]
	copy(dstSlice, srcSlice)
}

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func sizeOf(tp reflect.Type) uintptr {
	size, align := tp.Size(), uintptr(tp.Align())
	return (size + (align - 1)) / align * align
}

// appends to a slice, but guaranties to return a new one and not alter the original.
func appendNew[T any](sl []T, elems ...T) []T {
	sl2 := make([]T, len(sl), len(sl)+len(elems))
	copy(sl2, sl)
	sl2 = append(sl2, elems...)
	return sl2
}

// pagedSlice is a paged collection working with pages of length 32 slices.
// its primary purpose is pointer persistence, which is not given using simple slices.
//
// Implements [archetypes].
type pagedSlice[T any] struct {
	pages   [][]T
	len     int32
	lenLast int32
}

// Add adds a value to the paged slice.
func (p *pagedSlice[T]) Add(value T) {
	if p.len == 0 || p.lenLast == pageSize {
		p.pages = append(p.pages, make([]T, pageSize))
		p.lenLast = 0
	}
	p.pages[len(p.pages)-1][p.lenLast] = value
	p.len++
	p.lenLast++
}

// Get returns the value at the given index.
func (p *pagedSlice[T]) Get(index int32) *T {
	return &p.pages[index/pageSize][index%pageSize]
}

// Set sets the value at the given index.
func (p *pagedSlice[T]) Set(index int32, value T) {
	p.pages[index/pageSize][index%pageSize] = value
}

// Len returns the current number of items in the paged slice.
func (p *pagedSlice[T]) Len() int32 {
	return p.len
}
