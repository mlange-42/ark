package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnsafe(t *testing.T) {
	w := NewWorld(1024)
	u := w.Unsafe()

	assert.Equal(t, &w, u.world)
}

func TestUnsafeNewEntity(t *testing.T) {
	w := NewWorld(16)
	u := w.Unsafe()

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)

	e := u.NewEntity(posID, velID)

	assert.True(t, u.Has(e, posID))
	assert.True(t, u.Has(e, velID))
}

func TestUnsafeGet(t *testing.T) {
	w := NewWorld(16)
	u := w.Unsafe()

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)

	e := u.NewEntity(posID)

	assert.True(t, u.Has(e, posID))
	assert.False(t, u.Has(e, velID))

	assert.True(t, u.HasUnchecked(e, posID))
	assert.False(t, u.HasUnchecked(e, velID))

	pos := (*Position)(u.Get(e, posID))
	pos.X = 100

	pos2 := (*Position)(u.GetUnchecked(e, posID))

	assert.Equal(t, pos, pos2)
}
