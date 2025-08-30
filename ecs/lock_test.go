package ecs

import (
	"testing"
)

func TestLock(t *testing.T) {
	locks := newLock()

	if locks.IsLocked() {
		t.Errorf("expected IsLocked() to be false initially")
	}

	l1 := locks.Lock()
	if !locks.IsLocked() {
		t.Errorf("expected IsLocked() to be true after first lock")
	}
	if int(l1) != 0 {
		t.Errorf("expected first lock ID to be 0, got %d", l1)
	}

	l2 := locks.Lock()
	if !locks.IsLocked() {
		t.Errorf("expected IsLocked() to be true after second lock")
	}
	if int(l2) != 1 {
		t.Errorf("expected second lock ID to be 1, got %d", l2)
	}

	locks.Unlock(l1)
	if !locks.IsLocked() {
		t.Errorf("expected IsLocked() to remain true after unlocking l1")
	}

	expectPanicWithValue(t,
		"unbalanced unlock. Did you close a query that was already iterated?",
		func() { locks.Unlock(l1) })

	locks.Unlock(l2)
	if locks.IsLocked() {
		t.Errorf("expected IsLocked() to be false after unlocking all")
	}
}
