package ecs

// config provides configuration for an ECS [World].
type config struct {
	// Initial capacity for archetypes and the entity index.
	// The default value is 128.
	initialCapacity int
	// Initial capacity for archetypes with a relation component.
	// The default value is initialCapacity.
	initialCapacityRelations int
}

// newConfig creates a new default [World] configuration.
func newConfig(initialCapacity ...int) config {
	switch len(initialCapacity) {
	case 0:
		return config{
			initialCapacity:          1024,
			initialCapacityRelations: 128,
		}
	case 1:
		if initialCapacity[0] < 1 {
			panic("only positive values for the World's initialCapacity are allowed")
		}
		return config{
			initialCapacity:          initialCapacity[0],
			initialCapacityRelations: initialCapacity[0],
		}
	case 2:
		if initialCapacity[0] < 1 || initialCapacity[1] < 1 {
			panic("only positive values for the World's initialCapacity are allowed")
		}
		return config{
			initialCapacity:          initialCapacity[0],
			initialCapacityRelations: initialCapacity[1],
		}
	}
	panic("can only use a maximum of two values for the World's initialCapacity")
}
