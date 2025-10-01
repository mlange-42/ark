package ecs

import (
	"fmt"
	"math"
	"sync"
	"unsafe"
)

type number interface {
	int | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | cacheID | observerID
}

// entityPool is an implementation using implicit linked lists.
// Implements https://skypjack.github.io/2019-05-06-ecs-baf-part-3/
type entityPool struct {
	pointer   unsafe.Pointer
	entities  []Entity
	next      entityID
	available uint32
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
	p.pointer = unsafe.Pointer(&p.entities[0])
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

// bitPool is a pool of bits that makes it possible to obtain an un-set bit,
// and to recycle that bit for later use.
// This implementation uses an implicit list.
type bitPool struct {
	bits      []uint8
	length    uint8
	next      uint8
	available uint8
}

func newBitPool() bitPool {
	return bitPool{
		bits: make([]uint8, mask64TotalBits),
	}
}

// Get returns a fresh or recycled bit.
func (p *bitPool) Get() uint8 {
	if p.available == 0 {
		return p.getNew()
	}
	curr := p.next
	p.next, p.bits[p.next] = p.bits[p.next], p.next
	p.available--
	return p.bits[curr]
}

// Allocates and returns a new bit. For internal use.
func (p *bitPool) getNew() uint8 {
	if p.length >= mask64TotalBits {
		panic(fmt.Sprintf("run out of the maximum of %d bits. "+
			"This is likely caused by unclosed queries that lock the world. "+
			"Make sure that all queries finish their iteration or are closed manually", mask64TotalBits))
	}
	b := p.length
	p.bits[p.length] = b
	p.length++
	return b
}

// Recycle hands a bit back for recycling.
func (p *bitPool) Recycle(b uint8) {
	p.next, p.bits[b] = b, p.next
	p.available++
}

// Reset recycles all bits.
func (p *bitPool) Reset() {
	p.next = 0
	p.length = 0
	p.available = 0
}

// entityPool is an implementation using implicit linked lists.
// Implements https://skypjack.github.io/2019-05-06-ecs-baf-part-3/
type intPool[T number] struct {
	next              T
	pool              []T
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

type slicePools struct {
	relations slicePool[relationID]
	entities  slicePool[Entity]
	batches   slicePool[batchTable]
	tables    slicePool[tableID]
	ints      slicePool[uint32]
}

func newSlicePools() slicePools {
	return slicePools{
		relations: newSlicePool[relationID](8, 8),
		entities:  newSlicePool[Entity](8, 8),
		batches:   newSlicePool[batchTable](4, 128),
		tables:    newSlicePool[tableID](4, 128),
		ints:      newSlicePool[uint32](4, 128),
	}
}

type slicePool[E any] struct {
	free     [][]E
	sliceCap int
	mu       sync.Mutex
}

func newSlicePool[E any](size, sliceCap int) slicePool[E] {
	free := make([][]E, 0, size)
	for range size {
		free = append(free, make([]E, 0, sliceCap))
	}
	return slicePool[E]{
		free:     free,
		sliceCap: sliceCap,
	}
}

func (p *slicePool[E]) Get() []E {
	if len(p.free) == 0 {
		return make([]E, 0, p.sliceCap)
	}
	idx := len(p.free) - 1
	v := p.free[idx]
	p.free = p.free[:idx]
	return v
}

func (p *slicePool[E]) Has(size int) bool {
	if len(p.free) == 0 {
		return false
	}
	return cap(p.free[len(p.free)-1]) >= size
}

func (p *slicePool[E]) Recycle(s []E) {
	s = s[:0]
	p.free = append(p.free, s)
}
