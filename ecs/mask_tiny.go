//go:build tiny

package ecs

const maskTotalBits = mask64TotalBits

type bitMask = bitMask64

var newMask = newMask64
