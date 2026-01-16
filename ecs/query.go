package ecs

// UnsafeQuery is an unsafe query.
// It is significantly slower than type-safe generic queries like [Query2],
// and should only be used when component types are not known at compile time.
type UnsafeQuery struct {
	world     *World
	table     *table
	relations []relationID
	tables    []tableID
	filter    filter
	cursor    cursor
	lock      uint8
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
	return countQuery(&q.world.storage, &q.filter, q.relations, q.world.storage.allArchetypes)
}

// EntityAt returns the entity at a given index.
//
// The method is particularly useful for random sampling of entities from a query.
// However, performance depends on the number of archetypes in the world and in the query.
// In worlds with many archetypes, it is recommended to use a registered/cached filter.
//
// Do not use this to iterate a query! Use [Query.Next] instead.
//
// Panics if the index is out of range, as indicated by [Query.Count].
func (q *UnsafeQuery) EntityAt(index int) Entity {
	return entityAt(&q.world.storage, &q.filter, q.relations, q.world.storage.allArchetypes, uint32(index))
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
	if q.cursor.table < -1 {
		return
	}
	q.cursor.archetype = -2
	q.cursor.table = -2
	q.tables = nil
	q.table = nil
	q.world.unlockSafe(q.lock)
}

func (q *UnsafeQuery) nextTableOrArchetype() bool {
	if q.cursor.archetype >= 0 && q.nextTable() {
		return true
	}
	return q.nextArchetype()
}

func (q *UnsafeQuery) nextArchetype() bool {
	q.tables = nil
	maxArchIndex := int32(len(q.world.storage.archetypes) - 1)
	for q.cursor.archetype < maxArchIndex {
		q.cursor.archetype++
		archetype := &q.world.storage.archetypes[q.cursor.archetype]
		if !q.filter.matches(&archetype.mask) {
			continue
		}

		if !archetype.HasRelations() {
			table := &q.world.storage.tables[archetype.tables.tables[0]]
			if table.len > 0 {
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
		if table.len == 0 || !table.Matches(q.relations) {
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
	q.cursor.maxIndex = int64(q.table.len - 1)
}
