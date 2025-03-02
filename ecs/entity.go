package ecs

import "encoding/json"

type entityID uint32

var entityType = typeOf[Entity]()
var entitySize = sizeOf(entityType)

var wildcard = Entity{1, 0}

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

// MarshalJSON returns a JSON representation of the entity, for serialization purposes.
//
// The JSON representation of an entity is a two-element array of entity ID and generation.
func (e Entity) MarshalJSON() ([]byte, error) {
	arr := [2]uint32{uint32(e.id), e.gen}
	jsonValue, _ := json.Marshal(arr) // Ignore the error, as we can be sure this works.
	return jsonValue, nil
}

// UnmarshalJSON into an entity.
//
// For serialization purposes only. Do not use this to create entities!
func (e *Entity) UnmarshalJSON(data []byte) error {
	arr := [2]uint32{}
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	e.id = entityID(arr[0])
	e.gen = arr[1]

	return nil
}

func (e Entity) toRelation(id ID, out []RelationID) []RelationID {
	out = out[:0]
	out = append(out, RelationID{
		component: id,
		target:    Entity(e),
	})
	return out
}

// entityIndex denotes an entity's location by table and row index.
type entityIndex struct {
	table tableID
	row   uint32
}
