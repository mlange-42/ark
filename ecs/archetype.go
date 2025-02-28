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
	relationTables []map[entityID]*tableIDs
}

type tableIDs struct {
	tables []tableID
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
	if !a.HasRelations() {
		return a.tables
	}
	if len(relations) == 0 {
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
	index := slices.Index(a.tables, table)
	last := len(a.tables) - 1

	a.tables[index], a.tables[last] = a.tables[last], a.tables[index]
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
