package ecs

import (
	"fmt"
	"math"
	"reflect"
	"unsafe"

	"github.com/mlange-42/ark/ecs/stats"
)

type tableID uint32

// maxTableID is used as table ID for unused entities.
const maxTableID = math.MaxUint32

type table struct {
	zeroPointer unsafe.Pointer // pointer to the zero value, for fast zeroing
	components  []*column      // mapping from component IDs to columns
	ids         []ID           // components IDs in the same order as in the archetype
	columns     []column       // columns in dense order
	relationIDs []RelationID   // all relation IDs and targets of the table
	entities    column         // column for entities
	id          tableID
	archetype   archetypeID
	len         uint32
	cap         uint32
}

func newTable(id tableID, archetype *archetype, capacity uint32, reg *componentRegistry, targets []Entity, relationIDs []RelationID) table {
	entities := newColumn(0, entityType, entitySize, false, Entity{}, capacity)
	columns := make([]column, len(archetype.components))

	components := make([]*column, maskTotalBits)
	for i, id := range archetype.components {
		itemSize := uintptr(archetype.itemSizes[i])
		columns[i] = newColumn(uint32(i), reg.Types[id.id], itemSize, archetype.isRelation[i], targets[i], capacity)
		components[id.id] = &columns[i]
	}

	var zeroPointer unsafe.Pointer
	if archetype.zeroValue != nil {
		zeroPointer = unsafe.Pointer(&archetype.zeroValue[0])
	}

	return table{
		id:          id,
		archetype:   archetype.id,
		components:  components,
		entities:    entities,
		ids:         archetype.components,
		columns:     columns,
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
	return t.components[component.id].Get(index)
}

func (t *table) Has(component ID) bool {
	return t.components[component.id] != nil
}

func (t *table) GetEntity(index uintptr) Entity {
	return *(*Entity)(t.entities.Get(index))
}

func (t *table) GetRelation(component ID) Entity {
	return t.components[component.id].target
}

func (t *table) GetColumn(component ID) *column {
	return t.components[component.id]
}

func (t *table) Set(component ID, index uint32, comp unsafe.Pointer) {
	t.components[component.id].Set(index, comp)
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

		for i := range t.columns {
			column := &t.columns[i]
			sz := column.itemSize
			src := unsafe.Add(column.pointer, lastIndex*sz)
			dst := unsafe.Add(column.pointer, uintptr(index)*sz)
			copyPtr(src, dst, uintptr(sz))
			column.Zero(lastIndex, t.zeroPointer)
		}
	} else {
		for i := range t.columns {
			column := &t.columns[i]
			column.Zero(lastIndex, t.zeroPointer)
		}
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
	if len(relations) < len(t.relationIDs) {
		panic("relation targets must be fully specified")
	}
	for _, rel := range relations {
		column := t.components[rel.component.id]
		if column == nil {
			// TODO: is there a check anywhere that prevents adding/setting relations
			// on columns not in the table?
			continue
		}
		if !column.isRelation {
			panic(fmt.Sprintf("component %d is not a relation component", rel.component.id))
		}
		//if rel.target == wildcard {
		//	panic("relation targets must be fully specified, no wildcard allowed")
		//}
		if rel.target != column.target {
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
		if rel.target != t.components[rel.component.id].target {
			return false
		}
	}
	return true
}

func (t *table) Len() int {
	return int(t.len)
}

// Stats generates statistics for a table.
func (t *table) Stats(memPerEntity int, reg *componentRegistry) stats.Table {
	cap := int(t.cap)

	return stats.Table{
		Size:       int(t.Len()),
		Capacity:   cap,
		Memory:     cap * memPerEntity,
		MemoryUsed: t.Len() * memPerEntity,
	}
}

// UpdateStats updates statistics for a table.
func (t *table) UpdateStats(memPerEntity int, stats *stats.Table, reg *componentRegistry) {
	cap := int(t.cap)

	stats.Size = int(t.Len())
	stats.Capacity = cap
	stats.Memory = cap * memPerEntity
	stats.MemoryUsed = t.Len() * memPerEntity
}
