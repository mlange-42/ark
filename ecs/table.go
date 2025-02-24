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
	components  []int16
	entities    column
	ids         []ID
	columns     []column
	isRelation  []bool
	relations   []Entity
	relationIDs []relationID

	zeroValue   []byte
	zeroPointer unsafe.Pointer
}

func newTable(id tableID, archetype archetypeID, capacity uint32, reg *componentRegistry,
	ids []ID, componentsMap []int16, isRelation []bool, targets []Entity, relationIDs []relationID) table {

	entities := newColumn(entityType, capacity)
	columns := make([]column, len(ids))

	var maxSize uintptr = entitySize
	for i, id := range ids {
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
		components:  componentsMap,
		entities:    entities,
		ids:         ids,
		columns:     columns,
		isRelation:  isRelation,
		relations:   targets,
		zeroValue:   zeroValue,
		zeroPointer: zeroPointer,
		relationIDs: relationIDs,
	}
}

func (t *table) recycle(targets []Entity, relationIDs []relationID) {
	t.relations = targets
	t.relationIDs = relationIDs
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

func (t *table) AddAll(other *table) {
	t.entities.AddAll(&other.entities)
	for c := range t.columns {
		t.columns[c].AddAll(&other.columns[c])
	}
}

func (t *table) MatchesExact(relations []relationID) bool {
	if len(relations) != len(t.relationIDs) {
		panic("relation targets must be fully specified")
	}
	for _, rel := range relations {
		index := t.components[rel.component.id]
		if !t.isRelation[index] {
			panic(fmt.Sprintf("component %d is not a relation component", rel.component.id))
		}
		if rel.target == wildcard {
			panic("relation targets must be fully specified, no wildcard allowed")
		}
		if rel.target != t.relations[index] {
			return false
		}
	}
	return true
}

func (t *table) Matches(relations []relationID) bool {
	if len(relations) == 0 {
		return true
	}
	for _, rel := range relations {
		if rel.target == wildcard {
			continue
		}
		if rel.target != t.relations[t.components[rel.component.id]] {
			return false
		}
	}
	return true
}

func (t *table) Len() int {
	return t.entities.Len()
}
