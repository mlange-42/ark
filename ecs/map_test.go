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

	posMap.Add(e1, &Position{})
	velMap.Add(e1, &Velocity{})

	assert.True(t, posMap.Has(e1))
	assert.True(t, velMap.Has(e1))

	pos := posMap.Get(e1)
	pos.X = 100

	pos = posMap.Get(e1)
	assert.Equal(t, 100.0, pos.X)

	posMap.Remove(e1)
	assert.False(t, posMap.Has(e1))
}

func TestMapNewEntity(t *testing.T) {
	w := NewWorld(1024)

	posMap := NewMap[Position](&w)

	e := posMap.NewEntity(&Position{X: 1, Y: 2})

	pos := posMap.Get(e)
	assert.Equal(t, Position{X: 1, Y: 2}, *pos)
}

func TestMapRelation(t *testing.T) {
	w := NewWorld(1024)

	childMap := NewMap[ChildOf](&w)

	parent1 := w.NewEntity()
	parent2 := w.NewEntity()

	e := w.NewEntity()

	childMap.Add(e, &ChildOf{}, parent1)
	assert.Equal(t, childMap.GetRelation(e), parent1)
	assert.Equal(t, childMap.GetRelationUnchecked(e), parent1)

	childMap.SetRelation(e, parent2)
	assert.Equal(t, childMap.GetRelation(e), parent2)
	assert.Equal(t, childMap.GetRelationUnchecked(e), parent2)
}
