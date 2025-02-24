package ecs

import "slices"

type archetypeID uint32

type archetype struct {
	id             archetypeID
	mask           Mask
	components     []ID
	componentsMap  []int16
	isRelation     []bool
	tables         []tableID
	freeTables     []tableID
	numRelations   uint8
	relationTables []map[entityID][]tableID
}

func newArchetype(id archetypeID, mask *Mask, components []ID, tables []tableID, reg *componentRegistry) archetype {
	componentsMap := make([]int16, MaskTotalBits)
	for i := range MaskTotalBits {
		componentsMap[i] = -1
	}
	for i, id := range components {
		componentsMap[id.id] = int16(i)
	}

	numRelations := uint8(0)
	isRelation := make([]bool, len(components))
	//relationTables := make([]map[entityID][]tableID, len(components))
	for i, id := range components {
		if reg.IsRelation[id.id] {
			isRelation[i] = true
			//relationTables[]
			numRelations++
		}
	}
	return archetype{
		id:            id,
		mask:          *mask,
		components:    components,
		componentsMap: componentsMap,
		isRelation:    isRelation,
		tables:        tables,
		numRelations:  numRelations,
	}
}

func (a *archetype) HasRelations() bool {
	return a.numRelations > 0
}

func (a *archetype) GetTable(storage *storage, relations []relationID) (*table, bool) {
	if len(a.tables) == 0 {
		return nil, false
	}
	if !a.HasRelations() {
		return &storage.tables[a.tables[0]], true
	}
	for _, t := range a.tables {
		table := &storage.tables[t]
		if table.MatchesExact(relations) {
			return table, true
		}
	}
	return nil, false
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
	index := slices.Index(a.tables, table)
	last := len(a.tables) - 1

	a.tables[index], a.tables[last] = a.tables[last], a.tables[index]
	a.tables = a.tables[:last]

	a.freeTables = append(a.freeTables, table)
}

func (a *archetype) AddTable(table *table) {
	a.tables = append(a.tables, table.id)
}
