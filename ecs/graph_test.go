package ecs

import (
	"testing"
	"unsafe"
)

func TestGraph(t *testing.T) {
	g := newGraph()

	mask := newMask()
	node := g.Find(0, []ID{id(0), id(1)}, []ID{}, &mask)
	if g.nodes.Len() != 3 {
		t.Errorf("expected graph to have 3 nodes, got %d", g.nodes.Len())
	}
	if node.id != 2 {
		t.Errorf("expected node.id == 2, got %d", node.id)
	}
	if node.mask != newMask(id(0), id(1)) {
		t.Errorf("expected mask %v, got %v", newMask(id(0), id(1)), node.mask)
	}

	mask = node.mask
	node = g.Find(node.id, []ID{}, []ID{id(1)}, &mask)
	if node.id != 1 {
		t.Errorf("expected node.id == 1, got %d", node.id)
	}
	if node.mask != newMask(id(0)) {
		t.Errorf("expected mask %v, got %v", newMask(id(0)), node.mask)
	}

	mask = node.mask
	node = g.Find(node.id, []ID{id(2), id(3)}, []ID{id(0)}, &mask)
	if node.id != 4 {
		t.Errorf("expected node.id == 4, got %d", node.id)
	}
	if node.mask != newMask(id(2), id(3)) {
		t.Errorf("expected mask %v, got %v", newMask(id(2), id(3)), node.mask)
	}

	mask = node.mask
	node = g.Find(node.id, []ID{id(0)}, []ID{id(2), id(3)}, &mask)
	if node.id != 1 {
		t.Errorf("expected node.id == 1, got %d", node.id)
	}
	if node.mask != newMask(id(0)) {
		t.Errorf("expected mask %v, got %v", newMask(id(0)), node.mask)
	}

	mask = node.mask
	expectPanicWithValue(t,
		"entity does not have component with ID 3",
		func() { g.Find(node.id, []ID{}, []ID{id(3)}, &mask) })

	expectPanicWithValue(t,
		"entity already has component with ID 0, or it was added twice",
		func() { g.Find(node.id, []ID{id(0)}, []ID{}, &mask) })

	expectPanicWithValue(t,
		"component with ID 0 added and removed in the same exchange operation",
		func() { g.Find(node.id, []ID{id(0)}, []ID{id(0)}, &mask) })
}

func TestGraphNodePointers(t *testing.T) {
	g := newGraph()
	ptr := g.nodes.Get(0)

	node := g.nodes.Get(0)
	for i := range maskTotalBits {
		node = g.Find(node.id, []ID{id(i)}, nil, &bitMask{})
	}

	if unsafe.Pointer(ptr) != unsafe.Pointer(g.nodes.Get(0)) {
		t.Errorf("expected pointer to remain unchanged")
	}
}
