package ecs

import (
	"fmt"
	"math"
	"unsafe"
)

type tableID uint32

// maxTableID is used as table ID for unused entities.
const maxTableID = math.MaxUint32

type table struct {
	id          tableID
	archetype   archetypeID
	components  []int16      // mapping from component IDs to column indices
	entities    column       // column for entities
	ids         []ID         // components IDs in the same order as in the archetype
	columns     []column     // columns in dense order
	relationIDs []RelationID // all relation IDs and targets of the table

	zeroValue   []byte         // zero value with the size of the largest item type, for fast zeroing
	zeroPointer unsafe.Pointer // pointer to the zero value, for fast zeroing
}

func newTable(id tableID, archetype archetypeID, capacity uint32, reg *componentRegistry,
	ids []ID, componentsMap []int16, isRelation []bool, targets []Entity, relationIDs []RelationID) table {

	entities := newColumn(entityType, false, Entity{}, capacity)
	columns := make([]column, len(ids))

	var maxSize uintptr = entitySize
	for i, id := range ids {
		columns[i] = newColumn(reg.Types[id.id], isRelation[i], targets[i], capacity)
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
		components:  componentsMap,
		entities:    entities,
		ids:         ids,
		columns:     columns,
		zeroValue:   zeroValue,
		zeroPointer: zeroPointer,
		relationIDs: relationIDs,
	}
}

func (t *table) recycle(targets []Entity, relationIDs []RelationID) {
	t.relationIDs = relationIDs
	for i := range t.columns {
		t.columns[i].target = targets[i]
	}
}

func (t *table) HasRelations() bool {
	return len(t.relationIDs) > 0
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
	return t.columns[t.components[component.id]].target
}

func (t *table) GetColumn(component ID) *column {
	return &t.columns[t.components[component.id]]
}

func (t *table) Set(component ID, index uint32, comp unsafe.Pointer) {
	t.columns[t.components[component.id]].Set(index, comp)
}

func (t *table) SetEntity(index uint32, entity Entity) {
	t.entities.Set(index, unsafe.Pointer(&entity))
}

// Alloc allocates memory for the given number of entities.
func (t *table) Alloc(n uint32) {
	t.entities.Alloc(n)
	for i := range t.columns {
		t.columns[i].Alloc(n)
	}
}

func (t *table) Remove(index uint32) bool {
	swapped := t.entities.Remove(index, nil)
	for i := range t.columns {
		t.columns[i].Remove(index, t.zeroPointer)
	}
	return swapped
}

func (t *table) Reset() {
	t.entities.Reset(nil)
	for c := range t.columns {
		t.columns[c].Reset(t.zeroPointer)
	}
}

func (t *table) AddAll(other *table, count uint32) {
	t.entities.AddAll(&other.entities, count)
	for c := range t.columns {
		t.columns[c].AddAll(&other.columns[c], count)
	}
}

func (t *table) AddAllEntities(other *table, count uint32, allocColumns bool) {
	t.entities.AddAll(&other.entities, count)
	if allocColumns {
		for c := range t.columns {
			t.columns[c].Alloc(uint32(count))
		}
	}
}

func (t *table) MatchesExact(relations []RelationID) bool {
	if len(relations) != len(t.relationIDs) {
		panic("relation targets must be fully specified")
	}
	for _, rel := range relations {
		index := t.components[rel.component.id]
		if !t.columns[index].isRelation {
			panic(fmt.Sprintf("component %d is not a relation component", rel.component.id))
		}
		//if rel.target == wildcard {
		//	panic("relation targets must be fully specified, no wildcard allowed")
		//}
		if rel.target != t.columns[index].target {
			return false
		}
	}
	return true
}

func (t *table) Matches(relations []RelationID) bool {
	if len(relations) == 0 || !t.HasRelations() {
		return true
	}
	for _, rel := range relations {
		//if rel.target == wildcard {
		//	continue
		//}
		if rel.target != t.columns[t.components[rel.component.id]].target {
			return false
		}
	}
	return true
}

func (t *table) Len() int {
	return int(t.entities.Len())
}
