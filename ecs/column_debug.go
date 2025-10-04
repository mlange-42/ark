//go:build ark_debug

package ecs

import "unsafe"

// Get returns a pointer to the component at the given index.
func (c *column) Get(index uintptr) unsafe.Pointer {
	if c == nil {
		panic("entity does not have component the requested component type")
	}
	return unsafe.Add(c.pointer, index*c.itemSize)
}
