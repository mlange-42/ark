package ecs

import (
	"math"
	"reflect"

	"github.com/mlange-42/ark/ecs/stats"
)

// archetype ID/index type.
type archetypeID uint32

// maxArchetypeID is used as unassigned archetype ID.
const maxArchetypeID = math.MaxUint32

// archetype struct.
type archetype struct {
	*archetypeData
	componentsMap  []int16                  // mapping from component IDs to column indices; -1 indicates none
	relationTables []map[entityID]*tableIDs // lookup for relation targets of tables, indexed by column index
	tables         tableIDs                 // all active tables
	mask           bitMask                  // Bit mask for the archetype's components
	id             archetypeID              // ID of the archetype
	numRelations   uint8                    // number of relation components
}

type archetypeData struct {
	components   []ID                   // components IDs of the archetype in arbitrary order
	itemSizes    []uint32               // item size per component index
	isRelation   []bool                 // whether columns are relations components, indexed by column index
	freeTables   []tableID              // all inactive/free tables
	zeroValue    []byte                 // zero value with the size of the largest item type, for fast zeroing
	targetTables map[entityID]*tableIDs // all tables per target for cleanup
	node         nodeID                 // Node ID of the archetype
}

// tableIDs helper for faster search and remove operations.
type tableIDs struct {
	tables  []tableID          // List of tables
	indices map[tableID]uint32 // Mapping from table ID to the index in the list
}

// Creates a new tableIDs.
// The passed tables slice is used directly, so it should not be modified or stored afterwards.
func newTableIDs(tables ...tableID) tableIDs {
	indices := make(map[tableID]uint32, len(tables))
	for i, t := range tables {
		indices[t] = uint32(i)
	}
	return tableIDs{
		tables:  tables,
		indices: indices,
	}
}

// Append a table ID.
func (t *tableIDs) Append(id tableID) {
	t.tables = append(t.tables, id)
	t.indices[id] = uint32(len(t.tables) - 1)
}

// Remove a table ID.
func (t *tableIDs) Remove(id tableID) bool {
	index, ok := t.indices[id]
	if !ok {
		return false
	}

	last := uint32(len(t.tables) - 1)
	if index != last {
		t.tables[index], t.tables[last] = t.tables[last], t.tables[index]
		t.indices[t.tables[index]] = index
	}
	t.tables = t.tables[:last]
	delete(t.indices, id)

	return true
}

// Clear the list and the mapping.
func (t *tableIDs) Clear() {
	t.tables = t.tables[:0]
	t.indices = map[tableID]uint32{}
}

// newArchetype creates a new archetype.
func newArchetype(
	id archetypeID, node nodeID, mask *bitMask,
	components []ID, tables []tableID, reg *componentRegistry) (archetype, archetypeData) {
	componentsMap := make([]int16, maskTotalBits)
	for i := range maskTotalBits {
		componentsMap[i] = -1
	}

	sizes := make([]uint32, len(components))
	maxSize := entitySize
	for i, id := range components {
		componentsMap[id.id] = int16(i)
		tp := reg.Types[id.id]

		itemSize := tp.Size()
		sizes[i] = uint32(itemSize)
		if itemSize > maxSize {
			maxSize = itemSize
		}
	}
	var zeroValue []byte
	if maxSize > 0 {
		zeroValue = make([]byte, maxSize)
	}

	numRelations := uint8(0)
	isRelation := make([]bool, len(components))
	relationTables := make([]map[entityID]*tableIDs, len(components))
	for i, id := range components {
		if reg.IsRelation[id.id] {
			isRelation[i] = true
			relationTables[i] = map[entityID]*tableIDs{}
			numRelations++
		}
	}
	var targetTables map[entityID]*tableIDs
	if numRelations > 0 {
		targetTables = map[entityID]*tableIDs{}
	}
	archTables := newTableIDs(tables...)
	return archetype{
			id:             id,
			mask:           *mask,
			componentsMap:  componentsMap,
			tables:         archTables,
			numRelations:   numRelations,
			relationTables: relationTables,
		}, archetypeData{
			node:         node,
			components:   components,
			targetTables: targetTables,
			itemSizes:    sizes,
			isRelation:   isRelation,
			zeroValue:    zeroValue,
		}
}

