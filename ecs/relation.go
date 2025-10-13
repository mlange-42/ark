package ecs

import (
	"fmt"
	"reflect"
)

// relationType is the runtime type of RelationMarker
var relationType = reflect.TypeFor[RelationMarker]()

// RelationMarker is a marker for entity relation components.
// It must be embedded as first field of a component that represent an entity relationship
// (see the example).
//
// Entity relations allow for fast queries using entity relationships.
// E.g. to iterate over all entities that are the child of a certain parent entity.
type RelationMarker struct{}

// relationID is a pair of relation component type and relation target.
type relationID struct {
	target    Entity
	component ID
}

// Relation is the common type for specifying relationship targets.
// It can be created with [Rel], [RelIdx] and [RelID].
//
//   - [Rel] uses a generic type parameter to identify the component.
//     It is safe, but has some run-time overhead for component [ID] lookup on first usage.
//   - [RelIdx] uses an index to identify the component. It is fast but more error-prone.
//   - [RelID] uses a component ID. It is for use with the [Unsafe] API.
type Relation struct {
	componentType reflect.Type // Component type of the relation
	target        Entity       // Target entity of the relation
	component     ID           // Component ID of the relation
	index         uint8        // Component index of the relation in a mapper or query
}

// Rel creates a new [Relation] for a component type.
//
// It can be used as a safer but slower alternative to [RelIdx].
// Requires a component ID lookup when used the first time.
// On reuse, it is as fast as [RelIdx].
func Rel[C any](target Entity) Relation {
	return Relation{
		target:        target,
		componentType: reflect.TypeFor[C](),
		index:         255,
	}
}

// RelIdx creates a new [Relation] for a component index.
// The index refers to the position of the component in the generics
// of e.g. a [Map2] or [Filter2].
// For filters, components specified by [Filter2.With] are also covered by the index.
//
// It can be used as faster but less safe alternative to [Rel].
//
// Note that the index should not be confused with a component [ID] as obtained by [ComponentID]!
// For component IDs, use [RelID].
func RelIdx(index int, target Entity) Relation {
	return Relation{
		index:  uint8(index),
		target: target,
	}
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

// relationIDForUnsafe converts the Relation to a relationID.
//
// Relation must use ID or component type.
// Panics if used with an index Relation.
//
// Modifies the Relation to use an ID.
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

// id returns the component ID of this Relation.
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

// targetEntity returns the target [Entity] of this Relation.
func (r *Relation) targetEntity() Entity {
	return r.target
}

// Helper for converting relationSlice
type relationSlice []Relation

// ToRelations converts a slice of Relation items to relationIDs.
func (r relationSlice) ToRelations(world *World, mask *bitMask, ids []ID, out []relationID, copy bool) []relationID {
	if len(r) == 0 {
		return out
	}
	return r.toRelationsSlowPath(world, mask, ids, out, copy)
}

// toRelationsSlowPath is the slow path of ToRelations for more than zero relations.
func (r relationSlice) toRelationsSlowPath(world *World, mask *bitMask, ids []ID, out []relationID, copy bool) []relationID {
	if copy {
		temp := make([]relationID, 0, len(r)+len(out))
		temp = append(temp, out...)
		out = temp
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

// ToRelationIDsForUnsafe converts a slice of Relation items from the unsafe API to relationIDs.
func (r relationSlice) ToRelationIDsForUnsafe(world *World, out []relationID) []relationID {
	for _, rel := range r {
		out = append(out, rel.relationIDForUnsafe(world))
	}
	return out
}

// toRelation converts an entity and a component ID to relationIDs.
func toRelation(world *World, e Entity, id ID, out []relationID) []relationID {
	world.storage.checkRelationTarget(e)
	world.storage.checkRelationComponent(id)
	out = out[:0]
	out = append(out, relationID{target: e, component: id})
	return out
}

// Helper for converting relations
type relationEntities []Entity

// ToRelation converts a slice of entities to relationIDs.
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
