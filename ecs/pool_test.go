package ecs

import (
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
)

func TestEntityPoolConstructor(t *testing.T) {
	_ = newEntityPool(128, reservedEntities)
}

func TestEntityPool(t *testing.T) {
	p := newEntityPool(128, reservedEntities)

	expectedAll := []Entity{newEntity(0), newEntity(1), newEntity(2), newEntity(3), newEntity(4), newEntity(5), newEntity(6)}
	expectedAll[0].gen = math.MaxUint32
	expectedAll[1].gen = math.MaxUint32

	for range 5 {
		_ = p.Get()
	}
	expectSlicesEqual(t, expectedAll, p.entities, "Wrong initial entities")

	expectPanicsWithValue(t, "can't recycle reserved zero or wildcard entity", func() { p.Recycle(p.entities[0]) })
	expectPanicsWithValue(t, "can't recycle reserved zero or wildcard entity", func() { p.Recycle(p.entities[1]) })

	e0 := p.entities[reservedEntities]
	p.Recycle(e0)
	expectFalse(t, p.Alive(e0), "Dead entity should not be alive")

	e0Old := e0
	e0 = p.Get()
	expectedAll[reservedEntities].gen++
	expectTrue(t, p.Alive(e0), "Recycled entity of new generation should be alive")
	expectFalse(t, p.Alive(e0Old), "Recycled entity of old generation should not be alive")

	expectSlicesEqual(t, expectedAll, p.entities, "Wrong entities after get/recycle")

	e0Old = p.entities[reservedEntities]
	for i := range 5 {
		p.Recycle(p.entities[i+reservedEntities])
		expectedAll[i+1].gen++
	}

	expectEqual(t, 5, p.Cap())
	expectEqual(t, 130, p.TotalCap())

	expectFalse(t, p.Alive(e0Old), "Recycled entity of old generation should not be alive")

	for range 5 {
		_ = p.Get()
	}

	expectFalse(t, p.Alive(e0Old), "Recycled entity of old generation should not be alive")
	expectFalse(t, p.Alive(Entity{}), "Zero entity should not be alive")
	expectFalse(t, p.Alive(Entity{1, 0}), "Wildcard entity should not be alive")
}

func TestEntityPoolStochastic(t *testing.T) {
	n := 32
	p := newEntityPool(16, reservedEntities)

	for range n {
		p.Reset()
		expectEqual(t, 0, p.Len())
		expectEqual(t, 0, p.Available())

		alive := map[Entity]bool{}
		for range n {
			e := p.Get()
			alive[e] = true
		}

		for e, isAlive := range alive {
			expectEqual(t, isAlive, p.Alive(e), "Wrong alive state of entity %v after initialization", e)
			if rand.Float32() > 0.75 {
				continue
			}
			p.Recycle(e)
			alive[e] = false
		}
		for e, isAlive := range alive {
			expectEqual(t, isAlive, p.Alive(e), "Wrong alive state of entity %v after 1st removal. Entity is %v", e, p.entities[e.id])
		}
		for range n {
			e := p.Get()
			alive[e] = true
		}
		for e, isAlive := range alive {
			expectEqual(t, isAlive, p.Alive(e), "Wrong alive state of entity %v after 1st recycling. Entity is %v", e, p.entities[e.id])
		}
		expectEqual(t, uint32(0), p.available, "No more entities should be available")

		for e, isAlive := range alive {
			if !isAlive || rand.Float32() > 0.75 {
				continue
			}
			p.Recycle(e)
			alive[e] = false
		}
		for e, a := range alive {
			expectEqual(t, a, p.Alive(e), "Wrong alive state of entity %v after 2nd removal. Entity is %v", e, p.entities[e.id])
		}
	}
}

func TestBitPoolGet(t *testing.T) {
	p := newBitPool()

	allocated := make([]uint8, 64)
	for i := 0; i < 64; i++ {
		b := p.Get()
		allocated[i] = b
	}

	seen := make(map[uint8]bool)
	for _, b := range allocated {
		expectFalse(t, seen[b])
		seen[b] = true
	}

	expectPanics(t, func() {
		_ = p.Get()
	})

	p.Reset()
	expectEqual(t, ^uint64(0), p.free)
}

func TestBitPoolRecycle(t *testing.T) {
	p := newBitPool()

	b := p.Get()
	p.Recycle(b)

	b2 := p.Get()
	expectEqual(t, b, b2)
}

func TestBitPoolSafe(t *testing.T) {
	p := newBitPool()

	var wg sync.WaitGroup
	var count int32
	seen := make([]bool, 64)

	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b := p.GetSafe()
			atomic.AddInt32(&count, 1)
			seen[b] = true
		}()
	}
	wg.Wait()

	expectEqual(t, 64, count)

	for i := uint8(0); i < 64; i++ {
		expectTrue(t, seen[i], "bit %d was not allocated", i)
		wg.Add(1)
		go func(i uint8) {
			defer wg.Done()
			p.RecycleSafe(i)
		}(i)
	}
	wg.Wait()

	for i := 0; i < 64; i++ {
		_ = p.GetSafe()
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

func BenchmarkPoolAlive(b *testing.B) {
	pool := newEntityPool(1024, reservedEntities)

	for range 100 {
		_ = pool.Get()
	}
	entity := Entity{50, 0}

	for b.Loop() {
		_ = pool.Alive(entity)
	}
}
