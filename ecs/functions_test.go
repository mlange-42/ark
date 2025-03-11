package ecs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompResIDs(t *testing.T) {
	w := NewWorld(1024)

	posID := ComponentID[Position](&w)
	rotID := ComponentID[Velocity](&w)

	tPosID := TypeID(&w, reflect.TypeOf(Position{}))
	tRotID := TypeID(&w, reflect.TypeOf(Velocity{}))

	res1ID := ResourceID[Position](&w)
	res2ID := ResourceID[Velocity](&w)

	assert.Equal(t, posID, tPosID)
	assert.Equal(t, rotID, tRotID)

	assert.Equal(t, uint8(0), posID.id)
	assert.Equal(t, uint8(1), rotID.id)

	assert.Equal(t, uint8(0), res1ID.id)
	assert.Equal(t, uint8(1), res2ID.id)

	assert.Equal(t, []ID{id(0), id(1)}, ComponentIDs(&w))
	assert.Equal(t, []ResID{{id: 0}, {id: 1}}, ResourceIDs(&w))
}

func TestRegisterComponents(t *testing.T) {
	world := NewWorld(1024)

	ComponentID[Position](&world)

	assert.Equal(t, id(0), ComponentID[Position](&world))
	assert.Equal(t, id(1), ComponentID[Velocity](&world))

	world.lock()
	assert.PanicsWithValue(t,
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
	assert.True(t, ok)
	assert.Equal(t, info.Type, reflect.TypeOf(Position{}))

	info, ok = ComponentInfo(&w, ID{id: 3})
	assert.False(t, ok)
	assert.Equal(t, info, CompInfo{})

	resID := ResourceID[Velocity](&w)

	tp, ok := ResourceType(&w, resID)
	assert.True(t, ok)
	assert.Equal(t, tp, reflect.TypeOf(Velocity{}))

	tp, ok = ResourceType(&w, ResID{id: 3})
	assert.False(t, ok)
	assert.Equal(t, tp, nil)
}

func TestCompType(t *testing.T) {
	c := C[Position]()
	assert.Equal(t, typeOf[Position](), c.Type())
}

func TestResourceTypeID(t *testing.T) {
	w := NewWorld(1024)
	id1 := ResourceTypeID(&w, typeOf[Position]())
	id2 := ResourceTypeID(&w, typeOf[Velocity]())
	id3 := ResourceTypeID(&w, typeOf[Position]())

	assert.EqualValues(t, 0, id1.id)
	assert.EqualValues(t, 1, id2.id)
	assert.EqualValues(t, 0, id3.id)
}

func TestResourceShortcuts(t *testing.T) {
	w := NewWorld(1024)
	res := Position{1, 2}
	AddResource(&w, &res)

	res2 := GetResource[Position](&w)
	assert.Equal(t, res, *res2)
}

func BenchmarkComponentID(b *testing.B) {
	b.StopTimer()
	world := NewWorld(1024)
	id := ComponentID[Position](&world)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		id = ComponentID[Position](&world)
	}
	_ = id
}

func BenchmarkTypeID(b *testing.B) {
	b.StopTimer()
	world := NewWorld(1024)
	id := ComponentID[Position](&world)
	info, _ := ComponentInfo(&world, id)
	tp := info.Type

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		id = TypeID(&world, tp)
	}
	_ = id
}
