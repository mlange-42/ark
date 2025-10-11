package ecs

import "fmt"

// entityAtCache return the entity at a specific index in a query, for cached queries.
func entityAtCache(storage *storage, cache *cacheEntry, relations []relationID, index uint32) Entity {
	count := uint32(0)
	for _, tableID := range cache.tables.tables {
		table := &storage.tables[tableID]
		if table.Len() == 0 {
			continue
		}
		if !table.Matches(relations) {
			continue
		}
		len := uint32(table.Len())
		if count+len > index {
			return table.GetEntity(uintptr(index - count))
		}
		count += len
	}
	panic(fmt.Sprintf("entity index %d out of bounds for query with %d entities", index, count))
}

// entityAtCache return the entity at a specific index in a query, for uncached queries.
func entityAt(storage *storage, filter *filter, relations []relationID, archetypes []archetypeID, index uint32) Entity {
	count := uint32(0)
	for _, arch := range archetypes {
		archetype := &storage.archetypes[arch]
		if !filter.matches(&archetype.mask) {
			continue
		}

		if !archetype.HasRelations() {
			table := &storage.tables[archetype.tables.tables[0]]
			len := uint32(table.Len())
			if count+len > index {
				return table.GetEntity(uintptr(index - count))
			}
			count += len
			continue
		}

		tables := archetype.GetTables(relations)
		for _, tab := range tables {
			table := &storage.tables[tab]
			if !table.Matches(relations) {
				continue
			}
			len := uint32(table.Len())
			if count+len > index {
				return table.GetEntity(uintptr(index - count))
			}
			count += len
		}
	}
	panic(fmt.Sprintf("entity index %d out of bounds for query with %d entities", index, count))
}

// countQueryCache returns the number of entities in a query, for cached queries.
func countQueryCache(storage *storage, cache *cacheEntry, relations []relationID) int {
	count := 0
	for _, tableID := range cache.tables.tables {
		table := &storage.tables[tableID]
		if table.Len() == 0 {
			continue
		}
		if !table.Matches(relations) {
			continue
		}
		count += table.Len()
	}
	return count
}

// countQueryCache returns the number of entities in a query, for uncached queries.
func countQuery(storage *storage, filter *filter, relations []relationID, archetypes []archetypeID) int {
	count := 0
	for _, arch := range archetypes {
		archetype := &storage.archetypes[arch]
		if !filter.matches(&archetype.mask) {
			continue
		}

		if !archetype.HasRelations() {
			table := &storage.tables[archetype.tables.tables[0]]
			count += table.Len()
			continue
		}

		tables := archetype.GetTables(relations)
		for _, tab := range tables {
			table := &storage.tables[tab]
			if !table.Matches(relations) {
				continue
			}
			count += table.Len()
		}
	}
	return count
}
