package ecs

import "fmt"

// nodeID is the ID type for graph nodes.
type nodeID uint32

// node in the archetype graph.
type node struct {
	neighbors idMap
	mask      bitMask
	id        nodeID
	archetype archetypeID
}

// newNode creates a new node for the given archetype ID.
//
// The mask gets copied.
func newNode(id nodeID, archetype archetypeID, mask *bitMask) node {
	return node{
		id:        id,
		archetype: archetype,
		mask:      *mask,
		neighbors: newIDMap(),
	}
}

// GetArchetype returns the archetype of the node, and whether the node has an archetype.
func (n *node) GetArchetype() (archetypeID, bool) {
	return n.archetype, n.archetype != maxArchetypeID
}

// Archetype graph for faster lookup of transitions.
type graph struct {
	nodes []node
}

// newGraph creates a new empty graph.
func newGraph() graph {
	nodes := make([]node, 0, 128)
	nodes = append(nodes, newNode(0, 0, &bitMask{}))

	return graph{
		nodes: nodes,
	}
}

// Find the node for a given start node and components to add and remove.
//
// The bitMask argument gets modified and reflects the mask of the resulting node.
func (g *graph) Find(start nodeID, add []ID, remove []ID, outMask *bitMask) *node {
	startNode := &g.nodes[start]
	curr := startNode

	for _, id := range remove {
		if !outMask.Get(id.id) {
			panic(fmt.Sprintf("entity does not have component with ID %d", id.id))
		}
		outMask.Clear(id.id)
		if curr.neighbors.used.Get(id.id) {
			curr = &g.nodes[curr.neighbors.data[id.id]]
		} else {
			next := g.findOrCreate(outMask)
			next.neighbors.Set(id.id, curr.id)
			curr = &g.nodes[curr.id]
			curr.neighbors.Set(id.id, next.id)
			curr = next
		}
	}

	for _, id := range add {
		if outMask.Get(id.id) {
			panic(fmt.Sprintf("entity already has component with ID %d, or it was added twice", id.id))
		}
		if startNode.mask.Get(id.id) {
			panic(fmt.Sprintf("component with ID %d added and removed in the same exchange operation", id.id))
		}
		outMask.Set(id.id)
		if curr.neighbors.used.Get(id.id) {
			curr = &g.nodes[curr.neighbors.data[id.id]]
		} else {
			next := g.findOrCreate(outMask)
			next.neighbors.Set(id.id, curr.id)
			curr = &g.nodes[curr.id]
			curr.neighbors.Set(id.id, next.id)
			curr = next
		}
	}

	return curr
}

// FindAdd finds the node for a given start node and components to add.
//
// The bitMask argument gets modified and reflects the mask of the resulting node.
func (g *graph) FindAdd(start nodeID, add []ID, outMask *bitMask) *node {
	startNode := &g.nodes[start]
	curr := startNode

	for _, id := range add {
		if outMask.Get(id.id) {
			panic(fmt.Sprintf("entity already has component with ID %d, or it was added twice", id.id))
		}
		outMask.Set(id.id)
		if curr.neighbors.used.Get(id.id) {
			curr = &g.nodes[curr.neighbors.data[id.id]]
		} else {
			next := g.findOrCreate(outMask)
			next.neighbors.Set(id.id, curr.id)
			curr = &g.nodes[curr.id]
			curr.neighbors.Set(id.id, next.id)
			curr = next
		}
	}

	return curr
}

// FindRemove finds the node for a given start node and components to remove.
//
// The bitMask argument gets modified and reflects the mask of the resulting node.
func (g *graph) FindRemove(start nodeID, remove []ID, outMask *bitMask) *node {
	startNode := &g.nodes[start]
	curr := startNode

	for _, id := range remove {
		if !outMask.Get(id.id) {
			panic(fmt.Sprintf("entity does not have component with ID %d", id.id))
		}
		outMask.Clear(id.id)
		if curr.neighbors.used.Get(id.id) {
			curr = &g.nodes[curr.neighbors.data[id.id]]
		} else {
			next := g.findOrCreate(outMask)
			next.neighbors.Set(id.id, curr.id)
			curr = &g.nodes[curr.id]
			curr.neighbors.Set(id.id, next.id)
			curr = next
		}
	}

	return curr
}

// finds or creates a node if not reachable in the current graph.
func (g *graph) findOrCreate(mask *bitMask) *node {
	len := len(g.nodes)
	for i := range len {
		node := &g.nodes[i]
		if node.mask.Equals(mask) {
			return node
		}
	}
	g.nodes = append(g.nodes, newNode(nodeID(len), maxArchetypeID, mask))
	return &g.nodes[len]
}
