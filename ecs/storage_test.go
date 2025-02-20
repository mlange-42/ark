package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	storage := storage{}
	assert.Equal(t, 0, len(storage.archetypes))
	assert.Equal(t, 0, len(storage.tables))
}
