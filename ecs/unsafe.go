package ecs

import "unsafe"

// Unsafe provides access to Ark's unsafe ID-based API.
// Get an instance via [World.Unsafe].
//
// The unsafe API is significantly slower than the type-safe API,
// and should only be used when component types are not known at compile time.
type Unsafe struct {
	world           *World
	cachedRelations []relationID
}

// NewEntity creates a new entity with the given components.
func (u Unsafe) NewEntity(ids ...ID) Entity {
	entity, mask := u.world.newEntity(ids, nil)
	u.world.storage.observers.FireCreateEntityIfHas(entity, mask)
	return entity
}

// NewEntityRel creates a new entity with the given components and relation targets.
func (u Unsafe) NewEntityRel(ids []ID, relations ...Relation) Entity {
	u.cachedRelations = relationSlice(relations).ToRelationIDsForUnsafe(u.world, u.cachedRelations[:0])
	entity, mask := u.world.newEntity(ids, u.cachedRelations)
	u.world.storage.observers.FireCreateEntityIfHas(entity, mask)
	if len(relations) > 0 {
		u.world.storage.observers.FireCreateEntityRelIfHas(entity, mask)
	}
	return entity
}

// Get returns a pointer to the given component of an [Entity].
//
// ⚠️ Do not store the obtained pointer outside of the current context!
//
// Panics if the entity does not have the given component.
// Panics when called for a removed (and potentially recycled) entity.
func (u Unsafe) Get(entity Entity, comp ID) unsafe.Pointer {
	return u.world.storage.get(entity, comp)
}

// GetUnchecked returns a pointer to the given component of an [Entity].
// In contrast to [Unsafe.Get], it does not check whether the entity is alive.
//
// ⚠️ Do not store the obtained pointer outside of the current context!
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
func (u Unsafe) SetRelations(entity Entity, relations ...Relation) {
	u.cachedRelations = relationSlice(relations).ToRelationIDsForUnsafe(u.world, u.cachedRelations[:0])
	u.world.setRelations(entity, u.cachedRelations)
}

// Add the given components to an entity.
func (u Unsafe) Add(entity Entity, comp ...ID) {
	if !u.world.Alive(entity) {
		panic("can't add components to a dead entity")
	}
	oldMask, newMask := u.world.add(entity, comp, nil)
	u.world.storage.observers.FireAddIfHas(OnAddComponents, entity, oldMask, newMask)
}

// AddRel adds the given components and relation targets to an entity.
func (u Unsafe) AddRel(entity Entity, comps []ID, relations ...Relation) {
	if !u.world.Alive(entity) {
		panic("can't add components to a dead entity")
	}
	u.cachedRelations = relationSlice(relations).ToRelationIDsForUnsafe(u.world, u.cachedRelations[:0])
	oldMask, newMask := u.world.add(entity, comps, u.cachedRelations)
	u.world.storage.observers.FireAddIfHas(OnAddComponents, entity, oldMask, newMask)
	if len(relations) > 0 {
		u.world.storage.observers.FireAddIfHas(OnAddRelations, entity, oldMask, newMask)
	}
}

// Remove the given components from an entity.
func (u Unsafe) Remove(entity Entity, comp ...ID) {
	if !u.world.Alive(entity) {
		panic("can't remove components from a dead entity")
	}
	u.world.remove(entity, comp)
}

// Exchange the given components on entity.
func (u Unsafe) Exchange(entity Entity, add []ID, remove []ID, relations ...Relation) {
	if !u.world.Alive(entity) {
		panic("can't exchange components on a dead entity")
	}
	u.cachedRelations = relationSlice(relations).ToRelationIDsForUnsafe(u.world, u.cachedRelations[:0])
	oldMask, newMask := u.world.exchange(entity, add, remove, u.cachedRelations)

	if len(add) > 0 {
		u.world.storage.observers.FireAddIfHas(OnAddComponents, entity, oldMask, newMask)
		if len(relations) > 0 {
			u.world.storage.observers.FireAddIfHas(OnAddRelations, entity, oldMask, newMask)
		}
	}
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
	filter := NewFilter0(u.world)
	query := filter.Query()
	alive := make([]uint32, 0, query.Count())
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

	entities := make([]Entity, capacity)
	copy(entities, data.Entities)

	if capacity > 0 {
		u.world.storage.entityPool = entityPool{
			entities:  entities,
			next:      entityID(data.Next),
			available: data.Available,
			reserved:  entityID(reservedEntities),
		}
		u.world.storage.entityPool.pointer = unsafe.Pointer(&u.world.storage.entityPool.entities[0])
	}

	u.world.storage.entities = make([]entityIndex, capacity)
	u.world.storage.isTarget = make([]bool, capacity)

	table := &u.world.storage.tables[0]
	table.Extend(uint32(len(data.Alive)), u.world.storage.components)
	for _, idx := range data.Alive {
		entity := u.world.storage.entityPool.entities[idx]
		tableIdx := table.Add(entity, u.world.storage.components)
		u.world.storage.entities[entity.id] = entityIndex{table: table.id, row: tableIdx}
	}
}
