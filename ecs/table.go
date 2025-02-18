package ecs

import "unsafe"

type table struct {
	data    []column
	columns []unsafe.Pointer
}
