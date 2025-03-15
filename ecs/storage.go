package ecs

import (
	"unsafe"
)

type storage struct {
	entities           []entityIndex
	isTarget           []bool
	graph              graph
	archetypes         []archetype
	relationArchetypes []archetypeID
	tables             []table
	components         []componentStorage
	cache              cache
	entityPool         entityPool
	registry           componentRegistry
	config             config
}

type componentStorage struct {
	columns []*column
}

func newStorage(capacity ...int) storage {
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

	archetypes := make([]archetype, 0, 128)
	archetypes = append(archetypes, newArchetype(0, 0, &bitMask{}, []ID{}, []tableID{0}, &reg))
	tables := make([]table, 0, 128)
	tables = append(tables, newTable(0, &archetypes[0], uint32(config.initialCapacity), &reg, []Entity{}, []RelationID{}))
	return storage{
		config:     config,
		registry:   reg,
		cache:      newCache(),
		entities:   entities,
		isTarget:   isTarget,
		entityPool: newEntityPool(uint32(config.initialCapacity), reservedEntities),
		graph:      newGraph(),
		archetypes: archetypes,
		tables:     tables,
		components: make([]componentStorage, 0, maskTotalBits),
	}
}

func (s *storage) findOrCreateTable(oldTable *table, add []ID, remove []ID, relations []RelationID, outMask *bitMask) *table {
	startNode := s.archetypes[oldTable.archetype].node

	node := s.graph.Find(startNode, add, remove, outMask)
	var arch *archetype
	if archID, ok := node.GetArchetype(); ok {
		arch = &s.archetypes[archID]
	} else {
		arch = s.createArchetype(node)
		node.archetype = arch.id
	}

	var allRelations []RelationID
	if len(relations) > 0 {
		allRelations = appendNew(oldTable.relationIDs, relations...)
	} else {
		allRelations = oldTable.relationIDs
	}
	table, ok := arch.GetTable(s, allRelations)
	if !ok {
		table = s.createTable(arch, allRelations)
	}
	return table
}

func (s *storage) AddComponent(id uint8) {
	if len(s.components) != int(id) {
		panic("components can only be added to a storage sequentially")
	}
	s.components = append(s.components, componentStorage{columns: make([]*column, len(s.tables))})
}

// RemoveEntity removes the given entity from the world.
func (s *storage) RemoveEntity(entity Entity) {
	if !s.entityPool.Alive(entity) {
		panic("can't remove a dead entity")
	}
	index := &s.entities[entity.id]
	table := &s.tables[index.table]

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
		panic("can't get component of a dead entity")
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
	return s.getRelationUnchecked(entity, comp)
}

func (s *storage) getRelationUnchecked(entity Entity, comp ID) Entity {
	s.checkHasComponent(entity, comp)
	return s.tables[s.entities[entity.id].table].GetRelation(comp)
}

func (s *storage) registerTargets(relations []RelationID) {
	for _, rel := range relations {
		s.isTarget[rel.target.id] = true
	}
}

func (s *storage) registerFilter(filter *filter, relations []RelationID) cacheID {
	return s.cache.register(s, filter, relations)
}

func (s *storage) unregisterFilter(entry cacheID) {
	s.cache.unregister(entry)
}

func (s *storage) getRegisteredFilter(id cacheID) *cacheEntry {
	return s.cache.getEntry(id)
}

