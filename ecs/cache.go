package ecs

import "math"

type cacheID uint32

const maxCacheID = math.MaxUint32

// Cache entry for a filter.
type cacheEntry struct {
	filter    *filter         // The underlying filter.
	indices   map[tableID]int // Map of table indices for removal.
	relations []RelationID    // Entity relationships.
	tables    []tableID       // Tables matching the filter.
	id        cacheID         // Entry ID.
}

// cache provides filter caching to speed up queries.
//
// For registered filters, the relevant archetypes are tracked internally,
// so that there are no mask checks required during iteration.
// This is particularly helpful to avoid query iteration slowdown by a very high number of archetypes.
// If the number of archetypes exceeds approx. 50-100, uncached filters experience a slowdown.
// The relative slowdown increases with lower numbers of entities queried (noticeable below a few thousand entities).
// Cached filters avoid this slowdown.
//
// The overhead of tracking cached filters internally is very low, as updates are required only when new archetypes are created.
type cache struct {
	indices map[cacheID]int  // Mapping from filter IDs to indices in filters
	filters []cacheEntry     // The cached filters, indexed by indices
	intPool intPool[cacheID] // Pool for filter IDs
}

// newCache creates a new [cache].
func newCache() cache {
	return cache{
		intPool: newIntPool[cacheID](128),
		indices: map[cacheID]int{},
		filters: []cacheEntry{},
	}
}

func (c *cache) getEntry(id cacheID) *cacheEntry {
	return &c.filters[c.indices[id]]
}

// Register a filter.
func (c *cache) register(storage *storage, filter *filter, relations []RelationID) cacheID {
	// TODO: prevent duplicate registration
	id := c.intPool.Get()
	index := len(c.filters)
	c.filters = append(c.filters,
		cacheEntry{
			id:        id,
			filter:    filter,
			relations: relations,
			tables:    storage.getTableIDs(filter, relations),
			indices:   nil,
		})
	c.indices[id] = index
	return id
}

func (c *cache) unregister(id cacheID) {
	idx, ok := c.indices[id]
	if !ok {
		panic("no filter for id found to unregister")
	}
	delete(c.indices, id)

	last := len(c.filters) - 1
	if idx != last {
		c.filters[idx], c.filters[last] = c.filters[last], c.filters[idx]
		c.indices[c.filters[idx].id] = idx
	}
	c.filters[last] = cacheEntry{}
	c.filters = c.filters[:last]
}

// Adds a table.
//
// Iterates over all filters and adds the node to the resp. entry where the filter matches.
func (c *cache) addTable(storage *storage, table *table) {
	arch := &storage.archetypes[table.archetype]
	if !table.HasRelations() {
		for i := range c.filters {
			e := &c.filters[i]
			if !e.filter.matches(arch.mask) {
				continue
			}
			e.tables = append(e.tables, table.id)
		}
		return
	}

	for i := range c.filters {
		e := &c.filters[i]
		if !e.filter.matches(arch.mask) {
			continue
		}
		if !table.Matches(e.relations) {
			continue
		}
		e.tables = append(e.tables, table.id)
		if e.indices != nil {
			e.indices[table.id] = len(e.tables) - 1
		}
	}
}

// Removes a table.
//
// Can only be used for tables that have a relation target.
// Tables without a relation are never removed.
func (c *cache) removeTable(storage *storage, table *table) {
	//if !table.HasRelations() {
	//	// unreachable
	//	return
	//}
	arch := &storage.archetypes[table.archetype]
	for i := range c.filters {
		e := &c.filters[i]

		if e.indices == nil && e.filter.matches(arch.mask) {
			c.mapTables(storage, e)
		}

		if idx, ok := e.indices[table.id]; ok {
			lastIndex := len(e.tables) - 1
			swapped := idx != lastIndex
			if swapped {
				e.tables[idx], e.tables[lastIndex] = e.tables[lastIndex], e.tables[idx]
			}
			e.tables = e.tables[:lastIndex]
			if swapped {
				e.indices[e.tables[idx]] = idx
			}
			delete(e.indices, table.id)
		}
	}
}

func (c *cache) mapTables(storage *storage, e *cacheEntry) {
	e.indices = map[tableID]int{}
	for i, tableID := range e.tables {
		table := &storage.tables[tableID]
		if table.HasRelations() {
			e.indices[tableID] = i
		}
	}
}

func (c *cache) Reset() {
	for i := range c.filters {
		c.filters[i].tables = nil
	}
	c.indices = map[cacheID]int{}
	c.filters = c.filters[:0]
	c.intPool.Reset()
}
