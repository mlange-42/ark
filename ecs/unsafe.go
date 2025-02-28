package ecs

import "unsafe"

// Unsafe provides access to Ark's unsafe ID-based API.
// Get an instance via [World.Unsafe].
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

// TODO: Unsafe.NewEntity
// TODO: Unsafe.Add
// TODO: Unsafe.Remove
// TODO: Unsafe.Exchange
// TODO: Unsafe.GetRelation
// TODO: Unsafe.SetRelation
// TODO: Queries
