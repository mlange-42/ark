package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraph(t *testing.T) {
	g := newGraph()

	mask := bitMask{}
	node := g.Find(0, []ID{id(0), id(1)}, []ID{}, &mask)
	assert.Equal(t, 3, len(g.nodes))
	assert.EqualValues(t, 2, node.id)
	assert.Equal(t, newMask(id(0), id(1)), node.mask)

	mask = bitMask{}
	node = g.Find(node.id, []ID{}, []ID{id(1)}, &mask)
	assert.EqualValues(t, 1, node.id)
	assert.Equal(t, newMask(id(0)), node.mask)

	mask = bitMask{}
	node = g.Find(node.id, []ID{id(2), id(3)}, []ID{id(0)}, &mask)
	assert.EqualValues(t, 4, node.id)
	assert.Equal(t, newMask(id(2), id(3)), node.mask)

	mask = bitMask{}
	node = g.Find(node.id, []ID{id(0)}, []ID{id(2), id(3)}, &mask)
	assert.EqualValues(t, 1, node.id)
	assert.Equal(t, newMask(id(0)), node.mask)

	mask = bitMask{}
	assert.Panics(t, func() { g.Find(node.id, []ID{}, []ID{id(3)}, &mask) })
	assert.Panics(t, func() { g.Find(node.id, []ID{id(0)}, []ID{}, &mask) })
	assert.Panics(t, func() { g.Find(node.id, []ID{id(0)}, []ID{id(0)}, &mask) })
}
