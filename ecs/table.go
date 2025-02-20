package ecs

import "unsafe"

type table struct {
	components []int16
	entities   column
	columns    []column
	relations  []Entity
}

func newTable(capacity int, reg *registry, ids ...ID) table {
	components := make([]int16, MaskTotalBits)
	entities := newColumn(entityType, capacity)
	columns := make([]column, len(ids))

	for i := range MaskTotalBits {
		components[i] = -1
	}

	for i, id := range ids {
		components[id.id] = int16(i)
		columns[i] = newColumn(reg.Types[id.id], capacity)
	}
	return table{
		components: components,
		entities:   entities,
		columns:    columns,
		relations:  make([]Entity, len(ids)),
	}
}

func (t *table) Get(component ID, index uint32) unsafe.Pointer {
	return t.columns[t.components[component.id]].Get(index)
}

func (t *table) GetEntity(index uint32) Entity {
	return *(*Entity)(t.entities.Get(index))
}

func (t *table) GetRelation(component ID) Entity {
	return t.relations[t.components[component.id]]
}

func (t *table) GetColumn(component ID) *column {
	return &t.columns[t.components[component.id]]
}

func (t *table) GetEntities(component ID) *column {
	return &t.entities
}
