package ecs

import (
	"testing"
)

func TestArchetype(t *testing.T) {
	arch := archetype{}
	expectFalse(t, arch.HasRelations())
	expectEqual(t, 0, len(arch.tables.tables))
}

func TestTableIDs(t *testing.T) {
	ids := newTableIDs()

	ids.Append(2)
	ids.Append(3)
	ids.Append(4)
	ids.Append(5)

	expectEqual(t, 4, len(ids.tables))
	expectEqual(t, 4, len(ids.indices))

	ids.Remove(3)

	expectEqual(t, 3, len(ids.tables))
	expectEqual(t, 3, len(ids.indices))

	expectSlicesEqual(t, []tableID{2, 5, 4}, ids.tables)

	idx, ok := ids.indices[5]
	expectEqual(t, 1, idx)
	expectTrue(t, ok)

	idx, ok = ids.indices[3]
	expectFalse(t, ok)

	ids.Clear()
	expectEqual(t, 0, len(ids.tables))
	expectEqual(t, 0, len(ids.indices))
}
