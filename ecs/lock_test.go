package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLock(t *testing.T) {
	locks := lock{}

	assert.False(t, locks.IsLocked())

	l1 := locks.Lock()
	assert.True(t, locks.IsLocked())
	assert.Equal(t, 0, int(l1))

	l2 := locks.Lock()
	assert.True(t, locks.IsLocked())
	assert.Equal(t, 1, int(l2))

	locks.Unlock(l1)
	assert.True(t, locks.IsLocked())

	assert.PanicsWithValue(t, "unbalanced unlock. Did you close a query that was already iterated?",
		func() { locks.Unlock(l1) })

	locks.Unlock(l2)
	assert.False(t, locks.IsLocked())
}
