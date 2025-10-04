package ecs

import (
	"fmt"
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

func get[T any](storage *componentStorage, index *entityIndex) *T {
	col := storage.columns[index.table]
	if col == nil {
		return nil
	}
	return (*T)(col.Get(uintptr(index.row)))
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

// appends to a slice, but guaranties to return a new one and not alter the original.
func copyAppend[T any](sl []T, elems ...T) []T {
	sl2 := make([]T, len(sl), len(sl)+len(elems))
	copy(sl2, sl)
	sl2 = append(sl2, elems...)
	return sl2
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

func (s *storage) checkRelationComponent(id ID) {
	if !s.registry.IsRelation[id.id] {
		panic(fmt.Sprintf("component with ID %d is not a relation component", id.id))
	}
}

func (s *storage) checkRelationTarget(target Entity) {
	if !target.IsZero() && !s.entityPool.Alive(target) {
		panic("can't use a dead entity as relation target, except for the zero entity")
	}
}
