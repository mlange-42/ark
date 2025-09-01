package ecs

import (
	"testing"
)

func TestStorage(t *testing.T) {
	s := newStorage(16)
	expectEqual(t, 1, len(s.archetypes))
	expectEqual(t, 1, len(s.tables))

	s.AddComponent(0)
	s.AddComponent(1)

	expectPanicsWithValue(t, "components can only be added to a storage sequentially", func() { s.AddComponent(3) })
}
