package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArchetype(t *testing.T) {
	arch := archetype{}
	assert.False(t, arch.HasRelations())
	assert.Equal(t, 0, len(arch.tables.tables))
}
