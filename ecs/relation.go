package ecs

import "fmt"

var relationType = typeOf[Relation]()

// Relation is a marker for entity relation components.
// It must be embedded as first field of a component that represent an entity relation
// (see the example).
//
// Entity relations allow for fast queries using entity relationships.
// E.g. to iterate over all entities that are the child of a certain parent entity.
// Currently, each entity can only have a single relation component.
type Relation struct{}

type relationID struct {
	component ID
	target    Entity
}

// RelationIndex specifies an entity relation target by component index.
//
// Note that the index refers to the position of the component in the generics
// of e.g. a [Map2] or [Filter2].
// This should not be confused with component IDs as obtained by [ComponentID]!
type RelationIndex struct {
	index  uint8
	target Entity
}

// Rel creates a new RelationIndex.
func Rel(index int, target Entity) RelationIndex {
	return RelationIndex{
		index:  uint8(index),
		target: target,
	}
}

type relations []RelationIndex

func (r relations) toRelations(reg *componentRegistry, ids []ID, out []relationID) []relationID {
	out = out[:0]
	for _, rel := range r {
		id := ids[rel.index]
		if !reg.IsRelation[id.id] {
			panic(fmt.Sprintf("component at index %d is not a relation component", rel.index))
		}
		out = append(out, relationID{
			component: id,
			target:    rel.target,
		})
	}
	return out
}

type relationEntities []Entity

func (r relationEntities) toRelation(id ID, out []relationID) []relationID {
	out = out[:0]
	for _, rel := range r {
		out = append(out, relationID{
			component: id,
			target:    rel,
		})
	}
	return out
}
