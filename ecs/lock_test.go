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
