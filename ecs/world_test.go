package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWorld(t *testing.T) {
	w := NewWorld(1024)

	assert.Equal(t, 2, len(w.entities))
	assert.Equal(t, 1, len(w.storage.tables))
	assert.Equal(t, 1, len(w.storage.archetypes))
	assert.Equal(t, 1, len(w.storage.archetypes[0].tables))
}

func TestWorldNewEntity(t *testing.T) {
	w := NewWorld(8)

	assert.False(t, w.Alive(Entity{}))
	for i := range 10 {
		e := w.NewEntity()
		assert.EqualValues(t, e.id, i+2)
		assert.EqualValues(t, e.gen, 0)
		assert.True(t, w.Alive(e))
	}
	assert.Equal(t, 12, len(w.entities))

	idx := w.getEntityIndex(Entity{4, 0})
	assert.EqualValues(t, 0, idx.table)
	assert.EqualValues(t, 2, idx.row)
}

func TestWorldExchange(t *testing.T) {
	w := NewWorld(2)

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)

	e1 := w.NewEntity()
	e2 := w.NewEntity()
	e3 := w.NewEntity()

	w.exchange(e1, []ID{posID}, nil, nil, nil)
	w.exchange(e2, []ID{posID, velID}, nil, nil, nil)
	w.exchange(e3, []ID{posID, velID}, nil, nil, nil)

	assert.True(t, w.has(e1, posID))
	assert.False(t, w.has(e1, velID))

	assert.True(t, w.has(e2, posID))
	assert.True(t, w.has(e2, velID))

	pos := (*Position)(w.get(e1, posID))
	pos.X = 100

	pos = (*Position)(w.get(e1, posID))
	assert.Equal(t, pos.X, 100.0)

	w.exchange(e2, nil, []ID{posID}, nil, nil)
	assert.False(t, w.has(e2, posID))
	assert.True(t, w.has(e2, velID))
}

func TestWorldRemoveEntity(t *testing.T) {
	w := NewWorld(32)

	mapper := NewMap2[Position, Velocity](&w)

	entities := make([]Entity, 0, 100)
	for range 100 {
		e := mapper.NewEntity(&Position{}, &Velocity{})
		assert.True(t, w.Alive(e))
		entities = append(entities, e)
	}

	filter := NewFilter0(&w)
	query := filter.Query()
	cnt := 0
	for query.Next() {
		assert.True(t, w.Alive(query.Entity()))
		cnt++
	}
	assert.Equal(t, 100, cnt)

	for _, e := range entities {
		w.RemoveEntity(e)
		assert.False(t, w.Alive(e))
	}

	query = filter.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)

	e := w.NewEntity()
	w.RemoveEntity(e)
	assert.False(t, w.Alive(e))
}

func TestWorldRelations(t *testing.T) {
	w := NewWorld(16)

	parent1 := w.NewEntity()
	parent2 := w.NewEntity()

	mapper := NewMap3[Position, ChildOf, ChildOf2](&w)
	assert.True(t, w.storage.registry.IsRelation[1])
	assert.True(t, w.storage.registry.IsRelation[2])

	for range 10 {
		mapper.NewEntity(&Position{}, &ChildOf{}, &ChildOf2{}, Rel(1, parent1), Rel(2, parent2))
		mapper.NewEntity(&Position{}, &ChildOf{}, &ChildOf2{}, Rel(1, parent2), Rel(2, parent1))
	}

}
