package ecs

import (
	"fmt"
	"reflect"
)

var relationTp = reflect.TypeFor[RelationMarker]()

// RelationMarker is a marker for entity relation components.
// It must be embedded as first field of a component that represent an entity relationship
// (see the example).
//
// Entity relations allow for fast queries using entity relationships.
// E.g. to iterate over all entities that are the child of a certain parent entity.
type RelationMarker struct{}

type relationID struct {
	target    Entity
	component ID
}

// Relation is the common type for specifying relationship targets.
// It can be created with [Rel], [RelIdx] and [RelID].
//
//   - [Rel] is safe, but has some run-time overhead for component [ID] lookup.
//   - [RelIdx] is fast but more error-prone.
//   - [RelID] is used in the [Unsafe] API.
type Relation struct {
	target        Entity
	component     ID
	componentType reflect.Type
	index         uint8
}

func (r Relation) relationID() relationID {
	if r.index < 255 || r.componentType != nil {
		panic("only relations created with RelID can be used in the unsafe API")
	}
	return relationID{
		target:    r.target,
		component: r.component,
	}
}

func relID(id ID, target Entity) relationID {
	return relationID{
		target:    target,
		component: id,
	}
}

// id returns the component ID of this RelationID.
func (r Relation) id(ids []ID, world *World) ID {
	if r.index < 255 {
		return ids[r.index]
	}
	if r.componentType != nil {
		return TypeID(world, r.componentType)
	}
	return r.component
}

// targetEntity returns the target [Entity] of this RelationID.
func (r Relation) targetEntity() Entity {
	return r.target
}

// RelID creates a new [Relation] for a component ID.
//
// It is used in Ark's unsafe, ID-based API.
func RelID(id ID, target Entity) Relation {
	return Relation{
		target:    target,
		component: id,
		index:     255,
	}
}

// Rel creates a new [Relation] for a component type.
//
// It can be used as a safer but slower alternative to [RelIdx].
func Rel[C any](target Entity) Relation {
	return Relation{
		target:        target,
		componentType: reflect.TypeFor[C](),
		index:         255,
	}
}

// RelIdx creates a new [Relation] for a component index.
//
// It can be used as faster but less safe alternative to [Rel].
//
// Note that the index refers to the position of the component in the generics
// of e.g. a [Map2] or [Filter2].
// This should not be confused with component [ID] as obtained by [ComponentID]!
// For component IDs, use [RelationID].
func RelIdx(index int, target Entity) Relation {
	return Relation{
		index:  uint8(index),
		target: target,
	}
}

// Helper for converting relationSlice
type relationSlice []Relation

func (r relationSlice) toRelations(world *World, mask *bitMask, ids []ID, out []relationID) []relationID {
	// TODO: can this be made more efficient?
	if len(r) == 0 {
		return out
	}
	for _, rel := range r {
		id := rel.id(ids, world)
		world.storage.checkRelationTarget(rel.targetEntity())
		world.storage.checkRelationComponent(id)
		if !mask.Get(id) {
			panic(fmt.Sprintf("requested relation component with ID %d was not specified in the filter or map", id.id))
		}
		out = append(out, relationID{target: rel.target, component: id})
	}
	return out
}

func (r relationSlice) toRelationIDs(out []relationID) []relationID {
	for _, rel := range r {
		out = append(out, rel.relationID())
	}
	return out
}

func (e Entity) toRelation(world *World, id ID, out []relationID) []relationID {
	world.storage.checkRelationTarget(e)
	world.storage.checkRelationComponent(id)
	out = out[:0]
	out = append(out, relationID{target: e, component: id})
	return out
}

// Helper for converting relations
type relationEntities []Entity

func (r relationEntities) toRelation(world *World, id ID, out []relationID) []relationID {
	out = out[:0]
	if len(r) == 0 {
		return out
	}
	for _, rel := range r {
		world.storage.checkRelationTarget(rel)
		world.storage.checkRelationComponent(id)
		out = append(out, relationID{target: rel, component: id})
	}
	return out
}
