package ecs

import (
	"sort"
)

// ID is the component identifier.
type ID struct {
	id uint8
}

func id(id int) ID {
	return ID{uint8(id)}
}

// ids is a sorted list of component [ID]s.
type ids []ID

// newIDs creates a new list of component [ID]s.
//
// Safety: the caller must ensure that the IDs are sorted.
func newIDs(id ...ID) ids {
	return append([]ID(nil), id...)
}

// newIDsSorted creates a new list of component [ID]s and sorts them.
func newIDsSorted(id ...ID) ids {
	ids := ids(append([]ID(nil), id...))
	sort.Sort(ids)
	return ids
}

func (ids ids) Len() int           { return len(ids) }
func (ids ids) Less(i, j int) bool { return ids[i].id < ids[j].id }
func (ids ids) Swap(i, j int)      { ids[i], ids[j] = ids[j], ids[i] }

// Search performs binary search for a component [ID].
// It returns ths index of the ID, and whether it was present in the list.
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
