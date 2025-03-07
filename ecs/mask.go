//go:build !tiny

package ecs

import (
	"math/bits"
)

// MaskTotalBits is the size of a [bitMask] in bits.
// It is the maximum number of component types that may exist in any [World].
const MaskTotalBits = 256
const wordSize = 64

// bitMask is a 256 bit bit-mask.
// It is also a [Filter] for including certain components.
type bitMask struct {
	bits [4]uint64 // 4x 64 bits of the mask
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
	idx := bit.id / 64
	offset := bit.id - (64 * idx)
	mask := uint64(1 << offset)
	return b.bits[idx]&mask == mask
}

// Set sets the state of the bit at the given index.
func (b *bitMask) Set(bit ID, value bool) {
	idx := bit.id / 64
	offset := bit.id - (64 * idx)
	if value {
		b.bits[idx] |= (1 << offset)
	} else {
		b.bits[idx] &= ^(1 << offset)
	}
}

// Not returns the inversion of this mask.
func (b *bitMask) Not() bitMask {
	return bitMask{
		bits: [4]uint64{^b.bits[0], ^b.bits[1], ^b.bits[2], ^b.bits[3]},
	}
}

// IsZero returns whether no bits are set in the mask.
func (b *bitMask) IsZero() bool {
	return b.bits[0] == 0 && b.bits[1] == 0 && b.bits[2] == 0 && b.bits[3] == 0
}

// Reset the mask setting all bits to false.
func (b *bitMask) Reset() {
	b.bits = [4]uint64{0, 0, 0, 0}
}

// Contains reports if the other mask is a subset of this mask.
func (b *bitMask) Contains(other *bitMask) bool {
	return b.bits[0]&other.bits[0] == other.bits[0] &&
		b.bits[1]&other.bits[1] == other.bits[1] &&
		b.bits[2]&other.bits[2] == other.bits[2] &&
		b.bits[3]&other.bits[3] == other.bits[3]
}

// ContainsAny reports if any bit of the other mask is in this mask.
func (b *bitMask) ContainsAny(other *bitMask) bool {
	return b.bits[0]&other.bits[0] != 0 ||
		b.bits[1]&other.bits[1] != 0 ||
		b.bits[2]&other.bits[2] != 0 ||
		b.bits[3]&other.bits[3] != 0
}

// And returns the bitwise AND of two masks.
func (b *bitMask) And(other *bitMask) bitMask {
	return bitMask{
		bits: [4]uint64{
			b.bits[0] & other.bits[0],
			b.bits[1] & other.bits[1],
			b.bits[2] & other.bits[2],
			b.bits[3] & other.bits[3],
		},
	}
}

// Or returns the bitwise OR of two masks.
func (b *bitMask) Or(other *bitMask) bitMask {
	return bitMask{
		bits: [4]uint64{
			b.bits[0] | other.bits[0],
			b.bits[1] | other.bits[1],
			b.bits[2] | other.bits[2],
			b.bits[3] | other.bits[3],
		},
	}
}

// Xor returns the bitwise XOR of two masks.
func (b *bitMask) Xor(other *bitMask) bitMask {
	return bitMask{
		bits: [4]uint64{
			b.bits[0] ^ other.bits[0],
			b.bits[1] ^ other.bits[1],
			b.bits[2] ^ other.bits[2],
			b.bits[3] ^ other.bits[3],
		},
	}
}

// TotalBitsSet returns how many bits are set in this mask.
func (b *bitMask) TotalBitsSet() int {
	return bits.OnesCount64(b.bits[0]) + bits.OnesCount64(b.bits[1]) + bits.OnesCount64(b.bits[2]) + bits.OnesCount64(b.bits[3])
}

// Equals returns whether two masks are equal.
func (b *bitMask) Equals(other *bitMask) bool {
	return b.bits == other.bits
}

func (b *bitMask) toTypes(reg *registry) []ID {
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
