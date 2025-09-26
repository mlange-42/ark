package ecs

// Map is a mapper to access and manipulate components of an entity.
//
// Instances should be created during initialization and stored, e.g. in systems.
type Map[T any] struct {
	mask      bitMask
	world     *World
	storage   *componentStorage
	relations []relationID
	id        ID
	ids       [1]ID
}

// NewMap creates a new [Map].
func NewMap[T any](w *World) *Map[T] {
	id := ComponentID[T](w)
	return &Map[T]{
		mask:    newMask(id),
		world:   w,
		id:      id,
		ids:     [1]ID{id},
		storage: &w.storage.components[id.id],
	}
}

// NewEntity creates a new entity with the mapped component.
//
// If the mapped component is a relationship (see [RelationMarker]),
// a relation target entity must be provided.
func (m *Map[T]) NewEntity(comp *T, target ...Entity) Entity {
	return m.NewEntityFn(func(a *T) {
		*a = *comp
	}, target...)
}

// NewEntityFn creates a new entity with the mapped component and runs a callback instead of using a component for initialization.
// The callback can be nil.
//
// If the mapped component is a relationship (see [RelationMarker]),
// a relation target entity must be provided.
//
// ⚠️ Do not store the obtained pointer outside of the current context!
func (m *Map[T]) NewEntityFn(fn func(*T), target ...Entity) Entity {
	m.relations = relationEntities(target).toRelation(m.world, m.id, m.relations)
	entity, mask := m.world.newEntity(m.ids[:], m.relations)
	if fn != nil {
		fn(m.GetUnchecked(entity))
	}
	m.world.storage.observers.FireCreateEntity(entity, mask)
	return entity
}

// NewBatch creates a batch of new entities with the mapped component.
//
// If the mapped component is a relationship (see [RelationMarker]),
// a relation target entity must be provided.
func (m *Map[T]) NewBatch(count int, comp *T, target ...Entity) {
	m.NewBatchFn(count, func(entity Entity, a *T) {
		*a = *comp
	}, target...)
}

// NewBatchFn creates a batch of new entities with the mapped components, running the given initializer function on each.
// The initializer function can be nil.
//
// If the mapped component is a relationship (see [RelationMarker]),
// a relation target entity must be provided.
//
// ⚠️ Do not store the obtained pointers outside of the current context!
func (m *Map[T]) NewBatchFn(count int, fn func(Entity, *T), target ...Entity) {
	m.relations = relationEntities(target).toRelation(m.world, m.id, m.relations)
	tableID, start := m.world.newEntities(count, m.ids[:], m.relations)

	if fn != nil {
		table := &m.world.storage.tables[tableID]
		column := m.storage.columns[tableID]

		lock := m.world.lock()
		for i := range count {
			index := uintptr(start + i)
			fn(
				table.GetEntity(index),
				(*T)(column.Get(index)),
			)
		}
		m.world.unlock(lock)
	}

	if m.world.storage.observers.HasObservers(OnCreateEntity) {
		table := &m.world.storage.tables[tableID]
		earlyOut := true
		lock := m.world.lock()
		for i := range count {
			index := uintptr(start + i)
			if !m.world.storage.observers.doFireCreateEntity(table.GetEntity(index), &m.mask, earlyOut) {
				break
			}
			earlyOut = false
		}
		m.world.unlock(lock)
	}
}

// Get returns the mapped component for the given entity.
//
// Returns nil if the entity does not have the mapped component.
//
// ⚠️ Do not store the obtained pointer outside of the current context!
func (m *Map[T]) Get(entity Entity) *T {
	if !m.world.Alive(entity) {
		panic("can't get a component of a dead entity")
	}
	index := m.world.storage.entities[entity.id]
	column := m.storage.columns[index.table]
	if column == nil {
		return nil
	}
	return (*T)(column.Get(uintptr(index.row)))
}

// GetUnchecked returns the mapped component for the given entity.
// In contrast to [Map.Get], it does not check whether the entity is alive.
// Can be used as an optimization when it is certain that the entity is alive.
//
// Returns nil if the entity does not have the mapped component.
//
// ⚠️ Do not store the obtained pointer outside of the current context!
func (m *Map[T]) GetUnchecked(entity Entity) *T {
	index := m.world.storage.entities[entity.id]
	column := m.storage.columns[index.table]
	if column == nil {
		return nil
	}
	return (*T)(column.Get(uintptr(index.row)))
}

// Has return whether the given entity has the mapped component.
//
// Using [Map.Get] and checking for nil pointer may be faster
// than calling [Map.Has] and [Map.Get] subsequently.
func (m *Map[T]) Has(entity Entity) bool {
	if !m.world.Alive(entity) {
		panic("can't get a component of a dead entity")
	}
	return m.HasUnchecked(entity)
}

// HasUnchecked return whether the given entity has the mapped component.
// In contrast to [Map.Has], it does not check whether the entity is alive.
// Can be used as an optimization when it is certain that the entity is alive.
func (m *Map[T]) HasUnchecked(entity Entity) bool {
	index := m.world.storage.entities[entity.id]
	return m.storage.columns[index.table] != nil
}

