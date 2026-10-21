package ecs

import (
	"math"
	"math/rand"
	"testing"
)

// newTestEntityIndex returns an entity index with the reserved entities, as created by newStorage.
func newTestEntityIndex(capacity int) []entityIndex {
	entities := make([]entityIndex, reservedEntities, capacity+reservedEntities)
	for i := range reservedEntities {
		entities[i] = entityIndex{table: maxTableID, gen: math.MaxUint32}
	}
	return entities
}

// testAlive returns whether the entity is alive, as storage.isAlive does.
func testAlive(entities []entityIndex, e Entity) bool {
	return entities[e.id].gen == e.gen
}

func TestEntityPoolConstructor(t *testing.T) {
	_ = newEntityPool(reservedEntities)
}

func TestEntityPool(t *testing.T) {
	p := newEntityPool(reservedEntities)
	entities := newTestEntityIndex(128)

	created := []Entity{}
	for range 5 {
		created = append(created, p.Get(&entities))
	}
	expectSlicesEqual(t, []Entity{newEntity(2), newEntity(3), newEntity(4), newEntity(5), newEntity(6)}, created, "Wrong initial entities")
	expectEqual(t, 7, len(entities))

	expectPanicsWithValue(t, "can't recycle reserved zero or wildcard entity", func() { p.Recycle(entities, Entity{0, math.MaxUint32}) })
	expectPanicsWithValue(t, "can't recycle reserved zero or wildcard entity", func() { p.Recycle(entities, Entity{1, math.MaxUint32}) })

	e0 := created[0]
	p.Recycle(entities, e0)
	expectFalse(t, testAlive(entities, e0), "Dead entity should not be alive")

	e0Old := e0
	e0 = p.Get(&entities)
	expectEqual(t, Entity{e0Old.id, 1}, e0, "Recycled entity should have the same ID and the next generation")
	expectTrue(t, testAlive(entities, e0), "Recycled entity of new generation should be alive")
	expectFalse(t, testAlive(entities, e0Old), "Recycled entity of old generation should not be alive")
	expectEqual(t, 7, len(entities), "Recycling should not grow the entity index")

	p.Recycle(entities, e0)
	for _, e := range created[1:] {
		p.Recycle(entities, e)
	}

	expectEqual(t, 5, p.Cap(entities))
	expectEqual(t, 0, p.Len(entities))
	expectEqual(t, 5, p.Available())
	expectEqual(t, 130, cap(entities))

	expectFalse(t, testAlive(entities, e0Old), "Recycled entity of old generation should not be alive")

	for range 5 {
		_ = p.Get(&entities)
	}

	expectEqual(t, 5, p.Len(entities))
	expectEqual(t, 0, p.Available())
	expectFalse(t, testAlive(entities, e0Old), "Recycled entity of old generation should not be alive")
	expectFalse(t, testAlive(entities, Entity{}), "Zero entity should not be alive")
	expectFalse(t, testAlive(entities, Entity{1, 0}), "Wildcard entity should not be alive")
}

func TestEntityPoolStochastic(t *testing.T) {
	n := 32
	p := newEntityPool(reservedEntities)
	entities := newTestEntityIndex(16)

	for range n {
		entities = entities[:reservedEntities]
		p.Reset()
		expectEqual(t, 0, p.Len(entities))
		expectEqual(t, 0, p.Available())

		alive := map[Entity]bool{}
		for range n {
			e := p.Get(&entities)
			alive[e] = true
		}

		for e, isAlive := range alive {
			expectEqual(t, isAlive, testAlive(entities, e), "Wrong alive state of entity %v after initialization", e)
			if rand.Float32() > 0.75 {
				continue
			}
			p.Recycle(entities, e)
			alive[e] = false
		}
		for e, isAlive := range alive {
			expectEqual(t, isAlive, testAlive(entities, e), "Wrong alive state of entity %v after 1st removal. Entity is %v", e, entities[e.id])
		}
		for range n {
			e := p.Get(&entities)
			alive[e] = true
		}
		for e, isAlive := range alive {
			expectEqual(t, isAlive, testAlive(entities, e), "Wrong alive state of entity %v after 1st recycling. Entity is %v", e, entities[e.id])
		}
		expectEqual(t, uint32(0), p.available, "No more entities should be available")

		for e, isAlive := range alive {
			if !isAlive || rand.Float32() > 0.75 {
				continue
			}
			p.Recycle(entities, e)
			alive[e] = false
		}
		for e, a := range alive {
			expectEqual(t, a, testAlive(entities, e), "Wrong alive state of entity %v after 2nd removal. Entity is %v", e, entities[e.id])
		}
	}
}

func TestBitPool(t *testing.T) {
	p := newBitPool()

	for i := range mask64TotalBits {
		expectEqual(t, i, int(p.Get()))
	}

	expectPanics(t, func() { p.Get() })

	for i := range 10 {
		p.Recycle(uint8(i))
	}
	for i := 9; i >= 0; i-- {
		expectEqual(t, i, int(p.Get()))
	}

	expectPanics(t, func() { p.Get() })

	p.Reset()

	for i := range mask64TotalBits {
		expectEqual(t, i, int(p.Get()))
	}

	expectPanics(t, func() { p.Get() })

	for i := range 10 {
		p.Recycle(uint8(i))
	}
	for i := 9; i >= 0; i-- {
		expectEqual(t, i, int(p.Get()))
	}
}

func TestIntPool(t *testing.T) {
	p := newIntPool[int](16)

	for range 3 {
		for i := range 32 {
			expectEqual(t, i, p.Get())
		}

		expectEqual(t, 32, len(p.pool))

		p.Recycle(3)
		p.Recycle(4)
		expectEqual(t, 4, p.Get())
		expectEqual(t, 3, p.Get())

		p.Reset()
	}
}

func BenchmarkBitPool(b *testing.B) {
	p := newBitPool()

	var l uint8
	for b.Loop() {
		l = p.Get()
		p.Recycle(l)
	}
	_ = l
}
