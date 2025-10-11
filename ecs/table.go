package ecs

import (
	"fmt"
	"math"
	"reflect"
	"unsafe"

	"github.com/mlange-42/ark/ecs/stats"
)

// tableID is the type of table IDs.
type tableID uint32

// maxTableID is used as table ID for unused entities.
const maxTableID = math.MaxUint32

// table stores entities and components of an archetype,
// or of a certain combination of relations inside an archetype.
type table struct {
	entities    entityColumn   // column for entities
	zeroPointer unsafe.Pointer // pointer to the zero value, for fast zeroing
	components  []*column      // mapping from component IDs to columns
	ids         []ID           // components IDs in the same order as in the archetype
	columns     []column       // columns in dense order
	relationIDs []relationID   // all relation IDs and targets of the table
	id          tableID        // ID of the table
	archetype   archetypeID    // ID of the table's archetype
	len         uint32         // length of the table (number of rows)
	cap         uint32         // capacity of the table (number of rows)
	isFree      bool           // Whether the table is currently free
}

// newTable creates a new table.
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

// Recycle the table, using the given relations and targets, indexed by column index.
func (t *table) Recycle(targets []Entity, relationIDs []relationID) {
	t.relationIDs = relationIDs
	for i := range t.columns {
		t.columns[i].target = targets[i]
	}
	t.isFree = false
}

// HasRelations returns whether the table contains any relation components.
func (t *table) HasRelations() bool {
	return len(t.relationIDs) > 0
}

// Add an entity to the table.
// Returns the entity's new row index.
func (t *table) Add(entity Entity) uint32 {
	idx := t.len
	t.Alloc(1)
	t.entities.Set(idx, unsafe.Pointer(&entity))
	return idx
}

// Get a pointer to the given component at the given index.
func (t *table) Get(component ID, index uintptr) unsafe.Pointer {
	return t.components[component.id].Get(index)
}

// Has returns whether the table has a column for the given component.
func (t *table) Has(component ID) bool {
	return t.components[component.id] != nil
}

// GetEntity returns the entity at the given row index.
func (t *table) GetEntity(index uintptr) Entity {
	return t.entities.GetEntity(index)
}

// GetRelation returns the target entity for the given relation component.
func (t *table) GetRelation(component ID) Entity {
	return t.components[component.id].target
}

// Column returns the column pointer for the given component ID.
func (t *table) Column(component ID) *column {
	return t.components[component.id]
}

// Set the value of a component at the given row index from a column from another table.
func (t *table) Set(component ID, index uint32, src *column, srcIndex int) {
	t.components[component.id].Set(index, src, srcIndex)
}

// SetEntity sets the entity at the given row index.
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

// adjustCapacity changes the capacity of all columns.
// Does not check whether the change is necessary or feasible.
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
		copyPtr(src, dst, size)

		for i := range t.columns {
			column := &t.columns[i]

			if column.isTrivial {
				size := column.itemSize
				src := unsafe.Add(column.pointer, lastIndex*size)
				dst := unsafe.Add(column.pointer, uintptr(index)*size)
				copyPtr(src, dst, size)
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

// Reset the table.
// Clears all columns and sets the number of rows to zero.
func (t *table) Reset() {
	for c := range t.columns {
		t.columns[c].Reset(t.len, t.zeroPointer)
	}
	t.len = 0
}

// AddAll adds all entities with components from another table to this table.
func (t *table) AddAll(from *table, count uint32) {
	t.Alloc(count)
	t.entities.CopyToEnd(&from.entities, t.len, count)
	for c := range t.columns {
		t.columns[c].CopyToEnd(&from.columns[c], t.len, count)
	}
}

// AddAllEntities adds all entities (without components) from another table to this table.
func (t *table) AddAllEntities(from *table, count uint32) {
	t.Alloc(count)
	t.entities.CopyToEnd(&from.entities, t.len, count)
}

// MatchesExact returns whether this table matches the given relations exactly and exhaustively.
// Unspecified relations are not allowed.
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

// Matches returns whether this table matches the given relations.
// Unspecified relations are allowed.
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

// Len returns the length of the table (number of rows).
func (t *table) Len() int {
	return int(t.len)
}

// Stats generates statistics for a table.
func (t *table) Stats(memPerEntity int) stats.Table {
	cap := int(t.cap)

	return stats.Table{
		Size:       t.Len(),
		Capacity:   cap,
		Memory:     cap * memPerEntity,
		MemoryUsed: t.Len() * memPerEntity,
	}
}

// UpdateStats updates statistics for a table.
func (t *table) UpdateStats(memPerEntity int, stats *stats.Table) {
	cap := int(t.cap)

	stats.Size = t.Len()
	stats.Capacity = cap
	stats.Memory = cap * memPerEntity
	stats.MemoryUsed = t.Len() * memPerEntity
}
