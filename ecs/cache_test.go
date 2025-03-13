package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterCache(t *testing.T) {
	world := NewWorld()

	mapper2 := NewMap2[Position, Velocity](&world)
	mapper3 := NewMap3[Position, Velocity, Heading](&world)

	world.NewEntity()
	mapper2.NewEntity(&Position{}, &Velocity{})
	mapper3.NewEntity(&Position{}, &Velocity{}, &Heading{})

	filter2 := NewFilter2[Position, Velocity](&world).Register()
	filter3 := NewFilter3[Position, Velocity, Heading](&world).Register()

	assert.Equal(t, 0, int(filter2.cache))
	assert.Equal(t, 1, int(filter3.cache))

	assert.Equal(t, 2, len(world.storage.getTableIDs(&filter2.filter, filter2.relations)))
	assert.Equal(t, 1, len(world.storage.getTableIDs(&filter3.filter, filter3.relations)))

	assert.PanicsWithValue(t, "filter is already registered, can't register", func() { filter2.Register() })

	filter2.Unregister()
	filter3.Unregister()

	assert.PanicsWithValue(t, "filter is not registered, can't unregister", func() { filter2.Unregister() })
	assert.PanicsWithValue(t, "no filter for id found to unregister", func() { world.storage.unregisterFilter(100) })
}

func TestFilterCacheRelation(t *testing.T) {
	world := NewWorld()

	posMap := NewMap1[Position](&world)
	childMap := NewMap1[ChildOf](&world)
	child2Map := NewMap1[ChildOf2](&world)
	posChildMap := NewMap2[Position, ChildOf](&world)

	target1 := world.NewEntity()
	target2 := world.NewEntity()
	target3 := world.NewEntity()
	target4 := world.NewEntity()

	f1 := NewFilter1[ChildOf](&world).Register()
	f2 := NewFilter1[ChildOf](&world).Relations(RelIdx(0, target1)).Register()
	f3 := NewFilter1[ChildOf](&world).Relations(RelIdx(0, target2)).Register()

	c1 := world.storage.getRegisteredFilter(f1.cache)
	c2 := world.storage.getRegisteredFilter(f2.cache)
	c3 := world.storage.getRegisteredFilter(f3.cache)

	posMap.NewBatch(10, &Position{})

	assert.Equal(t, 0, len(c1.tables))
	assert.Equal(t, 0, len(c2.tables))
	assert.Equal(t, 0, len(c3.tables))

	e1 := childMap.NewEntity(&ChildOf{}, RelIdx(0, target1))
	assert.Equal(t, 1, len(c1.tables))
	assert.Equal(t, 1, len(c2.tables))

	childMap.NewEntity(&ChildOf{}, RelIdx(0, target3))
	assert.Equal(t, 2, len(c1.tables))
	assert.Equal(t, 1, len(c2.tables))

	child2Map.NewEntity(&ChildOf2{}, RelIdx(0, target2))

	world.RemoveEntity(e1)
	world.RemoveEntity(target1)
	assert.Equal(t, 1, len(c1.tables))
	assert.Equal(t, 0, len(c2.tables))

	childMap.NewEntity(&ChildOf{}, RelIdx(0, target2))
	posChildMap.NewEntity(&Position{}, &ChildOf{}, RelIdx(1, target2))
	posChildMap.NewEntity(&Position{}, &ChildOf{}, RelIdx(1, target3))
	posChildMap.NewEntity(&Position{}, &ChildOf{}, RelIdx(1, target4))

	assert.Equal(t, 5, len(c1.tables))
	assert.Equal(t, 2, len(c3.tables))

	world.RemoveEntities(NewFilter0(&world).Batch(), nil)
	assert.Equal(t, 0, len(c1.tables))
	assert.Equal(t, 0, len(c2.tables))
}
