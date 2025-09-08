//go:build !ark_debug

package ecs

import "unsafe"

// Next advances the query's cursor to the next entity.
func (q *UnsafeQuery) Next() bool {
	if int64(q.cursor.index) < q.cursor.maxIndex {
		q.cursor.index++
		return true
	}
	return q.nextTableOrArchetype()
}

// Entity returns the current entity.
func (q *UnsafeQuery) Entity() Entity {
	return q.table.GetEntity(q.cursor.index)
}

// Get returns the queried components of the current entity.
//
// ⚠️ Do not store the obtained pointer outside of the current context (i.e. the query loop)!
func (q *UnsafeQuery) Get(comp ID) unsafe.Pointer {
	return q.table.Get(comp, uintptr(q.cursor.index))
}
