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

// newColumn creates a new column for a given type and capacity.
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

// Len returns the number of components in the column.
func (c *column) Len() int {
	return int(c.len)
}

// Cap returns the current capacity of the column.
func (c *column) Cap() int {
	return c.data.Cap()
}

// Get returns a pointer to the component at the given index.
func (c *column) Get(index uint32) unsafe.Pointer {
	return unsafe.Pointer(uintptr(c.pointer) + uintptr(index)*c.itemSize)
}

// Add adds a component to the column.
func (c *column) Add(comp unsafe.Pointer) unsafe.Pointer {
	c.Extend(1)
	c.len++
	return c.Set(c.len-1, comp)
}

// Set overwrites the component at the given index.
func (c *column) Set(index uint32, comp unsafe.Pointer) unsafe.Pointer {
	dst := c.Get(index)
	if c.itemSize == 0 {
		return dst
	}

	copyPtr(comp, dst, c.itemSize)
	return dst
}

// Remove swap-removes the component at the given index.
// Returns whether a swap was necessary.
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

// Extend the column to be able to store the given number of additional components.
// Has no effect of the column's capacity is already sufficient.
// If the capacity needs to be increased, it will be doubled until it is sufficient.
func (c *column) Extend(by int) {
	required := c.Len() + by
	cap := c.Cap()
	if cap >= required {
		return
	}
	for cap < required {
		cap *= 2
	}
	old := c.data
	c.data = reflect.New(reflect.ArrayOf(cap, old.Type().Elem())).Elem()
	c.pointer = c.data.Addr().UnsafePointer()
	reflect.Copy(c.data, old)
}
