package ecs

import (
	"fmt"
	"math"
	"math/bits"
	"sync/atomic"
)

// Manages locks by mask bits.
//
// The number of simultaneous locks at a given time is limited to 64.
//
// The concurrency-safe methods access the locks atomically.
// The other methods use plain memory access, which is faster.
// They must not run concurrently with any other method, like before with the mutex.
type lock struct {
	locks uint64 // The actual locks, one bit per lock.
}

func newLock() lock {
	return lock{}
}

// Lock the world and get the Lock bit for later unlocking.
// This is not concurrency-safe.
func (m *lock) Lock() uint8 {
	l := freeLockBit(m.locks)
	m.locks |= uint64(1) << l
	return l
}

// Unlock unlocks the given lock bit.
// This is not concurrency-safe.
func (m *lock) Unlock(l uint8) {
	bit := uint64(1) << l
	if m.locks&bit == 0 {
		panic("unbalanced unlock. Did you close a query that was already iterated?")
	}
	m.locks &^= bit
}

// LockSafe locks the world and get the Lock bit for later unlocking.
// This is concurrency-safe.
func (m *lock) LockSafe() uint8 {
	for {
		locks := atomic.LoadUint64(&m.locks)
		l := freeLockBit(locks)
		if atomic.CompareAndSwapUint64(&m.locks, locks, locks|uint64(1)<<l) {
			return l
		}
	}
}

// UnlockSafe unlocks the given lock bit.
// This is concurrency-safe.
func (m *lock) UnlockSafe(l uint8) {
	bit := uint64(1) << l
	if atomic.AndUint64(&m.locks, ^bit)&bit == 0 {
		panic("unbalanced unlock. Did you close a query that was already iterated?")
	}
}

// IsLocked returns whether the world is locked by any queries.
func (m *lock) IsLocked() bool {
	return m.locks != 0
}

// Reset the locks.
func (m *lock) Reset() {
	m.locks = 0
}

// freeLockBit returns the lowest bit that is not set in the given locks.
func freeLockBit(locks uint64) uint8 {
	if locks == math.MaxUint64 {
		panic(fmt.Sprintf("run out of the maximum of %d bits. "+
			"This is likely caused by unclosed queries that lock the world. "+
			"Make sure that all queries finish their iteration or are closed manually", mask64TotalBits))
	}
	return uint8(bits.TrailingZeros64(^locks))
}
