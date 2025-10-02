package ecs

const (
	idMapChunkSize = 16
	idMapChunks    = maskTotalBits / idMapChunkSize
)

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
	return m.chunks[index/idMapChunkSize][index%idMapChunkSize], true
}

// GetPointer returns a pointer to the value at the given key and whether the key is present.
func (m *idMap[T]) GetPointer(index uint8) (*T, bool) {
	if !m.used.Get(index) {
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
	m.used.Set(index, true)
	m.chunkUsed[chunk]++
}

// Remove removes the value at the given key.
// It de-allocates empty chunks.
func (m *idMap[T]) Remove(index uint8) {
	chunk := index / idMapChunkSize
	m.used.Set(index, false)
	m.chunkUsed[chunk]--
	if m.chunkUsed[chunk] == 0 {
		m.chunks[chunk] = nil
	}
}

// idMap maps component IDs to values.
//
// It is a data structure meant for fast lookup while being memory-efficient.
// Access time is around 2ns, compared to 0.5ns for array access and 20ns for map[int]T.
//
// The range of keys is limited from 0 to 255 (or 63 with build tag ark_tiny).
type idMapSimple[T any] struct {
	zeroValue T
	values    []T
	used      bitMask
}

// newIDMap creates a new idMap
func newIDMapSimple[T any]() idMapSimple[T] {
	return idMapSimple[T]{
		values: make([]T, maskTotalBits),
		used:   bitMask{},
	}
}

// Get returns the value at the given key and whether the key is present.
func (m *idMapSimple[T]) Get(index uint8) (T, bool) {
	if !m.used.Get(index) {
		return m.zeroValue, false
	}
	return m.values[index], true
}

// GetPointer returns a pointer to the value at the given key and whether the key is present.
func (m *idMapSimple[T]) GetPointer(index uint8) (*T, bool) {
	if !m.used.Get(index) {
		return nil, false
	}
	return &m.values[index], true
}

// Set sets the value at the given key.
func (m *idMapSimple[T]) Set(index uint8, value T) {
	m.values[index] = value
	m.used.Set(index, true)
}

// Remove removes the value at the given key.
func (m *idMapSimple[T]) Remove(index uint8) {
	m.used.Set(index, false)
}
