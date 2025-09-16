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

// Relation is the common interface for specifying relationship targets.
// It is implemented by [Rel], [RelIdx] and [RelID].
//
//   - [Rel] is safe, but has some run-time overhead for component [ID] lookup.
//   - [RelIdx] is fast but more error-prone.
//   - [RelID] is used in the [Unsafe] API.
type Relation interface {
	id(ids []ID, world *World) ID
	targetEntity() Entity
}

// RelationID specifies an entity relation target by component [ID].
// Create with [RelID].
//
// It is used in Ark's unsafe, ID-based API.
type RelationID struct {
	component ID
	target    Entity
}

// RelID creates a new [Relation] for a component ID.
//
// It is used in Ark's unsafe, ID-based API.
func RelID(id ID, target Entity) RelationID {
	return RelationID{
		component: id,
		target:    target,
	}
}

// id returns the component ID of this RelationID.
func (r RelationID) id(ids []ID, world *World) ID {
	return r.component
}

// targetEntity returns the target [Entity] of this RelationID.
func (r RelationID) targetEntity() Entity {
	return r.target
}

// relationType specifies an entity relation target by component type.
// Create with [Rel].
//
// It can be used as a safer but slower alternative to [relationIndex].
type relationType[C any] struct {
	target Entity
}

// Rel creates a new [Relation] for a component type.
//
// It can be used as a safer but slower alternative to [RelIdx].
func Rel[C any](target Entity) Relation {
	return relationType[C]{
		target: target,
	}
}

// id returns the component ID of this RelationID.
func (r relationType[C]) id(ids []ID, world *World) ID {
	return ComponentID[C](world)
}

// targetEntity returns the target [Entity] of this RelationID.
func (r relationType[C]) targetEntity() Entity {
	return r.target
}

// relationIndex specifies an entity relation target by component index.
// Create with [RelIdx].
//
// It can be used as faster but more error-prone alternative to [Rel].
type relationIndex struct {
	index  uint8
	target Entity
}

// RelIdx creates a new [Relation] for a component index.
//
// It can be used as faster but less safe alternative to [Rel].
//
// Note that the index refers to the position of the component in the generics
// of e.g. a [Map2] or [Filter2].
// This should not be confused with component [ID] as obtained by [ComponentID]!
func RelIdx(index int, target Entity) Relation {
	return relationIndex{
		index:  uint8(index),
		target: target,
	}
}

// id returns the component ID of this RelationIndex.
func (r relationIndex) id(ids []ID, world *World) ID {
	return ids[r.index]
}

// targetEntity returns the target [Entity] of this RelationIndex.
func (r relationIndex) targetEntity() Entity {
	return r.target
}

// Helper for converting relations
type relations []Relation

func (r relations) toRelations(world *World, mask *bitMask, ids []ID, out []RelationID, startIdx uint8) []RelationID {
	// TODO: can this be made more efficient?
	out = out[:startIdx]
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
		out = append(out, RelationID{
			component: id,
			target:    rel.targetEntity(),
		})
	}
	return out
}

func (e Entity) toRelation(world *World, id ID, out []RelationID) []RelationID {
	world.storage.checkRelationTarget(e)
	world.storage.checkRelationComponent(id)
	out = out[:0]
	out = append(out, RelationID{
		component: id,
		target:    e,
	})
	return out
}

// Helper for converting relations
type relationEntities []Entity

func (r relationEntities) toRelation(world *World, id ID, out []RelationID) []RelationID {
	out = out[:0]
	if len(r) == 0 {
		return out
	}
	for _, rel := range r {
		world.storage.checkRelationTarget(rel)
		world.storage.checkRelationComponent(id)
		out = append(out, RelationID{
			component: id,
			target:    rel,
		})
	}
	return out
}
