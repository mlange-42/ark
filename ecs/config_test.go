package ecs

import (
	"testing"
)

func TestConfig(t *testing.T) {
	c := newConfig()
	expectEqual(t, 1024, c.initialCapacity)
	expectEqual(t, 128, c.initialCapacityRelations)

	c = newConfig(16)
	expectEqual(t, 16, c.initialCapacity)
	expectEqual(t, 16, c.initialCapacityRelations)

	c = newConfig(16, 8)
	expectEqual(t, 16, c.initialCapacity)
	expectEqual(t, 8, c.initialCapacityRelations)

	expectPanicsWithValue(t, "only positive values for the World's initialCapacity are allowed",
		func() { _ = newConfig(0) })
	expectPanicsWithValue(t, "only positive values for the World's initialCapacity are allowed",
		func() { _ = newConfig(1024, 0) })
	expectPanicsWithValue(t, "can only use a maximum of two values for the World's initialCapacity",
		func() { _ = newConfig(1024, 128, 32) })
}
