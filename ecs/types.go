package ecs

// ID is the component identifier.
// It is not relevant when using the default generic API.
type ID struct {
	id uint8
}

func id(id int) ID {
	return ID{uint8(id)}
}

// ResID is the resource identifier.
// It is not relevant when using the default generic API.
type ResID struct {
	id uint8
}
