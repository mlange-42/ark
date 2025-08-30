package ecs

import (
	"testing"
)

func TestArchetype(t *testing.T) {
	arch := archetype{}

	if arch.HasRelations() {
		t.Errorf("expected HasRelations() to be false, got true")
	}

	if len(arch.tables.tables) != 0 {
		t.Errorf("expected tables length to be 0, got %d", len(arch.tables.tables))
	}
}
