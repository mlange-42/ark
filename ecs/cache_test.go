package ecs

import (
	"testing"
)

func TestFilterCache(t *testing.T) {
	world := NewWorld()

	mapper2 := NewMap2[Position, Velocity](world)
	mapper3 := NewMap3[Position, Velocity, Heading](world)

	world.NewEntity()
	mapper2.NewEntity(&Position{}, &Velocity{})
	mapper3.NewEntity(&Position{}, &Velocity{}, &Heading{})

	filter2 := NewFilter2[Position, Velocity](world).Register()
	filter3 := NewFilter3[Position, Velocity, Heading](world).Register()

	expectEqual(t, 0, int(filter2.filter.cache))
	expectEqual(t, 1, int(filter3.filter.cache))

	expectEqual(t, 2, len(world.storage.getCacheTables(&filter2.filter, filter2.relations)))
	expectEqual(t, 1, len(world.storage.getCacheTables(&filter3.filter, filter3.relations)))

	expectPanicsWithValue(t, "filter is already registered, can't register", func() { filter2.Register() })

	filter2.Unregister()
	filter3.Unregister()

	expectPanicsWithValue(t, "filter is not registered, can't unregister", func() { filter2.Unregister() })
	expectPanicsWithValue(t, "no filter for id found to unregister", func() { world.storage.unregisterFilter(&filter{cache: 100}) })
}

func TestFilterCacheRelation(t *testing.T) {
	world := NewWorld()

	posMap := NewMap1[Position](world)
	childMap := NewMap1[ChildOf](world)
	child2Map := NewMap1[ChildOf2](world)
	posChildMap := NewMap2[Position, ChildOf](world)

	target1 := world.NewEntity()
	target2 := world.NewEntity()
	target3 := world.NewEntity()
	target4 := world.NewEntity()

	f1 := NewFilter1[ChildOf](world).Register()
	f2 := NewFilter1[ChildOf](world).Relations(RelIdx(0, target1)).Register()
	f3 := NewFilter1[ChildOf](world).Relations(RelIdx(0, target2)).Register()

	c1 := world.storage.getRegisteredFilter(f1.filter.cache)
	c2 := world.storage.getRegisteredFilter(f2.filter.cache)
	c3 := world.storage.getRegisteredFilter(f3.filter.cache)

	posMap.NewBatch(10, &Position{})

	expectEqual(t, 0, len(c1.tables.tables))
	expectEqual(t, 0, len(c2.tables.tables))
	expectEqual(t, 0, len(c3.tables.tables))

	e1 := childMap.NewEntity(&ChildOf{}, RelIdx(0, target1))
	expectEqual(t, 1, len(c1.tables.tables))
	expectEqual(t, 1, len(c2.tables.tables))

	childMap.NewEntity(&ChildOf{}, RelIdx(0, target3))
	expectEqual(t, 2, len(c1.tables.tables))
	expectEqual(t, 1, len(c2.tables.tables))

	child2Map.NewEntity(&ChildOf2{}, RelIdx(0, target2))

	world.RemoveEntity(e1)
	world.RemoveEntity(target1)
	expectEqual(t, 1, len(c1.tables.tables))
	expectEqual(t, 0, len(c2.tables.tables))

	childMap.NewEntity(&ChildOf{}, RelIdx(0, target2))
	posChildMap.NewEntity(&Position{}, &ChildOf{}, RelIdx(1, target2))
	posChildMap.NewEntity(&Position{}, &ChildOf{}, RelIdx(1, target3))
	posChildMap.NewEntity(&Position{}, &ChildOf{}, RelIdx(1, target4))

	expectEqual(t, 5, len(c1.tables.tables))
	expectEqual(t, 2, len(c3.tables.tables))

	world.RemoveEntities(NewFilter0(world).Batch(), nil)
	expectEqual(t, 0, len(c1.tables.tables))
	expectEqual(t, 0, len(c2.tables.tables))
}
