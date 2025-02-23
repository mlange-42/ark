package ecs

import (
	"math"
	"unsafe"
)

type tableID uint32

// maxTableID is used as table ID for unused entities.
const maxTableID = math.MaxUint32

type table struct {
	id         tableID
	archetype  archetypeID
	components []int16
	entities   column
	columns    []column
	relations  []Entity

	zeroValue   []byte
	zeroPointer unsafe.Pointer
}

func newTable(id tableID, archetype archetypeID, capacity uint32, reg *componentRegistry, ids []ID) table {
	components := make([]int16, MaskTotalBits)
	entities := newColumn(entityType, capacity)
	columns := make([]column, len(ids))

	for i := range MaskTotalBits {
		components[i] = -1
	}

	var maxSize uintptr = entitySize
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
		id:          id,
		archetype:   archetype,
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

func (t *table) Get(component ID, index uintptr) unsafe.Pointer {
	return t.columns[t.components[component.id]].Get(index)
}

func (t *table) Has(component ID) bool {
	return t.components[component.id] >= 0
}

func (t *table) GetEntity(index uintptr) Entity {
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

func (t *table) Set(component ID, index uint32, comp unsafe.Pointer) {
	t.columns[t.components[component.id]].Set(index, comp)
}

func (t *table) Remove(index uint32) bool {
	swapped := t.entities.Remove(index, t.zeroPointer)
	for i := range t.columns {
		t.columns[i].Remove(index, t.zeroPointer)
	}
	return swapped
}

func (t *table) Len() int {
	return t.entities.Len()
}
