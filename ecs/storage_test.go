package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	s := newStorage()
	assert.Equal(t, 1, len(s.archetypes))
	assert.Equal(t, 1, len(s.tables))

	s.AddComponent(0)
	s.AddComponent(1)

	assert.Panics(t, func() { s.AddComponent(3) })
}
