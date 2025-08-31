package ecs

import (
	"math/rand"
	"testing"
)

func all(ids ...ID) *bitMask256 {
	mask := newMask256(ids...)
	return &mask
}

func TestMask256(t *testing.T) {
	big := uint8(mask256TotalBits - 2)
	mask := newMask256(id(1), id(2), id(13), id(27), id8(big))
	if mask.TotalBitsSet() != 5 {
		t.Errorf("expected 5 bits to be set, got %d", mask.TotalBitsSet())
	}
	if !mask.Get(id(1)) {
		t.Errorf("expected bit %d to be set", id(1))
	}
	if !mask.Get(id(2)) {
		t.Errorf("expected bit %d to be set", id(2))
	}
	if !mask.Get(id(13)) {
		t.Errorf("expected bit %d to be set", id(13))
	}
	if !mask.Get(id(27)) {
		t.Errorf("expected bit %d to be set", id(27))
	}
	if !mask.Get(id8(big)) {
		t.Errorf("expected bit %d to be set", id8(big))
	}
	if mask.Get(id(0)) {
		t.Errorf("expected bit %d to be unset", id(0))
	}
	if mask.Get(id(3)) {
		t.Errorf("expected bit %d to be unset", id(3))
	}
	if mask.Get(id8(big - 1)) {
		t.Errorf("expected bit %d to be unset", id8(big-1))
	}
	if mask.Get(id8(big + 1)) {
		t.Errorf("expected bit %d to be unset", id8(big+1))
	}
	mask.Set(id(0), true)
	mask.Set(id(1), false)
	if !mask.Get(id(0)) {
		t.Errorf("expected bit %d to be set", id(0))
	}
	if mask.Get(id(1)) {
		t.Errorf("expected bit %d to be unset", id(1))
	}
	other1 := newMask256(id(1), id(2), id(32))
	other2 := newMask256(id(0), id(2))
	if mask.Contains(&other1) {
		t.Errorf("expected mask not to contain other1")
	}
	if !mask.Contains(&other2) {
		t.Errorf("expected mask to contain other2")
	}
	mask.Reset()
	if mask.TotalBitsSet() != 0 {
		t.Errorf("expected 0 bits to be set after reset, got %d", mask.TotalBitsSet())
	}
	mask = newMask256(id(1), id(2), id(13), id(27))
	other1 = newMask256(id(1), id(32))
	other2 = newMask256(id(0), id(32))
	if !mask.ContainsAny(&other1) {
		t.Errorf("expected mask to contain any of other1")
	}
	if mask.ContainsAny(&other2) {
		t.Errorf("expected mask not to contain any of other2")
	}
	if newMask256(id(1), id(33)) == newMask256(id(1), id(32)) {
		t.Errorf("expected masks not to be equal")
	}
	mask = newMask256(id(1), id(32))
	not := mask.Not()
	if !not.Get(id(0)) {
		t.Errorf("expected bit %d to be set in not", id(0))
	}
	if not.Get(id(1)) {
		t.Errorf("expected bit %d to be unset in not", id(1))
	}
	if not.Get(id(32)) {
		t.Errorf("expected bit %d to be unset in not", id(32))
	}
	if !mask.Equals(&mask) {
		t.Errorf("expected mask to be equal to itself")
	}
	if mask.Equals(&bitMask256{}) {
		t.Errorf("expected mask not to be equal to empty mask")
	}
	if mask.IsZero() {
		t.Errorf("expected mask not to be zero")
	}
	if !(&bitMask256{}).IsZero() {
		t.Errorf("expected empty mask to be zero")
	}
}

func TestBitMask256Copy(t *testing.T) {
	big := uint8(mask256TotalBits - 2)
	mask := newMask256(id(1), id(2), id(13), id(27), id8(big))
	mask2 := mask
	mask3 := &mask
	mask2.Set(id(1), false)
	mask3.Set(id(2), false)
	if !mask.Get(id(1)) {
		t.Errorf("expected bit %d to be set in mask", id(1))
	}
	if mask2.Get(id(1)) {
		t.Errorf("expected bit %d to be unset in mask2", id(1))
	}
	if !mask2.Get(id(2)) {
		t.Errorf("expected bit %d to be set in mask2", id(2))
	}
	if mask.Get(id(2)) {
		t.Errorf("expected bit %d to be unset in mask", id(2))
	}
	if mask3.Get(id(2)) {
		t.Errorf("expected bit %d to be unset in mask3", id(2))
	}
}

