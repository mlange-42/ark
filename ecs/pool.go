package ecs

import (
	"math"
	"unsafe"
)

// Reserved Entities.
// Add this to initial capacities of entity pool and lists,
// to avoid unexpected allocations.
const reservedEntities = 2

// entityPool is an implementation using implicit linked lists.
// Implements https://skypjack.github.io/2019-05-06-ecs-baf-part-3/
type entityPool struct {
	entities  []Entity
	next      entityID
	available uint32
	pointer   unsafe.Pointer
}

// newEntityPool creates a new, initialized Entity pool.
func newEntityPool(initialCapacity uint32) entityPool {
	entities := make([]Entity, 2, initialCapacity+reservedEntities)
	// The zero entity
	entities[0] = Entity{0, math.MaxUint32}
	// The wildcard entity
	entities[1] = Entity{1, math.MaxUint32}
	return entityPool{
		entities:  entities,
		next:      0,
		available: 0,
		pointer:   unsafe.Pointer(&entities[0]),
	}
}

// Get returns a fresh or recycled entity.
func (p *entityPool) Get() Entity {
	if p.available == 0 {
		return p.getNew()
	}
	curr := p.next
	p.next, p.entities[p.next].id = p.entities[p.next].id, p.next
	p.available--
	return p.entities[curr]
}

// Allocates and returns a new entity. For internal use.
func (p *entityPool) getNew() Entity {
	e := Entity{id: entityID(len(p.entities)), gen: 0}
	p.entities = append(p.entities, e)
	return e
}

// Recycle hands an entity back for recycling.
func (p *entityPool) Recycle(e Entity) {
	if e.id < 2 {
		panic("can't recycle reserved zero or wildcard entity")
	}
	p.entities[e.id].gen++
	p.next, p.entities[e.id].id = e.id, p.next
	p.available++
}

// Reset recycles all entities. Does NOT free the reserved memory.
func (p *entityPool) Reset() {
	p.entities = p.entities[:reservedEntities]
	p.next = 0
	p.available = 0
}

// Alive returns whether an entity is still alive, based on the entity's generations.
func (p *entityPool) Alive(e Entity) bool {
	return e.gen == (*Entity)(unsafe.Add(p.pointer, entitySize*uintptr(e.id))).gen
}

// Len returns the current number of used entities.
func (p *entityPool) Len() int {
	return len(p.entities) - reservedEntities - int(p.available)
}

// Cap returns the current capacity (used and recycled entities).
func (p *entityPool) Cap() int {
	return len(p.entities) - reservedEntities
}

// TotalCap returns the current capacity in terms of reserved memory.
func (p *entityPool) TotalCap() int {
	return cap(p.entities)
}

// Available returns the current number of available/recycled entities.
func (p *entityPool) Available() int {
	return int(p.available)
}
