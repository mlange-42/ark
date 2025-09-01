package ecs

import (
	"reflect"
	"testing"
)

func TestResources(t *testing.T) {
	res := newResources()

	posIDint, _ := res.registry.ComponentID(reflect.TypeOf(Position{}))
	rotIDint, _ := res.registry.ComponentID(reflect.TypeOf(Heading{}))
	posID := ResID{id: posIDint}
	rotID := ResID{id: rotIDint}

	expectFalse(t, res.Has(posID))
	expectNil(t, res.Get(posID))

	res.Add(posID, &Position{1, 2})

	expectTrue(t, res.Has(posID))
	pos, ok := res.Get(posID).(*Position)
	expectTrue(t, ok)
	expectEqual(t, Position{1, 2}, *pos)

	expectPanicsWithValue(t, "Resource of ID 0 was already added (type *ecs.Position)", func() { res.Add(posID, &Position{1, 2}) })

	pos, ok = res.Get(posID).(*Position)
	expectTrue(t, ok)
	expectEqual(t, Position{1, 2}, *pos)

	res.Add(rotID, &Heading{5})
	expectTrue(t, res.Has(rotID))
	res.Remove(rotID)
	expectFalse(t, res.Has(rotID))
	expectPanicsWithValue(t, "Resource of ID 1 is not present", func() { res.Remove(rotID) })
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
	expectTrue(t, ok)
	expectEqual(t, Position{1, 2}, *pos)

	rot, ok := res.Get(rotID).(*Heading)
	expectTrue(t, ok)
	expectEqual(t, Heading{5}, *rot)

	res.reset()

	expectFalse(t, res.Has(posID))
	expectFalse(t, res.Has(rotID))

	res.Add(posID, &Position{10, 20})
	res.Add(rotID, &Heading{50})

	pos, ok = res.Get(posID).(*Position)
	expectTrue(t, ok)
	expectEqual(t, Position{10, 20}, *pos)

	rot, ok = res.Get(rotID).(*Heading)
	expectTrue(t, ok)
	expectEqual(t, Heading{50}, *rot)
}
