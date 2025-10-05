package ecs

import (
	"fmt"
	"time"
	"unsafe"
)

type storage struct {
	entities           []entityIndex      // Entity positions in archetypes, indexed by entity ID
	isTarget           []bool             // Whether each entity is a target of a relationship
	graph              graph              // Graph for fast archetype traversal
	slices             slices             // Slices for internal re-use
	archetypes         []archetype        // All archetypes
	allArchetypes      []archetypeID      // list of all archetype IDs to simplify usage of componentIndex
	componentIndex     [][]archetypeID    // Archetypes indexed by components IDs; each archetype appears under all its component IDs
	relationArchetypes []archetypeID      // All archetypes with relationships
	tables             []table            // All tables
	components         []componentStorage // Component storages for fast random/world access
	cache              cache              // Filter cache
	entityPool         entityPool         // Entity pool for creation and recycling
	registry           componentRegistry  // Component registry
	locks              lock               // World locks
	observers          observerManager    // Observer/event manager
	config             config             // Storage configuration (initial capacities)
}

type componentStorage struct {
	columns []*column
}

type slices struct {
	relations []relationID
	batches   []batchTable
	tables    []tableID
	ints      []uint32

	entities        []Entity
	entitiesCleanup []Entity
}

func newSlices() slices {
	return slices{
		relations: make([]relationID, 0, 8),
		batches:   make([]batchTable, 0, 32),
		tables:    make([]tableID, 0, 32),
		ints:      make([]uint32, 0, 32),

		entities:        make([]Entity, 0, 16),
		entitiesCleanup: make([]Entity, 0, 256),
	}
}

func newStorage(numArchetypes int, capacity ...int) storage {
	config := newConfig(capacity...)

	reg := newComponentRegistry()
	entities := make([]entityIndex, reservedEntities, config.initialCapacity+reservedEntities)
	isTarget := make([]bool, reservedEntities, config.initialCapacity+reservedEntities)
	// Reserved zero and wildcard entities
	for i := range reservedEntities {
		entities[i] = entityIndex{table: maxTableID, row: 0}
	}
	componentsMap := make([]int16, maskTotalBits)
	for i := range maskTotalBits {
		componentsMap[i] = -1
	}

	archetypes := make([]archetype, 0, numArchetypes)
	archetypes = append(archetypes, newArchetype(0, 0, &bitMask{}, nil, []tableID{0}, &reg))
	tables := make([]table, 0, numArchetypes)
	tables = append(tables, newTable(0, &archetypes[0], uint32(config.initialCapacity), &reg, nil, nil))
	return storage{
		config:         config,
		registry:       reg,
		locks:          newLock(),
		observers:      newObserverManager(),
		cache:          newCache(),
		entities:       entities,
		isTarget:       isTarget,
		entityPool:     newEntityPool(uint32(config.initialCapacity), reservedEntities),
		graph:          newGraph(),
		slices:         newSlices(),
		archetypes:     archetypes,
		allArchetypes:  []archetypeID{0},
		componentIndex: make([][]archetypeID, 0, maskTotalBits),
		tables:         tables,
		components:     make([]componentStorage, 0, maskTotalBits),
	}
}

func (s *storage) findOrCreateTable(oldTable *table, add []ID, remove []ID, relations []relationID, outMask *bitMask) (*table, *archetype, bool) {
	startNode := s.archetypes[oldTable.archetype].node

	node := s.graph.Find(startNode, add, remove, outMask)
	var arch *archetype
	if archID, ok := node.GetArchetype(); ok {
		arch = &s.archetypes[archID]
	} else {
		arch = s.createArchetype(node)
		node.archetype = arch.id
	}

	relationRemoved := false
	// TODO: this could also use pooling
	var allRelations []relationID
	if len(remove) > 0 {
		// filter out removed relations
		allRelations = make([]relationID, 0, len(oldTable.relationIDs)+len(relations))
		for _, rel := range oldTable.relationIDs {
			if arch.mask.Get(rel.component.id) {
				allRelations = append(allRelations, rel)
			} else {
				relationRemoved = true
			}
		}
		allRelations = append(allRelations, relations...)
	} else {
		if len(relations) > 0 {
			allRelations = copyAppend(oldTable.relationIDs, relations...)
		} else {
			allRelations = oldTable.relationIDs
		}
	}
	table, ok := arch.GetTable(s, allRelations)
	if !ok {
		table = s.createTable(arch, allRelations)
	}
	return table, arch, relationRemoved
}

