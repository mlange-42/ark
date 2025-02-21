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
