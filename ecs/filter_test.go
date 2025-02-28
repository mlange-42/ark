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
		{NewFilter(id1, id2), All(id1, id2, id3), true},
		{NewFilter(id1, id2), All(id1), false},

		{NewFilter(id1, id2).Without(id3), All(id1, id2), true},
		{NewFilter(id1, id2).Without(id3), All(id1, id2, id3), false},
		{NewFilter(id1, id2).Without(id3), All(id1), false},

		{NewFilter(id1, id2).Exclusive(), All(id1, id2), true},
		{NewFilter(id1, id2).Exclusive(), All(id1, id2, id3), false},
		{NewFilter(id1, id2).Exclusive(), All(id1), false},
	}

	for _, test := range tests {
		assert.Equal(t, test.matches, test.filter.matches(&test.mask))
	}
}
