package ecs

var relationType = typeOf[RelationMarker]()

// RelationMarker is a marker for entity relation components.
// It must be embedded as first field of a component that represent an entity relation
// (see the example).
//
// Entity relations allow for fast queries using entity relationships.
// E.g. to iterate over all entities that are the child of a certain parent entity.
type RelationMarker struct{}

// Relation is the common interface for specifying relation targets.
// It is implemented by [RelationIndex], [RelationType] and [RelationID].
//
//   - [RelationType] is safe, but has some run-time overhead for component [ID] lookup.
//   - [RelationIndex] is fast but less safe.
//   - [RelationID] is used in the [Unsafe] API.
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

// RelID creates a new [RelationID].
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

// RelationType specifies an entity relation target by component type.
// Create with [Rel].
//
// It can be used as a safer but slower alternative to [RelationIndex].
type RelationType[C any] struct {
	target Entity
}

// Rel creates a new [RelationType].
func Rel[C any](target Entity) RelationType[C] {
	return RelationType[C]{
		target: target,
	}
}

// id returns the component ID of this RelationID.
func (r RelationType[C]) id(ids []ID, world *World) ID {
	return ComponentID[C](world)
}

// targetEntity returns the target [Entity] of this RelationID.
func (r RelationType[C]) targetEntity() Entity {
	return r.target
}

// RelationIndex specifies an entity relation target by component index.
// Create with [RelIdx].
//
// It can be used as faster but less safe alternative to [RelationType].
//
// Note that the index refers to the position of the component in the generics
// of e.g. a [Map2] or [Filter2].
// This should not be confused with component [ID] as obtained by [ComponentID]!
type RelationIndex struct {
	index  uint8
	target Entity
}

// RelIdx creates a new [RelationIndex].
func RelIdx(index int, target Entity) RelationIndex {
	return RelationIndex{
		index:  uint8(index),
		target: target,
	}
}

// id returns the component ID of this RelationIndex.
func (r RelationIndex) id(ids []ID, world *World) ID {
	return ids[r.index]
}

// targetEntity returns the target [Entity] of this RelationIndex.
func (r RelationIndex) targetEntity() Entity {
	return r.target
}

// Helper for converting relations
type relations []Relation

func (r relations) toRelations(world *World, ids []ID, base []RelationID, out []RelationID) []RelationID {
	out = out[:0]
	out = append(out, base...)
	for _, rel := range r {
		id := rel.id(ids, world)
		world.storage.checkRelationTarget(rel.targetEntity())
		world.storage.checkRelationComponent(id)
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
