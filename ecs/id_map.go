package ecs

const (
	idMapChunkSize = 16
)

type idMap struct {
	data []nodeID
	used bitMask
}

// newIDMap creates a new idMap
func newIDMap() idMap {
	return idMap{
		data: make([]nodeID, idMapChunkSize),
	}
}

// Get returns the value at the given key and whether the key is present.
func (m *idMap) Get(index uint8) (nodeID, bool) {
	if !m.used.Get(index) {
		return 0, false
	}
	return m.data[index], true
}

// Set sets the value at the given key.
func (m *idMap) Set(index uint8, value nodeID) {
	if len(m.data) <= int(index) {
		len := ((uint32(index) + idMapChunkSize) / idMapChunkSize) * idMapChunkSize
		data := make([]nodeID, len)
		copy(data, m.data)
		m.data = data
	}
	m.used.Set(index)
	m.data[index] = value
}
