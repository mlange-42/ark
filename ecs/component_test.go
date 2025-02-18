package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIDsSort(t *testing.T) {
	idsSorted := newSortedIDs(id(2), id(6), id(2), id(0))
	assert.Equal(t, ids([]ID{id(0), id(2), id(2), id(6)}), idsSorted)
}

func TestIDsSearch(t *testing.T) {
	n := 100
	arr := make([]ID, n)
	for i := range n {
		arr[i] = id(uint32(i + 5))
	}
	idsSorted := newIDs(arr...)

	tests := []struct {
		search uint32
		index  int
		found  bool
	}{
		{5, 0, true},
		{104, 99, true},
		{4, 0, false},
		{105, 100, false},
	}

	for _, test := range tests {
		idx, ok := idsSorted.Search(id(test.search))
		assert.Equal(t, test.found, ok)
		assert.Equal(t, test.index, idx)
	}
}
