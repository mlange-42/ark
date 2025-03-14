package ecs

// Resource provides access to a world resource.
//
// Create one with [NewResource].
type Resource[T any] struct {
	world *World
	id    ResID
}

// NewResource creates a new [Resource] mapper for a resource type.
// This does not add a resource to the world, but only creates a mapper for resource access!
//
// See also [World.Resources].
func NewResource[T any](w *World) Resource[T] {
	return Resource[T]{
		id:    ResourceID[T](w),
		world: w,
	}
}

// Add adds a resource to the world.
//
// Panics if there is already a resource of the given type.
func (g *Resource[T]) Add(res *T) {
	g.world.Resources().Add(g.id, res)
}

// Remove removes a resource from the world.
//
// Panics if there is no resource of the given type.
//
// See also [ecs.Resources.Remove].
func (g *Resource[T]) Remove() {
	g.world.Resources().Remove(g.id)
}

// Get gets the resource of the given type.
//
// Returns nil if there is no such resource.
func (g *Resource[T]) Get() *T {
	return g.world.Resources().Get(g.id).(*T)
}

// Has returns whether the world has the resource type.
func (g *Resource[T]) Has() bool {
	return g.world.Resources().Has(g.id)
}
