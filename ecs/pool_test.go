package ecs

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityPoolConstructor(t *testing.T) {
	_ = newEntityPool(128, reservedEntities)
}

func TestEntityPool(t *testing.T) {
	p := newEntityPool(128, reservedEntities)

	expectedAll := []Entity{newEntity(0), newEntity(1), newEntity(2), newEntity(3), newEntity(4), newEntity(5), newEntity(6)}
	expectedAll[0].gen = math.MaxUint32
	expectedAll[1].gen = math.MaxUint32

	for i := 0; i < 5; i++ {
		_ = p.Get()
	}
	assert.Equal(t, expectedAll, p.entities, "Wrong initial entities")

	assert.PanicsWithValue(t, "can't recycle reserved zero or wildcard entity", func() { p.Recycle(p.entities[0]) })
	assert.PanicsWithValue(t, "can't recycle reserved zero or wildcard entity", func() { p.Recycle(p.entities[1]) })

	e0 := p.entities[reservedEntities]
	p.Recycle(e0)
	assert.False(t, p.Alive(e0), "Dead entity should not be alive")

	e0Old := e0
	e0 = p.Get()
	expectedAll[reservedEntities].gen++
	assert.True(t, p.Alive(e0), "Recycled entity of new generation should be alive")
	assert.False(t, p.Alive(e0Old), "Recycled entity of old generation should not be alive")

	assert.Equal(t, expectedAll, p.entities, "Wrong entities after get/recycle")

	e0Old = p.entities[reservedEntities]
	for i := 0; i < 5; i++ {
		p.Recycle(p.entities[i+reservedEntities])
		expectedAll[i+1].gen++
	}

	assert.Equal(t, 5, p.Cap())
	assert.Equal(t, 130, p.TotalCap())

	assert.False(t, p.Alive(e0Old), "Recycled entity of old generation should not be alive")

	for i := 0; i < 5; i++ {
		_ = p.Get()
	}

	assert.False(t, p.Alive(e0Old), "Recycled entity of old generation should not be alive")
	assert.False(t, p.Alive(Entity{}), "Zero entity should not be alive")
	assert.False(t, p.Alive(Entity{1, 0}), "Wildcard entity should not be alive")
}

func TestEntityPoolStochastic(t *testing.T) {
	n := 32
	p := newEntityPool(16, reservedEntities)

	for range n {
		p.Reset()
		assert.Equal(t, 0, p.Len())
		assert.Equal(t, 0, p.Available())

		alive := map[Entity]bool{}
		for range n {
			e := p.Get()
			alive[e] = true
		}

		for e, isAlive := range alive {
			assert.Equal(t, isAlive, p.Alive(e), "Wrong alive state of entity %v after initialization", e)
			if rand.Float32() > 0.75 {
				continue
			}
			p.Recycle(e)
			alive[e] = false
		}
		for e, isAlive := range alive {
			assert.Equal(t, isAlive, p.Alive(e), "Wrong alive state of entity %v after 1st removal. Entity is %v", e, p.entities[e.id])
		}
		for range n {
			e := p.Get()
			alive[e] = true
		}
		for e, isAlive := range alive {
			assert.Equal(t, isAlive, p.Alive(e), "Wrong alive state of entity %v after 1st recycling. Entity is %v", e, p.entities[e.id])
		}
		assert.Equal(t, uint32(0), p.available, "No more entities should be available")

		for e, isAlive := range alive {
			if !isAlive || rand.Float32() > 0.75 {
				continue
			}
			p.Recycle(e)
			alive[e] = false
		}
		for e, a := range alive {
			assert.Equal(t, a, p.Alive(e), "Wrong alive state of entity %v after 2nd removal. Entity is %v", e, p.entities[e.id])
		}
	}
}

func TestBitPool(t *testing.T) {
	p := newBitPool()

	for i := 0; i < mask64TotalBits; i++ {
		assert.Equal(t, i, int(p.Get()))
	}

	assert.Panics(t, func() { p.Get() })

	for i := 0; i < 10; i++ {
		p.Recycle(uint8(i))
	}
	for i := 9; i >= 0; i-- {
		assert.Equal(t, i, int(p.Get()))
	}

	assert.Panics(t, func() { p.Get() })

	p.Reset()

	for i := 0; i < mask64TotalBits; i++ {
		assert.Equal(t, i, int(p.Get()))
	}

	assert.Panics(t, func() { p.Get() })

	for i := 0; i < 10; i++ {
		p.Recycle(uint8(i))
	}
	for i := 9; i >= 0; i-- {
		assert.Equal(t, i, int(p.Get()))
	}
}

func TestIntPool(t *testing.T) {
	p := newIntPool[int](16)

	for n := 0; n < 3; n++ {
		for i := 0; i < 32; i++ {
			assert.Equal(t, i, p.Get())
		}

		assert.Equal(t, 32, len(p.pool))

		p.Recycle(3)
		p.Recycle(4)
		assert.Equal(t, 4, p.Get())
		assert.Equal(t, 3, p.Get())

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
