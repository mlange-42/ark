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

// ResID is the resource identifier.
// It is not relevant when using the default generic API.
type ResID struct {
	id uint8
}

// Batch is like a filter for batch processing of entities.
// Create it using [Filter2.Batch] etc.
type Batch struct {
	filter    Filter
	relations []RelationID
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
	ID         ID
	Type       reflect.Type
	IsRelation bool
}
