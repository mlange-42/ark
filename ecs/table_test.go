package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTable(t *testing.T) {
	w := NewWorld()
	_ = ComponentID[Label](&w)
	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)

	table := newTable(8, &w.registry, posID, velID)

	assert.Equal(t, 2, len(table.columns))
	assert.Equal(t, 2, len(table.relations))
	assert.EqualValues(t, 0, table.components[posID.id])
	assert.EqualValues(t, 1, table.components[velID.id])
}
