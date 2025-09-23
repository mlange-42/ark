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
	entities    entityColumn   // column for entities
	zeroPointer unsafe.Pointer // pointer to the zero value, for fast zeroing
	components  []*column      // mapping from component IDs to columns
	ids         []ID           // components IDs in the same order as in the archetype
	columns     []column       // columns in dense order
	relationIDs []relationID   // all relation IDs and targets of the table
	id          tableID
	archetype   archetypeID
	len         uint32
	cap         uint32
	isFree      bool
}

func newTable(id tableID, archetype *archetype, capacity uint32, reg *componentRegistry, targets []Entity, relationIDs []relationID) table {
	entities := newEntityColumn(capacity)
	columns := make([]column, len(archetype.components))

	components := make([]*column, maskTotalBits)
	for i, id := range archetype.components {
		itemSize := uintptr(archetype.itemSizes[i])
		columns[i] = newColumn(uint32(i), reg.Types[id.id], itemSize, archetype.isRelation[i], reg.IsTrivial[id.id], targets[i], capacity)
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

func (t *table) Recycle(targets []Entity, relationIDs []relationID) {
	t.relationIDs = relationIDs
	for i := range t.columns {
		t.columns[i].target = targets[i]
	}
	t.isFree = false
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

func (t *table) Column(component ID) *column {
	return t.components[component.id]
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

func (t *table) Set(component ID, index uint32, src *column, srcIndex int) {
	t.components[component.id].Set(index, src, srcIndex)
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
// Has no effect if the table's capacity is already sufficient.
// If the capacity needs to be increased, it will be doubled until it is sufficient.
func (t *table) Extend(by uint32) {
	required := t.len + by
	if t.cap >= required {
		return
	}
	t.adjustCapacity(capPow2(required))
}

// CanShrink returns whether the table's capacity exceeds the next power-of-2 of what is required.
func (t *table) CanShrink(minCapacity uint32) bool {
	target := max(capPow2(t.len), minCapacity)
	return t.cap > target
}

// Shrink the table's capacity to the next power-of-2 of what is required.
func (t *table) Shrink(minCapacity uint32) bool {
	target := max(capPow2(t.len), minCapacity)
	if t.cap <= target {
		return false
	}
	t.adjustCapacity(target)
	return true
}

func (t *table) adjustCapacity(cap uint32) {
	t.cap = cap

	t.entities.data = reflect.New(reflect.ArrayOf(int(t.cap), entityType)).Elem()
	newPtr := t.entities.data.Addr().UnsafePointer()
	if t.len > 0 {
		copyPtr(t.entities.pointer, newPtr, uintptr(t.len)*entitySize)
	}
	t.entities.pointer = newPtr

	for i := range t.columns {
		column := &t.columns[i]
		old := column.data
		column.data = reflect.New(reflect.ArrayOf(int(t.cap), column.elemType)).Elem()
		if column.isTrivial {
			newPtr := column.data.Addr().UnsafePointer()
			if t.len > 0 {
				copyPtr(column.pointer, newPtr, uintptr(t.len)*column.itemSize)
			}
			column.pointer = newPtr
		} else {
			column.pointer = column.data.Addr().UnsafePointer()
			if t.len > 0 {
				reflect.Copy(column.data, old)
			}
		}
	}
}

// Remove swap-removes the entity at the given index.
// Returns whether a swap was necessary.
func (t *table) Remove(index uint32) bool {
	lastIndex := uintptr(t.len - 1)
	swapped := index != uint32(lastIndex)

	if swapped {
		size := entitySize
		src := unsafe.Add(t.entities.pointer, lastIndex*size)
		dst := unsafe.Add(t.entities.pointer, uintptr(index)*size)
		copyPtr(src, dst, uintptr(size))

		for i := range t.columns {
			column := &t.columns[i]

			if column.isTrivial {
				size := column.itemSize
				src := unsafe.Add(column.pointer, lastIndex*size)
				dst := unsafe.Add(column.pointer, uintptr(index)*size)
				copyPtr(src, dst, uintptr(size))
				column.Zero(lastIndex, t.zeroPointer)
				continue
			}

			copyValue(column.data, column.data, int(lastIndex), int(index))
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
	for c := range t.columns {
		t.columns[c].Reset(t.len, t.zeroPointer)
	}
	t.len = 0
}

func (t *table) AddAll(from *table, count uint32) {
	t.Alloc(count)
	t.entities.CopyToEnd(&from.entities, t.len, count)
	for c := range t.columns {
		t.columns[c].CopyToEnd(&from.columns[c], t.len, count)
	}
}

func (t *table) AddAllEntities(from *table, count uint32) {
	t.Alloc(count)
	t.entities.CopyToEnd(&from.entities, t.len, count)
}

func (t *table) MatchesExact(relations []relationID) bool {
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

func (t *table) Matches(relations []relationID) bool {
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
