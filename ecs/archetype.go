package ecs

import (
	"math"
	"reflect"
	"slices"

	"github.com/mlange-42/ark/ecs/stats"
)

type archetypeID uint32

// maxArchetypeID is used as unassigned archetype ID.
const maxArchetypeID = math.MaxUint32

type archetype struct {
	id             archetypeID
	node           nodeID
	mask           bitMask
	components     []ID                     // components IDs of the archetype in arbitrary order
	itemSizes      []uint32                 // item size per component ID
	componentsMap  []int16                  // mapping from component IDs to column indices; -1 indicates none
	isRelation     []bool                   // whether columns are relations components, indexed by column index
	relationTables []map[entityID]*tableIDs // lookup for relation targets of tables, indexed by column index
	tables         []tableID                // all active tables
	freeTables     []tableID                // all inactive/free tables
	numRelations   uint8                    // number of relation components

	zeroValue []byte // zero value with the size of the largest item type, for fast zeroing
}

type tableIDs struct {
	tables []tableID
}

func newArchetype(id archetypeID, node nodeID, mask *bitMask, components []ID, tables []tableID, reg *componentRegistry) archetype {
	componentsMap := make([]int16, maskTotalBits)
	for i := range maskTotalBits {
		componentsMap[i] = -1
	}

	sizes := make([]uint32, len(components))
	var maxSize uintptr = entitySize
	for i, id := range components {
		componentsMap[id.id] = int16(i)
		tp := reg.Types[id.id]

		itemSize := sizeOf(tp)
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
	return archetype{
		id:             id,
		node:           node,
		mask:           *mask,
		components:     components,
		itemSizes:      sizes,
		componentsMap:  componentsMap,
		isRelation:     isRelation,
		tables:         tables,
		numRelations:   numRelations,
		relationTables: relationTables,
		zeroValue:      zeroValue,
	}
}

func (a *archetype) HasRelations() bool {
	return a.numRelations > 0
}

func (a *archetype) GetTable(storage *storage, relations []RelationID) (*table, bool) {
	if len(a.tables) == 0 {
		return nil, false
	}
	if !a.HasRelations() {
		return &storage.tables[a.tables[0]], true
	}
	if len(relations) != int(a.numRelations) {
		panic("relations must be fully specified")
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

func (a *archetype) GetTables(relations []RelationID) []tableID {
	if !a.HasRelations() || len(relations) == 0 {
		return a.tables
	}
	index := a.componentsMap[relations[0].component.id]
	if tables, ok := a.relationTables[index][relations[0].target.id]; ok {
		return tables.tables
	}
	return nil
}

func (a *archetype) GetFreeTable() (tableID, bool) {
	if len(a.freeTables) == 0 {
		return 0, false
	}
	last := len(a.freeTables) - 1
	table := a.freeTables[last]

	a.freeTables = a.freeTables[:last]

	return table, true
}

func (a *archetype) FreeTable(table tableID) {
	// TODO: can we speed this up for large numbers of relation targets?
	index := slices.Index(a.tables, table)
	last := len(a.tables) - 1

	if index != last {
		a.tables[index], a.tables[last] = a.tables[last], a.tables[index]
	}
	a.tables = a.tables[:last]

	a.freeTables = append(a.freeTables, table)
}

func (a *archetype) AddTable(table *table) {
	a.tables = append(a.tables, table.id)
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
			tables.tables = append(tables.tables, table.id)
		} else {
			relations[target.id] = &tableIDs{tables: []tableID{table.id}}
		}
	}
}

func (a *archetype) RemoveTarget(entity Entity) {
	for i := range a.relationTables {
		if !a.isRelation[i] {
			continue
		}
		delete(a.relationTables[i], entity.id)
	}
}

func (a *archetype) Reset(storage *storage) {
	if !a.HasRelations() {
		storage.tables[a.tables[0]].Reset()
		return
	}

	for _, tab := range a.tables {
		table := &storage.tables[tab]
		table.Reset()
	}

	for i := len(a.tables) - 1; i >= 0; i-- {
		storage.cache.removeTable(storage, &storage.tables[a.tables[i]])
		a.FreeTable(a.tables[i])
	}

	for _, m := range a.relationTables {
		for key := range m {
			delete(m, key)
		}
	}
}

// Stats generates statistics for an archetype.
func (a *archetype) Stats(storage *storage) stats.Archetype {
	ids := a.components
	aTypes := make([]reflect.Type, len(ids))
	for j, id := range ids {
		aTypes[j], _ = storage.registry.ComponentType(id.id)
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
	tableStats := make([]stats.Table, len(a.tables))
	for i, id := range a.tables {
		table := &storage.tables[id]
		tableStats[i] = table.Stats(memPerEntity, &storage.registry)
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
		FreeTables:      len(a.freeTables),
		NumRelations:    int(a.numRelations),
		ComponentIDs:    intIDs,
		ComponentTypes:  aTypes,
		Memory:          memory,
		MemoryUsed:      memoryUsed,
		MemoryPerEntity: memPerEntity,
		Size:            count,
		Capacity:        cap,
		Tables:          tableStats,
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
	cntNew := int32(len(tables))
	if cntNew < cntOld {
		stats.Tables = stats.Tables[:cntNew]
		cntOld = cntNew
	}
	var i int32
	for i = 0; i < cntOld; i++ {
		tableStats := &stats.Tables[i]
		table := &storage.tables[tables[i]]
		table.UpdateStats(stats.MemoryPerEntity, tableStats, &storage.registry)
		cap += tableStats.Capacity
		count += tableStats.Size
		memory += tableStats.Memory
		memoryUsed += tableStats.MemoryUsed
	}
	for i = cntOld; i < cntNew; i++ {
		table := &storage.tables[tables[i]]
		tableStats := table.Stats(stats.MemoryPerEntity, &storage.registry)
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
