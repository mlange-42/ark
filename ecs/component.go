package ecs

import "reflect"

// ID is the component identifier.
type ID struct {
	id uint32
}

type componentInfo struct {
	id       ID
	typ      reflect.Type
	itemSize uintptr
}
