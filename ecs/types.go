package ecs

// ID is the component identifier.
// It is not relevant when using the default generic API.
type ID struct {
	id uint8
}

func id(id int) ID {
	return ID{uint8(id)}
}

func id8(id uint8) ID {
	return ID{id}
}

// ResID is the resource identifier.
// It is not relevant when using the default generic API.
type ResID struct {
	id uint8
}

// Batch is like a filter for batch processing of entities.
// Create it using [Filter2.Batch] etc.
type Batch struct {
	mask      Mask
	without   Mask
	relations []relationID
}
