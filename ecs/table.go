package ecs

import "unsafe"

type tableID uint32

type table struct {
	components []int16
	entities   column
	columns    []column
	relations  []Entity

	zeroValue   []byte
	zeroPointer unsafe.Pointer
}

func newTable(capacity uint32, reg *registry, ids ...ID) table {
	components := make([]int16, MaskTotalBits)
	entities := newColumn(entityType, capacity)
	columns := make([]column, len(ids))

	for i := range MaskTotalBits {
		components[i] = -1
	}

	var maxSize uintptr = 0
	for i, id := range ids {
		components[id.id] = int16(i)
		columns[i] = newColumn(reg.Types[id.id], capacity)
		if columns[i].itemSize > maxSize {
			maxSize = columns[i].itemSize
		}
	}
	var zeroValue []byte
	var zeroPointer unsafe.Pointer
	if maxSize > 0 {
		zeroValue = make([]byte, maxSize)
		zeroPointer = unsafe.Pointer(&zeroValue[0])
	}

	return table{
		components:  components,
		entities:    entities,
		columns:     columns,
		relations:   make([]Entity, len(ids)),
		zeroValue:   zeroValue,
		zeroPointer: zeroPointer,
	}
}

func (t *table) Add(entity Entity) uint32 {
	_, idx := t.entities.Add(unsafe.Pointer(&entity))
	for i := range t.columns {
		t.columns[i].Alloc(1)
	}
	return idx
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
