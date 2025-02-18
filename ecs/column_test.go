package ecs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColumnPointer(t *testing.T) {
	posType := reflect.TypeOf(Position{})
	column := newColumn(posType, 8)

	assert.Equal(t, uintptr(column.pointer), uintptr(column.data.Addr().UnsafePointer()))
}
