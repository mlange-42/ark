package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	out = inRelations.toRelations(&w, []ID{posID, childID, child2ID}, nil, out)

	assert.Equal(t, []RelationID{
		{component: childID, target: Entity{2, 0}},
		{component: child2ID, target: Entity{3, 0}},
	}, out)

	assert.Panics(t, func() {
		_ = inRelations.toRelations(&w, []ID{childID, child2ID, posID}, nil, out)
	})

	inRelations = relations{RelID(childID, Entity{2, 0}), RelID(child2ID, Entity{3, 0})}
	out = inRelations.toRelations(&w, []ID{posID, childID, child2ID}, nil, out)

	assert.Equal(t, []RelationID{
		{component: childID, target: Entity{2, 0}},
		{component: child2ID, target: Entity{3, 0}},
	}, out)

	inRelations = relations{Rel[ChildOf](Entity{2, 0}), Rel[ChildOf2](Entity{3, 0})}
	out = inRelations.toRelations(&w, []ID{posID, childID, child2ID}, nil, out)

	assert.Equal(t, []RelationID{
		{component: childID, target: Entity{2, 0}},
		{component: child2ID, target: Entity{3, 0}},
	}, out)

	inRelations = relations{Rel[ChildOf](Entity{2, 0})}
	out = inRelations.toRelations(&w, []ID{posID, childID, child2ID}, nil, out)
	assert.Equal(t, []RelationID{
		{component: childID, target: Entity{2, 0}},
	}, out)

	inRelations = relations{Rel[ChildOf2](Entity{3, 0})}
	out = inRelations.toRelations(&w, []ID{posID, childID, child2ID}, append([]RelationID(nil), out...), out)

	assert.Equal(t, []RelationID{
		{component: childID, target: Entity{2, 0}},
		{component: child2ID, target: Entity{3, 0}},
	}, out)
}
