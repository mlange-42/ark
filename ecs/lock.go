package ecs

import (
	"sync"
)

// Manages locks by mask bits.
//
// The number of simultaneous locks at a given time is limited to 64.
type lock struct {
	bitPool bitPool   // The bit pool for getting and recycling bits.
	locks   bitMask64 // The actual locks.
	mu      sync.Mutex
}

func newLock() lock {
	return lock{
		bitPool: newBitPool(),
	}
}

// Lock the world and get the Lock bit for later unlocking.
// This is not concurrency-safe.
func (m *lock) Lock() uint8 {
	lock := m.bitPool.Get()
	m.locks.Set(lock)
	return lock
}

// Unlock unlocks the given lock bit.
// This is not concurrency-safe.
func (m *lock) Unlock(l uint8) {
	if !m.locks.Get(l) {
		panic("unbalanced unlock. Did you close a query that was already iterated?")
	}
	m.locks.Clear(l)
	m.bitPool.Recycle(l)
}

// LockSafe locks the world and get the Lock bit for later unlocking.
// This is concurrency-safe.
func (m *lock) LockSafe() uint8 {
	m.mu.Lock()
	lock := m.bitPool.Get()
	m.locks.Set(lock)
	m.mu.Unlock()
	return lock
}

// UnlockSafe unlocks the given lock bit.
// This is concurrency-safe.
func (m *lock) UnlockSafe(l uint8) {
	m.mu.Lock()
	if !m.locks.Get(l) {
		m.mu.Unlock()
		panic("unbalanced unlock. Did you close a query that was already iterated?")
	}
	m.locks.Clear(l)
	m.bitPool.Recycle(l)
	m.mu.Unlock()
}

// IsLocked returns whether the world is locked by any queries.
func (m *lock) IsLocked() bool {
	return !m.locks.IsZero()
}

// Reset the locks and the pool.
func (m *lock) Reset() {
	m.locks = bitMask64{}
	m.bitPool.Reset()
}
