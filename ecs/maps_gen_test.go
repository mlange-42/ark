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
	}

	for _, e := range entities {
		mapper.Remove(e)
	}
}
