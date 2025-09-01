package ecs

import (
	"testing"
)

func TestRel(t *testing.T) {
	r := RelIdx(1, Entity{5, 0})
	expectEqual(t, relationIndex{1, Entity{5, 0}}, r.(relationIndex))
}

func TestRelID(t *testing.T) {
	r := RelID(ID{10}, Entity{5, 0})
	expectEqual(t, RelationID{ID{10}, Entity{5, 0}}, r)
}

func TestToRelations(t *testing.T) {
	w := NewWorld()

	childID := ComponentID[ChildOf](&w)
	child2ID := ComponentID[ChildOf2](&w)
	posID := ComponentID[Position](&w)

	inRelations := relations{RelIdx(1, Entity{2, 0}), RelIdx(2, Entity{3, 0})}
	var out []RelationID
	out = inRelations.toRelations(&w, []ID{posID, childID, child2ID}, out, 0)

	expectSlicesEqual(t, []RelationID{
		{component: childID, target: Entity{2, 0}},
		{component: child2ID, target: Entity{3, 0}},
	}, out)

	expectPanics(t, func() {
		_ = inRelations.toRelations(&w, []ID{childID, child2ID, posID}, out, 0)
	})

	inRelations = relations{RelID(childID, Entity{2, 0}), RelID(child2ID, Entity{3, 0})}
	out = inRelations.toRelations(&w, []ID{posID, childID, child2ID}, out, 0)

	expectSlicesEqual(t, []RelationID{
		{component: childID, target: Entity{2, 0}},
		{component: child2ID, target: Entity{3, 0}},
	}, out)

	inRelations = relations{Rel[ChildOf](Entity{2, 0}), Rel[ChildOf2](Entity{3, 0})}
	out = inRelations.toRelations(&w, []ID{posID, childID, child2ID}, out, 0)

	expectSlicesEqual(t, []RelationID{
		{component: childID, target: Entity{2, 0}},
		{component: child2ID, target: Entity{3, 0}},
	}, out)

	inRelations = relations{Rel[ChildOf](Entity{2, 0})}
	out = inRelations.toRelations(&w, []ID{posID, childID, child2ID}, out, 0)
	expectSlicesEqual(t, []RelationID{
		{component: childID, target: Entity{2, 0}},
	}, out)

	inRelations = relations{Rel[ChildOf2](Entity{3, 0})}
	out = inRelations.toRelations(&w, []ID{posID, childID, child2ID}, out, uint8(len(out)))

	expectSlicesEqual(t, []RelationID{
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

	expectEqual(t, e1, childMap.GetRelation(e3))
	expectEqual(t, e2, child2Map.GetRelation(e3))

	world.RemoveEntity(e1)
	expectEqual(t, Entity{}, childMap.GetRelation(e3))
	expectEqual(t, e2, child2Map.GetRelation(e3))

	world.RemoveEntity(e2)
	expectEqual(t, Entity{}, childMap.GetRelation(e3))
	expectEqual(t, Entity{}, child2Map.GetRelation(e3))
}

func TestStaleRelationTable(t *testing.T) {
	world := NewWorld()
	filter := NewFilter1[Position](&world)

	e1 := world.NewEntity()
	e2 := world.NewEntity()

	gen := NewMap3[Position, ChildOf, ChildOf2](&world)
	gen.NewEntity(
		&Position{},
		&ChildOf{}, &ChildOf2{},
		RelIdx(1, e1),
		RelIdx(2, e2),
	)

	world.RemoveEntity(e1)

	e4 := world.NewEntity()
	gen.NewEntity(
		&Position{},
		&ChildOf{}, &ChildOf2{},
		RelIdx(1, e4),
		RelIdx(2, e2),
	)

	query := filter.Query(Rel[ChildOf2](e2))
	expectEqual(t, 2, query.Count())
	query.Close()
}