func TestBitMask256(t *testing.T) {
	for i := range mask256TotalBits {
		mask := newMask256(id(i))
		if mask.TotalBitsSet() != 1 {
			t.Errorf("expected 1 bit to be set, got %d", mask.TotalBitsSet())
		}
		if !mask.Get(id(i)) {
			t.Errorf("expected bit %d to be set", id(i))
		}
	}
	mask := bitMask256{}
	if mask.TotalBitsSet() != 0 {
		t.Errorf("expected 0 bits to be set, got %d", mask.TotalBitsSet())
	}
	for i := range mask256TotalBits {
		mask.Set(id(i), true)
		if mask.TotalBitsSet() != i+1 {
			t.Errorf("expected %d bits to be set, got %d", i+1, mask.TotalBitsSet())
		}
		if !mask.Get(id(i)) {
			t.Errorf("expected bit %d to be set", id(i))
		}
	}
	big := int(mask256TotalBits - 10)
	mask = newMask256(id(1), id(2), id(13), id(27), id(big), id(big+1), id(big+2))
	if !mask.Contains(all(id(1), id(2), id(big), id(big+1))) {
		t.Errorf("expected mask to contain all")
	}
	if mask.Contains(all(id(1), id(2), id(big), id(big+5))) {
		t.Errorf("expected mask not to contain all")
	}
	if !mask.ContainsAny(all(id(6), id(big+2), id(big+6))) {
		t.Errorf("expected mask to contain any of all")
	}
	if mask.ContainsAny(all(id(6), id(big+3), id(big+5))) {
		t.Errorf("expected mask not to contain any of all")
	}
}

func TestMask256ToTypes(t *testing.T) {
	w := NewWorld(1024)
	id1 := ComponentID[Position](&w)
	id2 := ComponentID[Velocity](&w)
	mask := newMask256()
	comps := mask.toTypes(&w.storage.registry.registry)
	if len(comps) != 0 {
		t.Errorf("expected 0 components, got %d", len(comps))
	}
	mask = newMask256(id1, id2)
	comps = mask.toTypes(&w.storage.registry.registry)
	if len(comps) != 2 {
		t.Errorf("expected 2 components, got %d", len(comps))
	}
	if comps[0] != id1 {
		t.Errorf("expected component %d, got %d", id1, comps[0])
	}
	if comps[1] != id2 {
		t.Errorf("expected component %d, got %d", id2, comps[1])
	}
}

func BenchmarkMask256Get(b *testing.B) {
	mask := newMask256()
	for i := range mask256TotalBits {
		if rand.Float64() < 0.5 {
			mask.Set(id(i), true)
		}
	}
	idx := id(rand.Intn(mask256TotalBits))
	var v bool
	for b.Loop() {
		v = mask.Get(idx)
	}
	_ = v
}

func BenchmarkMask256Contains(b *testing.B) {
	mask := newMask256()
	for i := range mask256TotalBits {
		if rand.Float64() < 0.5 {
			mask.Set(id(i), true)
		}
	}
	filter := newMask256(id(rand.Intn(mask256TotalBits)))
	var v bool
	for b.Loop() {
		v = mask.Contains(&filter)
	}
	_ = v
}

func BenchmarkMask256ContainsAny(b *testing.B) {
	mask := newMask256()
	for i := range mask256TotalBits {
		if rand.Float64() < 0.5 {
			mask.Set(id(i), true)
		}
	}
	filter := newMask256(id(rand.Intn(mask256TotalBits)))
	var v bool
	for b.Loop() {
		v = mask.ContainsAny(&filter)
	}
	_ = v
}

func BenchmarkMask256Match(b *testing.B) {
	mask := newMask256(id(0), id(1), id(2))
	bits := newMask256(id(0), id(1), id(2))
	var v bool
	for b.Loop() {
		v = bits.Contains(&mask)
	}
	_ = v
}

func BenchmarkMask256Copy(b *testing.B) {
	mask := newMask256(id(0), id(1), id(2))
	var tempMask bitMask256
	for b.Loop() {
		tempMask = mask
	}
	_ = tempMask
}
