//go:build ark_tiny

package ecs

// maskTotalBits for 64 bit tiny ECS.
const maskTotalBits = mask64TotalBits

// bitMask for 64 bit tiny ECS.
type bitMask = bitMask64

// newMask constructor for 64 bit tiny ECS.
var newMask = newMask64
