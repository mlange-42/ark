package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTable(t *testing.T) {
	w := NewWorld(1024)
	_ = ComponentID[Label](&w)
	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)

	compMap := make([]int16, MaskTotalBits)
	compMap[1] = 0
	compMap[2] = 1
	table := newTable(0, 0, 8, &w.storage.registry, []ID{posID, velID}, compMap, make([]bool, 2), make([]Entity, 2), []relationID{})

	assert.Equal(t, 2, len(table.columns))
	assert.EqualValues(t, 0, table.components[posID.id])
	assert.EqualValues(t, 1, table.components[velID.id])
}
