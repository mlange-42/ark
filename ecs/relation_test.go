package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRel(t *testing.T) {
	r := Rel(1, Entity{5, 0})
	assert.Equal(t, RelationIndex{1, Entity{5, 0}}, r)
}

func TestRelID(t *testing.T) {
	r := RelID(ID{10}, Entity{5, 0})
	assert.Equal(t, RelationID{ID{10}, Entity{5, 0}}, r)
}

func TestToRelations(t *testing.T) {
	reg := newComponentRegistry()

	childID, _ := reg.ComponentID(typeOf[ChildOf]())
	child2ID, _ := reg.ComponentID(typeOf[ChildOf2]())
	posID, _ := reg.ComponentID(typeOf[Position]())

	relations := relations{Rel(1, Entity{2, 0}), Rel(2, Entity{3, 0})}
	var out []RelationID
	out = relations.toRelations(&reg, []ID{{posID}, {childID}, {child2ID}}, nil, out)

	assert.Equal(t, []RelationID{
		{component: ID{childID}, target: Entity{2, 0}},
		{component: ID{child2ID}, target: Entity{3, 0}},
	}, out)

	assert.Panics(t, func() {
		_ = relations.toRelations(&reg, []ID{{childID}, {child2ID}, {posID}}, nil, out)
	})
}
