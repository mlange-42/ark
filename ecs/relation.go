package ecs

var relationType = typeOf[Relation]()

// Relation is a marker for entity relation components.
// It must be embedded as first field of a component that represent an entity relation
// (see the example).
//
// Entity relations allow for fast queries using entity relationships.
// E.g. to iterate over all entities that are the child of a certain parent entity.
// Currently, each entity can only have a single relation component.
type Relation struct{}

type relation struct {
	component ID
	target    Entity
}