// HasRelations returns whether the archetype has any relation components.
func (a *archetype) HasRelations() bool {
	return a.numRelations > 0
}

// GetTable returns the table that matches the given relations exactly,
// and whether there is a matching table.
//
// Relations must be fully specified.
func (a *archetype) GetTable(storage *storage, relations []relationID) (*table, bool) {
	if len(a.tables.tables) == 0 {
		return nil, false
	}
	if !a.HasRelations() {
		return &storage.tables[a.tables.tables[0]], true
	}
	return a.getTableSlowPath(storage, relations)
}

// slow path for GetTable for archetypes with relations.
func (a *archetype) getTableSlowPath(storage *storage, relations []relationID) (*table, bool) {
	if uint8(len(relations)) < a.numRelations {
		panic("relation targets must be fully specified")
	}
	index := a.componentsMap[relations[0].component.id]
	tables, ok := a.relationTables[index][relations[0].target.id]
	if !ok {
		return nil, false
	}
	for _, t := range tables.tables {
		table := &storage.tables[t]
		if table.MatchesExact(relations) {
			return table, true
		}
	}
	return nil, false
}

// GetTables return all tables matching the first given relation, if any.
// Otherwise, returns all tables of the archetype.
//
// Relations do not need to be fully specified.
func (a *archetype) GetTables(relations []relationID) []tableID {
	if !a.HasRelations() || len(relations) == 0 {
		return a.tables.tables
	}
	index := a.componentsMap[relations[0].component.id]
	if tables, ok := a.relationTables[index][relations[0].target.id]; ok {
		return tables.tables
	}
	return nil
}

// GetFreeTable returns a free/unused table ID of the archetype,
// and whether there is any free table.
func (a *archetype) GetFreeTable() (tableID, bool) {
	if len(a.freeTables) == 0 {
		return 0, false
	}
	last := len(a.freeTables) - 1
	table := a.freeTables[last]

	a.freeTables = a.freeTables[:last]

	return table, true
}

// FreeTable frees the given table.
//
// Removes the table from the archetype's list of tables
// and from the archetype's relation tables.
//
// Does not clear the table's content.
func (a *archetype) FreeTable(table *table) {
	_ = a.tables.Remove(table.id)
	a.freeTables = append(a.freeTables, table.id)
	table.isFree = true

	// If there is only one relation, the resp. relationTables
	// entry is removed anyway.
	if a.numRelations <= 1 {
		return
	}

	// TODO: can/should we be more selective here?
	// For a potential solution, see https://github.com/mlange-42/ark/pull/264
	for _, m := range a.relationTables {
		for _, v := range m {
			_ = v.Remove(table.id)
		}
	}
	for _, v := range a.targetTables {
		_ = v.Remove(table.id)
	}
}

// FreeAllTables frees all tables of the archetype.
//
// Does not clear the tables' contents.
func (a *archetype) FreeAllTables(storage *storage) {
	for _, table := range a.tables.tables {
		storage.tables[table].isFree = true
	}
	a.freeTables = append(a.freeTables, a.tables.tables...)
	a.tables.Clear()

	for i := range a.relationTables {
		a.relationTables[i] = map[entityID]*tableIDs{}
	}
	a.targetTables = map[entityID]*tableIDs{}
}

// AddTable adds the given table to the archetype's list of tables,
// and to the relation tables of there are any relations.
func (a *archetype) AddTable(table *table) {
	a.tables.Append(table.id)
	if !a.HasRelations() {
		return
	}

	for i := range table.ids {
		column := &table.columns[i]
		if !column.isRelation {
			continue
		}
		target := column.target
		relations := a.relationTables[i]

		if tables, ok := relations[target.id]; ok {
			tables.Append(table.id)
		} else {
			tables := newTableIDs(table.id)
			relations[target.id] = &tables
		}

		if tables, ok := a.targetTables[target.id]; ok {
			if _, ok := tables.indices[table.id]; !ok {
				tables.Append(table.id)
			}
		} else {
			tables := newTableIDs(table.id)
			a.targetTables[target.id] = &tables
		}
	}
}

