package ecs

import "reflect"

// ID is the component identifier.
// It is not relevant when using the default generic API.
type ID struct {
	id uint8
}

func id(id int) ID {
	return ID{uint8(id)}
}

func id8(id uint8) ID {
	return ID{id}
}

// Index returns the internal component index of this component ID.
func (id ID) Index() uint8 {
	return id.id
}

// IDs is an immutable list of [ID] values.
type IDs struct {
	data []ID
}

func newIDs(ids []ID) IDs {
	return IDs{
		data: ids,
	}
}

// Get returns the ID at the given index.
func (ids *IDs) Get(index int) ID {
	return ids.data[index]
}

// Len returns the number of IDs.
func (ids *IDs) Len() int {
	return len(ids.data)
}

// ResID is the resource identifier.
// It is not relevant when using the default generic API.
type ResID struct {
	id uint8
}

// Index returns the internal component index of this resource ID.
func (id ResID) Index() uint8 {
	return id.id
}

// Batch is like a filter for batch processing of entities.
// Create it using [Filter2.Batch] etc.
//
// ⚠️ This should not be stored, but used immediately and re-generated
// each time a batch operation is called.
// Otherwise, changes to the origin filter or calls to [Filter2.Batch]
// with different relationship targets may modify stored instances.
type Batch struct {
	filter    *filter
	relations []relationID
	cache     cacheID
}

// EntityDump is a dump of the entire entity data of the world.
//
// See [World.DumpEntities] and [World.LoadEntities].
type EntityDump struct {
	Entities  []Entity // Entities in the World's entity pool.
	Alive     []uint32 // IDs of all alive entities in query iteration order.
	Next      uint32   // The next free entity of the World's entity pool.
	Available uint32   // The number of allocated and available entities in the World's entity pool.
}

// CompInfo provides information about a registered component.
// Returned by [ComponentInfo].
type CompInfo struct {
	Type       reflect.Type
	ID         ID
	IsRelation bool
}
