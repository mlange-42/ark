package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResource(t *testing.T) {
	w := NewWorld(1024)
	get := NewResource[Position](&w)

	assert.False(t, get.Has())
	get.Add(&Position{100, 101})

	assert.True(t, get.Has())
	res := get.Get()

	assert.Equal(t, Position{100, 101}, *res)

	get.Remove()
	assert.False(t, get.Has())
}

func ExampleResource() {
	world := NewWorld(1024)
	resAccess := NewResource[Position](&world)

	resAccess.Add(&Position{})

	res := resAccess.Get()
	res.X, res.Y = 10, 5
	// Output:
}
