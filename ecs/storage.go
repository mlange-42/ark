package ecs

type storage struct {
	archetypes []archetype
	tables     []table
}

func newStorage() storage {
	return storage{}
}
