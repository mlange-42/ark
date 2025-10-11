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
	componentType reflect.Type
	target        Entity
	component     ID
	index         uint8
}

func (r *Relation) relationIDForUnsafe(world *World) relationID {
	if r.index < 255 {
		panic("relations created with RelIdx can't be used in the unsafe API, use RelID or Rel instead")
	}
	if r.componentType != nil {
		r.component = TypeID(world, r.componentType)
		r.componentType = nil
	}
	return relationID{
		target:    r.target,
		component: r.component,
	}
}

// id returns the component ID of this RelationID.
func (r *Relation) id(ids []ID, world *World) ID {
	if r.index < 255 {
		return ids[r.index]
	}
	if r.componentType != nil {
		r.component = TypeID(world, r.componentType)
		r.componentType = nil
	}
	return r.component
}

// targetEntity returns the target [Entity] of this RelationID.
func (r *Relation) targetEntity() Entity {
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
// Required a component ID lookup when used the first time.
// Un reuse, it is as fast as [RelIdx].
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

func (r relationSlice) ToRelations(world *World, mask *bitMask, ids []ID, out []relationID, copyTo []relationID) []relationID {
	if len(r) == 0 {
		return out
	}
	return r.toRelationsSlowPath(world, mask, ids, out, copyTo)
}

func (r relationSlice) toRelationsSlowPath(world *World, mask *bitMask, ids []ID, out []relationID, copyTo []relationID) []relationID {
	if copyTo != nil {
		copyTo = append(copyTo, out...)
		out = copyTo
	}
	// Fast special case for a single relation (20% speedup)
	if len(r) == 1 {
		rel := r[0]
		id := rel.id(ids, world)
		world.storage.checkRelationTarget(rel.targetEntity())
		world.storage.checkRelationComponent(id)
		if !mask.Get(id.id) {
			panic(fmt.Sprintf("requested relation component with ID %d was not specified in the filter or map", id.id))
		}
		return append(out, relationID{target: rel.target, component: id})
	}
	// Slower with loop for more than one relation
	for _, rel := range r {
		id := rel.id(ids, world)
		world.storage.checkRelationTarget(rel.targetEntity())
		world.storage.checkRelationComponent(id)
		if !mask.Get(id.id) {
			panic(fmt.Sprintf("requested relation component with ID %d was not specified in the filter or map", id.id))
		}
		out = append(out, relationID{target: rel.target, component: id})
	}
	return out
}

func (r relationSlice) ToRelationIDsForUnsafe(world *World, out []relationID) []relationID {
	for _, rel := range r {
		out = append(out, rel.relationIDForUnsafe(world))
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

func (r relationEntities) ToRelation(world *World, id ID, out []relationID) []relationID {
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
