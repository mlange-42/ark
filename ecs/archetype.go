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
	mask           Mask
	components     []ID                     // components IDs of the archetype in arbitrary order
	componentsMap  []int16                  // mapping from component IDs to column indices; -1 indicates none
	isRelation     []bool                   // whether columns are relations components, indexed by column index
	relationTables []map[entityID]*tableIDs // lookup for relation targets of tables, indexed by column index
	tables         []tableID                // all active tables
	freeTables     []tableID                // all inactive/free tables
	numRelations   uint8                    // number of relation components
}

type tableIDs struct {
	tables []tableID
}

func newArchetype(id archetypeID, node nodeID, mask *Mask, components []ID, tables []tableID, reg *componentRegistry) archetype {
	componentsMap := make([]int16, MaskTotalBits)
	for i := range MaskTotalBits {
		componentsMap[i] = -1
	}
	for i, id := range components {
		componentsMap[id.id] = int16(i)
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
		componentsMap:  componentsMap,
		isRelation:     isRelation,
		tables:         tables,
		numRelations:   numRelations,
		relationTables: relationTables,
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
func (a *archetype) Stats(storage *storage) stats.Node {
	ids := a.components
	aCompCount := len(ids)
	aTypes := make([]reflect.Type, aCompCount)
	for j, id := range ids {
		aTypes[j], _ = storage.registry.ComponentType(id.id)
	}

	arches := a.tables
	var numArches int32
	cap := 0
	count := 0
	memory := 0
	var archStats []stats.Archetype
	if arches != nil {
		archStats = make([]stats.Archetype, numArches)
		for i, id := range arches {
			table := &storage.tables[id]
			archStats[i] = table.Stats(&storage.registry)
			stats := &archStats[i]
			cap += stats.Capacity
			count += stats.Size
			memory += stats.Memory
		}
	}

	memPerEntity := 0
	intIDs := make([]uint8, len(ids))
	for j, id := range ids {
		intIDs[j] = id.id
		memPerEntity += int(aTypes[j].Size())
	}

	return stats.Node{
		ArchetypeCount:       int(numArches),
		ActiveArchetypeCount: len(a.freeTables),
		HasRelation:          a.HasRelations(),
		Components:           aCompCount,
		ComponentIDs:         intIDs,
		ComponentTypes:       aTypes,
		Memory:               memory,
		MemoryPerEntity:      memPerEntity,
		Size:                 count,
		Capacity:             cap,
		Archetypes:           archStats,
	}
}

// UpdateStats updates statistics for an archetype.
func (a *archetype) UpdateStats(stats *stats.Node, storage *storage) {
	arches := a.tables

	cap := 0
	count := 0
	memory := 0

	cntOld := int32(len(stats.Archetypes))
	cntNew := int32(len(arches))
	var i int32
	for i = 0; i < cntOld; i++ {
		arch := &stats.Archetypes[i]
		table := &storage.tables[arches[i]]
		table.UpdateStats(stats, arch, &storage.registry)
		cap += arch.Capacity
		count += arch.Size
		memory += arch.Memory
	}
	for i = cntOld; i < cntNew; i++ {
		table := &storage.tables[arches[i]]
		arch := table.Stats(&storage.registry)
		stats.Archetypes = append(stats.Archetypes, arch)
		cap += arch.Capacity
		count += arch.Size
		memory += arch.Memory
	}

	stats.ArchetypeCount = int(cntNew)
	stats.ActiveArchetypeCount = len(a.freeTables)
	stats.Capacity = cap
	stats.Size = count
	stats.Memory = memory
}
