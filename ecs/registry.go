package ecs

import (
	"fmt"
	"math"
	"reflect"
)

// componentRegistry keeps track of type IDs.
type registry struct {
	Components map[reflect.Type]uint8 // Mapping from types to IDs.
	Types      []reflect.Type         // Mapping from IDs to types.
	IDs        []uint8                // List of IDs.
	Used       bitMask                // Mapping from IDs to used status.
}

// newComponentRegistry creates a new ComponentRegistry.
func newRegistry() registry {
	return registry{
		Components: map[reflect.Type]uint8{},
		Types:      make([]reflect.Type, maskTotalBits),
		Used:       bitMask{},
		IDs:        []uint8{},
	}
}

// ComponentID returns the ID for a component type, and registers it if not already registered.
// The second return value indicates if it is a newly created ID.
func (r *registry) ComponentID(tp reflect.Type) (uint8, bool) {
	if id, ok := r.Components[tp]; ok {
		return id, false
	}
	return r.registerComponent(tp, maskTotalBits), true
}

// ComponentType returns the type of a component by ID.
func (r *registry) ComponentType(id uint8) (reflect.Type, bool) {
	return r.Types[id], r.Used.Get(ID{id: id})
}

// Count returns the total number of reserved IDs. It is the maximum ID plus 1.
func (r *registry) Count() int {
	return len(r.Components)
}

// registerComponent registers a components and assigns an ID for it.
func (r *registry) registerComponent(tp reflect.Type, totalBits int) uint8 {
	val := len(r.Components)
	if val >= totalBits {
		panic(fmt.Sprintf("exceeded the maximum of %d component types or resource types", totalBits))
	}
	newID := uint8(val)
	id := id(val)
	r.Components[tp], r.Types[newID] = newID, tp
	r.Used.Set(id, true)
	r.IDs = append(r.IDs, newID)
	return newID
}

func (r *registry) unregisterLastComponent() {
	newID := uint8(len(r.Components) - 1)
	id := id8(newID)
	tp, _ := r.ComponentType(newID)
	delete(r.Components, tp)
	r.Types[newID] = nil
	r.Used.Set(id, false)
	r.IDs = r.IDs[:len(r.IDs)-1]
}

// componentRegistry keeps track of component IDs.
// In addition to [registry], it determines whether types
// are relation components and/or contain (or are) pointers.
type componentRegistry struct {
	registry
	IsRelation []bool
	IsTrivial  []bool
	Archetypes []int  // Number of archetypes for each component.
	generation uint32 // Generation to indicate changes to archetype count per component.
}

// newComponentRegistry creates a new ComponentRegistry.
func newComponentRegistry() componentRegistry {
	return componentRegistry{
		registry:   newRegistry(),
		IsRelation: make([]bool, maskTotalBits),
		IsTrivial:  make([]bool, maskTotalBits),
		Archetypes: make([]int, maskTotalBits),
		generation: 1,
	}
}

// ComponentID returns the ID for a component type, and registers it if not already registered.
// The second return value indicates if it is a newly created ID.
func (r *componentRegistry) ComponentID(tp reflect.Type) (uint8, bool) {
	if id, ok := r.Components[tp]; ok {
		return id, false
	}
	return r.registerComponent(tp, maskTotalBits), true
}

// registerComponent registers a components and assigns an ID for it.
func (r *componentRegistry) registerComponent(tp reflect.Type, totalBits int) uint8 {
	newID := r.registry.registerComponent(tp, totalBits)
	r.IsRelation[newID] = isRelation(tp)
	r.IsTrivial[newID] = isTrivial(tp)
	return newID
}

func (r *componentRegistry) unregisterLastComponent() {
	newID := uint8(len(r.Components) - 1)
	r.registry.unregisterLastComponent()
	r.IsRelation[newID] = false
}

func (r *componentRegistry) addArchetype(id uint8) {
	r.Archetypes[id]++
	r.generation++
}

func (r *componentRegistry) getGeneration() uint32 {
	return r.generation
}

// Returns the ID of the component present in the smallest number of archetypes.
func (r *componentRegistry) rareComponent(ids []ID) ID {
	minCount := math.MaxInt
	var rareID ID
	for _, id := range ids {
		count := r.Archetypes[id.id]
		if count < minCount {
			minCount = count
			rareID = id
		}
	}
	return rareID
}
