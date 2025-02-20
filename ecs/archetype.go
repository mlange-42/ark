package ecs

type archetypeID uint32

type archetype struct {
	id          archetypeID
	mask        Mask
	components  []ID
	tables      []*table
	hasRelation bool
}

func newArchetype(id archetypeID, mask *Mask, components []ID, tables []*table) archetype {
	return archetype{
		id:         id,
		mask:       *mask,
		components: components,
		tables:     tables,
	}
}

func (a *archetype) GetTable() (*table, bool) {
	if len(a.tables) == 0 {
		return nil, false
	}
	// TODO: consider relations
	return a.tables[0], true
}
