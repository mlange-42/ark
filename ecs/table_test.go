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
	table := newTable(0, &arch, 8, &w.storage.registry, make([]Entity, 2), []relationID{})

	expectEqual(t, 2, len(table.columns))
	expectEqual(t, 0, table.components[posID.id].index)
	expectEqual(t, 1, table.components[velID.id].index)

	for i := range 9 {
		table.Add(&w.storage, Entity{entityID(i + 2), 0})
	}
	expectEqual(t, 9, table.len)
	expectEqual(t, 16, table.cap)

	table2 := newTable(0, &arch, 8, &w.storage.registry, make([]Entity, 2), []relationID{})
	table2.AddAllEntities(&w.storage, &table, uint32(table.Len()))
	expectEqual(t, 9, table2.len)
	expectEqual(t, 16, table2.cap)
}

func TestTableMatches(t *testing.T) {
	w := NewWorld(1024)
	_ = ComponentID[Label](&w)
	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)
	childID := ComponentID[ChildOf](&w)
	child2ID := ComponentID[ChildOf2](&w)

	compMap := make([]int16, maskTotalBits)
	compMap[1] = 0
	compMap[2] = 1
	compMap[3] = 2

	arch := newArchetype(0, 0, &bitMask{}, []ID{posID, velID, childID}, []tableID{0}, &w.storage.registry)
	table := newTable(0, &arch, 8, &w.storage.registry,
		[]Entity{{}, {}, {2, 0}},
		[]relationID{{component: childID, target: Entity{2, 0}}},
	)

	expectTrue(t, table.MatchesExact([]relationID{{component: childID, target: Entity{2, 0}}}))
	expectTrue(t, table.MatchesExact([]relationID{{component: childID, target: Entity{2, 0}}, {component: child2ID, target: Entity{2, 0}}}))
	expectFalse(t, table.MatchesExact([]relationID{{component: childID, target: Entity{3, 0}}}))

	expectPanics(t, func() {
		table.MatchesExact([]relationID{})
	})
	expectPanics(t, func() {
		table.MatchesExact([]relationID{{component: posID, target: Entity{2, 0}}})
	})
}

func TestTableReset(t *testing.T) {
	w := NewWorld(1024)
	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)
	labelID := ComponentID[Label](&w)

	arch := newArchetype(0, 0, &bitMask{}, []ID{posID, velID, labelID}, []tableID{0}, &w.storage.registry)
	table := newTable(0, &arch, 8, &w.storage.registry, make([]Entity, 3), []relationID{})

	table.Reset()

	for i := range 75 {
		table.Add(&w.storage, Entity{entityID(i + 2), 0})
	}
	expectEqual(t, 75, table.len)
	table.Reset()
	expectEqual(t, 0, table.len)
}
