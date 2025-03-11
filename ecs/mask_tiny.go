//go:build tiny

package ecs

import (
	"math/bits"
)

// maskTotalBits is the size of a [bitMask] in bits.
// It is the maximum number of component types that may exist in any [World].
const maskTotalBits = 64

// bitMask is a 64 bit bit-mask.
// It is also a [Filter] for including certain components.
type bitMask struct {
	bits uint64 // 4x 64 bits of the mask
}

// newMask creates a new Mask from a list of IDs.
// Matches all entities that have the respective components, and potentially further components.
func newMask(ids ...ID) bitMask {
	var mask bitMask
	for _, id := range ids {
		mask.Set(id, true)
	}
	return mask
}

// Get reports whether the bit at the given index [ID] is set.
func (b *bitMask) Get(bit ID) bool {
	mask := uint64(1 << bit.id)
	return b.bits&mask == mask
}

// Set sets the state of the bit at the given index.
func (b *bitMask) Set(bit ID, value bool) {
	if value {
		b.bits |= (1 << bit.id)
	} else {
		b.bits &= ^(1 << bit.id)
	}
}

// Not returns the inversion of this mask.
func (b *bitMask) Not() bitMask {
	return bitMask{
		bits: ^b.bits,
	}
}

// IsZero returns whether no bits are set in the mask.
func (b *bitMask) IsZero() bool {
	return b.bits == 0
}

// Reset the mask setting all bits to false.
func (b *bitMask) Reset() {
	b.bits = 0
}

// Contains reports if the other mask is a subset of this mask.
func (b *bitMask) Contains(other *bitMask) bool {
	return b.bits&other.bits == other.bits
}

// ContainsAny reports if any bit of the other mask is in this mask.
func (b *bitMask) ContainsAny(other *bitMask) bool {
	return b.bits&other.bits != 0
}

// TotalBitsSet returns how many bits are set in this mask.
func (b *bitMask) TotalBitsSet() int {
	return bits.OnesCount64(b.bits)
}

// Equals returns whether two masks are equal.
func (b *bitMask) Equals(other *bitMask) bool {
	return b.bits == other.bits
}

func (b *bitMask) toTypes(reg *registry) []ID {
	if b.bits == 0 {
		return []ID{}
	}

	count := int(b.TotalBitsSet())
	types := make([]ID, count)
	totalIDs := reg.Count()

	idx := 0
	for j := 0; j < totalIDs; j++ {
		id := ID{id: uint8(j)}
		if b.Get(id) {
			types[idx] = id
			idx++
		}
	}
	return types
}
