package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	w := NewWorld(1024)

	posMap := NewMap[Position](&w)
	velMap := NewMap[Velocity](&w)

	e1 := w.NewEntity()

	posMap.Add(e1)
	velMap.Add(e1)

	assert.True(t, posMap.Has(e1))
	assert.True(t, velMap.Has(e1))

	pos := posMap.Get(e1)
	pos.X = 100

	pos = posMap.Get(e1)
	assert.Equal(t, 100.0, pos.X)

	posMap.Remove(e1)
	assert.False(t, posMap.Has(e1))
}
