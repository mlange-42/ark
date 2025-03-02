package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraph(t *testing.T) {
	g := newGraph()

	mask := NewMask()
	node := g.Find(0, &mask, []ID{id(0), id(1)}, []ID{})
	assert.Equal(t, 3, len(g.nodes))
	assert.EqualValues(t, 2, node.id)
	assert.Equal(t, NewMask(id(0), id(1)), node.mask)

	node = g.Find(node.id, &node.mask, []ID{}, []ID{id(1)})
	assert.EqualValues(t, 1, node.id)
	assert.Equal(t, NewMask(id(0)), node.mask)

	node = g.Find(node.id, &node.mask, []ID{id(2), id(3)}, []ID{id(0)})
	assert.EqualValues(t, 4, node.id)
	assert.Equal(t, NewMask(id(2), id(3)), node.mask)

	node = g.Find(node.id, &node.mask, []ID{id(0)}, []ID{id(2), id(3)})
	assert.EqualValues(t, 1, node.id)
	assert.Equal(t, NewMask(id(0)), node.mask)

	assert.Panics(t, func() { g.Find(node.id, &node.mask, []ID{}, []ID{id(3)}) })
	assert.Panics(t, func() { g.Find(node.id, &node.mask, []ID{id(0)}, []ID{}) })
	assert.Panics(t, func() { g.Find(node.id, &node.mask, []ID{id(0)}, []ID{id(0)}) })
}
