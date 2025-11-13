package ecs

import (
	"reflect"
	"testing"
)

func TestCompResIDs(t *testing.T) {
	w := NewWorld(1024)

	posID := ComponentID[Position](w)
	rotID := ComponentID[Velocity](w)

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

	expectSlicesEqual(t, []ID{id(0), id(1)}, ComponentIDs(w))
	expectSlicesEqual(t, []ResID{{id: 0}, {id: 1}}, ResourceIDs(w))
}

func TestRegisterComponents(t *testing.T) {
	world := NewWorld(1024)

	ComponentID[Position](world)

	expectEqual(t, id(0), ComponentID[Position](world))
	expectEqual(t, id(1), ComponentID[Velocity](world))

	world.lock()
	expectPanicsWithValue(t,
		"attempt to register a new component in a locked world",
		func() {
			ComponentID[Heading](world)
		})
}

func TestComponentInfo(t *testing.T) {
	w := NewWorld(1024)
	_ = ComponentID[Velocity](w)
	posID := ComponentID[Position](w)

	info, ok := ComponentInfo(w, posID)
	expectTrue(t, ok)
	expectEqual(t, info.Type, reflect.TypeOf(Position{}))

	info, ok = ComponentInfo(w, ID{id: 3})
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
	id := ComponentID[Position](world)

	for b.Loop() {
		id = ComponentID[Position](world)
	}
	_ = id
}

func BenchmarkTypeID(b *testing.B) {

	world := NewWorld(1024)
	id := ComponentID[Position](world)
	info, _ := ComponentInfo(world, id)
	tp := info.Type

	for b.Loop() {
		id = TypeID(world, tp)
	}
	_ = id
}
