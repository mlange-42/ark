package ecs

// Manages locks by mask bits.
//
// The number of simultaneous locks at a given time is limited to [MaskTotalBits].
type lock struct {
	locks   bitMask // The actual locks.
	bitPool bitPool // The bit pool for getting and recycling bits.
}

// Lock the world and get the Lock bit for later unlocking.
func (m *lock) Lock() uint8 {
	lock := m.bitPool.Get()
	m.locks.Set(id8(lock), true)
	return lock
}

// Unlock unlocks the given lock bit.
func (m *lock) Unlock(l uint8) {
	if !m.locks.Get(id8(l)) {
		panic("unbalanced unlock. Did you close a query that was already iterated?")
	}
	m.locks.Set(id8(l), false)
	m.bitPool.Recycle(l)
}

// IsLocked returns whether the world is locked by any queries.
func (m *lock) IsLocked() bool {
	return !m.locks.IsZero()
}

// Reset the locks and the pool.
func (m *lock) Reset() {
	m.locks = bitMask{}
	m.bitPool.Reset()
}
