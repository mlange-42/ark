package ecs

import (
	"reflect"
	"sort"
)

// ID is the component identifier.
type ID struct {
	id uint32
}

func id(id uint32) ID {
	return ID{id}
}

type componentInfo struct {
	id  ID
	typ reflect.Type
}

type ids []ID

func newIDs(id ...ID) ids {
	return append([]ID(nil), id...)
}

func newSortedIDs(id ...ID) ids {
	ids := ids(append([]ID(nil), id...))
	sort.Sort(ids)
	return ids
}
func (ids ids) Len() int           { return len(ids) }
func (ids ids) Less(i, j int) bool { return ids[i].id < ids[j].id }
func (ids ids) Swap(i, j int)      { ids[i], ids[j] = ids[j], ids[i] }

func (ids ids) Search(id ID) (int, bool) {
	// Define f(-1) == false and f(n) == true.
	// Invariant: f(i-1) == false, f(j) == true.
	n := len(ids)
	i, j := 0, n
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if ids[h].id < id.id {
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}
	// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
	return i, i < n && ids[i].id == id.id
}
