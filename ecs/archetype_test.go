package ecs

import (
	"testing"
)

func TestArchetype(t *testing.T) {
	arch := archetype{}
	expectFalse(t, arch.HasRelations())
	expectEqual(t, 0, len(arch.tables.tables))
}
