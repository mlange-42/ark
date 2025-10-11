package ecs

import (
	"math"
	"reflect"
	"unsafe"
)

func capPow2(required uint32) uint32 {
	if required == 0 {
		return 1
	}
	required--
	required |= required >> 1
	required |= required >> 2
	required |= required >> 4
	required |= required >> 8
	required |= required >> 16
	return required + 1
}

func get[T any](storage *componentStorage, table tableID, row uint32) *T {
	col := storage.columns[table]
	if col == nil {
		return nil
	}
	return (*T)(col.Get(uintptr(row)))
}

// copyPtr copies from one pointer to another.
func copyPtr(src, dst unsafe.Pointer, itemSize uintptr) {
	dstSlice := (*[math.MaxInt32]byte)(dst)[:itemSize:itemSize]
	srcSlice := (*[math.MaxInt32]byte)(src)[:itemSize:itemSize]
	copy(dstSlice, srcSlice)
}

func copyValue(src, dst reflect.Value, from, to int) {
	dst.Index(to).Set(src.Index(from))
}

func copyRange(src, dst reflect.Value, start, count int) {
	srcSlice := src.Slice(0, count)
	dstSlice := dst.Slice(start, start+count)
	reflect.Copy(dstSlice, srcSlice)
}

func zeroValueAt(v reflect.Value, index int) {
	elem := v.Index(index)
	elem.SetZero()
}

// isRelation determines whether a type is a relation component.
func isRelation(tp reflect.Type) bool {
	if tp.Kind() != reflect.Struct || tp.NumField() == 0 {
		return false
	}
	field := tp.Field(0)
	return field.Type == relationTp && field.Name == relationTp.Name()
}

// isTrivial checks if a type is "trivial" (contains no pointers, slices, maps, strings, or channels).
// It also returns false if the type itself is one of these.
func isTrivial(tp reflect.Type) bool {
	// Base case: If the type is invalid, return false
	if tp == nil {
		return false
	}

	// Check if the type itself is a pointer, slice, map, or channel
	switch tp.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Interface, reflect.String:
		return false
	}

	// If it's a struct, check its fields recursively
	if tp.Kind() == reflect.Struct {
		for i := range tp.NumField() {
			field := tp.Field(i).Type
			if !isTrivial(field) {
				return false
			}
		}
	}

	// If it's an array, check its elements recursively
	if tp.Kind() == reflect.Array {
		elemType := tp.Elem()
		if !isTrivial(elemType) {
			return false
		}
	}

	// If none of the above conditions matched, it's trivial
	return true
}

// pagedSlice is a paged collection working with pages of length 32 slices.
// its primary purpose is pointer persistence, which is not given using simple slices.
//
// Implements [archetypes].
type pagedSlice[T any] struct {
	pages    [][]T
	len      uint32
	lenLast  uint32
	pageSize uint32
}

func newPagesSlice[T any](pageSize uint32) pagedSlice[T] {
	return pagedSlice[T]{
		pageSize: pageSize,
	}
}

// Add adds a value to the paged slice.
func (p *pagedSlice[T]) Add(value T) {
	if p.len == 0 || p.lenLast == p.pageSize {
		p.pages = append(p.pages, make([]T, p.pageSize))
		p.lenLast = 0
	}
	p.pages[len(p.pages)-1][p.lenLast] = value
	p.len++
	p.lenLast++
}

// Get returns the value at the given index.
func (p *pagedSlice[T]) Get(index uint32) *T {
	return &p.pages[index/p.pageSize][index%p.pageSize]
}

// Set sets the value at the given index.
func (p *pagedSlice[T]) Set(index uint32, value T) {
	p.pages[index/p.pageSize][index%p.pageSize] = value
}

// Len returns the current number of items in the paged slice.
func (p *pagedSlice[T]) Len() uint32 {
	return p.len
}
