package ecs

import (
	"testing"
)

func TestNewTable(t *testing.T) {
	w := NewWorld(1024)
	_ = ComponentID[Label](&w)
	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)

	arch := newArchetype(0, 0, &bitMask{}, []ID{posID, velID}, []tableID{0}, &w.storage.registry)
	table := newTable(0, &arch, 8, &w.storage.registry, make([]Entity, 2), []RelationID{})

	expectEqual(t, len(table.columns), 2)
	expectEqual(t, table.components[posID.id].index, 0)
	expectEqual(t, table.components[velID.id].index, 1)

	for i := 0; i < 9; i++ {
		table.Add(Entity{entityID(i + 2), 0})
	}
	expectEqual(t, table.len, 9)
	expectEqual(t, table.cap, 16)

	table2 := newTable(0, &arch, 8, &w.storage.registry, make([]Entity, 2), []RelationID{})
	table2.AddAllEntities(&table, uint32(table.Len()))
	expectEqual(t, table2.len, 9)
	expectEqual(t, table2.cap, 16)
}

func TestTableMatches(t *testing.T) {
	w := NewWorld(1024)
	_ = ComponentID[Label](&w)
	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)
	childID := ComponentID[ChildOf](&w)

	compMap := make([]int16, maskTotalBits)
	compMap[1] = 0
	compMap[2] = 1
	compMap[3] = 2

	arch := newArchetype(0, 0, &bitMask{}, []ID{posID, velID, childID}, []tableID{0}, &w.storage.registry)
	table := newTable(0, &arch, 8, &w.storage.registry,
		[]Entity{{}, {}, {2, 0}},
		[]RelationID{{component: childID, target: Entity{2, 0}}},
	)

	expectTrue(t, table.MatchesExact([]RelationID{{component: childID, target: Entity{2, 0}}}))
	expectFalse(t, table.MatchesExact([]RelationID{{component: childID, target: Entity{3, 0}}}))

	expectPanic(t, func() {
		table.MatchesExact([]RelationID{})
	})
	expectPanic(t, func() {
		table.MatchesExact([]RelationID{{component: posID, target: Entity{2, 0}}})
	})
}

func TestTableReset(t *testing.T) {
	w := NewWorld(1024)
	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)
	labelID := ComponentID[Label](&w)

	arch := newArchetype(0, 0, &bitMask{}, []ID{posID, velID, labelID}, []tableID{0}, &w.storage.registry)
	table := newTable(0, &arch, 8, &w.storage.registry, make([]Entity, 3), []RelationID{})

	table.Reset()

	for i := 0; i < 75; i++ {
		table.Add(Entity{entityID(i + 2), 0})
	}
	expectEqual(t, table.len, 75)
	table.Reset()
	expectEqual(t, table.len, 0)
}
