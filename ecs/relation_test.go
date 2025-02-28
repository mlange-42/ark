package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRel(t *testing.T) {
	r := Rel(1, Entity{5, 0})
	assert.Equal(t, RelationIndex{1, Entity{5, 0}}, r)
}

func TestRelID(t *testing.T) {
	r := RelID(ID{10}, Entity{5, 0})
	assert.Equal(t, RelationID{ID{10}, Entity{5, 0}}, r)
}
