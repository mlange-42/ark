package ecs

type entityID uint32

var entityType = typeOf[Entity]()
var entitySize = sizeOf(entityType)

var wildcard = Entity{1, 0}

// Wildcard entity for relation targets.
func Wildcard() Entity {
	return wildcard
}

// Entity identifier.
type Entity struct {
	id  entityID // Entity ID
	gen uint32   // Entity generation
}

func newEntity(id entityID) Entity {
	return Entity{id, 0}
}

// IsZero returns whether this entity is the reserved zero entity.
func (e Entity) IsZero() bool {
	return e.id == 0
}

// isWildcard returns whether this entity is the reserved wildcard entity.
func (e Entity) isWildcard() bool {
	return e.id == 1
}

type entityIndex struct {
	table tableID
	row   uint32
}

type entities []entityIndex
