package ecs

// ID is the component identifier.
type ID struct {
	id uint8
}

func id(id int) ID {
	return ID{uint8(id)}
}

// ResID is the resource identifier type.
type ResID struct {
	id uint8
}
