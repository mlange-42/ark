package ecs

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"
)

// entityID is the type for entity IDs.
type entityID uint32

// runtime type of entities
var entityType = reflect.TypeFor[Entity]()

// memory size of an Entity
var entitySize = entityType.Size()

// memory size of an entityIndex
var entityIndexSize = reflect.TypeFor[entityIndex]().Size()

//var wildcard = Entity{1, 0}

// Entity is an identifier for entities.
//
// Can be stored safely in components, resources or elsewhere.
// For stored entities, it may be necessary to check their alive status with [World.Alive].
//
// ⚠️ Always store entities by value, never by pointer!
//
// In Ark, entities are returned to a pool when they are removed from the world.
// These entities can be recycled, with the same ID ([Entity.ID]), but an incremented generation ([Entity.Gen]).
// This allows to determine whether an entity hold by the user is still alive, despite it was potentially recycled.
type Entity struct {
	id  entityID // Entity ID
	gen uint32   // Entity generation
}

// newEntity creates a new entity with the given ID.
func newEntity(id entityID) Entity {
	return Entity{id, 0}
}

// ID returns the entity's ID, primarily for debugging purposes.
func (e Entity) ID() uint32 {
	return uint32(e.id)
}

// Gen returns the entity's generation, primarily for debugging purposes.
func (e Entity) Gen() uint32 {
	return e.gen
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

// MarshalBinary returns a binary representation of the entity, for serialization and networking purposes.
func (e *Entity) MarshalBinary() []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint32(buf[0:4], uint32(e.id))
	binary.BigEndian.PutUint32(buf[4:8], e.gen)
	return buf
}

// UnmarshalBinary into an entity.
//
// For serialization and networking purposes only. Do not use this to create entities!
func (e *Entity) UnmarshalBinary(data []byte) error {
	if len(data) != 8 {
		return fmt.Errorf("invalid data length: expected 8 bytes, got %d", len(data))
	}
	e.id = entityID(binary.BigEndian.Uint32(data[0:4]))
	e.gen = binary.BigEndian.Uint32(data[4:8])
	return nil
}

// entityIndex denotes an entity's location by table and row index.
type entityIndex struct {
	table tableID
	row   uint32
}
