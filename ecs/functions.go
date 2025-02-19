package ecs

import "reflect"

// ComponentID returns the [ID] for a component type via generics.
// Registers the type if it is not already registered.
//
// The number of unique component types per [World] is limited to 256 ([MaskTotalBits]).
//
// Panics if called on a locked world and the type is not registered yet.
//
// Note that type aliases are not considered separate component types.
// Type re-definitions, however, are separate types.
//
// ⚠️ Warning: Using IDs that are outside of the range of registered IDs anywhere in [World] or other places will result in undefined behavior!
func ComponentID[T any](w *World) ID {
	tp := reflect.TypeOf((*T)(nil)).Elem()

	id, newID := w.registry.ComponentID(tp)
	if newID {
		// TODO: check lock and unroll
		//if w.IsLocked() {
		//	w.registry.unregisterLastComponent()
		//	panic("attempt to register a new component in a locked world")
		//}
	}
	return id
}
