package ecs

type archetypeID uint32

type archetype struct {
	id            archetypeID
	mask          Mask
	components    []ID
	componentsMap []int16
	isRelation    []bool
	tables        []*table
	numRelations  uint8
}

func newArchetype(id archetypeID, mask *Mask, components []ID, tables []*table, reg *componentRegistry) archetype {
	componentsMap := make([]int16, MaskTotalBits)
	for i := range MaskTotalBits {
		componentsMap[i] = -1
	}
	for i, id := range components {
		componentsMap[id.id] = int16(i)
	}

	numRelations := uint8(0)
	isRelation := make([]bool, len(components))
	for _, id := range components {
		if reg.IsRelation.Get(id) {
			isRelation[id.id] = true
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

func (a *archetype) GetTable(relations []relation) (*table, bool) {
	if len(a.tables) == 0 {
		return nil, false
	}
	if !a.HasRelations() {
		return a.tables[0], true
	}
	for _, t := range a.tables {
		if t.MatchesExact(relations) {
			return t, true
		}
	}
	return nil, false
}
