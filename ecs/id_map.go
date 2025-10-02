package ecs

import "math/bits"

const (
	idMapChunkSize = 16
	idMapChunks    = maskTotalBits / idMapChunkSize

	// idMapChunkMask is idMapChunkSize - 1
	idMapChunkMask = idMapChunkSize - 1
)

// idMapChunkShift is log2(idMapChunkSize)
var idMapChunkShift = bits.TrailingZeros(uint(idMapChunkSize))

// idMap maps component IDs to values.
//
// Is is a data structure meant for fast lookup while being memory-efficient.
// Access time is around 2ns, compared to 0.5ns for array access and 20ns for map[int]T.
//
// The memory footprint is reduced by using chunks, and only allocating chunks if they contain a key.
//
// The range of keys is limited from 0 to 255 (63 with build tag ark_tiny).
type idMap[T any] struct {
	zeroValue T
	chunks    [][]T
	chunkUsed []uint8
	used      bitMask
}

// newIDMap creates a new idMap
func newIDMap[T any]() idMap[T] {
	return idMap[T]{
		chunks:    make([][]T, idMapChunks),
		used:      bitMask{},
		chunkUsed: make([]uint8, idMapChunks),
	}
}

// Get returns the value at the given key and whether the key is present.
func (m *idMap[T]) Get(index uint8) (T, bool) {
	if !m.used.Get(index) {
		return m.zeroValue, false
	}
	chunk := index >> idMapChunkShift
	offset := index & idMapChunkMask
	return m.chunks[chunk][offset], true
}

// Set sets the value at the given key.
func (m *idMap[T]) Set(index uint8, value T) {
	chunk := index >> idMapChunkShift
	offset := index & idMapChunkMask

	if m.chunks[chunk] == nil {
		m.chunks[chunk] = make([]T, idMapChunkSize)
	}
	m.chunks[chunk][offset] = value
	m.used.Set(index, true)
	m.chunkUsed[chunk]++
}

// Remove removes the value at the given key.
// It de-allocates empty chunks.
func (m *idMap[T]) Remove(index uint8) {
	chunk := index >> idMapChunkShift
	m.used.Set(index, false)
	m.chunkUsed[chunk]--
	if m.chunkUsed[chunk] == 0 {
		m.chunks[chunk] = nil
	}
}
