package ecs

import "fmt"

type nodeID uint32

type node struct {
	id        nodeID
	archetype archetypeID
	mask      Mask
	neighbors idMap[nodeID]
}

func newNode(id nodeID, archetype archetypeID, mask Mask) node {
	return node{
		id:        id,
		archetype: archetype,
		mask:      mask,
		neighbors: newIDMap[nodeID](),
	}
}

func (n *node) GetArchetype() (archetypeID, bool) {
	return n.archetype, n.archetype != maxArchetypeID
}

// Archetype graph for faster lookup of transitions.
type graph struct {
	nodes []node
}

func newGraph() graph {
	return graph{
		nodes: []node{
			newNode(0, 0, NewMask()),
		},
	}
}

func (g *graph) Find(start nodeID, add []ID, remove []ID, outMask *Mask) *node {
	startNode := &g.nodes[start]
	curr := startNode
	*outMask = startNode.mask

	for _, id := range remove {
		if !outMask.Get(id) {
			panic(fmt.Sprintf("entity does not have component with ID %d", id.id))
		}
		outMask.Set(id, false)
		if next, ok := curr.neighbors.Get(id.id); ok {
			curr = &g.nodes[next]
		} else {
			next := g.findOrCreate(outMask)
			next.neighbors.Set(id.id, curr.id)
			curr.neighbors.Set(id.id, next.id)
			curr = next
		}
	}

	for _, id := range add {
		if outMask.Get(id) {
			panic(fmt.Sprintf("entity already has component with ID %d, or it was added twice", id.id))
		}
		if startNode.mask.Get(id) {
			panic(fmt.Sprintf("component with ID %d added and removed in the same exchange operation", id.id))
		}

		outMask.Set(id, true)
		if next, ok := curr.neighbors.Get(id.id); ok {
			curr = &g.nodes[next]
		} else {
			next := g.findOrCreate(outMask)
			next.neighbors.Set(id.id, curr.id)
			curr.neighbors.Set(id.id, next.id)
			curr = next
		}
	}
	return curr
}

func (g *graph) findOrCreate(mask *Mask) *node {
	for i := range g.nodes {
		if g.nodes[i].mask.Equals(mask) {
			return &g.nodes[i]
		}
	}
	idx := len(g.nodes)
	g.nodes = append(g.nodes, newNode(nodeID(idx), maxArchetypeID, *mask))
	return &g.nodes[idx]
}