func (s *storage) findOrCreateTableAdd(oldTable *table, add []ID, relations []relationID, outMask *bitMask) (*table, *archetype) {
	startNode := s.archetypes[oldTable.archetype].node

	node := s.graph.FindAdd(startNode, add, outMask)
	var arch *archetype
	if archID, ok := node.GetArchetype(); ok {
		arch = &s.archetypes[archID]
	} else {
		arch = s.createArchetype(node)
		node.archetype = arch.id
	}

	// TODO: this could also use pooling
	var allRelations []relationID
	if len(relations) > 0 {
		allRelations = copyAppend(oldTable.relationIDs, relations...)
	} else {
		allRelations = oldTable.relationIDs
	}

	table, ok := arch.GetTable(s, allRelations)
	if !ok {
		table = s.createTable(arch, allRelations)
	}
	return table, arch
}

func (s *storage) findOrCreateTableRemove(oldTable *table, remove []ID, outMask *bitMask) (*table, *archetype, bool) {
	startNode := s.archetypes[oldTable.archetype].node

	node := s.graph.FindRemove(startNode, remove, outMask)
	var arch *archetype
	if archID, ok := node.GetArchetype(); ok {
		arch = &s.archetypes[archID]
	} else {
		arch = s.createArchetype(node)
		node.archetype = arch.id
	}

	relationRemoved := false
	// TODO: this could also use pooling
	allRelations := make([]relationID, 0, len(oldTable.relationIDs))
	for _, rel := range oldTable.relationIDs {
		if arch.mask.Get(rel.component.id) {
			allRelations = append(allRelations, rel)
		} else {
			relationRemoved = true
		}
	}
	table, ok := arch.GetTable(s, allRelations)
	if !ok {
		table = s.createTable(arch, allRelations)
	}
	return table, arch, relationRemoved
}

func (s *storage) AddComponent(id uint8) {
	if len(s.components) != int(id) {
		panic("components can only be added to a storage sequentially")
	}
	s.components = append(s.components, componentStorage{columns: make([]*column, len(s.tables))})
	s.componentIndex = append(s.componentIndex, []archetypeID{})
}

// RemoveEntity removes the given entity from the world.
func (s *storage) RemoveEntity(entity Entity) {
	if !s.entityPool.Alive(entity) {
		panic("can't remove a dead entity")
	}
	index := &s.entities[entity.id]
	table := &s.tables[index.table]

	hasEntityObs := s.observers.HasObservers(OnRemoveEntity)
	hasRelationObs := table.HasRelations() && s.observers.HasObservers(OnRemoveRelations)
	if hasEntityObs || hasRelationObs {
		l := s.lock()
		if hasEntityObs {
			mask := &s.archetypes[table.archetype].mask
			s.observers.FireRemoveEntity(entity, mask, true)
		}
		if hasRelationObs {
			mask := &s.archetypes[table.archetype].mask
			s.observers.FireRemoveEntityRel(entity, mask, true)
		}
		s.unlock(l)
	}

	swapped := table.Remove(index.row)

	s.entityPool.Recycle(entity)

	if swapped {
		swapEntity := table.GetEntity(uintptr(index.row))
		s.entities[swapEntity.id].row = index.row
	}
	index.table = maxTableID

	if s.isTarget[entity.id] {
		s.cleanupArchetypes(entity)
		s.isTarget[entity.id] = false
	}
}

func (s *storage) Reset() {
	s.entities = s.entities[:reservedEntities]
	s.entityPool.Reset()
	s.isTarget = s.isTarget[:reservedEntities]
	s.cache.Reset()
	s.locks.Reset()
	s.observers.Reset()

	for i := range s.archetypes {
		s.archetypes[i].Reset(s)
	}
}

