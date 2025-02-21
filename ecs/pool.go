package ecs

import (
	"math"
	"unsafe"
)

// entityPool is an implementation using implicit linked lists.
// Implements https://skypjack.github.io/2019-05-06-ecs-baf-part-3/
type entityPool struct {
	entities  []Entity
	next      entityID
	available uint32
	pointer   unsafe.Pointer
	reserved  entityID
}

// newEntityPool creates a new, initialized Entity pool.
func newEntityPool(initialCapacity uint32, reserved uint32) entityPool {
	entities := make([]Entity, reserved, initialCapacity+reserved)
	// Reserved zero and wildcard entities.
	for i := range reserved {
		entities[i] = Entity{entityID(i), math.MaxUint32}
	}
	return entityPool{
		entities:  entities,
		next:      0,
		available: 0,
		pointer:   unsafe.Pointer(&entities[0]),
		reserved:  entityID(reserved),
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
	if e.id < p.reserved {
		panic("can't recycle reserved zero or wildcard entity")
	}
	p.entities[e.id].gen++
	p.next, p.entities[e.id].id = e.id, p.next
	p.available++
}

// Reset recycles all entities. Does NOT free the reserved memory.
func (p *entityPool) Reset() {
	p.entities = p.entities[:p.reserved]
	p.next = 0
	p.available = 0
}

// Alive returns whether an entity is still alive, based on the entity's generations.
func (p *entityPool) Alive(e Entity) bool {
	return e.gen == (*Entity)(unsafe.Add(p.pointer, entitySize*uintptr(e.id))).gen
}

// Len returns the current number of used entities.
func (p *entityPool) Len() int {
	return len(p.entities) - int(p.reserved) - int(p.available)
}

// Cap returns the current capacity (used and recycled entities).
func (p *entityPool) Cap() int {
	return len(p.entities) - int(p.reserved)
}

// TotalCap returns the current capacity in terms of reserved memory.
func (p *entityPool) TotalCap() int {
	return cap(p.entities)
}

// Available returns the current number of available/recycled entities.
func (p *entityPool) Available() int {
	return int(p.available)
}
