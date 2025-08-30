package ecs

import (
	"reflect"
	"testing"
)

func TestCompResIDs(t *testing.T) {
	w := NewWorld(1024)

	posID := ComponentID[Position](&w)
	rotID := ComponentID[Velocity](&w)

	tPosID := TypeID(&w, reflect.TypeOf(Position{}))
	tRotID := TypeID(&w, reflect.TypeOf(Velocity{}))

	res1ID := ResourceID[Position](&w)
	res2ID := ResourceID[Velocity](&w)

	if posID != tPosID {
		t.Errorf("expected posID == tPosID, got %v vs %v", posID, tPosID)
	}
	if rotID != tRotID {
		t.Errorf("expected rotID == tRotID, got %v vs %v", rotID, tRotID)
	}

	if posID.id != 0 {
		t.Errorf("expected posID.id == 0, got %d", posID.id)
	}
	if rotID.id != 1 {
		t.Errorf("expected rotID.id == 1, got %d", rotID.id)
	}

	if res1ID.id != 0 {
		t.Errorf("expected res1ID.id == 0, got %d", res1ID.id)
	}
	if res2ID.id != 1 {
		t.Errorf("expected res2ID.id == 1, got %d", res2ID.id)
	}

	expectedCompIDs := []ID{id(0), id(1)}
	if got := ComponentIDs(&w); !reflect.DeepEqual(got, expectedCompIDs) {
		t.Errorf("expected ComponentIDs %v, got %v", expectedCompIDs, got)
	}

	expectedResIDs := []ResID{{id: 0}, {id: 1}}
	if got := ResourceIDs(&w); !reflect.DeepEqual(got, expectedResIDs) {
		t.Errorf("expected ResourceIDs %v, got %v", expectedResIDs, got)
	}
}

func TestRegisterComponents(t *testing.T) {
	world := NewWorld(1024)

	ComponentID[Position](&world)

	if ComponentID[Position](&world) != id(0) {
		t.Errorf("expected Position ID == 0")
	}
	if ComponentID[Velocity](&world) != id(1) {
		t.Errorf("expected Velocity ID == 1")
	}

	world.lock()
	expectPanicWithValue(t,
		"attempt to register a new component in a locked world",
		func() {
			ComponentID[Heading](&world)
		})
}

func TestComponentInfo(t *testing.T) {
	w := NewWorld(1024)
	_ = ComponentID[Velocity](&w)
	posID := ComponentID[Position](&w)

	info, ok := ComponentInfo(&w, posID)
	if !ok {
		t.Errorf("expected ComponentInfo to return ok=true for posID")
	}
	if info.Type != reflect.TypeOf(Position{}) {
		t.Errorf("expected info.Type == Position, got %v", info.Type)
	}

	info, ok = ComponentInfo(&w, ID{id: 3})
	if ok {
		t.Errorf("expected ComponentInfo to return ok=false for unknown ID")
	}
	if info != (CompInfo{}) {
		t.Errorf("expected empty CompInfo for unknown ID, got %v", info)
	}

	resID := ResourceID[Velocity](&w)
	tp, ok := ResourceType(&w, resID)
	if !ok {
		t.Errorf("expected ResourceType to return ok=true for resID")
	}
	if tp != reflect.TypeOf(Velocity{}) {
		t.Errorf("expected ResourceType == Velocity, got %v", tp)
	}

	tp, ok = ResourceType(&w, ResID{id: 3})
	if ok {
		t.Errorf("expected ResourceType to return ok=false for unknown ResID")
	}
	if tp != nil {
		t.Errorf("expected nil type for unknown ResID, got %v", tp)
	}
}

func TestCompType(t *testing.T) {
	c := C[Position]()
	if c.Type() != reflect.TypeFor[Position]() {
		t.Errorf("expected component type to be Position, got %v", c.Type())
	}
}

func TestResourceTypeID(t *testing.T) {
	w := NewWorld(1024)
	id1 := ResourceTypeID(&w, reflect.TypeFor[Position]())
	id2 := ResourceTypeID(&w, reflect.TypeFor[Velocity]())
	id3 := ResourceTypeID(&w, reflect.TypeFor[Position]())

	if id1.id != 0 {
		t.Errorf("expected id1.id == 0, got %d", id1.id)
	}
	if id2.id != 1 {
		t.Errorf("expected id2.id == 1, got %d", id2.id)
	}
	if id3.id != 0 {
		t.Errorf("expected id3.id == 0, got %d", id3.id)
	}
}

func TestResourceShortcuts(t *testing.T) {
	w := NewWorld(1024)
	res := Position{1, 2}
	AddResource(&w, &res)

	res2 := GetResource[Position](&w)
	if *res2 != res {
		t.Errorf("expected GetResource to return %v, got %v", res, *res2)
	}
}

func BenchmarkComponentID(b *testing.B) {
	world := NewWorld(1024)
	id := ComponentID[Position](&world)

	for b.Loop() {
		id = ComponentID[Position](&world)
	}
	_ = id
}

func BenchmarkTypeID(b *testing.B) {
	world := NewWorld(1024)
	id := ComponentID[Position](&world)
	info, _ := ComponentInfo(&world, id)
	tp := info.Type

	for b.Loop() {
		id = TypeID(&world, tp)
	}
	_ = id
}
