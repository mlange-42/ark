package ecs

import (
	"math/rand"
	"testing"
)

func all256(ids ...ID) *bitMask256 {
	mask := newMask256(ids...)
	return &mask
}

func TestMask256(t *testing.T) {
	big := uint8(mask256TotalBits - 2)
	mask := newMask256(id(1), id(2), id(13), id(27), id8(big))

	expectEqual(t, 5, mask.TotalBitsSet())

	expectTrue(t, mask.Get(id(1)))
	expectTrue(t, mask.Get(id(2)))
	expectTrue(t, mask.Get(id(13)))
	expectTrue(t, mask.Get(id(27)))
	expectTrue(t, mask.Get(id8(big)))

	expectFalse(t, mask.Get(id(0)))
	expectFalse(t, mask.Get(id(3)))
	expectFalse(t, mask.Get(id8(big-1)))
	expectFalse(t, mask.Get(id8(big+1)))

	mask.Set(id(0), true)
	mask.Set(id(1), false)

	expectTrue(t, mask.Get(id(0)))
	expectFalse(t, mask.Get(id(1)))

	other1 := newMask256(id(1), id(2), id(32))
	other2 := newMask256(id(0), id(2))

	expectFalse(t, mask.Contains(&other1))
	expectTrue(t, mask.Contains(&other2))

	mask.Reset()
	expectEqual(t, 0, mask.TotalBitsSet())

	mask = newMask256(id(1), id(2), id(13), id(27))
	other1 = newMask256(id(1), id(32))
	other2 = newMask256(id(0), id(32))

	expectTrue(t, mask.ContainsAny(&other1))
	expectFalse(t, mask.ContainsAny(&other2))

	expectEqual(t, newMask256(id(1), id(32)), newMask256(id(1), id(32)))
	expectNotEqual(t, newMask256(id(1), id(33)), newMask256(id(1), id(32)))

	mask = newMask256(id(1), id(32))
	not := mask.Not()

	expectTrue(t, not.Get(id(0)))
	expectFalse(t, not.Get(id(1)))
	expectFalse(t, not.Get(id(32)))

	expectTrue(t, mask.Equals(&mask))
	expectFalse(t, mask.Equals(&bitMask256{}))

	expectFalse(t, mask.IsZero())
	expectTrue(t, (&bitMask256{}).IsZero())

}

func TestBitMask256Copy(t *testing.T) {
	big := uint8(mask256TotalBits - 2)

	mask := newMask256(id(1), id(2), id(13), id(27), id8(big))
	mask2 := mask
	mask3 := &mask

	mask2.Set(id(1), false)
	mask3.Set(id(2), false)

	expectTrue(t, mask.Get(id(1)))
	expectFalse(t, mask2.Get(id(1)))

	expectTrue(t, mask2.Get(id(2)))
	expectFalse(t, mask.Get(id(2)))
	expectFalse(t, mask3.Get(id(2)))
}

func TestBitMask256(t *testing.T) {
	for i := range mask256TotalBits {
		mask := newMask256(id(i))
		expectEqual(t, 1, mask.TotalBitsSet())
		expectTrue(t, mask.Get(id(i)))
	}
	mask := bitMask256{}
	expectEqual(t, 0, mask.TotalBitsSet())

	for i := range mask256TotalBits {
		mask.Set(id(i), true)
		expectEqual(t, i+1, mask.TotalBitsSet())
		expectTrue(t, mask.Get(id(i)))
	}

	big := int(mask256TotalBits - 10)

	mask = newMask256(id(1), id(2), id(13), id(27), id(big), id(big+1), id(big+2))

	expectTrue(t, mask.Contains(all256(id(1), id(2), id(big), id(big+1))))
	expectFalse(t, mask.Contains(all256(id(1), id(2), id(big), id(big+5))))

	expectTrue(t, mask.ContainsAny(all256(id(6), id(big+2), id(big+6))))
	expectFalse(t, mask.ContainsAny(all256(id(6), id(big+3), id(big+5))))

	mask = newMask256()
	for i := range 256 {
		expectFalse(t, mask.Get(ID{uint8(i)}))
		mask.Set(ID{uint8(i)}, true)
		expectTrue(t, mask.Get(ID{uint8(i)}))
		mask.Set(ID{uint8(i)}, false)
		expectFalse(t, mask.Get(ID{uint8(i)}))
	}
}

func TestMask256ToTypes(t *testing.T) {
	w := NewWorld(1024)

	id1 := ComponentID[Position](&w)
	id2 := ComponentID[Velocity](&w)

	mask := newMask256()
	comps := mask.toTypes(&w.storage.registry.registry)
	expectSlicesEqual(t, []ID{}, comps)

	mask = newMask256(id1, id2)
	comps = mask.toTypes(&w.storage.registry.registry)
	expectSlicesEqual(t, []ID{id1, id2}, comps)
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
