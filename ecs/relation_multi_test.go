package ecs

import "testing"

// Dedicated tests for entities with more than one relation component (#219).

func TestMultipleRelationsMap2SetRelations(t *testing.T) {
	w := NewWorld(16)

	parentA := w.NewEntity()
	parentB := w.NewEntity()
	parentC := w.NewEntity()
	parentD := w.NewEntity()

	mapper := NewMap2[ChildOf, ChildOf2](w)

	child := mapper.NewEntity(
		&ChildOf{}, &ChildOf2{},
		RelIdx(0, parentA),
		RelIdx(1, parentB),
	)
	expectEqual(t, parentA, mapper.GetRelation(child, 0))
	expectEqual(t, parentB, mapper.GetRelation(child, 1))

	mapper.SetRelations(child, RelIdx(0, parentC))
	expectEqual(t, parentC, mapper.GetRelation(child, 0))
	expectEqual(t, parentB, mapper.GetRelation(child, 1))

	mapper.SetRelations(child, RelIdx(1, parentD))
	expectEqual(t, parentC, mapper.GetRelation(child, 0))
	expectEqual(t, parentD, mapper.GetRelation(child, 1))

	mapper.SetRelations(child, Rel[ChildOf](parentA), Rel[ChildOf2](parentB))
	expectEqual(t, parentA, mapper.GetRelation(child, 0))
	expectEqual(t, parentB, mapper.GetRelation(child, 1))
}

func TestMultipleRelationsMap2NewBatch(t *testing.T) {
	n := 12
	w := NewWorld(16)

	parent1 := w.NewEntity()
	parent2 := w.NewEntity()
	parent3 := w.NewEntity()

	mapper := NewMap2[ChildOf, ChildOf2](w)
	filter := NewFilter2[ChildOf, ChildOf2](w)

	mapper.NewBatchFn(n, nil, RelIdx(0, parent1), RelIdx(1, parent1))
	mapper.NewBatchFn(n, nil, RelIdx(0, parent2), RelIdx(1, parent2))
	mapper.NewBatchFn(n, nil, RelIdx(0, parent3), RelIdx(1, parent3))

	query := filter.Query()
	expectEqual(t, 3*n, query.Count())
	query.Close()

	query = filter.Query(RelIdx(0, parent2), RelIdx(1, parent2))
	expectEqual(t, n, query.Count())
	query.Close()

	query = filter.Query(RelIdx(0, parent1))
	expectEqual(t, n, query.Count())
	cnt := 0
	for query.Next() {
		entity := query.Entity()
		expectEqual(t, parent1, mapper.GetRelation(entity, 0))
		expectEqual(t, parent1, mapper.GetRelation(entity, 1))
		cnt++
	}
	expectEqual(t, n, cnt)
}

func TestMultipleRelationsMap2SetRelationsBatch(t *testing.T) {
	n := 16
	w := NewWorld(16)

	parent1 := w.NewEntity()
	parent2 := w.NewEntity()
	parent3 := w.NewEntity()
	parent4 := w.NewEntity()

	mapper := NewMap2[ChildOf, ChildOf2](w)
	filter := NewFilter2[ChildOf, ChildOf2](w)

	mapper.NewBatchFn(n, nil, RelIdx(0, parent1), RelIdx(1, parent1))
	mapper.NewBatchFn(n, nil, RelIdx(0, parent2), RelIdx(1, parent2))

	mapper.SetRelationsBatch(
		filter.Batch(RelIdx(1, parent1)),
		func(entity Entity) {
			expectEqual(t, parent1, mapper.GetRelation(entity, 0))
			expectEqual(t, parent4, mapper.GetRelation(entity, 1))
		},
		RelIdx(1, parent4),
	)

	query := filter.Query(RelIdx(1, parent4))
	expectEqual(t, n, query.Count())
	query.Close()

	query = filter.Query(RelIdx(1, parent1))
	expectEqual(t, 0, query.Count())
	query.Close()

	mapper.SetRelationsBatch(
		filter.Batch(RelIdx(0, parent2)),
		nil,
		RelIdx(0, parent3),
	)

	query = filter.Query(RelIdx(0, parent3), RelIdx(1, parent2))
	expectEqual(t, n, query.Count())
	query.Close()

	query = filter.Query(RelIdx(0, parent2))
	expectEqual(t, 0, query.Count())
	query.Close()
}

func TestMultipleRelationsMap3IndependentTargets(t *testing.T) {
	w := NewWorld(8)

	p1 := w.NewEntity()
	p2 := w.NewEntity()
	p3 := w.NewEntity()

	mapper := NewMap3[Position, ChildOf, ChildOf2](w)
	child := mapper.NewEntity(
		&Position{X: 1},
		&ChildOf{}, &ChildOf2{},
		RelIdx(1, p1),
		RelIdx(2, p2),
	)

	childOf := NewMap1[ChildOf](w)
	childOf2 := NewMap1[ChildOf2](w)
	expectEqual(t, p1, childOf.GetRelation(child, 0))
	expectEqual(t, p2, childOf2.GetRelation(child, 0))

	mapper.SetRelations(child, RelIdx(1, p3))
	expectEqual(t, p3, childOf.GetRelation(child, 0))
	expectEqual(t, p2, childOf2.GetRelation(child, 0))

	pos, _, _ := mapper.Get(child)
	expectEqual(t, 1.0, pos.X)
}
