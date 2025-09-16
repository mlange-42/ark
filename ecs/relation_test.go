package ecs

import (
	"testing"
)

func all(ids ...ID) *bitMask {
	mask := newMask(ids...)
	return &mask
}

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

	ids := []ID{posID, childID, child2ID}
	out = inRelations.toRelations(&w, all(ids...), ids, out[:0])

	expectSlicesEqual(t, []RelationID{
		{component: childID, target: Entity{2, 0}},
		{component: child2ID, target: Entity{3, 0}},
	}, out)

	expectPanicsWithValue(t, "component with ID 2 is not a relation component",
		func() {
			ids := []ID{childID, child2ID, posID}
			_ = inRelations.toRelations(&w, all(ids...), ids, out[:0])
		})

	expectPanicsWithValue(t, "requested relation component with ID 0 was not specified in the filter or map",
		func() {
			ids := []ID{posID, childID, child2ID}
			_ = inRelations.toRelations(&w, all(posID, child2ID), ids, out[:0])
		})

	inRelations = relations{RelID(childID, Entity{2, 0}), RelID(child2ID, Entity{3, 0})}

	out = inRelations.toRelations(&w, all(ids...), ids, out[:0])

	expectSlicesEqual(t, []RelationID{
		{component: childID, target: Entity{2, 0}},
		{component: child2ID, target: Entity{3, 0}},
	}, out)

	inRelations = relations{Rel[ChildOf](Entity{2, 0}), Rel[ChildOf2](Entity{3, 0})}
	out = inRelations.toRelations(&w, all(ids...), ids, out[:0])

	expectSlicesEqual(t, []RelationID{
		{component: childID, target: Entity{2, 0}},
		{component: child2ID, target: Entity{3, 0}},
	}, out)

	inRelations = relations{Rel[ChildOf](Entity{2, 0})}
	out = inRelations.toRelations(&w, all(ids...), ids, out[:0])
	expectSlicesEqual(t, []RelationID{
		{component: childID, target: Entity{2, 0}},
	}, out)

	inRelations = relations{Rel[ChildOf2](Entity{3, 0})}
	out = inRelations.toRelations(&w, all(ids...), ids, out)

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
	filter := NewFilter2[Position, ChildOf2](&world)

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

func TestRelationChecks(t *testing.T) {
	world := NewWorld()
	builder1 := NewMap1[Position](&world)
	builder2 := NewMap2[Position, ChildOf](&world)

	parent := world.NewEntity()
	builder1.NewEntity(&Position{})
	e1 := builder2.NewEntity(&Position{}, &ChildOf{}, Rel[ChildOf](parent))
	e2 := builder2.NewEntity(&Position{}, &ChildOf{}, Rel[ChildOf](parent))

	expectPanicsWithValue(t, "requested relation component with ID 2 was not specified in the filter or map", func() {
		builder2.NewEntity(&Position{}, &ChildOf{}, Rel[ChildOf2](parent))
	})

	exchange1 := NewExchange1[ChildOf2](&world).Removes(C[ChildOf]())
	exchange1.Exchange(e1, &ChildOf2{}, Rel[ChildOf2](parent))

	expectPanicsWithValue(t, "requested relation component with ID 1 was not specified in the filter or map", func() {
		exchange1.Exchange(e2, &ChildOf2{}, Rel[ChildOf](parent))
	})

	filter1 := NewFilter1[Position](&world)
	expectPanicsWithValue(t, "requested relation component with ID 1 was not specified in the filter or map", func() {
		filter1.Query(Rel[ChildOf](parent))
	})
	filter2 := NewFilter1[Position](&world).With(C[ChildOf]())
	query := filter2.Query(Rel[ChildOf](parent))
	query.Close()
}

func BenchmarkToRelations0(b *testing.B) {
	world := NewWorld()

	ids := []ID{
		ComponentID[Position](&world),
		ComponentID[ChildOf](&world),
		ComponentID[ChildOf2](&world),
	}
	mask := newMask(ids...)

	rels := relations{}
	relations := []RelationID{}

	for b.Loop() {
		_ = rels.toRelations(&world, &mask, ids, relations)
	}
}

func BenchmarkToRelations1(b *testing.B) {
	world := NewWorld()
	parent := world.NewEntity()

	ids := []ID{
		ComponentID[Position](&world),
		ComponentID[ChildOf](&world),
		ComponentID[ChildOf2](&world),
	}
	mask := newMask(ids...)

	rels := relations{
		RelIdx(1, parent),
	}
	relations := []RelationID{}

	for b.Loop() {
		_ = rels.toRelations(&world, &mask, ids, relations)
	}
}
