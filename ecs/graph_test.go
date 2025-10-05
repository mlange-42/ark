package ecs

import (
	"testing"
)

func TestGraph(t *testing.T) {
	g := newGraph()

	mask := newMask()
	node := g.Find(0, []ID{id(0), id(1)}, []ID{}, &mask)
	expectEqual(t, 3, len(g.nodes))
	expectEqual(t, 2, node.id)
	expectEqual(t, newMask(id(0), id(1)), node.mask)

	mask = node.mask
	node = g.Find(node.id, []ID{}, []ID{id(1)}, &mask)
	expectEqual(t, 1, node.id)
	expectEqual(t, newMask(id(0)), node.mask)

	mask = node.mask
	node = g.Find(node.id, []ID{id(2), id(3)}, []ID{id(0)}, &mask)
	expectEqual(t, 4, node.id)
	expectEqual(t, newMask(id(2), id(3)), node.mask)

	mask = node.mask
	node = g.Find(node.id, []ID{id(0)}, []ID{id(2), id(3)}, &mask)
	expectEqual(t, 1, node.id)
	expectEqual(t, newMask(id(0)), node.mask)

	mask = node.mask
	expectPanicsWithValue(t,
		"entity does not have component with ID 3",
		func() { g.Find(node.id, []ID{}, []ID{id(3)}, &mask) })
	expectPanicsWithValue(t,
		"entity already has component with ID 0, or it was added twice",
		func() { g.Find(node.id, []ID{id(0)}, []ID{}, &mask) })
	expectPanicsWithValue(t,
		"component with ID 0 added and removed in the same exchange operation",
		func() { g.Find(node.id, []ID{id(0)}, []ID{id(0)}, &mask) })
}

func BenchmarkGraphFind1(b *testing.B) {
	g := newGraph()

	id1 := ID{0}

	add := []ID{{1}}

	mask1 := newMask()
	node := g.Find(0, []ID{id1}, nil, &mask1)

	for b.Loop() {
		n := g.Find(node.id, add, nil, &mask1)
		g.Find(n.id, nil, add, &mask1)
	}
}

func BenchmarkGraphFind10(b *testing.B) {
	g := newGraph()

	id1 := ID{0}

	add := []ID{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {9}, {10}}

	mask1 := newMask()
	node := g.Find(0, []ID{id1}, nil, &mask1)

	for b.Loop() {
		n := g.Find(node.id, add, nil, &mask1)
		g.Find(n.id, nil, add, &mask1)
	}
}
