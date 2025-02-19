package ecs

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIDsSort(t *testing.T) {
	idsSorted := newIDsSorted(id(2), id(6), id(2), id(0))
	assert.Equal(t, ids([]ID{id(0), id(2), id(2), id(6)}), idsSorted)
}

func TestIDsSearch(t *testing.T) {
	n := 100
	arr := make([]ID, n)
	for i := range n {
		arr[i] = id(i + 5)
	}
	idsSorted := newIDs(arr...)

	tests := []struct {
		search int
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

func benchmarkIDsSearch(b *testing.B, n int) {
	arr := make([]ID, n)
	for i := range n {
		arr[i] = id(i)
	}
	idsSorted := newIDs(arr...)
	searchFor := id(int(float32(n) * 0.6))

	for b.Loop() {
		_, _ = idsSorted.Search(searchFor)
	}
}

func benchmarkIDsSearchLinear(b *testing.B, n int) {
	arr := make([]ID, n)
	for i := range n {
		arr[i] = id(i)
	}
	idsSorted := newIDs(arr...)
	searchFor := id(int(float32(n) * 0.5))

	for b.Loop() {
		_, _ = idsSorted.SearchLinear(searchFor)
	}
}

func benchmarkIDsContains(b *testing.B, k, n int) {
	numQueries := 1000

	allIDs := make([]ID, n)
	for i := range n {
		allIDs[i] = id(i)
	}
	archIDs := newIDs(allIDs...)

	queries := make([]ids, numQueries)
	for i := range numQueries {
		rand.Shuffle(n, func(i, j int) {
			allIDs[i], allIDs[j] = allIDs[j], allIDs[i]
		})
		queries[i] = newIDsSorted(allIDs[:k]...)
	}

	for b.Loop() {
		for i := range numQueries {
			_ = archIDs.Contains(queries[i]...)
		}
	}
}

func BenchmarkIDsSearch_2(b *testing.B) {
	benchmarkIDsSearch(b, 2)
}

func BenchmarkIDsSearch_8(b *testing.B) {
	benchmarkIDsSearch(b, 8)
}

func BenchmarkIDsSearch_64(b *testing.B) {
	benchmarkIDsSearch(b, 64)
}

func BenchmarkIDsSearch_256(b *testing.B) {
	benchmarkIDsSearch(b, 256)
}

func BenchmarkIDsSearch_1024(b *testing.B) {
	benchmarkIDsSearch(b, 1024)
}

func BenchmarkIDsSearchLinear_2(b *testing.B) {
	benchmarkIDsSearchLinear(b, 2)
}

func BenchmarkIDsSearchLinear_8(b *testing.B) {
	benchmarkIDsSearchLinear(b, 8)
}

func BenchmarkIDsSearchLinear_64(b *testing.B) {
	benchmarkIDsSearchLinear(b, 64)
}

func BenchmarkIDsSearchLinear_256(b *testing.B) {
	benchmarkIDsSearchLinear(b, 256)
}

func BenchmarkIDsSearchLinear_1024(b *testing.B) {
	benchmarkIDsSearchLinear(b, 1024)
}

func BenchmarkIDsContains_1in8_1000(b *testing.B) {
	benchmarkIDsContains(b, 1, 8)
}

func BenchmarkIDsContains_2in8_1000(b *testing.B) {
	benchmarkIDsContains(b, 2, 8)
}

func BenchmarkIDsContains_4in8_1000(b *testing.B) {
	benchmarkIDsContains(b, 4, 8)
}

func BenchmarkIDsContains_8in8_1000(b *testing.B) {
	benchmarkIDsContains(b, 8, 8)
}

func BenchmarkIDsContains_1in32_1000(b *testing.B) {
	benchmarkIDsContains(b, 1, 32)
}

func BenchmarkIDsContains_2in32_1000(b *testing.B) {
	benchmarkIDsContains(b, 2, 32)
}

func BenchmarkIDsContains_4in32_1000(b *testing.B) {
	benchmarkIDsContains(b, 4, 32)
}

func BenchmarkIDsContains_8in32_1000(b *testing.B) {
	benchmarkIDsContains(b, 8, 32)
}
