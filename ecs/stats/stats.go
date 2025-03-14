// Package stats provides the structs returned by ecs.World.Stats().
package stats

import (
	"fmt"
	"reflect"
	"strings"
)

// World statistics.
type World struct {
	// Component types, indexed by component ID.
	ComponentTypes []reflect.Type
	// Archetype statistics.
	Archetypes []Archetype
	// Entity statistics.
	Entities Entities
	// Memory reserved for entities and components, in bytes.
	Memory int
	// Memory actually used for alive entities and their components components, in bytes.
	MemoryUsed int
	// Number of cached filters.
	CachedFilters int
	// Locked state of the world.
	Locked bool
}

// Entities statistics.
type Entities struct {
	// Currently used/alive entities.
	Used int
	// Recycled/available entities.
	Recycled int
	// Current total number of entities in the pool (used + recycled).
	Total int
	// Current capacity of entity pool and entity list.
	Capacity int
}

// Archetype statistics.
type Archetype struct {
	// Component IDs.
	ComponentIDs []uint8
	// Component types for ComponentIDs.
	ComponentTypes []reflect.Type
	// Table statistics.
	Tables []Table
	// Number of entities in the tables of this archetype.
	Size int
	// Sum of capacity of the tables in this archetype.
	Capacity int
	// Number of relation components in the archetype.
	NumRelations int
	// Memory reserved for entities and components, in bytes.
	Memory int
	// Memory actually used for alive entities and their components components, in bytes.
	MemoryUsed int
	// Memory for components per entity, in bytes.
	MemoryPerEntity int
	// Number of free tables.
	FreeTables int
}

// Table statistics.
type Table struct {
	// Number of entities in the table.
	Size int
	// Capacity of the table.
	Capacity int
	// Memory reserved for entities and components, in bytes.
	Memory int
	// Memory actually used for alive entities and their components components, in bytes.
	MemoryUsed int
}

func (w *World) String() string {
	b := strings.Builder{}
	fmt.Fprintf(
		&b, "World     -- Components: %d, Archetypes: %d, Filters: %d, Memory: %.1f kB, Locked: %t\n",
		len(w.ComponentTypes), len(w.Archetypes), w.CachedFilters, float64(w.Memory)/1024.0, w.Locked,
	)
	typeNames := make([]string, len(w.ComponentTypes))
	for i, tp := range w.ComponentTypes {
		typeNames[i] = tp.Name()
	}
	fmt.Fprintf(&b, "             Components: %s\n", strings.Join(typeNames, ", "))
	fmt.Fprint(&b, w.Entities.String())
	for i := range w.Archetypes {
		fmt.Fprint(&b, w.Archetypes[i].String())
	}
	return b.String()
}

func (e *Entities) String() string {
	return fmt.Sprintf("Entities  -- Used: %d, Recycled: %d, Total: %d, Capacity: %d\n", e.Used, e.Recycled, e.Total, e.Capacity)
}

func (a *Archetype) String() string {
	typeNames := make([]string, len(a.ComponentTypes))
	for i, tp := range a.ComponentTypes {
		typeNames[i] = tp.Name()
	}
	return fmt.Sprintf(
		"Archetype -- Tables: %4d, Comps: %2d, Entities: %6d, Cap: %6d, Mem: %7.1f kB, Per entity: %4d B\n             Components: %s\n",
		len(a.Tables), len(a.ComponentIDs), a.Size, a.Capacity, float64(a.Memory)/1024.0, a.MemoryPerEntity, strings.Join(typeNames, ", "),
	)
}

func (t *Table) String() string {
	return fmt.Sprintf("Table     -- Entities: %6d, Cap: %6d, Mem: %7.1f kB\n", t.Size, t.Capacity, float64(t.Memory)/1024.0)
}
