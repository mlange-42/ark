package ecs

import (
	"math"
	"math/rand"
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
	if !equalSlices(expectedAll, p.entities) {
		t.Errorf("Wrong initial entities")
	}
	expectPanicWithValue(t, "can't recycle reserved zero or wildcard entity", func() { p.Recycle(p.entities[0]) })
	expectPanicWithValue(t, "can't recycle reserved zero or wildcard entity", func() { p.Recycle(p.entities[1]) })
	e0 := p.entities[reservedEntities]
	p.Recycle(e0)
	if p.Alive(e0) {
		t.Errorf("Dead entity should not be alive")
	}
	e0Old := e0
	e0 = p.Get()
	expectedAll[reservedEntities].gen++
	if !p.Alive(e0) {
		t.Errorf("Recycled entity of new generation should be alive")
	}
	if p.Alive(e0Old) {
		t.Errorf("Recycled entity of old generation should not be alive")
	}
	if !equalSlices(expectedAll, p.entities) {
		t.Errorf("Wrong entities after get/recycle")
	}
	e0Old = p.entities[reservedEntities]
	for i := range 5 {
		p.Recycle(p.entities[i+reservedEntities])
		expectedAll[i+1].gen++
	}
	if p.Cap() != 5 {
		t.Errorf("Expected capacity 5, got %d", p.Cap())
	}
	if p.TotalCap() != 130 {
		t.Errorf("Expected total capacity 130, got %d", p.TotalCap())
	}
	if p.Alive(e0Old) {
		t.Errorf("Recycled entity of old generation should not be alive")
	}
	for range 5 {
		_ = p.Get()
	}
	if p.Alive(e0Old) {
		t.Errorf("Recycled entity of old generation should not be alive")
	}
	if p.Alive(Entity{}) {
		t.Errorf("Zero entity should not be alive")
	}
	if p.Alive(Entity{1, 0}) {
		t.Errorf("Wildcard entity should not be alive")
	}
}

func TestEntityPoolStochastic(t *testing.T) {
	n := 32
	p := newEntityPool(16, reservedEntities)
	for range n {
		p.Reset()
		if p.Len() != 0 {
			t.Errorf("Expected length 0, got %d", p.Len())
		}
		if p.Available() != 0 {
			t.Errorf("Expected available 0, got %d", p.Available())
		}
		alive := map[Entity]bool{}
		for range n {
			e := p.Get()
			alive[e] = true
		}
		for e, isAlive := range alive {
			if p.Alive(e) != isAlive {
				t.Errorf("Wrong alive state of entity %v after initialization", e)
			}
			if rand.Float32() > 0.75 {
				continue
			}
			p.Recycle(e)
			alive[e] = false
		}
		for e, isAlive := range alive {
			if p.Alive(e) != isAlive {
				t.Errorf("Wrong alive state of entity %v after 1st removal. Entity is %v", e, p.entities[e.id])
			}
		}
		for range n {
			e := p.Get()
			alive[e] = true
		}
		for e, isAlive := range alive {
			if p.Alive(e) != isAlive {
				t.Errorf("Wrong alive state of entity %v after 1st recycling. Entity is %v", e, p.entities[e.id])
			}
		}
		if p.available != 0 {
			t.Errorf("No more entities should be available")
		}
		for e, isAlive := range alive {
			if !isAlive || rand.Float32() > 0.75 {
				continue
			}
			p.Recycle(e)
			alive[e] = false
		}
		for e, a := range alive {
			if p.Alive(e) != a {
				t.Errorf("Wrong alive state of entity %v after 2nd removal. Entity is %v", e, p.entities[e.id])
			}
		}
	}
}

func TestBitPool(t *testing.T) {
	p := newBitPool()
	for i := range mask64TotalBits {
		if p.Get() != uint8(i) {
			t.Errorf("Expected %d, got %d", i, p.Get())
		}
	}
	expectPanic(t, func() { p.Get() })
	for i := range 10 {
		p.Recycle(uint8(i))
	}
	for i := 9; i >= 0; i-- {
		if p.Get() != uint8(i) {
			t.Errorf("Expected %d, got %d", i, p.Get())
		}
	}
	expectPanic(t, func() { p.Get() })
	p.Reset()
	for i := range mask64TotalBits {
		if p.Get() != uint8(i) {
			t.Errorf("Expected %d, got %d", i, p.Get())
		}
	}
	expectPanic(t, func() { p.Get() })
	for i := range 10 {
		p.Recycle(uint8(i))
	}
	for i := 9; i >= 0; i-- {
		if p.Get() != uint8(i) {
			t.Errorf("Expected %d, got %d", i, p.Get())
		}
	}
}

func TestIntPool(t *testing.T) {
	p := newIntPool[int](16)
	for range 3 {
		for i := range 32 {
			if p.Get() != i {
				t.Errorf("Expected %d, got %d", i, p.Get())
			}
		}
		if len(p.pool) != 32 {
			t.Errorf("Expected pool length 32, got %d", len(p.pool))
		}
		p.Recycle(3)
		p.Recycle(4)
		if p.Get() != 4 {
			t.Errorf("Expected 4, got %d", p.Get())
		}
		if p.Get() != 3 {
			t.Errorf("Expected 3, got %d", p.Get())
		}
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
