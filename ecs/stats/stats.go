// Package stats provides the structs returned by ecs.World.Stats().
package stats

import (
	"fmt"
	"reflect"
	"strings"
)

// World provide statistics for an [ecs.World].
type World struct {
	// Entity statistics.
	Entities Entities
	// Total number of components.
	ComponentCount int
	// Component types, indexed by component ID.
	ComponentTypes []reflect.Type
	// Locked state of the world.
	Locked bool
	// Archetype statistics.
	Archetypes []Archetype
	// Memory used by entities and components.
	Memory int
	// Number of cached filters.
	CachedFilters int
}

// Entities provide statistics about [ecs.World] entities.
type Entities struct {
	// Currently used/alive entities.
	Used int
	// Current capacity of the entity pool.
	Total int
	// Recycled/available entities.
	Recycled int
	// Current capacity of the entities list.
	Capacity int
}

// Archetype statistics.
type Archetype struct {
	// Number of free tables.
	FreeTables int
	// Whether the archetype contains relation tables.
	HasRelation bool
	// Number of components.
	Components int
	// Component IDs.
	ComponentIDs []uint8
	// Component types for ComponentIDs.
	ComponentTypes []reflect.Type
	// Memory for components per entity, in bytes.
	MemoryPerEntity int
	// Total reserved memory for entities and components, in bytes.
	Memory int
	// Number of entities in the tables of this archetype.
	Size int
	// Sum of capacity of the tables in this archetype.
	Capacity int
	// Table statistics.
	Tables []Table
}

// Table statistics.
type Table struct {
	// Number of entities in the table.
	Size int
	// Capacity of the table.
	Capacity int
	// Total reserved memory for entities and components, in bytes.
	Memory int
}

func (s *World) String() string {
	b := strings.Builder{}

	fmt.Fprintf(
		&b, "World -- Components: %d, Archetypes: %d, Filters: %d, Memory: %.1f kB, Locked: %t\n",
		s.ComponentCount, len(s.Archetypes), s.CachedFilters, float64(s.Memory)/1024.0, s.Locked,
	)

	typeNames := make([]string, len(s.ComponentTypes))
	for i, tp := range s.ComponentTypes {
		typeNames[i] = tp.Name()
	}
	fmt.Fprintf(&b, "  Components: %s\n", strings.Join(typeNames, ", "))
	fmt.Fprint(&b, s.Entities.String())

	for i := range s.Archetypes {
		fmt.Fprint(&b, s.Archetypes[i].String())
	}

	return b.String()
}

func (s *Entities) String() string {
	return fmt.Sprintf("Entities -- Used: %d, Recycled: %d, Total: %d, Capacity: %d\n", s.Used, s.Recycled, s.Total, s.Capacity)
}

func (s *Archetype) String() string {
	typeNames := make([]string, len(s.ComponentTypes))
	for i, tp := range s.ComponentTypes {
		typeNames[i] = tp.Name()
	}

	return fmt.Sprintf(
		"Archetype -- Components: %2d, Entities: %6d, Capacity: %6d, Memory: %7.1f kB, Per entity: %4d B\n  Components: %s\n",
		s.Components, s.Size, s.Capacity, float64(s.Memory)/1024.0, s.MemoryPerEntity, strings.Join(typeNames, ", "),
	)
}
