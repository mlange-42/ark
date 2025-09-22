package ecs

import (
	"reflect"
	"unsafe"
)

// column storage for components in an table.
type column struct {
	pointer    unsafe.Pointer // pointer to the first element
	data       reflect.Value  // data buffer
	itemSize   uintptr        // memory size of items
	target     Entity         // target entity if for a relation component
	index      uint32         // index of the column in the containing table
	isRelation bool           // whether this column is for a relation component
	elemType   reflect.Type   // element type of the column
	isTrivial  bool           // Whether the column's type is trivial , i.e. without pointers.
}

// newColumn creates a new column for a given type and capacity.
func newColumn(index uint32, tp reflect.Type, itemSize uintptr, isRelation bool, isTrivial bool, target Entity, capacity uint32) column {
	// TODO: should we use a slice instead of an array here?
	data := reflect.New(reflect.ArrayOf(int(capacity), tp)).Elem()
	pointer := data.Addr().UnsafePointer()

	return column{
		index:      index,
		data:       data,
		pointer:    pointer,
		itemSize:   itemSize,
		isRelation: isRelation,
		target:     target,
		elemType:   tp,
		isTrivial:  isTrivial,
	}
}

// Get returns a pointer to the component at the given index.
func (c *column) Get(index uintptr) unsafe.Pointer {
	return unsafe.Add(c.pointer, index*c.itemSize)
}

// CopyToEnd copies from the given column to the end of this column.
// Column length must be increased before.
func (c *column) CopyToEnd(from *column, ownLen uint32, count uint32) {
	start := ownLen - count
	if c.isTrivial {
		src := from.Get(0)
		dst := c.Get(uintptr(start))
		copyPtr(src, dst, c.itemSize*uintptr(count))
		return
	}
	copyRange(from.data, c.data, int(start), int(count))
}

// Set overwrites the component at the given index.
func (c *column) Set(index uint32, src *column, srcIndex int) {
	if c.itemSize == 0 {
		return
	}
	if c.isTrivial {
		comp := src.Get(uintptr(srcIndex))
		dst := c.Get(uintptr(index))
		copyPtr(comp, dst, c.itemSize)
		return
	}
	copyValue(src.data, c.data, srcIndex, int(index))
}

// Zero resets the memory at the given index.
func (c *column) Zero(index uintptr, zero unsafe.Pointer) {
	if c.itemSize == 0 {
		return
	}
	if c.isTrivial {
		dst := unsafe.Add(c.pointer, index*c.itemSize)
		copyPtr(zero, dst, uintptr(c.itemSize))
	} else {
		// TODO: Do we really need this?
		// Tests indicate stuff get GC'd also with copyPtr.
		zeroValueAt(c.data, int(index))
	}
}

// Zero resets a block of storage in one buffer.
func (c *column) ZeroRange(start, len uint32, zero unsafe.Pointer) {
	size := uint32(c.itemSize)
	if size == 0 {
		return
	}
	var i uint32
	for i = range len {
		dst := unsafe.Add(c.pointer, (i+start)*size)
		copyPtr(zero, dst, c.itemSize)
	}
}

func (c *column) Reset(ownLen uint32, zero unsafe.Pointer) {
	if ownLen == 0 {
		return
	}
	if ownLen <= 64 && c.isTrivial { // A coarse estimate where manually zeroing is faster
		c.ZeroRange(0, ownLen, zero)
	} else {
		c.data.SetZero()
	}
}

// entityColumn storage for entities in an table.
type entityColumn struct {
	pointer unsafe.Pointer // pointer to the first element
	data    reflect.Value  // data buffer
}

// newColumn creates a new column for a given type and capacity.
func newEntityColumn(capacity uint32) entityColumn {
	// TODO: should we use a slice instead of an array here?
	data := reflect.New(reflect.ArrayOf(int(capacity), entityType)).Elem()
	pointer := data.Addr().UnsafePointer()

	return entityColumn{
		data:    data,
		pointer: pointer,
	}
}

// Get returns a pointer to the component at the given index.
func (c *entityColumn) Get(index uintptr) unsafe.Pointer {
	return unsafe.Add(c.pointer, index*entitySize)
}

// CopyToEnd copies from the given column to the end of this column.
// Column length must be increased before.
func (c *entityColumn) CopyToEnd(from *entityColumn, ownLen uint32, count uint32) {
	start := ownLen - count
	src := from.Get(0)
	dst := c.Get(uintptr(start))
	copyPtr(src, dst, entitySize*uintptr(count))
}

// Set overwrites the component at the given index.
func (c *entityColumn) Set(index uint32, comp unsafe.Pointer) unsafe.Pointer {
	dst := c.Get(uintptr(index))

	copyPtr(comp, dst, entitySize)
	return dst
}
