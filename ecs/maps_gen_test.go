package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap2(t *testing.T) {
	w := NewWorld(4)

	mapper := NewMap2[Position, Velocity](&w)

	entities := []Entity{}
	for i := range 12 {
		v := float64(i)
		e := w.NewEntity()
		mapper.Add(e, &Position{v + 1, v + 2}, &Velocity{v + 3, v + 4})
		entities = append(entities, e)
	}

	for i, e := range entities {
		v := float64(i)
		pos, vel := mapper.Get(e)
		assert.Equal(t, Position{v + 1, v + 2}, *pos)
		assert.Equal(t, Velocity{v + 3, v + 4}, *vel)
		assert.True(t, mapper.HasAll(e))
	}

	for _, e := range entities {
		mapper.Remove(e)
	}

	e := mapper.NewEntity(&Position{101, 102}, &Velocity{103, 104})
	pos, vel := mapper.Get(e)
	assert.Equal(t, Position{101, 102}, *pos)
	assert.Equal(t, Velocity{103, 104}, *vel)
}

func TestMap2NewBatch(t *testing.T) {
	w := NewWorld(16)

	mapper := NewMap2[Position, Velocity](&w)

	for range 12 {
		_ = mapper.NewEntity(&Position{1, 2}, &Velocity{3, 4})
	}
	mapper.NewBatch(24, &Position{5, 6}, &Velocity{7, 8})

	filter := NewFilter2[Position, Velocity](&w)
	query := filter.Query()
	cnt := 0
	var lastEntity Entity
	for query.Next() {
		pos, vel := query.Get()
		if cnt < 12 {
			assert.Equal(t, Position{1, 2}, *pos)
			assert.Equal(t, Velocity{3, 4}, *vel)
		} else {
			assert.Equal(t, Position{5, 6}, *pos)
			assert.Equal(t, Velocity{7, 8}, *vel)
		}
		lastEntity = query.Entity()
		cnt++
	}
	assert.True(t, mapper.HasAll(lastEntity))
	assert.Equal(t, 36, cnt)
}

func TestMap2NewBatchFn(t *testing.T) {
	w := NewWorld(16)

	mapper := NewMap2[Position, Velocity](&w)

	for range 12 {
		_ = mapper.NewEntity(&Position{1, 2}, &Velocity{3, 4})
	}
	mapper.NewBatchFn(24, func(entity Entity, pos *Position, vel *Velocity) {
		pos.X = 5
		pos.Y = 6
		vel.X = 7
		vel.Y = 8
	})

	filter := NewFilter2[Position, Velocity](&w)
	query := filter.Query()
	cnt := 0
	var lastEntity Entity
	for query.Next() {
		pos, vel := query.Get()
		if cnt < 12 {
			assert.Equal(t, Position{1, 2}, *pos)
			assert.Equal(t, Velocity{3, 4}, *vel)
		} else {
			assert.Equal(t, Position{5, 6}, *pos)
			assert.Equal(t, Velocity{7, 8}, *vel)
		}
		lastEntity = query.Entity()
		cnt++
	}
	assert.True(t, mapper.HasAll(lastEntity))
	assert.Equal(t, 36, cnt)
}
