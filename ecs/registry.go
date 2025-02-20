package ecs

import (
	"fmt"
	"reflect"
)

// componentRegistry keeps track of type IDs.
type registry struct {
	Components map[reflect.Type]ID // Mapping from types to IDs.
	Types      []reflect.Type      // Mapping from IDs to types.
	IDs        []ID                // List of IDs.
	Used       Mask                // Mapping from IDs tu used status.
}

// newComponentRegistry creates a new ComponentRegistry.
func newRegistry() registry {
	return registry{
		Components: map[reflect.Type]ID{},
		Types:      make([]reflect.Type, MaskTotalBits),
		Used:       Mask{},
		IDs:        []ID{},
	}
}

// ComponentID returns the ID for a component type, and registers it if not already registered.
// The second return value indicates if it is a newly created ID.
func (r *registry) ComponentID(tp reflect.Type) (ID, bool) {
	if id, ok := r.Components[tp]; ok {
		return id, false
	}
	return r.registerComponent(tp, MaskTotalBits), true
}

// ComponentType returns the type of a component by ID.
func (r *registry) ComponentType(id ID) (reflect.Type, bool) {
	return r.Types[id.id], r.Used.Get(ID{id: id.id})
}

// Count returns the total number of reserved IDs. It is the maximum ID plus 1.
func (r *registry) Count() int {
	return len(r.Components)
}

// Reset clears the registry.
func (r *registry) Reset() {
	for t := range r.Components {
		delete(r.Components, t)
	}
	for i := range r.Types {
		r.Types[i] = nil
	}
	r.Used.Reset()
	r.IDs = r.IDs[:0]
}

// registerComponent registers a components and assigns an ID for it.
func (r *registry) registerComponent(tp reflect.Type, totalBits int) ID {
	val := len(r.Components)
	if val >= totalBits {
		panic(fmt.Sprintf("exceeded the maximum of %d component types or resource types", totalBits))
	}
	id := id(val)
	r.Components[tp], r.Types[id.id] = id, tp
	r.Used.Set(id, true)
	r.IDs = append(r.IDs, id)
	return id
}
