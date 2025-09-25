package ecs

import (
	"math/bits"
)

// maskTotalBits is the size of a [bitMask] in bits.
// It is the maximum number of component types that may exist in any [World].
const mask256TotalBits = 256
const wordSize = 64

// bitMask is a 256 bit bit-mask.
type bitMask256 struct {
	bits [4]uint64 // 4x 64 bits of the mask
}

// newMask creates a new Mask from a list of IDs.
// Matches all entities that have the respective components, and potentially further components.
func newMask256(ids ...ID) bitMask256 {
	var mask bitMask256
	for _, id := range ids {
		mask.Set(id.id, true)
	}
	return mask
}

// Get reports whether the bit at the given index [ID] is set.
func (b *bitMask256) Get(bit uint8) bool {
	idx := bit >> 6
	mask := uint64(1) << (bit & 63)
	return b.bits[idx]&mask == mask
}

// Set sets the state of the bit at the given index.
func (b *bitMask256) Set(bit uint8, value bool) {
	idx := bit >> 6
	mask := uint64(1) << (bit & 63)
	if value {
		b.bits[idx] |= mask
	} else {
		b.bits[idx] &^= mask // faster than b.bits[idx] &= ^mask
	}
}

// Not returns the inversion of this mask.
func (b *bitMask256) Not() bitMask256 {
	return bitMask256{
		bits: [4]uint64{^b.bits[0], ^b.bits[1], ^b.bits[2], ^b.bits[3]},
	}
}

// OrI calculates the OR of this mask and other in-place.
func (b *bitMask256) OrI(other *bitMask256) {
	b.bits[0] |= other.bits[0]
	b.bits[1] |= other.bits[1]
	b.bits[2] |= other.bits[2]
	b.bits[3] |= other.bits[3]
}

// IsZero returns whether no bits are set in the mask.
func (b *bitMask256) IsZero() bool {
	return b.bits[0] == 0 && b.bits[1] == 0 && b.bits[2] == 0 && b.bits[3] == 0
}

// Reset the mask setting all bits to false.
func (b *bitMask256) Reset() {
	b.bits = [4]uint64{0, 0, 0, 0}
}

// SetAll sets all bits to 1.
func (b *bitMask256) SetAll() {
	b.bits = [4]uint64{^uint64(0), ^uint64(0), ^uint64(0), ^uint64(0)}
}

// Contains reports if the other mask is a subset of this mask.
func (b *bitMask256) Contains(other *bitMask256) bool {
	b0, b1, b2, b3 := b.bits[0], b.bits[1], b.bits[2], b.bits[3]
	o0, o1, o2, o3 := other.bits[0], other.bits[1], other.bits[2], other.bits[3]
	return b0&o0 == o0 &&
		b1&o1 == o1 &&
		b2&o2 == o2 &&
		b3&o3 == o3
}

// ContainsAny reports if any bit of the other mask is in this mask.
func (b *bitMask256) ContainsAny(other *bitMask256) bool {
	b0, b1, b2, b3 := b.bits[0], b.bits[1], b.bits[2], b.bits[3]
	o0, o1, o2, o3 := other.bits[0], other.bits[1], other.bits[2], other.bits[3]
	return b0&o0 != 0 ||
		b1&o1 != 0 ||
		b2&o2 != 0 ||
		b3&o3 != 0
}

// TotalBitsSet returns how many bits are set in this mask.
func (b *bitMask256) TotalBitsSet() int {
	return bits.OnesCount64(b.bits[0]) + bits.OnesCount64(b.bits[1]) + bits.OnesCount64(b.bits[2]) + bits.OnesCount64(b.bits[3])
}

// Equals returns whether two masks are equal.
func (b *bitMask256) Equals(other *bitMask256) bool {
	return b.bits == other.bits
}

func (b *bitMask256) toTypes(reg *registry) []ID {
	count := int(b.TotalBitsSet())
	types := make([]ID, count)

	totalIDs := reg.Count()
	bins := totalIDs/wordSize + 1
	bits := totalIDs % wordSize

	idx := 0
	for i := range bins {
		if b.bits[i] == 0 {
			continue
		}
		cnt := wordSize
		if i == bins-1 {
			cnt = bits
		}
		for j := range cnt {
			id := ID{id: uint8(i*wordSize + j)}
			if b.Get(id.id) {
				types[idx] = id
				idx++
			}
		}
	}
	return types
}
