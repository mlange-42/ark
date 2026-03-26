package ecs

import (
	"reflect"
	"unsafe"
)

// column storage for components in an table.
type column struct {
	data       reflect.Value  // data buffer
	pointer    unsafe.Pointer // pointer to the first element
	itemSize   uintptr        // memory size of items
	elemType   reflect.Type   // element type of the column
	typePtr    unsafe.Pointer // pointer to the element type's rtype
	target     Entity         // target entity if for a relation component
	index      uint32         // index of the column in the containing table
	isRelation bool           // whether this column is for a relation component
	isTrivial  bool           // Whether the column's type is trivial , i.e. without pointers.
}

// newColumn creates a new column for a given type and capacity.
func newColumn(index uint32, tp reflect.Type, itemSize uintptr, isRelation bool, isTrivial bool, target Entity, capacity uint32) column {
	// TODO: should we use a slice instead of an array here?
	data := reflect.MakeSlice(reflect.SliceOf(tp), int(capacity), int(capacity))
	pointer := data.Index(0).Addr().UnsafePointer()

	return column{
		pointer:    pointer,
		itemSize:   itemSize,
		target:     target,
		index:      index,
		data:       data,
		isRelation: isRelation,
		elemType:   tp,
		typePtr:    rtypePtr(tp),
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
func (c *column) Set(index uint32, src *column, srcIndex uint32) {
	if c.itemSize == 0 {
		return
	}
	if c.isTrivial {
		comp := src.Get(uintptr(srcIndex))
		dst := c.Get(uintptr(index))
		copyPtr(comp, dst, c.itemSize)
		return
	}
	copyValue(src.data, c.data, int(srcIndex), int(index))
}

// Zero resets the memory at the given index.
func (c *column) Zero(index uintptr) {
	if c.itemSize == 0 {
		return
	}

	dst := unsafe.Add(c.pointer, index*c.itemSize)
	if c.isTrivial {
		memclrNoHeapPointers(dst, c.itemSize)
	} else {
		//zeroValueAt(c.data, int(index))
		// Fast GC-safe zero
		typedmemclr(c.typePtr, dst)
	}
}

// ZeroRange resets a block of storage in one buffer.
func (c *column) ZeroRange(start, length uint32) {
	if length == 0 || c.itemSize == 0 {
		return
	}

	elemSize := c.itemSize
	base := uintptr(start) * elemSize

	if c.isTrivial {
		// Single bulk memclr for trivial (pointer-free) types
		total := uintptr(length) * elemSize
		memclrNoHeapPointers(unsafe.Add(c.pointer, base), total)
		return
	}

	// Non-trivial: per-element GC-safe zeroing
	ptr := unsafe.Add(c.pointer, base)
	for i := uint32(0); i < length; i++ {
		typedmemclr(c.typePtr, ptr)
		ptr = unsafe.Add(ptr, elemSize)
	}
}

// Reset the column. Zeroes the memory.
func (c *column) Reset(ownLen uint32) {
	c.ZeroRange(0, ownLen)
}

// entityColumn storage for entities in an table.
type entityColumn struct {
	data    reflect.Value  // data buffer
	pointer unsafe.Pointer // pointer to the first element
}

// newColumn creates a new column for a given type and capacity.
func newEntityColumn(capacity uint32) entityColumn {
	// TODO: should we use a slice instead of an array here?
	data := reflect.MakeSlice(reflect.SliceOf(entityType), int(capacity), int(capacity))
	pointer := data.Index(0).Addr().UnsafePointer()

	return entityColumn{
		data:    data,
		pointer: pointer,
	}
}

// Get returns a pointer to the entity at the given index.
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