func (s *storage) get(entity Entity, component ID) unsafe.Pointer {
	if !s.entityPool.Alive(entity) {
		panic("can't get component of a dead entity")
	}
	return s.getUnchecked(entity, component)
}

func (s *storage) getUnchecked(entity Entity, component ID) unsafe.Pointer {
	s.checkHasComponent(entity, component)
	index := s.entities[entity.id]
	return s.tables[index.table].Get(component, uintptr(index.row))
}

func (s *storage) has(entity Entity, component ID) bool {
	if !s.entityPool.Alive(entity) {
		panic("can't check component of a dead entity")
	}
	return s.hasUnchecked(entity, component)
}

func (s *storage) hasUnchecked(entity Entity, component ID) bool {
	index := s.entities[entity.id]
	return s.tables[index.table].Has(component)
}

func (s *storage) getRelation(entity Entity, comp ID) Entity {
	if !s.entityPool.Alive(entity) {
		panic("can't get relation for a dead entity")
	}
	s.checkHasComponent(entity, comp)
	return s.tables[s.entities[entity.id].table].GetRelation(comp)
}

func (s *storage) getRelationUnchecked(entity Entity, comp ID) Entity {
	s.checkHasComponent(entity, comp)
	return s.tables[s.entities[entity.id].table].GetRelation(comp)
}

func (s *storage) registerTargets(relations []relationID) {
	for _, rel := range relations {
		s.isTarget[rel.target.id] = true
	}
}

func (s *storage) registerFilter(filter *filter, relations []relationID) {
	s.cache.register(s, filter, relations)
}

func (s *storage) unregisterFilter(filter *filter) {
	s.cache.unregister(filter)
}

func (s *storage) getRegisteredFilter(id cacheID) *cacheEntry {
	return s.cache.getEntry(id)
}

func (s *storage) createEntity(table tableID) (Entity, uint32) {
	entity := s.entityPool.Get()

	idx := s.tables[table].Add(entity)
	if int(entity.id) == len(s.entities) {
		s.entities = append(s.entities, entityIndex{table: table, row: idx})
		s.isTarget = append(s.isTarget, false)
	} else {
		s.entities[entity.id] = entityIndex{table: table, row: idx}
	}
	return entity, idx
}

func (s *storage) createEntities(table *table, count int) {
	startIdx := table.Len()
	table.Alloc(uint32(count))

	len := len(s.entities)
	for i := range count {
		index := uint32(startIdx + i)
		entity := s.entityPool.Get()
		table.SetEntity(index, entity)

		if int(entity.id) == len {
			s.entities = append(s.entities, entityIndex{table: table.id, row: index})
			s.isTarget = append(s.isTarget, false)
			len++
		} else {
			s.entities[entity.id] = entityIndex{table: table.id, row: index}
		}
	}
}

func (s *storage) createArchetype(node *node) *archetype {
	comps := node.mask.toTypes(&s.registry.registry)
	index := len(s.archetypes)
	s.archetypes = append(s.archetypes, newArchetype(archetypeID(index), node.id, &node.mask, comps, nil, &s.registry))
	archetype := &s.archetypes[index]

	s.allArchetypes = append(s.allArchetypes, archetype.id)
	for _, id := range archetype.components {
		s.componentIndex[id.id] = append(s.componentIndex[id.id], archetype.id)
		s.registry.addArchetype(id.id)
	}
	if archetype.HasRelations() {
		s.relationArchetypes = append(s.relationArchetypes, archetype.id)
	}

	return archetype
}

