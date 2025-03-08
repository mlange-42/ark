package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTable(t *testing.T) {
	w := NewWorld(1024)
	_ = ComponentID[Label](&w)
	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)

	arch := newArchetype(0, 0, &bitMask{}, []ID{posID, velID}, []tableID{0}, &w.storage.registry)
	table := newTable(0, &arch, 8, &w.storage.registry, make([]Entity, 2), []RelationID{})

	assert.Equal(t, 2, len(table.columns))
	assert.EqualValues(t, 0, table.components[posID.id])
	assert.EqualValues(t, 1, table.components[velID.id])

	for i := range 9 {
		table.Add(Entity{entityID(i + 2), 0})
	}
	assert.EqualValues(t, 9, table.len)
	assert.EqualValues(t, 16, table.cap)

	table2 := newTable(0, &arch, 8, &w.storage.registry, make([]Entity, 2), []RelationID{})
	table2.AddAllEntities(&table, uint32(table.Len()))
	assert.EqualValues(t, 9, table2.len)
	assert.EqualValues(t, 16, table2.cap)
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

	assert.True(t, table.MatchesExact([]RelationID{{component: childID, target: Entity{2, 0}}}))
	assert.False(t, table.MatchesExact([]RelationID{{component: childID, target: Entity{3, 0}}}))

	assert.Panics(t, func() {
		table.MatchesExact([]RelationID{})
	})
	assert.Panics(t, func() {
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

	for i := range 75 {
		table.Add(Entity{entityID(i + 2), 0})
	}
	assert.EqualValues(t, 75, table.len)
	table.Reset()
	assert.EqualValues(t, 0, table.len)
}
