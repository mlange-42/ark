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
	if res.Has(posID) {
		t.Errorf("expected resource to not exist")
	}
	if res.Get(posID) != nil {
		t.Errorf("expected resource to be nil")
	}
	res.Add(posID, &Position{1, 2})
	if !res.Has(posID) {
		t.Errorf("expected resource to exist")
	}
	pos, ok := res.Get(posID).(*Position)
	if !ok {
		t.Errorf("expected resource to be of type *Position")
	}
	if *pos != (Position{1, 2}) {
		t.Errorf("expected resource to be %v, got %v", Position{1, 2}, *pos)
	}
	expectPanicWithValue(t, "Resource of ID 0 was already added (type *ecs.Position)", func() { res.Add(posID, &Position{1, 2}) })
	pos, ok = res.Get(posID).(*Position)
	if !ok {
		t.Errorf("expected resource to be of type *Position")
	}
	if *pos != (Position{1, 2}) {
		t.Errorf("expected resource to be %v, got %v", Position{1, 2}, *pos)
	}
	res.Add(rotID, &Heading{5})
	if !res.Has(rotID) {
		t.Errorf("expected resource to exist")
	}
	res.Remove(rotID)
	if res.Has(rotID) {
		t.Errorf("expected resource to not exist")
	}
	expectPanicWithValue(t, "Resource of ID 1 is not present", func() { res.Remove(rotID) })
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
	if !ok {
		t.Errorf("expected resource to be of type *Position")
	}
	if *pos != (Position{1, 2}) {
		t.Errorf("expected resource to be %v, got %v", Position{1, 2}, *pos)
	}
	rot, ok := res.Get(rotID).(*Heading)
	if !ok {
		t.Errorf("expected resource to be of type *Heading")
	}
	if *rot != (Heading{5}) {
		t.Errorf("expected resource to be %v, got %v", Heading{5}, *rot)
	}
	res.reset()
	if res.Has(posID) {
		t.Errorf("expected resource to not exist")
	}
	if res.Has(rotID) {
		t.Errorf("expected resource to not exist")
	}
	res.Add(posID, &Position{10, 20})
	res.Add(rotID, &Heading{50})
	pos, ok = res.Get(posID).(*Position)
	if !ok {
		t.Errorf("expected resource to be of type *Position")
	}
	if *pos != (Position{10, 20}) {
		t.Errorf("expected resource to be %v, got %v", Position{10, 20}, *pos)
	}
	rot, ok = res.Get(rotID).(*Heading)
	if !ok {
		t.Errorf("expected resource to be of type *Heading")
	}
	if *rot != (Heading{50}) {
		t.Errorf("expected resource to be %v, got %v", Heading{50}, *rot)
	}
}
