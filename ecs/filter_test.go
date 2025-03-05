package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	id1 := ID{0}
	id2 := ID{1}
	id3 := ID{2}

	tests := []struct {
		filter  Filter
		mask    Mask
		matches bool
	}{
		{NewFilter(nil, id1, id2), NewMask(id1, id2, id3), true},
		{NewFilter(nil, id1, id2), NewMask(id1), false},

		{NewFilter(nil, id1, id2).Without(id3), NewMask(id1, id2), true},
		{NewFilter(nil, id1, id2).Without(id3), NewMask(id1, id2, id3), false},
		{NewFilter(nil, id1, id2).Without(id3), NewMask(id1), false},

		{NewFilter(nil, id1, id2).Exclusive(), NewMask(id1, id2), true},
		{NewFilter(nil, id1, id2).Exclusive(), NewMask(id1, id2, id3), false},
		{NewFilter(nil, id1, id2).Exclusive(), NewMask(id1), false},
	}

	for _, test := range tests {
		assert.Equal(t, test.matches, test.filter.matches(&test.mask))
	}
}

func BenchmarkFilterCopy(b *testing.B) {
	f := NewFilter(nil, id(1))

	var ff Filter
	for b.Loop() {
		ff = f
	}
	_ = ff
}
