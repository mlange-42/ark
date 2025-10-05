package ecs

import (
	"testing"
)

func TestIDMap(t *testing.T) {
	big1 := uint8(maskTotalBits - 20)
	big2 := uint8(maskTotalBits - 3)

	m := newIDMap()

	e0 := nodeID(0)
	e1 := nodeID(1)
	e121 := nodeID(big1)
	e200 := nodeID(big2)

	m.Set(0, e0)
	m.Set(1, e1)
	m.Set(big1, e121)
	m.Set(big2, e200)

	e, ok := m.Get(0)
	expectTrue(t, ok)
	expectEqual(t, e0, e)

	e, ok = m.Get(1)
	expectTrue(t, ok)
	expectEqual(t, e1, e)

	e, ok = m.Get(big1)
	expectTrue(t, ok)
	expectEqual(t, e121, e)

	e, ok = m.Get(big2)
	expectTrue(t, ok)
	expectEqual(t, e200, e)

	e, ok = m.Get(15)
	expectFalse(t, ok)
	expectEqual(t, 0, e)
}

func BenchmarkIdMapping_IDMap(b *testing.B) {

	entities := [maskTotalBits]nodeID{}
	m := newIDMap()

	for i := range maskTotalBits {
		entities[i] = nodeID(i)
		m.Set(uint8(i), entities[i])
	}

	var v nodeID
	for i := 0; b.Loop(); i++ {
		v, _ = m.Get(uint8(i % maskTotalBits))
	}
	_ = v
}

func BenchmarkIdMapping_Array(b *testing.B) {

	entities := [maskTotalBits]nodeID{}
	m := [maskTotalBits]nodeID{}

	for i := range maskTotalBits {
		entities[i] = nodeID(i)
		m[i] = entities[i]
	}

	var v nodeID
	for i := 0; b.Loop(); i++ {
		v = m[i%maskTotalBits]
	}
	_ = v
}

func BenchmarkIdMapping_HashMap(b *testing.B) {

	entities := [maskTotalBits]nodeID{}
	m := make(map[uint8]nodeID, maskTotalBits)

	for i := range maskTotalBits {
		entities[i] = nodeID(i)
		m[uint8(i)] = entities[i]
	}

	var v nodeID
	for i := 0; b.Loop(); i++ {
		v = m[uint8(i%maskTotalBits)]
	}
	_ = v
}
