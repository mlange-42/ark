package ecs

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func all(ids ...ID) *bitMask256 {
	mask := newMask256(ids...)
	return &mask
}

func TestMask256(t *testing.T) {
	big := uint8(mask256TotalBits - 2)
	mask := newMask256(id(1), id(2), id(13), id(27), id8(big))

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

	other1 := newMask256(id(1), id(2), id(32))
	other2 := newMask256(id(0), id(2))

	assert.False(t, mask.Contains(&other1))
	assert.True(t, mask.Contains(&other2))

	mask.Reset()
	assert.Equal(t, 0, mask.TotalBitsSet())

	mask = newMask256(id(1), id(2), id(13), id(27))
	other1 = newMask256(id(1), id(32))
	other2 = newMask256(id(0), id(32))

	assert.True(t, mask.ContainsAny(&other1))
	assert.False(t, mask.ContainsAny(&other2))

	assert.Equal(t, newMask256(id(1), id(32)), newMask256(id(1), id(32)))
	assert.NotEqual(t, newMask256(id(1), id(33)), newMask256(id(1), id(32)))

	mask = newMask256(id(1), id(32))
	not := mask.Not()

	assert.True(t, not.Get(id(0)))
	assert.False(t, not.Get(id(1)))
	assert.False(t, not.Get(id(32)))
}

func TestBitMask256Copy(t *testing.T) {
	big := uint8(mask256TotalBits - 2)

	mask := newMask256(id(1), id(2), id(13), id(27), id8(big))
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
	for i := 0; i < mask256TotalBits; i++ {
		mask := newMask256(id(i))
		assert.Equal(t, 1, mask.TotalBitsSet())
		assert.True(t, mask.Get(id(i)))
	}
	mask := bitMask256{}
	assert.Equal(t, 0, mask.TotalBitsSet())

	for i := 0; i < mask256TotalBits; i++ {
		mask.Set(id(i), true)
		assert.Equal(t, i+1, mask.TotalBitsSet())
		assert.True(t, mask.Get(id(i)))
	}

	big := int(mask256TotalBits - 10)

	mask = newMask256(id(1), id(2), id(13), id(27), id(big), id(big+1), id(big+2))

	assert.True(t, mask.Contains(all(id(1), id(2), id(big), id(big+1))))
	assert.False(t, mask.Contains(all(id(1), id(2), id(big), id(big+5))))

	assert.True(t, mask.ContainsAny(all(id(6), id(big+2), id(big+6))))
	assert.False(t, mask.ContainsAny(all(id(6), id(big+3), id(big+5))))
}

func TestMask256ToTypes(t *testing.T) {
	w := NewWorld(1024)

	id1 := ComponentID[Position](&w)
	id2 := ComponentID[Velocity](&w)

	mask := newMask256()
	comps := mask.toTypes(&w.storage.registry.registry)
	assert.Equal(t, []ID{}, comps)

	mask = newMask256(id1, id2)
	comps = mask.toTypes(&w.storage.registry.registry)
	assert.Equal(t, []ID{id1, id2}, comps)
}

func BenchmarkMask256Get(b *testing.B) {
	mask := newMask256()
	for i := 0; i < mask256TotalBits; i++ {
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
	for i := 0; i < mask256TotalBits; i++ {
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
	for i := 0; i < mask256TotalBits; i++ {
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
