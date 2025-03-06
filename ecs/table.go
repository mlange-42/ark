package ecs

import (
	"fmt"
	"math"
	"reflect"
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
	len         uint32
	cap         uint32

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
		cap:         capacity,
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
	idx := t.len
	t.Alloc(1)
	t.entities.Set(idx, unsafe.Pointer(&entity))
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
	t.Extend(n)
	t.len += n
}

// Extend the table to be able to store the given number of additional entities.
// Has no effect of the table's capacity is already sufficient.
// If the capacity needs to be increased, it will be doubled until it is sufficient.
func (t *table) Extend(by uint32) {
	required := t.len + by
	if t.cap >= required {
		return
	}
	for t.cap < required {
		t.cap *= 2
	}

	old := t.entities.data
	t.entities.data = reflect.New(reflect.ArrayOf(int(t.cap), old.Type().Elem())).Elem()
	t.entities.pointer = t.entities.data.Addr().UnsafePointer()
	reflect.Copy(t.entities.data, old)

	for i := range t.columns {
		column := &t.columns[i]
		old := column.data
		column.data = reflect.New(reflect.ArrayOf(int(t.cap), old.Type().Elem())).Elem()
		column.pointer = column.data.Addr().UnsafePointer()
		reflect.Copy(column.data, old)
	}
}

// Remove swap-removes the entity at the given index.
// Returns whether a swap was necessary.
func (t *table) Remove(index uint32) bool {
	lastIndex := uintptr(t.len - 1)
	swapped := index != uint32(lastIndex)

	if swapped {
		sz := t.entities.itemSize
		src := unsafe.Add(t.entities.pointer, lastIndex*sz)
		dst := unsafe.Add(t.entities.pointer, uintptr(index)*sz)
		copyPtr(src, dst, uintptr(sz))
	}

	for i := range t.columns {
		column := &t.columns[i]
		if swapped {
			sz := column.itemSize
			src := unsafe.Add(column.pointer, lastIndex*sz)
			dst := unsafe.Add(column.pointer, uintptr(index)*sz)
			copyPtr(src, dst, uintptr(sz))
		}
		column.Zero(lastIndex, t.zeroPointer)
	}

	t.len--
	return swapped
}

func (t *table) Reset() {
	t.entities.Reset(t.len, nil)
	for c := range t.columns {
		t.columns[c].Reset(t.len, t.zeroPointer)
	}
	t.len = 0
}

func (t *table) AddAll(other *table, count uint32) {
	t.Alloc(count)
	t.entities.SetLast(&other.entities, t.len, count)
	for c := range t.columns {
		t.columns[c].SetLast(&other.columns[c], t.len, count)
	}
}

func (t *table) AddAllEntities(other *table, count uint32) {
	t.Alloc(count)
	t.entities.SetLast(&other.entities, t.len, count)
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
	return int(t.len)
}
