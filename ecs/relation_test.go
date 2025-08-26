package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRel(t *testing.T) {
	r := RelIdx(1, Entity{5, 0})
	assert.Equal(t, relationIndex{1, Entity{5, 0}}, r)
}

func TestRelID(t *testing.T) {
	r := RelID(ID{10}, Entity{5, 0})
	assert.Equal(t, RelationID{ID{10}, Entity{5, 0}}, r)
}

func TestToRelations(t *testing.T) {
	w := NewWorld()

	childID := ComponentID[ChildOf](&w)
	child2ID := ComponentID[ChildOf2](&w)
	posID := ComponentID[Position](&w)

	inRelations := relations{RelIdx(1, Entity{2, 0}), RelIdx(2, Entity{3, 0})}
	var out []RelationID
	out = inRelations.toRelations(&w, []ID{posID, childID, child2ID}, out, 0)

	assert.Equal(t, []RelationID{
		{component: childID, target: Entity{2, 0}},
		{component: child2ID, target: Entity{3, 0}},
	}, out)

	assert.Panics(t, func() {
		_ = inRelations.toRelations(&w, []ID{childID, child2ID, posID}, out, 0)
	})

	inRelations = relations{RelID(childID, Entity{2, 0}), RelID(child2ID, Entity{3, 0})}
	out = inRelations.toRelations(&w, []ID{posID, childID, child2ID}, out, 0)

	assert.Equal(t, []RelationID{
		{component: childID, target: Entity{2, 0}},
		{component: child2ID, target: Entity{3, 0}},
	}, out)

	inRelations = relations{Rel[ChildOf](Entity{2, 0}), Rel[ChildOf2](Entity{3, 0})}
	out = inRelations.toRelations(&w, []ID{posID, childID, child2ID}, out, 0)

	assert.Equal(t, []RelationID{
		{component: childID, target: Entity{2, 0}},
		{component: child2ID, target: Entity{3, 0}},
	}, out)

	inRelations = relations{Rel[ChildOf](Entity{2, 0})}
	out = inRelations.toRelations(&w, []ID{posID, childID, child2ID}, out, 0)
	assert.Equal(t, []RelationID{
		{component: childID, target: Entity{2, 0}},
	}, out)

	inRelations = relations{Rel[ChildOf2](Entity{3, 0})}
	out = inRelations.toRelations(&w, []ID{posID, childID, child2ID}, out, uint8(len(out)))

	assert.Equal(t, []RelationID{
		{component: childID, target: Entity{2, 0}},
		{component: child2ID, target: Entity{3, 0}},
	}, out)
}

func TestRemoveRelationTarget(t *testing.T) {
	world := NewWorld()

	e1 := world.NewEntity()
	e2 := world.NewEntity()

	childMap := NewMap[ChildOf](&world)
	child2Map := NewMap[ChildOf2](&world)

	gen2Map := NewMap3[Position, ChildOf, ChildOf2](&world)
	e3 := gen2Map.NewEntity(
		&Position{},
		&ChildOf{}, &ChildOf2{},
		RelIdx(1, e1),
		RelIdx(2, e2),
	)
	_ = e3

	assert.Equal(t, e1, childMap.GetRelation(e3))
	assert.Equal(t, e2, child2Map.GetRelation(e3))

	world.RemoveEntity(e1)
	assert.Equal(t, Entity{}, childMap.GetRelation(e3))
	assert.Equal(t, e2, child2Map.GetRelation(e3))

	world.RemoveEntity(e2)
	assert.Equal(t, Entity{}, childMap.GetRelation(e3))
	assert.Equal(t, Entity{}, child2Map.GetRelation(e3))
}

func TestMultiRelation_RemoveAndReAdd_NoDuplicates(t *testing.T) {
	world := NewWorld()

	f := func(t *testing.T, world *World) {
		e1 := world.NewEntity()
		e2 := world.NewEntity()

		gen := NewMap3[Position, ChildOf, ChildOf2](world)
		e3 := gen.NewEntity(
			&Position{},
			&ChildOf{}, &ChildOf2{},
			RelIdx(1, e1), RelIdx(2, e2),
		)

		filter := NewFilter1[Position](world)

		check := func(prefix string, rel Relation, wanted ...Entity) {
			var entities = make(map[Entity]bool)
			for query := filter.Query(rel); query.Next(); {
				entity := query.Entity()
				require.False(t, entities[entity], "%s: entity %v returned multiple times", prefix, entity)
				entities[entity] = true
			}
			for _, w := range wanted {
				require.True(t, entities[w], "%s: entity %v not returned", prefix, w)
				delete(entities, w)
			}
			require.Empty(t, entities, "%s: unexpected entities returned: %v", prefix, entities)
		}

		check("initial, query by ChildOf", Rel[ChildOf](e1), e3)
		check("initial, query by ChildOf2", Rel[ChildOf2](e2), e3)

		world.RemoveEntity(e1)

		check("after removing e1, query by ChildOf2", Rel[ChildOf2](e2), e3)

		e4 := world.NewEntity()
		e5 := gen.NewEntity(
			&Position{},
			&ChildOf{}, &ChildOf2{},
			RelIdx(1, e4), RelIdx(2, e2),
		)
		check("after adding e5, query by ChildOf", Rel[ChildOf](e4), e5)
		check("after adding e5, query by ChildOf2", Rel[ChildOf2](e2), e3, e5)
	}

	t.Run("fresh world", func(t *testing.T) { f(t, &world) })
	world.Reset()
	t.Run("after reset", func(t *testing.T) { f(t, &world) })
}
