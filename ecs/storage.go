package ecs

import "fmt"

type storage struct {
	registry        componentRegistry
	archetypes      []archetype
	tables          []table
	initialCapacity uint32
	components      []componentStorage
}

type componentStorage struct {
	columns []*column
}

func newStorage(capacity uint32) storage {
	reg := newComponentRegistry()
	tables := make([]table, 0, 128)
	tables = append(tables, newTable(0, 0, capacity, &reg, []ID{}, []int16{}, []Entity{}, 0))
	archetypes := make([]archetype, 0, 128)
	archetypes = append(archetypes, newArchetype(0, &Mask{}, []ID{}, []*table{&tables[0]}, &reg))
	return storage{
		registry:        reg,
		archetypes:      archetypes,
		tables:          tables,
		initialCapacity: capacity,
		components:      make([]componentStorage, 0, MaskTotalBits),
	}
}

func (s *storage) findOrCreateTable(mask *Mask, relations []relationID) *table {
	// TODO: use archetype graph
	var arch *archetype
	for i := range s.archetypes {
		if s.archetypes[i].mask.Equals(mask) {
			arch = &s.archetypes[i]
			break
		}
	}
	if arch == nil {
		arch = s.createArchetype(mask)
	}
	table, ok := arch.GetTable(relations)
	if !ok {
		table = s.createTable(arch, relations)
	}
	return table
}

func (s *storage) AddComponent(id uint8) {
	if len(s.components) != int(id) {
		panic("components can only be added to a storage sequentially")
	}
	s.components = append(s.components, componentStorage{columns: make([]*column, len(s.tables))})
}

func (s *storage) createArchetype(mask *Mask) *archetype {
	comps := mask.toTypes(&s.registry.registry)
	index := len(s.archetypes)
	s.archetypes = append(s.archetypes, newArchetype(archetypeID(index), mask, comps, nil, &s.registry))
	return &s.archetypes[index]
}

func (s *storage) createTable(archetype *archetype, relations []relationID) *table {
	index := len(s.tables)

	targets := make([]Entity, len(archetype.components))
	numRelations := uint8(0)
	for _, rel := range relations {
		idx := archetype.componentsMap[rel.component.id]
		targets[idx] = rel.target
		numRelations++
	}
	if numRelations != archetype.numRelations {
		panic("relations must be fully specified")
	}

	s.tables = append(s.tables, newTable(tableID(index), archetype.id, s.initialCapacity, &s.registry, archetype.components, archetype.componentsMap, targets, numRelations))
	table := &s.tables[index]
	archetype.tables = append(archetype.tables, table)
	for i := range s.components {
		id := ID{id: uint8(i)}
		comps := &s.components[i]
		if archetype.mask.Get(id) {
			comps.columns = append(comps.columns, table.GetColumn(id))
		} else {
			comps.columns = append(comps.columns, nil)
		}
	}
	return table
}

func (s *storage) getExchangeMask(mask *Mask, add []ID, rem []ID) {
	for _, comp := range rem {
		if !mask.Get(comp) {
			panic(fmt.Sprintf("entity does not have a component of type %v, can't remove", s.registry.Types[comp.id]))
		}
		mask.Set(comp, false)
	}
	for _, comp := range add {
		if mask.Get(comp) {
			panic(fmt.Sprintf("entity already has component of type %v, can't add", s.registry.Types[comp.id]))
		}
		mask.Set(comp, true)
	}
}
