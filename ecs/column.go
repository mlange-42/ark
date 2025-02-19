package ecs

import (
	"reflect"
	"unsafe"
)

// column storage for components in an archetype.
type column struct {
	data     reflect.Value
	pointer  unsafe.Pointer
	itemSize uintptr
	len      uint32
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
		len:      0,
	}
}

func (c *column) Len() int {
	return int(c.len)
}

func (c *column) Cap() int {
	return c.data.Cap()
}

func (c *column) Get(index uint32) unsafe.Pointer {
	return unsafe.Pointer(uintptr(c.pointer) + uintptr(index)*c.itemSize)
}

func (c *column) Add(comp unsafe.Pointer) unsafe.Pointer {
	c.Extend(1)
	c.len++
	return c.Set(c.len-1, comp)
}

func (c *column) Set(index uint32, comp unsafe.Pointer) unsafe.Pointer {
	dst := c.Get(index)
	if c.itemSize == 0 {
		return dst
	}

	copyPtr(comp, dst, c.itemSize)
	return dst
}

func (c *column) Remove(index uint32) bool {
	lastIndex := uint32(c.len - 1)
	swapped := index != lastIndex

	if swapped && c.itemSize != 0 {
		src := unsafe.Add(c.pointer, lastIndex*uint32(c.itemSize))
		dst := unsafe.Add(c.pointer, index*uint32(c.itemSize))
		copyPtr(src, dst, c.itemSize)
	}
	c.len--
	// TODO: zero the last element?
	return swapped
}

func (c *column) Extend(by int) {
	required := c.Len() + by
	cap := c.Cap()
	if cap >= required {
		return
	}
	for cap < required {
		cap *= 2
	}

	if c.itemSize == 0 {
		return
	}
	old := c.data
	c.data = reflect.New(reflect.ArrayOf(cap, old.Type().Elem())).Elem()
	c.pointer = c.data.Addr().UnsafePointer()
	reflect.Copy(c.data, old)
}