// RemoveTarget removes a relation target entity from the archetype's relation tables.
func (a *archetype) RemoveTarget(entity Entity) {
	for i := range a.relationTables {
		if !a.isRelation[i] {
			continue
		}
		delete(a.relationTables[i], entity.id)
	}
	delete(a.targetTables, entity.id)
}

// Reset the archetype. Clears and frees all tables.
func (a *archetype) Reset(storage *storage) {
	if !a.HasRelations() {
		storage.tables[a.tables.tables[0]].Reset()
		return
	}

	for i := len(a.tables.tables) - 1; i >= 0; i-- {
		table := &storage.tables[a.tables.tables[i]]
		table.Reset()
	}

	a.FreeAllTables(storage)
}

// Stats generates statistics for an archetype.
func (a *archetype) Stats(storage *storage) stats.Archetype {
	ids := a.components
	aTypes := make([]reflect.Type, len(ids))
	aTypeNames := make([]string, len(ids))
	for j, id := range ids {
		tp, _ := storage.registry.ComponentType(id.id)
		aTypes[j] = tp
		aTypeNames[j] = tp.Name()
	}

	memPerEntity := int(entitySize)
	intIDs := make([]uint8, len(ids))
	for j, id := range ids {
		intIDs[j] = id.id
		memPerEntity += int(a.itemSizes[j])
	}

	cap := 0
	count := 0
	memory := 0
	memoryUsed := 0
	tableStats := make([]stats.Table, len(a.tables.tables))
	for i, id := range a.tables.tables {
		table := &storage.tables[id]
		tableStats[i] = table.Stats(memPerEntity)
		stats := &tableStats[i]
		cap += stats.Capacity
		count += stats.Size
		memory += stats.Memory
		memoryUsed += stats.MemoryUsed
	}
	for _, id := range a.freeTables {
		table := &storage.tables[id]
		cap += int(table.cap)
		memory += memPerEntity * int(table.cap)
	}

	return stats.Archetype{
		FreeTables:         len(a.freeTables),
		NumRelations:       int(a.numRelations),
		ComponentIDs:       intIDs,
		ComponentTypes:     aTypes,
		ComponentTypeNames: aTypeNames,
		Memory:             memory,
		MemoryUsed:         memoryUsed,
		MemoryPerEntity:    memPerEntity,
		Size:               count,
		Capacity:           cap,
		Tables:             tableStats,
	}
}

// UpdateStats updates statistics for an archetype.
func (a *archetype) UpdateStats(stats *stats.Archetype, storage *storage) {
	tables := a.tables

	cap := 0
	count := 0
	memory := 0
	memoryUsed := 0

	cntOld := int32(len(stats.Tables))
	cntNew := int32(len(tables.tables))
	if cntNew < cntOld {
		stats.Tables = stats.Tables[:cntNew]
		cntOld = cntNew
	}
	var i int32
	for i := range cntOld {
		tableStats := &stats.Tables[i]
		table := &storage.tables[tables.tables[i]]
		table.UpdateStats(stats.MemoryPerEntity, tableStats)
		cap += tableStats.Capacity
		count += tableStats.Size
		memory += tableStats.Memory
		memoryUsed += tableStats.MemoryUsed
	}
	for i = cntOld; i < cntNew; i++ {
		table := &storage.tables[tables.tables[i]]
		tableStats := table.Stats(stats.MemoryPerEntity)
		stats.Tables = append(stats.Tables, tableStats)
		cap += tableStats.Capacity
		count += tableStats.Size
		memory += tableStats.Memory
		memoryUsed += tableStats.MemoryUsed
	}
	for _, id := range a.freeTables {
		table := &storage.tables[id]
		cap += int(table.cap)
		memory += stats.MemoryPerEntity * int(table.cap)
	}

	stats.FreeTables = len(a.freeTables)
	stats.Capacity = cap
	stats.Size = count
	stats.Memory = memory
	stats.MemoryUsed = memoryUsed
}