func (s *storage) createTable(archetype *archetype, relations []relationID) *table {
	targets := make([]Entity, len(archetype.components))

	if uint8(len(relations)) < archetype.numRelations {
		panic("relation targets must be fully specified")
	}
	for _, rel := range relations {
		idx := archetype.componentsMap[rel.component.id]
		targets[idx] = rel.target
	}
	for i := range relations {
		rel := &relations[i]
		s.checkRelationComponent(rel.component)
		s.checkRelationTarget(rel.target)
	}

	var newTableID tableID
	recycled := false
	if id, ok := archetype.GetFreeTable(); ok {
		newTableID = id
		s.tables[newTableID].Recycle(targets, relations)
		recycled = true
	} else {
		newTableID = tableID(len(s.tables))
		cap := s.config.initialCapacity
		if archetype.HasRelations() {
			cap = s.config.initialCapacityRelations
		}
		s.tables = append(s.tables, newTable(
			newTableID, archetype, uint32(cap), &s.registry,
			targets, relations))
	}
	archetype.AddTable(&s.tables[newTableID])

	table := &s.tables[newTableID]
	if !recycled {
		for i := range s.components {
			id := ID{id: uint8(i)}
			comps := &s.components[i]
			if archetype.mask.Get(id.id) {
				comps.columns = append(comps.columns, table.GetColumn(id))
			} else {
				comps.columns = append(comps.columns, nil)
			}
		}
	}

	s.cache.addTable(s, table)
	return table
}

// Removes empty archetypes that have a target relation to the given entity.
func (s *storage) cleanupArchetypes(target Entity) {
	newRelations := s.slices.relations
	for _, arch := range s.relationArchetypes {
		archetype := &s.archetypes[arch]
		len := len(archetype.tables.tables)
		for i := len - 1; i >= 0; i-- {
			table := &s.tables[archetype.tables.tables[i]]

			foundTarget := false
			for _, rel := range table.relationIDs {
				if rel.target.id == target.id {
					newRelations = append(newRelations, relationID{component: rel.component, target: Entity{}})
					foundTarget = true
				}
			}
			if !foundTarget {
				continue
			}

			if table.Len() > 0 {
				allRelations := s.getExchangeTargetsUnchecked(table, newRelations)
				newTable, ok := archetype.GetTable(s, allRelations)
				if !ok {
					newTable = s.createTable(archetype, allRelations)
					// Get the old table again, as pointers may have changed.
					table = &s.tables[table.id]
				}
				s.moveEntities(table, newTable, uint32(table.Len()))
			}
			archetype.FreeTable(table)
			s.cache.removeTable(table)

			newRelations = newRelations[:0]
		}
		archetype.RemoveTarget(target)
	}
	s.slices.relations = newRelations[:0]
}

// moveEntities moves all entities from src to dst.
func (s *storage) moveEntities(src, dst *table, count uint32) {
	oldLen := dst.Len()
	dst.AddAll(src, count)

	newLen := dst.Len()
	newTable := dst.id
	for i := oldLen; i < newLen; i++ {
		entity := dst.GetEntity(uintptr(i))
		s.entities[entity.id] = entityIndex{table: newTable, row: uint32(i)}
	}
	src.Reset()
}

func (s *storage) getExchangeTargetsUnchecked(oldTable *table, relations []relationID) []relationID {
	targets := s.slices.entities
	for i := range oldTable.columns {
		targets = append(targets, oldTable.columns[i].target)
	}
	for _, rel := range relations {
		column := oldTable.components[rel.component.id]
		// TODO: check this!
		// As rel.target is always the zero entity, and the zero entity can't be removed,
		// this should not be possible.
		//if rel.target == targets[column.index] {
		//	continue
		//}
		targets[column.index] = rel.target
	}

	result := make([]relationID, 0, len(oldTable.relationIDs))
	for i, e := range targets {
		if !oldTable.columns[i].isRelation {
			continue
		}
		id := oldTable.ids[i]
		result = append(result, relationID{component: id, target: e})
	}
	s.slices.entities = targets[:0]
	return result
}

