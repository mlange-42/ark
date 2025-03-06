package ecs

import (
	"reflect"
	"unsafe"
)

// column storage for components in an archetype.
type column struct {
	data       reflect.Value  // data buffer
	pointer    unsafe.Pointer // pointer to the first element
	isRelation bool           // whether this column is for a relation component
	target     Entity         // target entity if for a relation component
	itemSize   uintptr        // memory size of items
}

// newColumn creates a new column for a given type and capacity.
func newColumn(tp reflect.Type, isRelation bool, target Entity, capacity uint32) column {
	// TODO: should be use a slice instead of an array here?
	data := reflect.New(reflect.ArrayOf(int(capacity), tp)).Elem()
	pointer := data.Addr().UnsafePointer()

	return column{
		data:       data,
		pointer:    pointer,
		itemSize:   sizeOf(tp),
		isRelation: isRelation,
		target:     target,
	}
}

// Get returns a pointer to the component at the given index.
func (c *column) Get(index uintptr) unsafe.Pointer {
	return unsafe.Add(c.pointer, index*c.itemSize)
}

func (c *column) SetLast(other *column, ownLen uint32, count uint32) {
	start := ownLen - count
	src := other.Get(0)
	dst := c.Get(uintptr(start))
	copyPtr(src, dst, c.itemSize*uintptr(count))
}

// Set overwrites the component at the given index.
func (c *column) Set(index uint32, comp unsafe.Pointer) unsafe.Pointer {
	dst := c.Get(uintptr(index))
	if c.itemSize == 0 {
		return dst
	}

	copyPtr(comp, dst, uintptr(c.itemSize))
	return dst
}

// Zero resets the memory at the given index.
func (c *column) Zero(index uintptr, zero unsafe.Pointer) {
	if c.itemSize == 0 {
		return
	}
	dst := unsafe.Add(c.pointer, index*c.itemSize)
	copyPtr(zero, dst, uintptr(c.itemSize))
}

// Zero resets a block of storage in one buffer.
func (c *column) ZeroRange(start, len uint32, zero unsafe.Pointer) {
	size := uint32(c.itemSize)
	if size == 0 {
		return
	}
	var i uint32
	for i = 0; i < len; i++ {
		dst := unsafe.Add(c.pointer, (i+start)*size)
		copyPtr(zero, dst, c.itemSize)
	}
}

func (c *column) Reset(ownLen uint32, zero unsafe.Pointer) {
	if ownLen == 0 {
		return
	}
	if zero == nil {
		return
	}
	if ownLen <= 64 { // A coarse estimate where manually zeroing is faster
		c.ZeroRange(0, ownLen, zero)
	} else {
		c.data.SetZero()
	}
}
