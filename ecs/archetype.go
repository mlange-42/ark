package ecs

type archetypeID uint32

type archetype struct {
	id          archetypeID
	mask        Mask
	components  []ID
	isRelation  []bool
	tables      []*table
	hasRelation bool
}

func newArchetype(id archetypeID, mask *Mask, components []ID, tables []*table, reg *componentRegistry) archetype {
	hasRelation := false
	isRelation := make([]bool, len(components))
	for _, id := range components {
		if reg.IsRelation.Get(id) {
			hasRelation = true
			isRelation[id.id] = true
		}
	}
	return archetype{
		id:          id,
		mask:        *mask,
		components:  components,
		isRelation:  isRelation,
		tables:      tables,
		hasRelation: hasRelation,
	}
}

func (a *archetype) GetTable() (*table, bool) {
	if len(a.tables) == 0 {
		return nil, false
	}
	// TODO: consider relations
	return a.tables[0], true
}
