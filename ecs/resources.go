package ecs

import (
	"fmt"
	"reflect"
)

// Resources manage a world's resources. Access it using [World.Resources].
//
// Although this type provides an ID-based API, the recommended usage is via [Resource].
type Resources struct {
	resources []any
	registry  registry
}

// newResources creates a new Resources manager.
func newResources() Resources {
	return Resources{
		registry:  newRegistry(),
		resources: make([]any, maskTotalBits),
	}
}

// Add a resource to the world.
// The resource should always be a pointer.
//
// Panics if there is already a resource of the given type.
//
// See [Resource.Add] for the recommended type-safe way.
func (r *Resources) Add(id ResID, res any) {
	if r.resources[id.id] != nil {
		panic(fmt.Sprintf("Resource of ID %d was already added (type %v)", id.id, reflect.TypeOf(res)))
	}
	r.resources[id.id] = res
}

// Remove a resource from the world.
//
// Panics if there is no resource of the given type.
//
// See [Resource.Remove] for the recommended type-safe way.
func (r *Resources) Remove(id ResID) {
	if r.resources[id.id] == nil {
		panic(fmt.Sprintf("Resource of ID %d is not present", id.id))
	}
	r.resources[id.id] = nil
}

// Get returns a pointer to the resource of the given type.
//
// Returns nil if there is no such resource.
//
// See [Resource.Get] for the recommended type-safe way.
func (r *Resources) Get(id ResID) interface{} {
	return r.resources[id.id]
}

// Has returns whether the world has the given resource.
//
// See [Resource.Has] for the recommended type-safe way.
func (r *Resources) Has(id ResID) bool {
	return r.resources[id.id] != nil
}

// reset removes all resources.
func (r *Resources) reset() {
	for i := range r.resources {
		r.resources[i] = nil
	}
}
