package ecs

import (
	"math/bits"
	"sync/atomic"
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
		mask.Set(id.id, true)
	}
	return mask
}

// Get reports whether the bit at the given index [ID] is set.
func (b *bitMask64) Get(bit uint8) bool {
	mask := uint64(1 << bit)
	return b.bits&mask == mask
}

// Set sets the state of the bit at the given index.
func (b *bitMask64) Set(bit uint8, value bool) {
	if value {
		b.bits |= (1 << bit)
	} else {
		b.bits &^= (1 << bit)
	}
}

func (b *bitMask64) SetTrue(bit uint8) {
	b.bits |= (1 << bit)
}

func (b *bitMask64) SetFalse(bit uint8) {
	b.bits &= ^(1 << bit)
}

func (b *bitMask64) SetTrueSafe(bit uint8) {
	atomic.OrUint64(&b.bits, 1<<bit)
}

func (b *bitMask64) SetFalseSafe(bit uint8) {
	mask := ^(uint64(1) << bit)
	atomic.AndUint64(&b.bits, mask)
}

// Not returns the inversion of this mask.
func (b *bitMask64) Not() bitMask64 {
	return bitMask64{
		bits: ^b.bits,
	}
}

// OrI calculates the OR of this mask and other in-place.
func (b *bitMask64) OrI(other *bitMask64) {
	b.bits |= other.bits
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

	count := b.TotalBitsSet()
	types := make([]ID, count)
	totalIDs := reg.Count()

	idx := 0
	for j := range totalIDs {
		id := ID{id: uint8(j)}
		if b.Get(id.id) {
			types[idx] = id
			idx++
		}
	}
	return types
}
