package ecs

import (
	"reflect"
	"testing"
)

func TestCompResIDs(t *testing.T) {
	w := NewWorld(1024)

	posID := w.ComponentID[Position]()
	rotID := w.ComponentID[Velocity]()

	tPosID := TypeID(w, reflect.TypeOf(Position{}))
	tRotID := TypeID(w, reflect.TypeOf(Velocity{}))

	res1ID := ResourceID[Position](w)
	res2ID := ResourceID[Velocity](w)

	expectEqual(t, posID, tPosID)
	expectEqual(t, rotID, tRotID)

	expectEqual(t, uint8(0), posID.id)
	expectEqual(t, uint8(1), rotID.id)

	expectEqual(t, uint8(0), res1ID.id)
	expectEqual(t, uint8(1), res2ID.id)

	expectSlicesEqual(t, []ID{id(0), id(1)}, w.ComponentIDs())
	expectSlicesEqual(t, []ResID{{id: 0}, {id: 1}}, ResourceIDs(w))
}

func TestRegisterComponents(t *testing.T) {
	world := NewWorld(1024)

	world.ComponentID[Position]()

	expectEqual(t, id(0), world.ComponentID[Position]())
	expectEqual(t, id(1), world.ComponentID[Velocity]())

	world.lock()
	expectPanicsWithValue(t,
		"attempt to register a new component in a locked world",
		func() {
			world.ComponentID[Heading]()
		})
}

func TestComponentInfo(t *testing.T) {
	w := NewWorld(1024)
	_ = w.ComponentID[Velocity]()
	posID := w.ComponentID[Position]()

	info, ok := w.ComponentInfo(posID)
	expectTrue(t, ok)
	expectEqual(t, info.Type, reflect.TypeOf(Position{}))

	info, ok = w.ComponentInfo(ID{id: 3})
	expectFalse(t, ok)
	expectEqual(t, info, CompInfo{})

	resID := ResourceID[Velocity](w)

	tp, ok := ResourceType(w, resID)
	expectTrue(t, ok)
	expectEqual(t, tp, reflect.TypeOf(Velocity{}))

	tp, ok = ResourceType(w, ResID{id: 3})
	expectFalse(t, ok)
	expectEqual(t, tp, nil)
}

func TestCompType(t *testing.T) {
	c := C[Position]()
	expectEqual(t, reflect.TypeFor[Position](), c.Type())
}

func TestResourceTypeID(t *testing.T) {
	w := NewWorld(1024)
	id1 := ResourceTypeID(w, reflect.TypeFor[Position]())
	id2 := ResourceTypeID(w, reflect.TypeFor[Velocity]())
	id3 := ResourceTypeID(w, reflect.TypeFor[Position]())

	expectEqual(t, 0, id1.id)
	expectEqual(t, 1, id2.id)
	expectEqual(t, 0, id3.id)
}

func TestResourceShortcuts(t *testing.T) {
	w := NewWorld(1024)
	res := Position{1, 2}
	AddResource(w, &res)

	res2 := GetResource[Position](w)
	expectEqual(t, res, *res2)

	res3 := GetResource[Velocity](w)
	expectNil(t, res3)
}

func BenchmarkComponentID(b *testing.B) {

	world := NewWorld(1024)
	id := world.ComponentID[Position]()

	for b.Loop() {
		id = world.ComponentID[Position]()
	}
	_ = id
}

func BenchmarkTypeID(b *testing.B) {

	world := NewWorld(1024)
	id := world.ComponentID[Position]()
	info, _ := world.ComponentInfo(id)
	tp := info.Type

	for b.Loop() {
		id = TypeID(world, tp)
	}
	_ = id
}
