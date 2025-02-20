package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWorld(t *testing.T) {
	w := NewWorld(1024)

	assert.Equal(t, 1, len(w.entities))
	assert.Equal(t, 1, len(w.storage.tables))
	assert.Equal(t, 1, len(w.storage.archetypes))
	assert.Equal(t, 1, len(w.storage.archetypes[0].tables))
}

func TestWorldNewEntity(t *testing.T) {
	w := NewWorld(1024)

	assert.False(t, w.Alive(Entity{}))
	for i := range 10 {
		e := w.NewEntity()
		assert.EqualValues(t, e.id, i+1)
		assert.EqualValues(t, e.gen, 0)
		assert.True(t, w.Alive(e))
	}
	assert.Equal(t, 11, len(w.entities))

	idx := w.getEntityIndex(Entity{3, 0})
	assert.EqualValues(t, 0, idx.table)
	assert.EqualValues(t, 2, idx.row)
}

func TestWorldExchange(t *testing.T) {
	w := NewWorld(1024)

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)

	e1 := w.NewEntity()
	e2 := w.NewEntity()
	e3 := w.NewEntity()

	w.exchange(e1, []ID{posID}, nil, nil)
	w.exchange(e2, []ID{posID, velID}, nil, nil)
	w.exchange(e3, []ID{posID, velID}, nil, nil)

	assert.True(t, w.has(e1, posID))
	assert.False(t, w.has(e1, velID))

	assert.True(t, w.has(e2, posID))
	assert.True(t, w.has(e2, velID))

	pos := (*Position)(w.get(e1, posID))
	pos.X = 100

	pos = (*Position)(w.get(e1, posID))
	assert.Equal(t, pos.X, 100.0)

	w.exchange(e2, nil, []ID{posID}, nil)
	assert.False(t, w.has(e2, posID))
	assert.True(t, w.has(e2, velID))
}
