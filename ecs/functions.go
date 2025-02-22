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
	id := w.componentID(tp)
	return id
}

// Comp is a helper to pass component types to functions and methods.
// Use function [C] to create one.
type Comp struct {
	tp reflect.Type
}

// C creates a [Comp] instance for the given type.
func C[T any]() Comp {
	return Comp{typeOf[T]()}
}

// ResourceID returns the [ResID] for a resource type via generics.
// Registers the type if it is not already registered.
//
// The number of resources per [World] is limited to [MaskTotalBits].
func ResourceID[T any](w *World) ResID {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	return w.resourceID(tp)
}
