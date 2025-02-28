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
