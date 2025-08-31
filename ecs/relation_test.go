package ecs

import (
	"testing"
)

func TestRel(t *testing.T) {
	r := RelIdx(1, Entity{5, 0})
	if r != (relationIndex{1, Entity{5, 0}}) {
		t.Errorf("expected %v, got %v", relationIndex{1, Entity{5, 0}}, r)
	}
}

func TestRelID(t *testing.T) {
	r := RelID(ID{10}, Entity{5, 0})
	if r != (RelationID{ID{10}, Entity{5, 0}}) {
		t.Errorf("expected %v, got %v", RelationID{ID{10}, Entity{5, 0}}, r)
	}
}

func TestToRelations(t *testing.T) {
	w := NewWorld()
	childID := ComponentID[ChildOf](&w)
	child2ID := ComponentID[ChildOf2](&w)
	posID := ComponentID[Position](&w)
	inRelations := relations{RelIdx(1, Entity{2, 0}), RelIdx(2, Entity{3, 0})}
	var out []RelationID
	out = inRelations.toRelations(&w, []ID{posID, childID, child2ID}, out, 0)
	if !equalRelations(out, []RelationID{
		{component: childID, target: Entity{2, 0}},
		{component: child2ID, target: Entity{3, 0}},
	}) {
		t.Errorf("expected %v, got %v", []RelationID{
			{component: childID, target: Entity{2, 0}},
			{component: child2ID, target: Entity{3, 0}},
		}, out)
	}
	expectPanic(t, func() {
		_ = inRelations.toRelations(&w, []ID{childID, child2ID, posID}, out, 0)
	})
	inRelations = relations{RelID(childID, Entity{2, 0}), RelID(child2ID, Entity{3, 0})}
	out = inRelations.toRelations(&w, []ID{posID, childID, child2ID}, out, 0)
	if !equalRelations(out, []RelationID{
		{component: childID, target: Entity{2, 0}},
		{component: child2ID, target: Entity{3, 0}},
	}) {
		t.Errorf("expected %v, got %v", []RelationID{
			{component: childID, target: Entity{2, 0}},
			{component: child2ID, target: Entity{3, 0}},
		}, out)
	}
	inRelations = relations{Rel[ChildOf](Entity{2, 0}), Rel[ChildOf2](Entity{3, 0})}
	out = inRelations.toRelations(&w, []ID{posID, childID, child2ID}, out, 0)
	if !equalRelations(out, []RelationID{
		{component: childID, target: Entity{2, 0}},
		{component: child2ID, target: Entity{3, 0}},
	}) {
		t.Errorf("expected %v, got %v", []RelationID{
			{component: childID, target: Entity{2, 0}},
			{component: child2ID, target: Entity{3, 0}},
		}, out)
	}
	inRelations = relations{Rel[ChildOf](Entity{2, 0})}
	out = inRelations.toRelations(&w, []ID{posID, childID, child2ID}, out, 0)
	if !equalRelations(out, []RelationID{
		{component: childID, target: Entity{2, 0}},
	}) {
		t.Errorf("expected %v, got %v", []RelationID{
			{component: childID, target: Entity{2, 0}},
		}, out)
	}
	inRelations = relations{Rel[ChildOf2](Entity{3, 0})}
	out = inRelations.toRelations(&w, []ID{posID, childID, child2ID}, out, uint8(len(out)))
	if !equalRelations(out, []RelationID{
		{component: childID, target: Entity{2, 0}},
		{component: child2ID, target: Entity{3, 0}},
	}) {
		t.Errorf("expected %v, got %v", []RelationID{
			{component: childID, target: Entity{2, 0}},
			{component: child2ID, target: Entity{3, 0}},
		}, out)
	}
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
	if childMap.GetRelation(e3) != e1 {
		t.Errorf("expected %v, got %v", e1, childMap.GetRelation(e3))
	}
	if child2Map.GetRelation(e3) != e2 {
		t.Errorf("expected %v, got %v", e2, child2Map.GetRelation(e3))
	}
	world.RemoveEntity(e1)
	if childMap.GetRelation(e3) != (Entity{}) {
		t.Errorf("expected %v, got %v", Entity{}, childMap.GetRelation(e3))
	}
	if child2Map.GetRelation(e3) != e2 {
		t.Errorf("expected %v, got %v", e2, child2Map.GetRelation(e3))
	}
	world.RemoveEntity(e2)
	if childMap.GetRelation(e3) != (Entity{}) {
		t.Errorf("expected %v, got %v", Entity{}, childMap.GetRelation(e3))
	}
	if child2Map.GetRelation(e3) != (Entity{}) {
		t.Errorf("expected %v, got %v", Entity{}, child2Map.GetRelation(e3))
	}
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
	if query.Count() != 2 {
		t.Errorf("expected 2, got %d", query.Count())
	}
	query.Close()
}
