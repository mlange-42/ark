package ecs

import "fmt"

type storage struct {
	registry        registry
	archetypes      []archetype
	tables          []table
	initialCapacity uint32
}

func newStorage(capacity uint32) storage {
	reg := newRegistry()
	tables := make([]table, 0, 128)
	tables = append(tables, newTable(0, 0, capacity, &reg))
	archetypes := make([]archetype, 0, 128)
	archetypes = append(archetypes, newArchetype(0, &Mask{}, []ID{}, []*table{&tables[0]}))
	return storage{
		registry:        reg,
		archetypes:      archetypes,
		tables:          tables,
		initialCapacity: capacity,
	}
}

func (s *storage) findOrCreateTable(mask *Mask) *table {
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
	table, ok := arch.GetTable()
	if !ok {
		table = s.createTable(arch)
	}
	return table
}

func (s *storage) createArchetype(mask *Mask) *archetype {
	comps := mask.toTypes(&s.registry)
	index := len(s.archetypes)
	s.archetypes = append(s.archetypes, newArchetype(archetypeID(index), mask, comps, nil))
	return &s.archetypes[index]
}

func (s *storage) createTable(archetype *archetype) *table {
	index := len(s.tables)
	s.tables = append(s.tables, newTable(tableID(index), archetype.id, s.initialCapacity, &s.registry, archetype.components...))
	return &s.tables[index]
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
