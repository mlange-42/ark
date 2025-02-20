package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntity(t *testing.T) {
	e := newEntity(100)
	assert.EqualValues(t, 100, e.id)
	assert.EqualValues(t, 0, e.gen)
}

func TestEntityIndex(t *testing.T) {
	index := entityIndex{}
	assert.EqualValues(t, 0, index.table)
	assert.EqualValues(t, 0, index.row)
}
