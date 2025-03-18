package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraph(t *testing.T) {
	g := newGraph()

	mask := newMask()
	node := g.Find(0, []ID{id(0), id(1)}, []ID{}, &mask)
	assert.EqualValues(t, 3, g.nodes.Len())
	assert.EqualValues(t, 2, node.id)
	assert.Equal(t, newMask(id(0), id(1)), node.mask)

	mask = node.mask
	node = g.Find(node.id, []ID{}, []ID{id(1)}, &mask)
	assert.EqualValues(t, 1, node.id)
	assert.Equal(t, newMask(id(0)), node.mask)

	mask = node.mask
	node = g.Find(node.id, []ID{id(2), id(3)}, []ID{id(0)}, &mask)
	assert.EqualValues(t, 4, node.id)
	assert.Equal(t, newMask(id(2), id(3)), node.mask)

	mask = node.mask
	node = g.Find(node.id, []ID{id(0)}, []ID{id(2), id(3)}, &mask)
	assert.EqualValues(t, 1, node.id)
	assert.Equal(t, newMask(id(0)), node.mask)

	mask = node.mask
	assert.PanicsWithValue(t,
		"entity does not have component with ID 3",
		func() { g.Find(node.id, []ID{}, []ID{id(3)}, &mask) })
	assert.PanicsWithValue(t,
		"entity already has component with ID 0, or it was added twice",
		func() { g.Find(node.id, []ID{id(0)}, []ID{}, &mask) })
	assert.PanicsWithValue(t,
		"component with ID 0 added and removed in the same exchange operation",
		func() { g.Find(node.id, []ID{id(0)}, []ID{id(0)}, &mask) })
}
