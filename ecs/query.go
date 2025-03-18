package ecs

import "unsafe"

// UnsafeQuery is an unsafe query.
// It is significantly slower than type-safe generic queries like [Query2],
// and should only be used when component types are not known at compile time.
type UnsafeQuery struct {
	world     *World
	table     *table
	relations []RelationID
	tables    []tableID
	filter    filter
	cursor    cursor
	lock      uint8
}

// Next advances the query's cursor to the next entity.
func (q *UnsafeQuery) Next() bool {
	q.cursor.checkQueryNext()
	if int64(q.cursor.index) < q.cursor.maxIndex {
		q.cursor.index++
		return true
	}
	return q.nextTableOrArchetype()
}

// Entity returns the current entity.
func (q *UnsafeQuery) Entity() Entity {
	q.cursor.checkQueryGet()
	return q.table.GetEntity(q.cursor.index)
}

// Get returns the queried components of the current entity.
//
// ⚠️ Do not store the obtained pointer outside of the current context (i.e. the query loop)!
func (q *UnsafeQuery) Get(comp ID) unsafe.Pointer {
	q.cursor.checkQueryGet()
	return q.table.Get(comp, uintptr(q.cursor.index))
}

// Has returns whether the current entity has the given component.
func (q *UnsafeQuery) Has(comp ID) bool {
	return q.table.Has(comp)
}

// GetRelation returns the entity relation target of the component at the given index.
func (q *UnsafeQuery) GetRelation(comp ID) Entity {
	return q.table.GetRelation(comp)
}

// Count returns the number of entities matching this query.
func (q *UnsafeQuery) Count() int {
	return countQuery(&q.world.storage, &q.filter, q.relations)
}

// IDs returns the IDs of all component of the current [Entity]n.
func (q *UnsafeQuery) IDs() IDs {
	return newIDs(q.table.ids)
}

// Close closes the Query and unlocks the world.
//
// Automatically called when iteration completes.
// Needs to be called only if breaking out of the query iteration or not iterating at all.
func (q *UnsafeQuery) Close() {
	q.cursor.archetype = -2
	q.cursor.table = -2
	q.tables = nil
	q.table = nil
	q.world.unlock(q.lock)
}

func (q *UnsafeQuery) nextTableOrArchetype() bool {
	if q.cursor.archetype >= 0 && q.nextTable() {
		return true
	}
	return q.nextArchetype()
}

func (q *UnsafeQuery) nextArchetype() bool {
	maxArchIndex := int32(len(q.world.storage.archetypes) - 1)
	for q.cursor.archetype < maxArchIndex {
		q.cursor.archetype++
		archetype := &q.world.storage.archetypes[q.cursor.archetype]
		if !q.filter.matches(archetype.mask) {
			continue
		}

		if !archetype.HasRelations() {
			table := &q.world.storage.tables[archetype.tables[0]]
			if table.Len() > 0 {
				q.setTable(0, table)
				return true
			}
			continue
		}

		q.tables = archetype.GetTables(q.relations)
		q.cursor.table = -1
		if q.nextTable() {
			return true
		}
	}
	q.Close()
	return false
}

func (q *UnsafeQuery) nextTable() bool {
	maxTableIndex := int32(len(q.tables) - 1)
	for q.cursor.table < maxTableIndex {
		q.cursor.table++
		table := &q.world.storage.tables[q.tables[q.cursor.table]]
		if table.Len() == 0 {
			continue
		}
		if !table.Matches(q.relations) {
			continue
		}
		q.setTable(q.cursor.table, table)
		return true
	}
	return false
}

func (q *UnsafeQuery) setTable(index int32, table *table) {
	q.cursor.table = index
	q.table = table
	q.cursor.index = 0
	q.cursor.maxIndex = int64(q.table.Len() - 1)
}
