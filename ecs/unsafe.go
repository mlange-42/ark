package ecs

import "unsafe"

// Unsafe provides access to Ark's unsafe ID-based API.
// Get an instance via [World.Unsafe].
//
// The unsafe API is significantly slower than the type-safe API,
// and should only be used when component types are not known at compile time.
type Unsafe struct {
	world *World
}

// NewEntity creates a new entity with the given components.
func (u Unsafe) NewEntity(ids ...ID) Entity {
	return u.world.newEntityWith(ids, nil, nil)
}

// NewEntityRel creates a new entity with the given components and relation targets.
func (u Unsafe) NewEntityRel(ids []ID, relations ...RelationID) Entity {
	return u.world.newEntityWith(ids, nil, relations)
}

// Get returns a pointer to the given component of an [Entity].
//
// ⚠️ Important: The obtained pointer should not be stored persistently!
//
// Panics if the entity does not have the given component.
// Panics when called for a removed (and potentially recycled) entity.
func (u Unsafe) Get(entity Entity, comp ID) unsafe.Pointer {
	return u.world.storage.get(entity, comp)
}

// GetUnchecked returns a pointer to the given component of an [Entity].
// In contrast to [Unsafe.Get], it does not check whether the entity is alive.
//
// ⚠️ Important: The obtained pointer should not be stored persistently!
//
// Panics if the entity does not have the given component.
func (u Unsafe) GetUnchecked(entity Entity, comp ID) unsafe.Pointer {
	return u.world.storage.getUnchecked(entity, comp)
}

// Has returns whether an [Entity] has the given component.
//
// Panics when called for a removed (and potentially recycled) entity.
func (u Unsafe) Has(entity Entity, comp ID) bool {
	return u.world.storage.has(entity, comp)
}

// HasUnchecked returns whether an [Entity] has the given component.
// In contrast to [Unsafe.Has], it does not check whether the entity is alive.
//
// Panics when called for a removed (and potentially recycled) entity.
func (u Unsafe) HasUnchecked(entity Entity, comp ID) bool {
	return u.world.storage.hasUnchecked(entity, comp)
}

// GetRelation returns the relation target for the entity and the mapped component.
func (u Unsafe) GetRelation(entity Entity, comp ID) Entity {
	return u.world.storage.getRelation(entity, comp)
}

// GetRelationUnchecked returns the relation target for the entity and the mapped component.
// In contrast to [Unsafe.GetRelation], it does not check whether the entity is alive.
// Can be used as an optimization when it is certain that the entity is alive.
func (u Unsafe) GetRelationUnchecked(entity Entity, comp ID) Entity {
	return u.world.storage.getRelationUnchecked(entity, comp)
}

// SetRelations sets relation targets for an entity.
func (u Unsafe) SetRelations(entity Entity, relations ...RelationID) {
	u.world.setRelations(entity, relations)
}

// Add the given components to an entity.
func (u Unsafe) Add(entity Entity, comp ...ID) {
	if !u.world.Alive(entity) {
		panic("can't add components to a dead entity")
	}
	u.world.exchange(entity, comp, nil, nil, nil)
}

// AddRel adds the given components and relation targets to an entity.
func (u Unsafe) AddRel(entity Entity, comps []ID, relations ...RelationID) {
	if !u.world.Alive(entity) {
		panic("can't add components to a dead entity")
	}
	u.world.exchange(entity, comps, nil, nil, relations)
}

// Remove the given components from an entity.
func (u Unsafe) Remove(entity Entity, comp ...ID) {
	if !u.world.Alive(entity) {
		panic("can't remove components from a dead entity")
	}
	u.world.exchange(entity, nil, comp, nil, nil)
}

// Exchange the given components on entity.
func (u Unsafe) Exchange(entity Entity, add []ID, remove []ID, relations ...RelationID) {
	if !u.world.Alive(entity) {
		panic("can't exchange components on a dead entity")
	}
	u.world.exchange(entity, add, remove, nil, relations)
}

// IDs returns all component IDs of an entity.
func (u Unsafe) IDs(entity Entity) IDs {
	if !u.world.Alive(entity) {
		panic("can't get component IDs of a dead entity")
	}
	index := u.world.storage.entities[entity.id]
	return newIDs(u.world.storage.tables[index.table].ids)
}

// DumpEntities dumps entity information into an [EntityDump] object.
// This dump can be used with [Unsafe.LoadEntities] to set the World's entity state.
//
// For world serialization with components and resources, see module [github.com/mlange-42/ark-serde].
func (u Unsafe) DumpEntities() EntityDump {
	alive := []uint32{}

	filter := NewFilter(u.world)
	query := filter.Query()
	for query.Next() {
		alive = append(alive, uint32(query.Entity().id))
	}

	data := EntityDump{
		Entities:  append([]Entity{}, u.world.storage.entityPool.entities...),
		Alive:     alive,
		Next:      uint32(u.world.storage.entityPool.next),
		Available: u.world.storage.entityPool.available,
	}

	return data
}

// LoadEntities resets all entities to the state saved with [Unsafe.DumpEntities].
//
// Use this only on an empty world! Can be used after [World.Reset].
//
// The resulting world will have the same entities (in terms of ID, generation and alive state)
// as the original world. This is necessary for proper serialization of entity relations.
// However, the entities will not have any components.
//
// Panics if the world has any dead or alive entities.
//
// For world serialization with components and resources, see module [github.com/mlange-42/ark-serde].
func (u Unsafe) LoadEntities(data *EntityDump) {
	u.world.checkLocked()

	if len(u.world.storage.entityPool.entities) > 2 || u.world.storage.entityPool.available > 0 {
		panic("can set entity data only on a fresh or reset world")
	}

	capacity := len(data.Entities)

	entities := make([]Entity, 0, capacity)
	entities = append(entities, data.Entities...)

	if len(data.Entities) > 0 {
		u.world.storage.entityPool = entityPool{
			entities:  entities,
			next:      entityID(data.Next),
			available: data.Available,
			reserved:  entityID(reservedEntities),
		}
		u.world.storage.entityPool.pointer = unsafe.Pointer(&u.world.storage.entityPool.entities[0])
	}

	u.world.storage.entities = make([]entityIndex, len(data.Entities), capacity)
	u.world.storage.isTarget = make([]bool, len(data.Entities), capacity)

	table := &u.world.storage.tables[0]
	for _, idx := range data.Alive {
		entity := u.world.storage.entityPool.entities[idx]
		tableIdx := table.Add(entity)
		u.world.storage.entities[entity.id] = entityIndex{table: table.id, row: tableIdx}
	}
}
