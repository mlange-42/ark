package ecs

import (
	"testing"
)

func TestIDMap(t *testing.T) {
	big1 := uint8(maskTotalBits - 20)
	big2 := uint8(maskTotalBits - 3)

	m := newIDMap[*Entity]()

	e0 := Entity{0, 0}
	e1 := Entity{1, 0}
	e121 := Entity{entityID(big1), 0}
	e200 := Entity{entityID(big2), 0}

	m.Set(0, &e0)
	m.Set(1, &e1)
	m.Set(big1, &e121)
	m.Set(big2, &e200)

	e, ok := m.Get(0)
	if !ok || *e != e0 {
		t.Errorf("expected Get(0) to return %v, got %v", e0, e)
	}

	e, ok = m.Get(1)
	if !ok || *e != e1 {
		t.Errorf("expected Get(1) to return %v, got %v", e1, e)
	}

	e, ok = m.Get(big1)
	if !ok || *e != e121 {
		t.Errorf("expected Get(big1) to return %v, got %v", e121, e)
	}

	e, ok = m.Get(big2)
	if !ok || *e != e200 {
		t.Errorf("expected Get(big2) to return %v, got %v", e200, e)
	}

	e, ok = m.Get(15)
	if ok || e != nil {
		t.Errorf("expected Get(15) to return nil, got %v", e)
	}

	m.Remove(0)
	m.Remove(1)

	e, ok = m.Get(0)
	if ok || e != nil {
		t.Errorf("expected Get(0) after removal to return nil, got %v", e)
	}

	if m.chunks[0] != nil {
		t.Errorf("expected m.chunks[0] to be nil after removals")
	}
}

func TestIDMapPointers(t *testing.T) {
	big1 := uint8(maskTotalBits - 20)
	big2 := uint8(maskTotalBits - 3)

	m := newIDMap[Entity]()

	e0 := Entity{0, 0}
	e1 := Entity{1, 0}
	e121 := Entity{entityID(big1), 0}
	e200 := Entity{entityID(big2), 0}

	m.Set(0, e0)
	m.Set(1, e1)
	m.Set(big1, e121)
	m.Set(big2, e200)

	e, ok := m.GetPointer(0)
	if !ok || *e != e0 {
		t.Errorf("expected GetPointer(0) to return %v, got %v", e0, *e)
	}

	e, ok = m.GetPointer(1)
	if !ok || *e != e1 {
		t.Errorf("expected GetPointer(1) to return %v, got %v", e1, *e)
	}

	e, ok = m.GetPointer(big1)
	if !ok || *e != e121 {
		t.Errorf("expected GetPointer(big1) to return %v, got %v", e121, *e)
	}

	e, ok = m.GetPointer(big2)
	if !ok || *e != e200 {
		t.Errorf("expected GetPointer(big2) to return %v, got %v", e200, *e)
	}

	e, ok = m.GetPointer(15)
	if ok || e != nil {
		t.Errorf("expected GetPointer(15) to return nil, got %v", e)
	}

	m.Remove(0)
	m.Remove(1)

	e, ok = m.GetPointer(0)
	if ok || e != nil {
		t.Errorf("expected GetPointer(0) after removal to return nil, got %v", e)
	}

	if m.chunks[0] != nil {
		t.Errorf("expected m.chunks[0] to be nil after removals")
	}
}

func BenchmarkIdMapping_IDMap(b *testing.B) {
	entities := [maskTotalBits]Entity{}
	m := newIDMap[*Entity]()

	for i := range maskTotalBits {
		entities[i] = Entity{entityID(i), 0}
		m.Set(uint8(i), &entities[i])
	}

	var ptr *Entity
	for i := 0; b.Loop(); i++ {
		ptr, _ = m.Get(uint8(i % maskTotalBits))
	}
	_ = ptr
}

func BenchmarkIdMapping_Array(b *testing.B) {
	entities := [maskTotalBits]Entity{}
	m := [maskTotalBits]*Entity{}

	for i := range maskTotalBits {
		entities[i] = Entity{entityID(i), 0}
		m[i] = &entities[i]
	}

	var ptr *Entity
	for i := 0; b.Loop(); i++ {
		ptr = m[i%maskTotalBits]
	}
	_ = ptr
}

func BenchmarkIdMapping_HashMap(b *testing.B) {
	entities := [maskTotalBits]Entity{}
	m := make(map[uint8]*Entity, maskTotalBits)

	for i := range maskTotalBits {
		entities[i] = Entity{entityID(i), 0}
		m[uint8(i)] = &entities[i]
	}

	var ptr *Entity
	for i := 0; b.Loop(); i++ {
		ptr = m[uint8(i%maskTotalBits)]
	}
	_ = ptr
}

func BenchmarkIdMapping_HashMapEntity(b *testing.B) {
	entities := [maskTotalBits]Entity{}
	m := make(map[Entity]*Entity, maskTotalBits)

	for i := range maskTotalBits {
		entities[i] = Entity{entityID(i), 0}
		m[Entity{entityID(i), 0}] = &entities[i]
	}

	var ptr *Entity
	for i := 0; b.Loop(); i++ {
		ptr = m[Entity{entityID(i % maskTotalBits), 0}]
	}
	_ = ptr
}
