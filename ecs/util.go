package ecs

import (
	"math"
	"reflect"
	"unsafe"
)

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
