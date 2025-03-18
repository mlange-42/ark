package ecs

import "fmt"

type nodeID uint32

type node struct {
	neighbors idMap[*node]
	mask      bitMask
	id        nodeID
	archetype archetypeID
}

func newNode(id nodeID, archetype archetypeID, mask bitMask) node {
	return node{
		id:        id,
		archetype: archetype,
		mask:      mask,
		neighbors: newIDMap[*node](),
	}
}

func (n *node) GetArchetype() (archetypeID, bool) {
	return n.archetype, n.archetype != maxArchetypeID
}

// Archetype graph for faster lookup of transitions.
type graph struct {
	nodes pagedSlice[node]
}

func newGraph() graph {
	nodes := pagedSlice[node]{}
	nodes.Add(newNode(0, 0, newMask()))

	return graph{
		nodes: nodes,
	}
}

func (g *graph) Find(start nodeID, add []ID, remove []ID, outMask *bitMask) *node {
	startNode := g.nodes.Get(int32(start))
	curr := startNode

	for _, id := range remove {
		if !outMask.Get(id) {
			panic(fmt.Sprintf("entity does not have component with ID %d", id.id))
		}
		outMask.Set(id, false)
		if next, ok := curr.neighbors.Get(id.id); ok {
			curr = next
		} else {
			next := g.findOrCreate(outMask)
			next.neighbors.Set(id.id, curr)
			curr.neighbors.Set(id.id, next)
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
			curr = next
		} else {
			next := g.findOrCreate(outMask)
			next.neighbors.Set(id.id, curr)
			curr.neighbors.Set(id.id, next)
			curr = next
		}
	}
	return curr
}

func (g *graph) findOrCreate(mask *bitMask) *node {
	len := g.nodes.Len()
	for i := range len {
		node := g.nodes.Get(i)
		if node.mask.Equals(mask) {
			return node
		}
	}
	g.nodes.Add(newNode(nodeID(len), maxArchetypeID, *mask))
	return g.nodes.Get(len)
}
