package ecs

type entityID uint32

var entityType = typeOf[Entity]()

// Entity identifier.
type Entity struct {
	id  entityID // Entity ID
	gen uint32   // Entity generation
}

func newEntity(id entityID) Entity {
	return Entity{id, 0}
}

type entityIndex struct {
	table int
	row   int
}
