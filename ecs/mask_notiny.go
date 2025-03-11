//go:build !tiny

package ecs

const maskTotalBits = mask256TotalBits

type bitMask = bitMask256

var newMask = newMask256
