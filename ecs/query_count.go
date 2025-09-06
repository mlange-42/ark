package ecs

func countQueryCache(storage *storage, cache *cacheEntry, relations []RelationID) int {
	count := 0
	for _, tableID := range cache.tables {
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

func countQuery(storage *storage, filter *filter, relations []RelationID) int {
	count := 0
	for arch := range storage.archetypes {
		archetype := &storage.archetypes[arch]
		if !filter.matches(archetype.mask) {
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

func countQueryComponent(storage *storage, filter *filter, relations []RelationID, rareComp uint8) int {
	count := 0
	archetypes := storage.archetypesMap[rareComp]
	for _, arch := range archetypes {
		archetype := &storage.archetypes[arch]
		if !filter.matches(archetype.mask) {
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
