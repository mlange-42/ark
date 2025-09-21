package ecs

import (
	"math/bits"
)

// maskTotalBits is the size of a [bitMask] in bits.
// It is the maximum number of component types that may exist in any [World].
const mask64TotalBits = 64

// bitMask is a 64 bit bit-mask.
type bitMask64 struct {
	bits uint64 // 64 bits of the mask
}

// newMask creates a new Mask from a list of IDs.
// Matches all entities that have the respective components, and potentially further components.
func newMask64(ids ...ID) bitMask64 {
	var mask bitMask64
	for _, id := range ids {
		mask.Set(id, true)
	}
	return mask
}

// Get reports whether the bit at the given index [ID] is set.
func (b *bitMask64) Get(bit ID) bool {
	mask := uint64(1 << bit.id)
	return b.bits&mask == mask
}

// Set sets the state of the bit at the given index.
func (b *bitMask64) Set(bit ID, value bool) {
	if value {
		b.bits |= (1 << bit.id)
	} else {
		b.bits &= ^(1 << bit.id)
	}
}

func (b *bitMask64) SetTrue(bit ID) {
	b.bits |= (1 << bit.id)
}

func (b *bitMask64) SetFalse(bit ID) {
	b.bits &= ^(1 << bit.id)
}

// Not returns the inversion of this mask.
func (b *bitMask64) Not() bitMask64 {
	return bitMask64{
		bits: ^b.bits,
	}
}

// IsZero returns whether no bits are set in the mask.
func (b *bitMask64) IsZero() bool {
	return b.bits == 0
}

// Reset the mask setting all bits to false.
func (b *bitMask64) Reset() {
	b.bits = 0
}

// Contains reports if the other mask is a subset of this mask.
func (b *bitMask64) Contains(other *bitMask64) bool {
	return b.bits&other.bits == other.bits
}

// ContainsAny reports if any bit of the other mask is in this mask.
func (b *bitMask64) ContainsAny(other *bitMask64) bool {
	return b.bits&other.bits != 0
}

// TotalBitsSet returns how many bits are set in this mask.
func (b *bitMask64) TotalBitsSet() int {
	return bits.OnesCount64(b.bits)
}

// Equals returns whether two masks are equal.
func (b *bitMask64) Equals(other *bitMask64) bool {
	return b.bits == other.bits
}

func (b *bitMask64) toTypes(reg *registry) []ID {
	if b.bits == 0 {
		return []ID{}
	}

	count := int(b.TotalBitsSet())
	types := make([]ID, count)
	totalIDs := reg.Count()

	idx := 0
	for j := range totalIDs {
		id := ID{id: uint8(j)}
		if b.Get(id) {
			types[idx] = id
			idx++
		}
	}
	return types
}
