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

	t2 := table{id: 2}
	t3 := table{id: 3}
	t4 := table{id: 4}
	t5 := table{id: 5}

	ids.Append(&t2)
	ids.Append(&t3)
	ids.Append(&t4)
	ids.Append(&t5)

	expectEqual(t, 4, len(ids.tables))
	expectEqual(t, 4, len(ids.indices))

	ids.Remove(3)

	expectEqual(t, 3, len(ids.tables))
	expectEqual(t, 3, len(ids.indices))

	expectSlicesEqual(t, []*table{&t2, &t5, &t4}, ids.tables)

	idx, ok := ids.indices[5]
	expectEqual(t, 1, idx)
	expectTrue(t, ok)

	_, ok = ids.indices[3]
	expectFalse(t, ok)

	ids.Clear()
	expectEqual(t, 0, len(ids.tables))
	expectEqual(t, 0, len(ids.indices))
}
