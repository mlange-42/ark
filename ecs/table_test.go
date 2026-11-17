package ecs

import (
	"testing"
)

func TestNewTable(t *testing.T) {
	w := NewWorld(1024)
	_ = w.ComponentID[Label]()
	posID := w.ComponentID[Position]()
	velID := w.ComponentID[Velocity]()

	arch, data := newArchetype(0, 0, &bitMask{}, []ID{posID, velID}, []tableID{0}, &w.storage.registry)
	arch.archetypeData = &data
	table := newTable(0, &arch, 8, &w.storage.registry, make([]Entity, 2), []relationID{})

	expectEqual(t, 2, len(table.columns))
	expectEqual(t, 0, table.components[posID.id].index)
	expectEqual(t, 1, table.components[velID.id].index)

	expectTrue(t, table.components[posID.id].isTrivial)
	expectTrue(t, table.components[velID.id].isTrivial)

	for i := range 9 {
		table.Add(Entity{entityID(i + 2), 0})
	}
	expectEqual(t, 9, table.len)
	expectEqual(t, 16, table.cap)

	table2 := newTable(0, &arch, 8, &w.storage.registry, make([]Entity, 2), []relationID{})
	table2.AddAllEntities(&table, uint32(table.Len()))
	expectEqual(t, 9, table2.len)
	expectEqual(t, 16, table2.cap)
}

func TestTableMatches(t *testing.T) {
	w := NewWorld(1024)
	_ = w.ComponentID[Label]()
	posID := w.ComponentID[Position]()
	velID := w.ComponentID[Velocity]()
	childID := w.ComponentID[ChildOf]()
	child2ID := w.ComponentID[ChildOf2]()

	compMap := make([]int16, maskTotalBits)
	compMap[1] = 0
	compMap[2] = 1
	compMap[3] = 2

	arch, data := newArchetype(0, 0, &bitMask{}, []ID{posID, velID, childID}, []tableID{0}, &w.storage.registry)
	arch.archetypeData = &data
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
	posID := w.ComponentID[Position]()
	velID := w.ComponentID[Velocity]()
	labelID := w.ComponentID[Label]()

	arch, data := newArchetype(0, 0, &bitMask{}, []ID{posID, velID, labelID}, []tableID{0}, &w.storage.registry)
	arch.archetypeData = &data
	table := newTable(0, &arch, 8, &w.storage.registry, make([]Entity, 3), []relationID{})

	table.Reset()

	for i := range 75 {
		table.Add(Entity{entityID(i + 2), 0})
	}
	expectEqual(t, 75, table.len)
	table.Reset()
	expectEqual(t, 0, table.len)
}
