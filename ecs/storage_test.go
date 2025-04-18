package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	s := newStorage(16)
	assert.Equal(t, 1, len(s.archetypes))
	assert.Equal(t, 1, len(s.tables))

	s.AddComponent(0)
	s.AddComponent(1)

	assert.PanicsWithValue(t, "components can only be added to a storage sequentially", func() { s.AddComponent(3) })
}
