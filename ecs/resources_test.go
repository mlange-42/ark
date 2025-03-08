package ecs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResources(t *testing.T) {
	res := newResources()

	posIDint, _ := res.registry.ComponentID(reflect.TypeOf(Position{}))
	rotIDint, _ := res.registry.ComponentID(reflect.TypeOf(Heading{}))
	posID := ResID{id: posIDint}
	rotID := ResID{id: rotIDint}

	assert.False(t, res.Has(posID))
	assert.Nil(t, res.Get(posID))

	res.Add(posID, &Position{1, 2})

	assert.True(t, res.Has(posID))
	pos, ok := res.Get(posID).(*Position)
	assert.True(t, ok)
	assert.Equal(t, Position{1, 2}, *pos)

	assert.PanicsWithValue(t, "Resource of ID 0 was already added (type *ecs.Position)", func() { res.Add(posID, &Position{1, 2}) })

	pos, ok = res.Get(posID).(*Position)
	assert.True(t, ok)
	assert.Equal(t, Position{1, 2}, *pos)

	res.Add(rotID, &Heading{5})
	assert.True(t, res.Has(rotID))
	res.Remove(rotID)
	assert.False(t, res.Has(rotID))
	assert.PanicsWithValue(t, "Resource of ID 1 is not present", func() { res.Remove(rotID) })
}

func TestResourcesReset(t *testing.T) {
	res := newResources()

	posIDint, _ := res.registry.ComponentID(reflect.TypeOf(Position{}))
	rotIDint, _ := res.registry.ComponentID(reflect.TypeOf(Heading{}))
	posID := ResID{id: posIDint}
	rotID := ResID{id: rotIDint}

	res.Add(posID, &Position{1, 2})
	res.Add(rotID, &Heading{5})

	pos, ok := res.Get(posID).(*Position)
	assert.True(t, ok)
	assert.Equal(t, Position{1, 2}, *pos)

	rot, ok := res.Get(rotID).(*Heading)
	assert.True(t, ok)
	assert.Equal(t, Heading{5}, *rot)

	res.reset()

	assert.False(t, res.Has(posID))
	assert.False(t, res.Has(rotID))

	res.Add(posID, &Position{10, 20})
	res.Add(rotID, &Heading{50})

	pos, ok = res.Get(posID).(*Position)
	assert.True(t, ok)
	assert.Equal(t, Position{10, 20}, *pos)

	rot, ok = res.Get(rotID).(*Heading)
	assert.True(t, ok)
	assert.Equal(t, Heading{50}, *rot)
}
