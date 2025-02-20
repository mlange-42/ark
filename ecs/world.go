package ecs

// World is the central type holding entity and component data, as well as resources.
type World struct {
	registry registry
	storage  storage
}

// NewWorld creates a new [World].
func NewWorld() World {
	return World{
		registry: newRegistry(),
		storage:  newStorage(),
	}
}
