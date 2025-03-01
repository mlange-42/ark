//go:build !tiny

package ecs

import (
	"math/bits"
)

// MaskTotalBits is the size of a [Mask] in bits.
// It is the maximum number of component types that may exist in any [World].
const MaskTotalBits = 256
const wordSize = 64

// Mask is a 256 bit bit-mask.
// It is also a [Filter] for including certain components.
type Mask struct {
	bits [4]uint64 // 4x 64 bits of the mask
}

// NewMask creates a new Mask from a list of IDs.
// Matches all entities that have the respective components, and potentially further components.
func NewMask(ids ...ID) Mask {
	var mask Mask
	for _, id := range ids {
		mask.Set(id, true)
	}
	return mask
}

// Get reports whether the bit at the given index [ID] is set.
func (b *Mask) Get(bit ID) bool {
	idx := bit.id / 64
	offset := bit.id - (64 * idx)
	mask := uint64(1 << offset)
	return b.bits[idx]&mask == mask
}

// Set sets the state of the bit at the given index.
func (b *Mask) Set(bit ID, value bool) {
	idx := bit.id / 64
	offset := bit.id - (64 * idx)
	if value {
		b.bits[idx] |= (1 << offset)
	} else {
		b.bits[idx] &= ^(1 << offset)
	}
}

// Not returns the inversion of this mask.
func (b *Mask) Not() Mask {
	return Mask{
		bits: [4]uint64{^b.bits[0], ^b.bits[1], ^b.bits[2], ^b.bits[3]},
	}
}

// IsZero returns whether no bits are set in the mask.
func (b *Mask) IsZero() bool {
	return b.bits[0] == 0 && b.bits[1] == 0 && b.bits[2] == 0 && b.bits[3] == 0
}

// Reset the mask setting all bits to false.
func (b *Mask) Reset() {
	b.bits = [4]uint64{0, 0, 0, 0}
}

// Contains reports if the other mask is a subset of this mask.
func (b *Mask) Contains(other *Mask) bool {
	return b.bits[0]&other.bits[0] == other.bits[0] &&
		b.bits[1]&other.bits[1] == other.bits[1] &&
		b.bits[2]&other.bits[2] == other.bits[2] &&
		b.bits[3]&other.bits[3] == other.bits[3]
}

// ContainsAny reports if any bit of the other mask is in this mask.
func (b *Mask) ContainsAny(other *Mask) bool {
	return b.bits[0]&other.bits[0] != 0 ||
		b.bits[1]&other.bits[1] != 0 ||
		b.bits[2]&other.bits[2] != 0 ||
		b.bits[3]&other.bits[3] != 0
}

// And returns the bitwise AND of two masks.
func (b *Mask) And(other *Mask) Mask {
	return Mask{
		bits: [4]uint64{
			b.bits[0] & other.bits[0],
			b.bits[1] & other.bits[1],
			b.bits[2] & other.bits[2],
			b.bits[3] & other.bits[3],
		},
	}
}

// Or returns the bitwise OR of two masks.
func (b *Mask) Or(other *Mask) Mask {
	return Mask{
		bits: [4]uint64{
			b.bits[0] | other.bits[0],
			b.bits[1] | other.bits[1],
			b.bits[2] | other.bits[2],
			b.bits[3] | other.bits[3],
		},
	}
}

// Xor returns the bitwise XOR of two masks.
func (b *Mask) Xor(other *Mask) Mask {
	return Mask{
		bits: [4]uint64{
			b.bits[0] ^ other.bits[0],
			b.bits[1] ^ other.bits[1],
			b.bits[2] ^ other.bits[2],
			b.bits[3] ^ other.bits[3],
		},
	}
}

// TotalBitsSet returns how many bits are set in this mask.
func (b *Mask) TotalBitsSet() int {
	return bits.OnesCount64(b.bits[0]) + bits.OnesCount64(b.bits[1]) + bits.OnesCount64(b.bits[2]) + bits.OnesCount64(b.bits[3])
}

// Equals returns whether two masks are equal.
func (b *Mask) Equals(other *Mask) bool {
	return b.bits == other.bits
}

func (b *Mask) toTypes(reg *registry) []ID {
	count := int(b.TotalBitsSet())
	types := make([]ID, count)

	totalIDs := reg.Count()
	bins := totalIDs/wordSize + 1
	bits := totalIDs % wordSize

	idx := 0
	for i := 0; i < bins; i++ {
		if b.bits[i] == 0 {
			continue
		}
		cnt := wordSize
		if i == bins-1 {
			cnt = bits
		}
		for j := 0; j < cnt; j++ {
			id := ID{id: uint8(i*wordSize + j)}
			if b.Get(id) {
				types[idx] = id
				idx++
			}
		}
	}
	return types
}
