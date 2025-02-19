package ecs

type World struct {
	registry registry
}

func NewWorld() World {
	return World{
		registry: newRegistry(),
	}
}
