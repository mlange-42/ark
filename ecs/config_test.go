package ecs

import (
	"testing"
)

func TestConfig(t *testing.T) {
	c := newConfig()
	if c.initialCapacity != 1024 {
		t.Errorf("expected initialCapacity to be 1024, got %d", c.initialCapacity)
	}
	if c.initialCapacityRelations != 128 {
		t.Errorf("expected initialCapacityRelations to be 128, got %d", c.initialCapacityRelations)
	}

	c = newConfig(16)
	if c.initialCapacity != 16 {
		t.Errorf("expected initialCapacity to be 16, got %d", c.initialCapacity)
	}
	if c.initialCapacityRelations != 16 {
		t.Errorf("expected initialCapacityRelations to be 16, got %d", c.initialCapacityRelations)
	}

	c = newConfig(16, 8)
	if c.initialCapacity != 16 {
		t.Errorf("expected initialCapacity to be 16, got %d", c.initialCapacity)
	}
	if c.initialCapacityRelations != 8 {
		t.Errorf("expected initialCapacityRelations to be 8, got %d", c.initialCapacityRelations)
	}

	expectPanicWithValue(t, "only positive values for the World's initialCapacity are allowed", func() {
		_ = newConfig(0)
	})
	expectPanicWithValue(t, "only positive values for the World's initialCapacity are allowed", func() {
		_ = newConfig(1024, 0)
	})
	expectPanicWithValue(t, "can only use a maximum of two values for the World's initialCapacity", func() {
		_ = newConfig(1024, 128, 32)
	})
}
