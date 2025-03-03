package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	c := newConfig()
	assert.Equal(t, 1024, c.initialCapacity)
	assert.Equal(t, 128, c.initialCapacityRelations)

	c = newConfig(16)
	assert.Equal(t, 16, c.initialCapacity)
	assert.Equal(t, 16, c.initialCapacityRelations)

	c = newConfig(16, 8)
	assert.Equal(t, 16, c.initialCapacity)
	assert.Equal(t, 8, c.initialCapacityRelations)

	assert.PanicsWithValue(t, "only positive values for the World's initialCapacity are allowed",
		func() { _ = newConfig(0) })
	assert.PanicsWithValue(t, "only positive values for the World's initialCapacity are allowed",
		func() { _ = newConfig(1024, 0) })
	assert.PanicsWithValue(t, "can only use a maximum of two values for the World's initialCapacity",
		func() { _ = newConfig(1024, 128, 32) })
}
