package ecs

import (
	"fmt"
)

type number interface {
	int | int8 | uint8 | int16 | uint16 | int32 | uint32 | cacheID | observerID
}

// entityPool recycles entity IDs, using implicit linked lists.
// Implements https://skypjack.github.io/2019-05-06-ecs-baf-part-3/
//
// The pool has no storage of its own. It works on the entity index of the storage,
// see [entityIndex]. For a removed entity, the row field of its index entry
// holds the ID of the next free entity, and the gen field the generation
// it gets when it is recycled.
type entityPool struct {
	next      entityID
	available uint32
	reserved  entityID
}

// newEntityPool creates a new, initialized Entity pool.
func newEntityPool(reserved uint32) entityPool {
	return entityPool{
		reserved: entityID(reserved),
	}
}

// Get returns a fresh or recycled entity.
// For a fresh entity, it appends an index entry to entities.
// The caller must set table and row of the entity's index entry.
func (p *entityPool) Get(entities *[]entityIndex) Entity {
	if p.available == 0 {
		e := Entity{id: entityID(len(*entities)), gen: 0}
		*entities = append(*entities, entityIndex{table: maxTableID})
		return e
	}
	curr := p.next
	index := &(*entities)[curr]
	p.next = entityID(index.row)
	p.available--
	return Entity{id: curr, gen: index.gen}
}

// Recycle hands an entity back for recycling.
// Marks the entity as dead in its index entry.
func (p *entityPool) Recycle(entities []entityIndex, e Entity) {
	if e.id < p.reserved {
		panic("can't recycle reserved zero or wildcard entity")
	}
	index := &entities[e.id]
	index.table = maxTableID
	index.row = uint32(p.next)
	index.gen++
	p.next = e.id
	p.available++
}

// Reset recycles all entities. The caller must truncate the entity index to the reserved entities.
func (p *entityPool) Reset() {
	p.next = 0
	p.available = 0
}

// Len returns the current number of used entities.
func (p *entityPool) Len(entities []entityIndex) int {
	return len(entities) - int(p.reserved) - int(p.available)
}

// Cap returns the current capacity (used and recycled entities).
func (p *entityPool) Cap(entities []entityIndex) int {
	return len(entities) - int(p.reserved)
}

// Available returns the current number of available/recycled entities.
func (p *entityPool) Available() int {
	return int(p.available)
}

// bitPool is a pool of bits that makes it possible to obtain an un-set bit,
// and to recycle that bit for later use.
// This implementation uses an implicit list.
type bitPool struct {
	free []uint8
}

// newBitPool creates a new empty bitPool.
func newBitPool() bitPool {
	free := make([]uint8, mask64TotalBits)
	for i := range mask64TotalBits {
		free[i] = uint8(63 - i)
	}
	return bitPool{
		free: free,
	}
}

// Get returns a fresh or recycled bit.
func (p *bitPool) Get() uint8 {
	if len(p.free) == 0 {
		panic(fmt.Sprintf("run out of the maximum of %d bits. "+
			"This is likely caused by unclosed queries that lock the world. "+
			"Make sure that all queries finish their iteration or are closed manually", mask64TotalBits))
	}
	b := p.free[len(p.free)-1]
	p.free = p.free[:len(p.free)-1]
	return b
}

// Recycle hands a bit back for recycling.
func (p *bitPool) Recycle(b uint8) {
	p.free = append(p.free, b)
}

// Reset recycles all bits.
func (p *bitPool) Reset() {
	p.free = p.free[:mask64TotalBits]
	for i := range mask64TotalBits {
		p.free[i] = uint8(63 - i)
	}
}

// entityPool is an implementation using implicit linked lists.
// Implements https://skypjack.github.io/2019-05-06-ecs-baf-part-3/
type intPool[T number] struct {
	pool              []T
	next              T
	available         uint32
	capacityIncrement uint32
}

// newEntityPool creates a new, initialized Entity pool.
func newIntPool[T number](capacityIncrement uint32) intPool[T] {
	return intPool[T]{
		pool:              make([]T, 0, capacityIncrement),
		next:              0,
		available:         0,
		capacityIncrement: capacityIncrement,
	}
}

// Get returns a fresh or recycled entity.
func (p *intPool[T]) Get() T {
	if p.available == 0 {
		return p.getNew()
	}
	curr := p.next
	p.next, p.pool[p.next] = p.pool[p.next], p.next
	p.available--
	return p.pool[curr]
}

// Allocates and returns a new entity. For internal use.
func (p *intPool[T]) getNew() T {
	e := T(len(p.pool))
	if len(p.pool) == cap(p.pool) {
		old := p.pool
		p.pool = make([]T, len(p.pool), len(p.pool)+int(p.capacityIncrement))
		copy(p.pool, old)
	}
	p.pool = append(p.pool, e)
	return e
}

// Recycle hands an entity back for recycling.
func (p *intPool[T]) Recycle(e T) {
	p.next, p.pool[e] = e, p.next
	p.available++
}

// Reset recycles all entities. Does NOT free the reserved memory.
func (p *intPool[T]) Reset() {
	p.pool = p.pool[:0]
	p.next = 0
	p.available = 0
}
