package ecs

import (
	"testing"
)

func TestFilter(t *testing.T) {
	id1 := ID{0}
	id2 := ID{1}
	id3 := ID{2}

	tests := []struct {
		filter  UnsafeFilter
		mask    bitMask
		matches bool
	}{
		{NewUnsafeFilter(nil, id1, id2).Without(), newMask(id1, id2, id3), true},
		{NewUnsafeFilter(nil, id1, id2), newMask(id1), false},

		{NewUnsafeFilter(nil, id1, id2).Without(id3), newMask(id1, id2), true},
		{NewUnsafeFilter(nil, id1, id2).Without(id3), newMask(id1, id2, id3), false},
		{NewUnsafeFilter(nil, id1, id2).Without(id3), newMask(id1), false},

		{NewUnsafeFilter(nil, id1, id2).Exclusive(), newMask(id1, id2), true},
		{NewUnsafeFilter(nil, id1, id2).Exclusive(), newMask(id1, id2, id3), false},
		{NewUnsafeFilter(nil, id1, id2).Exclusive(), newMask(id1), false},
	}

	for _, test := range tests {
		expectEqual(t, test.matches, test.filter.matches(&test.mask))
	}
}

func BenchmarkFilterCopy(b *testing.B) {
	f := NewUnsafeFilter(nil, id(1))

	var ff UnsafeFilter
	for b.Loop() {
		ff = f
	}
	_ = ff
}

func BenchmarkFilterMatches(b *testing.B) {
	filterMask := newMask(id(1), id(2))
	filter := filter{
		mask: filterMask,
	}
	mask := newMask(id(1), id(2), id(3))

	for b.Loop() {
		_ = filter.matches(&mask)
	}
}

func BenchmarkFilterWithoutMatches(b *testing.B) {
	filterMask := newMask(id(1), id(2))
	withoutMask := newMask(id(1), id(2))
	filter := filter{
		mask:       filterMask,
		without:    withoutMask,
		hasWithout: true,
	}
	mask := newMask(id(1), id(2), id(3))

	for b.Loop() {
		_ = filter.matches(&mask)
	}
}
