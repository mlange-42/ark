package ecs

import (
	"testing"
)

func TestStorage(t *testing.T) {
	s := newStorage(16)
	if len(s.archetypes) != 1 {
		t.Errorf("expected 1 archetype, got %d", len(s.archetypes))
	}
	if len(s.tables) != 1 {
		t.Errorf("expected 1 table, got %d", len(s.tables))
	}
	s.AddComponent(0)
	s.AddComponent(1)
	expectPanicWithValue(t, "components can only be added to a storage sequentially", func() { s.AddComponent(3) })
}
