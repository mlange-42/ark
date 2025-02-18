package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeOf(t *testing.T) {
	posType := typeOf[Position]()
	assert.Equal(t, "Position", posType.Name())
}
