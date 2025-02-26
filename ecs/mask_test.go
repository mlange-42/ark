package ecs

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func all(ids ...ID) *Mask {
	mask := All(ids...)
	return &mask
}

func TestMask(t *testing.T) {
	big := uint8(MaskTotalBits - 2)
	mask := All(id(1), id(2), id(13), id(27), id8(big))

	assert.Equal(t, 5, mask.TotalBitsSet())

	assert.True(t, mask.Get(id(1)))
	assert.True(t, mask.Get(id(2)))
	assert.True(t, mask.Get(id(13)))
	assert.True(t, mask.Get(id(27)))
	assert.True(t, mask.Get(id8(big)))

	assert.False(t, mask.Get(id(0)))
	assert.False(t, mask.Get(id(3)))
	assert.False(t, mask.Get(id8(big-1)))
	assert.False(t, mask.Get(id8(big+1)))

	mask.Set(id(0), true)
	mask.Set(id(1), false)

	assert.True(t, mask.Get(id(0)))
	assert.False(t, mask.Get(id(1)))

	other1 := All(id(1), id(2), id(32))
	other2 := All(id(0), id(2))

	assert.False(t, mask.Contains(&other1))
	assert.True(t, mask.Contains(&other2))

	mask.Reset()
	assert.Equal(t, 0, mask.TotalBitsSet())

	mask = All(id(1), id(2), id(13), id(27))
	other1 = All(id(1), id(32))
	other2 = All(id(0), id(32))

	assert.True(t, mask.ContainsAny(&other1))
	assert.False(t, mask.ContainsAny(&other2))
}

func TestBitMaskLogic(t *testing.T) {
	big := uint8(MaskTotalBits - 2)

	assert.Equal(t, All(id(5)), all(id(0), id(5)).And(all(id(5), id8(big))))
	assert.Equal(t, All(id(0), id(5), id8(big)), all(id(0), id(5)).Or(all(id(5), id8(big))))
	assert.Equal(t, All(id(0), id8(big)), all(id(0), id(5)).Xor(all(id(5), id8(big))))
}

func TestBitMaskCopy(t *testing.T) {
	big := uint8(MaskTotalBits - 2)

	mask := All(id(1), id(2), id(13), id(27), id8(big))
	mask2 := mask
	mask3 := &mask

	mask2.Set(id(1), false)
	mask3.Set(id(2), false)

	assert.True(t, mask.Get(id(1)))
	assert.False(t, mask2.Get(id(1)))

	assert.True(t, mask2.Get(id(2)))
	assert.False(t, mask.Get(id(2)))
	assert.False(t, mask3.Get(id(2)))
}

func TestBitMask256(t *testing.T) {
	for i := 0; i < MaskTotalBits; i++ {
		mask := All(id(i))
		assert.Equal(t, 1, mask.TotalBitsSet())
		assert.True(t, mask.Get(id(i)))
	}
	mask := Mask{}
	assert.Equal(t, 0, mask.TotalBitsSet())

	for i := 0; i < MaskTotalBits; i++ {
		mask.Set(id(i), true)
		assert.Equal(t, i+1, mask.TotalBitsSet())
		assert.True(t, mask.Get(id(i)))
	}

	big := int(MaskTotalBits - 10)

	mask = All(id(1), id(2), id(13), id(27), id(big), id(big+1), id(big+2))

	assert.True(t, mask.Contains(all(id(1), id(2), id(big), id(big+1))))
	assert.False(t, mask.Contains(all(id(1), id(2), id(big), id(big+5))))

	assert.True(t, mask.ContainsAny(all(id(6), id(big+2), id(big+6))))
	assert.False(t, mask.ContainsAny(all(id(6), id(big+3), id(big+5))))
}

func TestMaskToTypes(t *testing.T) {
	w := NewWorld(1024)

	id1 := ComponentID[Position](&w)
	id2 := ComponentID[Velocity](&w)

	mask := All()
	comps := mask.toTypes(&w.storage.registry.registry)
	assert.Equal(t, []ID{}, comps)

	mask = All(id1, id2)
	comps = mask.toTypes(&w.storage.registry.registry)
	assert.Equal(t, []ID{id1, id2}, comps)
}

func BenchmarkMaskGet(b *testing.B) {
	mask := All()
	for i := 0; i < MaskTotalBits; i++ {
		if rand.Float64() < 0.5 {
			mask.Set(id(i), true)
		}
	}
	idx := id(rand.Intn(MaskTotalBits))

	var v bool
	for b.Loop() {
		v = mask.Get(idx)
	}
	_ = v
}

func BenchmarkMaskContains(b *testing.B) {
	mask := All()
	for i := 0; i < MaskTotalBits; i++ {
		if rand.Float64() < 0.5 {
			mask.Set(id(i), true)
		}
	}
	filter := All(id(rand.Intn(MaskTotalBits)))

	var v bool
	for b.Loop() {
		v = mask.Contains(&filter)
	}
	_ = v
}

func BenchmarkMaskContainsAny(b *testing.B) {
	mask := All()
	for i := 0; i < MaskTotalBits; i++ {
		if rand.Float64() < 0.5 {
			mask.Set(id(i), true)
		}
	}
	filter := All(id(rand.Intn(MaskTotalBits)))

	var v bool
	for b.Loop() {
		v = mask.ContainsAny(&filter)
	}
	_ = v
}

func BenchmarkMaskMatch(b *testing.B) {
	mask := All(id(0), id(1), id(2))
	bits := All(id(0), id(1), id(2))
	var v bool
	for b.Loop() {
		v = mask.Matches(&bits)
	}
	_ = v
}

func BenchmarkMaskCopy(b *testing.B) {
	mask := All(id(0), id(1), id(2))
	var tempMask Mask
	for b.Loop() {
		tempMask = mask
	}
	_ = tempMask
}
