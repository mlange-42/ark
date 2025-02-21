package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntity(t *testing.T) {
	e := newEntity(100)
	assert.EqualValues(t, 100, e.id)
	assert.EqualValues(t, 0, e.gen)
}

func TestEntityIndex(t *testing.T) {
	index := entityIndex{}
	assert.EqualValues(t, 0, index.table)
	assert.EqualValues(t, 0, index.row)
}

func TestReservedEntities(t *testing.T) {
	w := NewWorld(1024)

	zero := Entity{}
	wildcard := Entity{1, 0}

	assert.False(t, w.Alive(zero))
	assert.False(t, w.Alive(wildcard))
	assert.False(t, w.Alive(Wildcard()))

	assert.True(t, zero.IsZero())
	assert.False(t, wildcard.IsZero())
	assert.True(t, wildcard.isWildcard())
}