func (s *storage) getExchangeTargets(oldTable *table, relations []relationID, mask *bitMask) ([]relationID, bool) {
	changed := false
	targets := s.slices.entities
	for i := range oldTable.columns {
		targets = append(targets, oldTable.columns[i].target)
	}
	for _, rel := range relations {
		// Validity of the target is checked when creating a new table.
		// Whether the component is a relation is checked when creating a new table.
		column := oldTable.components[rel.component.id]
		if column == nil {
			tp, _ := s.registry.ComponentType(rel.component.id)
			panic(fmt.Sprintf("entity has no component of type %s to set relation target for", tp.Name()))
		}
		if rel.target == targets[column.index] {
			continue
		} else if mask != nil {
			mask.Set(rel.component.id)
		}
		targets[column.index] = rel.target
		changed = true
	}
	if !changed {
		return nil, false
	}

	result := make([]relationID, 0, len(oldTable.relationIDs))
	for i, e := range targets {
		if !oldTable.columns[i].isRelation {
			continue
		}
		id := oldTable.ids[i]
		result = append(result, relationID{component: id, target: e})
	}
	s.slices.entities = targets[:0]
	return result, true
}

// the returned slice comes from the pool and should be recycled.
func (s *storage) getTables(batch *Batch) []tableID {
	tables := s.slices.tables

	if batch.filter.cache != maxCacheID {
		cache := s.getRegisteredFilter(batch.filter.cache)
		for _, tableID := range cache.tables.tables {
			table := &s.tables[tableID]
			if table.Len() == 0 {
				continue
			}
			if !table.Matches(batch.relations) {
				continue
			}
			tables = append(tables, tableID)
		}
		return tables
	}

	for i := range s.archetypes {
		archetype := &s.archetypes[i]
		if !batch.filter.matches(&archetype.mask) {
			continue
		}

		if !archetype.HasRelations() {
			table := &s.tables[archetype.tables.tables[0]]
			tables = append(tables, table.id)
			continue
		}

		tableIDs := archetype.GetTables(batch.relations)
		for _, tab := range tableIDs {
			table := &s.tables[tab]
			if !table.Matches(batch.relations) {
				continue
			}
			tables = append(tables, tab)
		}
	}
	return tables
}

func (s *storage) getTableIDs(filter *filter, relations []relationID) []tableID {
	tables := []tableID{}

	for i := range s.archetypes {
		archetype := &s.archetypes[i]
		if !filter.matches(&archetype.mask) {
			continue
		}

		if !archetype.HasRelations() {
			tables = append(tables, archetype.tables.tables[0])
			continue
		}

		tableIDs := archetype.GetTables(relations)
		for _, tab := range tableIDs {
			table := &s.tables[tab]
			if !table.Matches(relations) {
				continue
			}
			tables = append(tables, tab)
		}
	}
	return tables
}

// Shrink reduces memory usage by shrinking the capacity of archetype tables.
// See [World.Shrink] for details.
func (s *storage) Shrink(stopAfter time.Duration) bool {
	start := time.Now()
	var tableIdx int
	anyFound := false
	for tableIdx = range s.tables {
		table := &s.tables[tableIdx]

		if !table.HasRelations() {
			if table.Shrink(uint32(s.config.initialCapacity)) {
				anyFound = true
			}
		} else {
			if table.Shrink(uint32(s.config.initialCapacityRelations)) {
				anyFound = true
			}
			if !table.isFree && table.Len() == 0 {
				s.archetypes[table.archetype].FreeTable(table)
				anyFound = true
			}
		}

		if anyFound && (stopAfter == 0 || time.Since(start) >= stopAfter) {
			break
		}
	}

	tableIdx++
	for tableIdx < len(s.tables) {
		table := &s.tables[tableIdx]

		if !table.HasRelations() {
			if table.CanShrink(uint32(s.config.initialCapacity)) {
				return true
			}
		} else {
			if table.CanShrink(uint32(s.config.initialCapacityRelations)) {
				return true
			}
			if !table.isFree && table.Len() == 0 {
				return true
			}
		}

		tableIdx++
	}

	return false
}

// lock the world and get the lock bit for later unlocking.
func (s *storage) lock() uint8 {
	return s.locks.Lock()
}

// unlock unlocks the given lock bit.
func (s *storage) unlock(l uint8) {
	s.locks.Unlock(l)
}
