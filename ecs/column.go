package ecs

import (
	"reflect"
	"unsafe"
)

type column struct {
	data     reflect.Value
	pointer  unsafe.Pointer
	itemSize uintptr
}

func newColumn(tp reflect.Type, capacity int) column {
	size, align := tp.Size(), uintptr(tp.Align())
	size = (size + (align - 1)) / align * align

	data := reflect.New(reflect.ArrayOf(capacity, tp)).Elem()
	pointer := data.Addr().UnsafePointer()

	return column{
		data:     data,
		pointer:  pointer,
		itemSize: size,
	}
}

func (c *column) Get(index uint32) unsafe.Pointer {
	return unsafe.Pointer(uintptr(c.pointer) + uintptr(index)*c.itemSize)
}

func (c *column) Set(index uint32, comp unsafe.Pointer) unsafe.Pointer {
	dst := c.Get(index)
	if c.itemSize == 0 {
		return dst
	}

	copyPtr(comp, dst, c.itemSize)
	return dst
}

func (c *column) Extend(capacity int) {
	if c.itemSize == 0 {
		return
	}
	old := c.data
	c.data = reflect.New(reflect.ArrayOf(capacity, old.Type().Elem())).Elem()
	c.pointer = c.data.Addr().UnsafePointer()
	reflect.Copy(c.data, old)
}
