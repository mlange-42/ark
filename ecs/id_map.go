package ecs

const (
	idMapChunkSize = 16
	idMapChunks    = MaskTotalBits / idMapChunkSize
)

// idMap maps component IDs to values.
//
// Is is a data structure meant for fast lookup while being memory-efficient.
// Access time is around 2ns, compared to 0.5ns for array access and 20ns for map[int]T.
//
// The memory footprint is reduced by using chunks, and only allocating chunks if they contain a key.
//
// The range of keys is limited from 0 to [MaskTotalBits]-1.
type idMap[T any] struct {
	zeroValue T
	chunks    [][]T
	chunkUsed []uint8
	used      Mask
}

// newIDMap creates a new idMap
func newIDMap[T any]() idMap[T] {
	return idMap[T]{
		chunks:    make([][]T, idMapChunks),
		used:      Mask{},
		chunkUsed: make([]uint8, idMapChunks),
	}
}

// Get returns the value at the given key and whether the key is present.
func (m *idMap[T]) Get(index uint8) (T, bool) {
	if !m.used.Get(id8(index)) {
		return m.zeroValue, false
	}
	return m.chunks[index/idMapChunkSize][index%idMapChunkSize], true
}

// Get returns a pointer to the value at the given key and whether the key is present.
func (m *idMap[T]) GetPointer(index uint8) (*T, bool) {
	if !m.used.Get(id8(index)) {
		return nil, false
	}
	return &m.chunks[index/idMapChunkSize][index%idMapChunkSize], true
}

// Set sets the value at the given key.
func (m *idMap[T]) Set(index uint8, value T) {
	chunk := index / idMapChunkSize
	if m.chunks[chunk] == nil {
		m.chunks[chunk] = make([]T, idMapChunkSize)
	}
	m.chunks[chunk][index%idMapChunkSize] = value
	m.used.Set(id8(index), true)
	m.chunkUsed[chunk]++
}

// Remove removes the value at the given key.
// It de-allocates empty chunks.
func (m *idMap[T]) Remove(index uint8) {
	chunk := index / idMapChunkSize
	m.used.Set(id8(index), false)
	m.chunkUsed[chunk]--
	if m.chunkUsed[chunk] == 0 {
		m.chunks[chunk] = nil
	}
}
