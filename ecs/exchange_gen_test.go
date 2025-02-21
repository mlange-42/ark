package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExchange2(t *testing.T) {
	w := NewWorld(1024)

	builder := NewMap2[Position, Velocity](&w)
	headMap := NewMap[Heading](&w)
	ex := NewExchange1[Heading](&w).Removes(C[Velocity](), C[Position]())

	e := builder.NewEntity(&Position{}, &Velocity{})

	ex.Exchange(e, &Heading{})
	assert.True(t, headMap.Has(e))
}
