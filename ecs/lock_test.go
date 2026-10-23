package ecs

import (
	"testing"
)

func TestLock(t *testing.T) {
	locks := newLock()

	expectFalse(t, locks.IsLocked())

	l1 := locks.Lock()
	expectTrue(t, locks.IsLocked())
	expectEqual(t, 0, int(l1))

	l2 := locks.Lock()
	expectTrue(t, locks.IsLocked())
	expectEqual(t, 1, int(l2))

	locks.Unlock(l1)
	expectTrue(t, locks.IsLocked())

	expectPanicsWithValue(t, "unbalanced unlock. Did you close a query that was already iterated?",
		func() { locks.Unlock(l1) })

	locks.Unlock(l2)
	expectFalse(t, locks.IsLocked())
}

func TestLockSafe(t *testing.T) {
	locks := newLock()

	expectFalse(t, locks.IsLocked())

	l1 := locks.LockSafe()
	expectTrue(t, locks.IsLocked())
	expectEqual(t, 0, int(l1))

	l2 := locks.LockSafe()
	expectTrue(t, locks.IsLocked())
	expectEqual(t, 1, int(l2))

	locks.UnlockSafe(l1)
	expectTrue(t, locks.IsLocked())

	expectPanicsWithValue(t, "unbalanced unlock. Did you close a query that was already iterated?",
		func() { locks.UnlockSafe(l1) })

	locks.UnlockSafe(l2)
	expectFalse(t, locks.IsLocked())
}

func TestLockMaxLocks(t *testing.T) {
	locks := newLock()

	for i := range mask64TotalBits {
		if i%2 == 0 {
			expectEqual(t, i, int(locks.Lock()))
		} else {
			expectEqual(t, i, int(locks.LockSafe()))
		}
	}
	expectPanics(t, func() { locks.Lock() })
	expectPanics(t, func() { locks.LockSafe() })

	// Freed bits are re-used, lowest first.
	locks.Unlock(10)
	locks.UnlockSafe(3)
	expectEqual(t, 3, int(locks.LockSafe()))
	expectEqual(t, 10, int(locks.Lock()))
	expectPanics(t, func() { locks.Lock() })

	locks.Reset()
	expectFalse(t, locks.IsLocked())
	for i := range mask64TotalBits {
		expectEqual(t, i, int(locks.LockSafe()))
	}
}

func BenchmarkLockSafe(b *testing.B) {
	locks := newLock()

	for b.Loop() {
		l := locks.LockSafe()
		locks.UnlockSafe(l)
	}
}