func (s *storage) createEntity(table tableID) (Entity, uint32) {
	entity := s.entityPool.Get()

	idx := s.tables[table].Add(entity)
	len := len(s.entities)
	if int(entity.id) == len {
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
	if archetype.HasRelations() {
		s.relationArchetypes = append(s.relationArchetypes, archetype.id)
	}
	return archetype
}

func (s *storage) createTable(archetype *archetype, relations []RelationID) *table {
	targets := make([]Entity, len(archetype.components))
	numRelations := uint8(0)
	for _, rel := range relations {
		idx := archetype.componentsMap[rel.component.id]
		targets[idx] = rel.target
		numRelations++
	}
	if numRelations != archetype.numRelations {
		panic("relations must be fully specified")
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
		s.tables[newTableID].recycle(targets, relations)
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
			if archetype.mask.Get(id) {
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
	newRelations := []RelationID{}
	for _, arch := range s.relationArchetypes {
		archetype := &s.archetypes[arch]
		len := len(archetype.tables)
		for i := len - 1; i >= 0; i-- {
			table := &s.tables[archetype.tables[i]]

			foundTarget := false
			for _, rel := range table.relationIDs {
				if rel.target.id == target.id {
					newRelations = append(newRelations, RelationID{component: rel.component, target: Entity{}})
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
					newTable = s.createTable(archetype, newRelations)
				}
				s.moveEntities(table, newTable, uint32(table.Len()))
			}
			archetype.FreeTable(table.id)
			s.cache.removeTable(s, table)

			newRelations = newRelations[:0]
		}
		archetype.RemoveTarget(target)
	}
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

func (s *storage) getExchangeTargetsUnchecked(oldTable *table, relations []RelationID) []RelationID {
	targets := make([]Entity, len(oldTable.columns))
	for i := range oldTable.columns {
		targets[i] = oldTable.columns[i].target
	}
	for _, rel := range relations {
		column := oldTable.components[rel.component.id]
		if rel.target == targets[column.index] {
			continue
		}
		targets[column.index] = rel.target
	}

	result := make([]RelationID, 0, len(oldTable.relationIDs))
	for i, e := range targets {
		if !oldTable.columns[i].isRelation {
			continue
		}
		id := oldTable.ids[i]
		result = append(result, RelationID{component: id, target: e})
	}

	return result
}

func (s *storage) getExchangeTargets(oldTable *table, relations []RelationID) ([]RelationID, bool) {
	changed := false
	targets := make([]Entity, len(oldTable.columns))
	for i := range oldTable.columns {
		targets[i] = oldTable.columns[i].target
	}
	for _, rel := range relations {
		// Validity of the target is checked when creating a new table.
		// Whether the component is a relation is checked when creating a new table.
		column := oldTable.components[rel.component.id]
		if rel.target == targets[column.index] {
			continue
		}
		targets[column.index] = rel.target
		changed = true
	}
	if !changed {
		return nil, false
	}

	result := make([]RelationID, 0, len(oldTable.relationIDs))
	for i, e := range targets {
		if !oldTable.columns[i].isRelation {
			continue
		}
		id := oldTable.ids[i]
		result = append(result, RelationID{component: id, target: e})
	}

	return result, true
}

func (s *storage) getTables(batch *Batch) []*table {
	tables := []*table{}

	if batch.cache != maxCacheID {
		cache := s.getRegisteredFilter(batch.cache)
		for _, tableID := range cache.tables {
			table := &s.tables[tableID]
			if table.Len() == 0 {
				continue
			}
			if !table.Matches(batch.relations) {
				continue
			}
			tables = append(tables, table)
		}
		return tables
	}

	for i := range s.archetypes {
		archetype := &s.archetypes[i]
		if !batch.filter.matches(&archetype.mask) {
			continue
		}

		if !archetype.HasRelations() {
			table := &s.tables[archetype.tables[0]]
			tables = append(tables, table)
			continue
		}

		tableIDs := archetype.GetTables(batch.relations)
		for _, tab := range tableIDs {
			table := &s.tables[tab]
			if !table.Matches(batch.relations) {
				continue
			}
			tables = append(tables, table)
		}
	}
	return tables
}

func (s *storage) getTableIDs(filter *filter, relations []RelationID) []tableID {
	tables := []tableID{}

	for i := range s.archetypes {
		archetype := &s.archetypes[i]
		if !filter.matches(&archetype.mask) {
			continue
		}

		if !archetype.HasRelations() {
			tables = append(tables, archetype.tables[0])
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
