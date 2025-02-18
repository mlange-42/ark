package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntity(t *testing.T) {
	e := newEntity(100)
	assert.Equal(t, 100, e.id)
	assert.Equal(t, 0, e.gen)
}