// Add the mapped component to the given entity.
//
// If the mapped component is a relationship (see [RelationMarker]),
// a relation target entity must be provided.
func (m *Map[T]) Add(entity Entity, comp *T, target ...Entity) {
	m.AddFn(entity, func(a *T) {
		*a = *comp
	}, target...)
}

// AddFn adds the mapped component to the given entity and runs a callback instead of using a component for initialization.
// The callback can be nil.
//
// If the mapped component is a relationship (see [RelationMarker]),
// a relation target entity must be provided.
//
// ⚠️ Do not store the obtained pointer outside of the current context!
func (m *Map[T]) AddFn(entity Entity, fn func(*T), target ...Entity) {
	if !m.world.Alive(entity) {
		panic("can't add a component to a dead entity")
	}
	m.relations = relationEntities(target).toRelation(m.world, m.id, m.relations)
	oldMask, newMask := m.world.exchange(entity, m.ids[:], nil, m.relations)
	if fn != nil {
		fn(m.GetUnchecked(entity))
	}

	m.world.storage.observers.FireAdd(entity, oldMask, newMask)
}

// Set the mapped component of the given entity to the given values.
// Requires the entity to already have the mapped component.
//
// Triggers an OnSetComponents event, which is not the case when assigning by
// pointer dereference or assigning to fields. This can be used to notify observers
// about component updates.
//
// This is not a component operation, so it can be performed on a locked world.
func (m *Map[T]) Set(entity Entity, comp *T) {
	if !m.world.Alive(entity) {
		panic("can't set component of a dead entity")
	}
	m.world.storage.checkHasComponent(entity, m.ids[0])

	index := &m.world.storage.entities[entity.id]
	*(*T)(m.storage.columns[index.table].Get(uintptr(index.row))) = *comp

	if m.world.storage.observers.HasObservers(OnSetComponents) {
		newMask := &m.world.storage.archetypes[m.world.storage.tables[index.table].archetype].mask
		m.world.storage.observers.doFireSet(entity, &m.mask, newMask)
	}
}

// AddBatch adds the mapped component to all entities matching the given batch filter.
//
// If the mapped component is a relationship (see [RelationMarker]),
// a relation target entity must be provided.
func (m *Map[T]) AddBatch(batch *Batch, comp *T, target ...Entity) {
	m.AddBatchFn(batch, func(_ Entity, a *T) {
		*a = *comp
	}, target...)
}

// AddBatchFn adds the mapped component to all entities matching the given batch filter,
// running the given function on each. The function can be nil.
//
// If the mapped component is a relationship (see [RelationMarker]),
// a relation target entity must be provided.
//
// ⚠️ Do not store the obtained pointers outside of the current context!
func (m *Map[T]) AddBatchFn(batch *Batch, fn func(Entity, *T), target ...Entity) {
	m.relations = relationEntities(target).toRelation(m.world, m.id, m.relations)

	var process func(tableID tableID, start, len uint32)
	if fn != nil {
		process = func(tableID tableID, start, len uint32) {
			table := &m.world.storage.tables[tableID]
			column := m.storage.columns[tableID]

			for i := range len {
				index := uintptr(start + i)
				fn(
					table.GetEntity(index),
					(*T)(column.Get(index)),
				)
			}
		}
	}
	m.world.exchangeBatch(batch, m.ids[:], nil, m.relations, process)
}

// Remove the mapped component from the given entity.
func (m *Map[T]) Remove(entity Entity) {
	if !m.world.Alive(entity) {
		panic("can't remove a component from a dead entity")
	}
	m.world.exchange(entity, nil, m.ids[:], nil)
}

// RemoveBatch removes the mapped component from all entities matching the given batch filter,
// running the given function on each. The function can be nil.
func (m *Map[T]) RemoveBatch(batch *Batch, fn func(entity Entity)) {
	removeBatch(m.world, batch, m.ids[:], fn)
}

// GetRelation returns the relation target for the entity and the mapped component.
func (m *Map[T]) GetRelation(entity Entity) Entity {
	return m.world.storage.getRelation(entity, m.id)
}

// GetRelationUnchecked returns the relation target for the entity and the mapped component.
// In contrast to [Map.GetRelation], it does not check whether the entity is alive.
// Can be used as an optimization when it is certain that the entity is alive.
func (m *Map[T]) GetRelationUnchecked(entity Entity) Entity {
	return m.world.storage.getRelationUnchecked(entity, m.id)
}

// SetRelation sets the relation target for the entity and the mapped component.
func (m *Map[T]) SetRelation(entity Entity, target Entity) {
	m.relations = target.toRelation(m.world, m.id, m.relations)
	m.world.setRelations(entity, m.relations)
}

// SetRelationBatch sets the relation target for all entities matching the given batch filter.
func (m *Map[T]) SetRelationBatch(batch *Batch, target Entity, fn func(entity Entity)) {
	m.relations = target.toRelation(m.world, m.id, m.relations)
	setRelationsBatch(m.world, batch, fn, m.relations)
}
