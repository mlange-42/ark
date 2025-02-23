package ecs

import "fmt"

type storage struct {
	registry   componentRegistry
	entities   entities
	isTarget   []bool
	entityPool entityPool

	archetypes         []archetype
	relationArchetypes []archetypeID
	tables             []table
	components         []componentStorage
	initialCapacity    uint32
}

type componentStorage struct {
	columns []*column
}

func newStorage(capacity uint32) storage {
	reg := newComponentRegistry()
	entities := make([]entityIndex, reservedEntities, capacity+reservedEntities)
	isTarget := make([]bool, reservedEntities, capacity+reservedEntities)
	// Reserved zero and wildcard entities
	for i := range reservedEntities {
		entities[i] = entityIndex{table: maxTableID, row: 0}
	}

	tables := make([]table, 0, 128)
	tables = append(tables, newTable(0, 0, capacity, &reg, []ID{}, make([]int16, MaskTotalBits), []bool{}, []Entity{}, []relationID{}))
	archetypes := make([]archetype, 0, 128)
	archetypes = append(archetypes, newArchetype(0, &Mask{}, []ID{}, []tableID{0}, &reg))
	return storage{
		registry:        reg,
		entities:        entities,
		isTarget:        isTarget,
		entityPool:      newEntityPool(capacity, reservedEntities),
		archetypes:      archetypes,
		tables:          tables,
		initialCapacity: capacity,
		components:      make([]componentStorage, 0, MaskTotalBits),
	}
}

func (s *storage) findOrCreateTable(oldTable *table, mask *Mask, relations []relationID) *table {
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
	allRelation := appendNew(oldTable.relationIDs, relations...)
	table, ok := arch.GetTable(s, allRelation)
	if !ok {
		table = s.createTable(arch, allRelation)
	}
	return table
}

func (s *storage) AddComponent(id uint8) {
	if len(s.components) != int(id) {
		panic("components can only be added to a storage sequentially")
	}
	s.components = append(s.components, componentStorage{columns: make([]*column, len(s.tables))})
}

// RemoveEntity removes the given entity from the world.
func (s *storage) RemoveEntity(entity Entity) {
	if !s.entityPool.Alive(entity) {
		panic("can't remove a dead entity")
	}
	index := &s.entities[entity.id]
	table := &s.tables[index.table]

	swapped := table.Remove(index.row)

	s.entityPool.Recycle(entity)

	if swapped {
		swapEntity := table.GetEntity(uintptr(index.row))
		s.entities[swapEntity.id].row = index.row
	}
	index.table = maxTableID

	if s.isTarget[entity.id] {
		s.cleanupArchetypes(entity)
		s.isTarget[entity.id] = false
	}
}

func (s *storage) getRelation(entity Entity, comp ID) Entity {
	if !s.entityPool.Alive(entity) {
		panic("can't get relation for a dead entity")
	}
	return s.tables[s.entities[entity.id].table].GetRelation(comp)
}

func (s *storage) registerTargets(relations []relationID) {
	for _, rel := range relations {
		s.isTarget[rel.target.id] = true
	}
}

func (s *storage) createEntity(table tableID) (Entity, uint32) {
	entity := s.entityPool.Get()

	idx := s.tables[table].Add(entity)
	len := len(s.entities)
	if int(entity.id) == len {
		s.entities = append(s.entities, entityIndex{table: table, row: idx})
		s.isTarget = append(s.isTarget, false)
	} else {
		s.entities[entity.id] = entityIndex{table: table, row: idx}
	}
	return entity, idx
}

func (s *storage) createArchetype(mask *Mask) *archetype {
	comps := mask.toTypes(&s.registry.registry)
	index := len(s.archetypes)
	s.archetypes = append(s.archetypes, newArchetype(archetypeID(index), mask, comps, nil, &s.registry))
	archetype := &s.archetypes[index]
	if archetype.HasRelations() {
		s.relationArchetypes = append(s.relationArchetypes, archetype.id)
	}
	return archetype
}

func (s *storage) createTable(archetype *archetype, relations []relationID) *table {
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

	var newTableID tableID
	if id, ok := archetype.GetFreeTable(); ok {
		newTableID = id
		s.tables[newTableID].recycle(targets, relations)
	} else {
		newTableID = tableID(len(s.tables))
		s.tables = append(s.tables, newTable(
			newTableID, archetype.id, s.initialCapacity, &s.registry,
			archetype.components, archetype.componentsMap,
			archetype.isRelation, targets, relations))
	}

	table := &s.tables[newTableID]
	archetype.tables = append(archetype.tables, newTableID)
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

// Removes empty archetypes that have a target relation to the given entity.
func (s *storage) cleanupArchetypes(target Entity) {
	newRelations := []relationID{}
	for _, arch := range s.relationArchetypes {
		archetype := &s.archetypes[arch]
		for _, t := range archetype.tables {
			table := &s.tables[t]

			foundTarget := false
			for _, rel := range table.relationIDs {
				if rel.target.id == target.id {
					newRelations = append(newRelations, relationID{component: rel.component, target: Entity{}})
					foundTarget = true
				}
			}
			if !foundTarget {
				continue
			}

			allRelations := s.getExchangeTargetsUnchecked(table, newRelations)
			newTable, ok := archetype.GetTable(s, allRelations)
			if !ok {
				newTable = s.createTable(archetype, newRelations)
			}
			s.moveEntities(table, newTable)
			archetype.FreeTable(newTable.id)

			newRelations = newRelations[:0]
		}
	}
}

// moveEntities moves all entities from src to dst.
func (s *storage) moveEntities(src, dst *table) {
	oldLen := dst.Len()
	dst.AddAll(src)

	newLen := dst.Len()
	newTable := dst.id
	for i := oldLen; i < newLen; i++ {
		entity := dst.GetEntity(uintptr(i))
		s.entities[entity.id] = entityIndex{table: newTable, row: uint32(i)}
	}
	src.Reset()
}

func (s *storage) getExchangeTargetsUnchecked(oldTable *table, relations []relationID) []relationID {
	targets := append([]Entity(nil), oldTable.relations...)
	for _, rel := range relations {
		index := oldTable.components[rel.component.id]
		if rel.target == targets[index] {
			continue
		}
		targets[index] = rel.target
	}

	result := make([]relationID, 0, len(oldTable.relationIDs))
	for i, e := range targets {
		if !oldTable.isRelation[i] {
			continue
		}
		id := oldTable.ids[i]
		result = append(result, relationID{component: id, target: e})
	}

	return result
}

func (s *storage) getExchangeTargets(oldTable *table, relations []relationID) ([]relationID, bool) {
	changed := false
	targets := append([]Entity(nil), oldTable.relations...)
	for _, rel := range relations {
		if !rel.target.IsZero() && !s.entityPool.Alive(rel.target) {
			panic("can't make a dead entity a relation target")
		}
		index := oldTable.components[rel.component.id]
		if !oldTable.isRelation[index] {
			panic(fmt.Sprintf("component %d is not a relation component", rel.component.id))
		}
		if rel.target == targets[index] {
			continue
		}
		targets[index] = rel.target
		changed = true
	}
	if !changed {
		return nil, false
	}

	result := make([]relationID, 0, len(oldTable.relationIDs))
	for i, e := range targets {
		if !oldTable.isRelation[i] {
			continue
		}
		id := oldTable.ids[i]
		result = append(result, relationID{component: id, target: e})
	}

	return result, true
}
