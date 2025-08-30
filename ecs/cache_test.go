package ecs

import (
	"testing"
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

	if filter2.cache != 0 {
		t.Errorf("expected filter2.cache to be 0, got %d", filter2.cache)
	}
	if filter3.cache != 1 {
		t.Errorf("expected filter3.cache to be 1, got %d", filter3.cache)
	}

	if got := len(world.storage.getTableIDs(&filter2.filter, filter2.relations)); got != 2 {
		t.Errorf("expected 2 table IDs for filter2, got %d", got)
	}
	if got := len(world.storage.getTableIDs(&filter3.filter, filter3.relations)); got != 1 {
		t.Errorf("expected 1 table ID for filter3, got %d", got)
	}

	expectPanicWithValue(t, "filter is already registered, can't register", func() {
		filter2.Register()
	})

	filter2.Unregister()
	filter3.Unregister()

	expectPanicWithValue(t, "filter is not registered, can't unregister", func() {
		filter2.Unregister()
	})

	expectPanicWithValue(t, "no filter for id found to unregister", func() {
		world.storage.unregisterFilter(100)
	})
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

	if len(c1.tables) != 0 {
		t.Errorf("expected c1.tables to be empty, got %d", len(c1.tables))
	}
	if len(c2.tables) != 0 {
		t.Errorf("expected c2.tables to be empty, got %d", len(c2.tables))
	}
	if len(c3.tables) != 0 {
		t.Errorf("expected c3.tables to be empty, got %d", len(c3.tables))
	}

	e1 := childMap.NewEntity(&ChildOf{}, RelIdx(0, target1))
	if len(c1.tables) != 1 {
		t.Errorf("expected c1.tables to have 1 entry, got %d", len(c1.tables))
	}
	if len(c2.tables) != 1 {
		t.Errorf("expected c2.tables to have 1 entry, got %d", len(c2.tables))
	}

	childMap.NewEntity(&ChildOf{}, RelIdx(0, target3))
	if len(c1.tables) != 2 {
		t.Errorf("expected c1.tables to have 2 entries, got %d", len(c1.tables))
	}
	if len(c2.tables) != 1 {
		t.Errorf("expected c2.tables to still have 1 entry, got %d", len(c2.tables))
	}

	child2Map.NewEntity(&ChildOf2{}, RelIdx(0, target2))

	world.RemoveEntity(e1)
	world.RemoveEntity(target1)
	if len(c1.tables) != 1 {
		t.Errorf("expected c1.tables to have 1 entry after removal, got %d", len(c1.tables))
	}
	if len(c2.tables) != 0 {
		t.Errorf("expected c2.tables to be empty after removal, got %d", len(c2.tables))
	}

	childMap.NewEntity(&ChildOf{}, RelIdx(0, target2))
	posChildMap.NewEntity(&Position{}, &ChildOf{}, RelIdx(1, target2))
	posChildMap.NewEntity(&Position{}, &ChildOf{}, RelIdx(1, target3))
	posChildMap.NewEntity(&Position{}, &ChildOf{}, RelIdx(1, target4))

	if len(c1.tables) != 5 {
		t.Errorf("expected c1.tables to have 5 entries, got %d", len(c1.tables))
	}
	if len(c3.tables) != 2 {
		t.Errorf("expected c3.tables to have 2 entries, got %d", len(c3.tables))
	}

	world.RemoveEntities(NewFilter0(&world).Batch(), nil)
	if len(c1.tables) != 0 {
		t.Errorf("expected c1.tables to be empty after batch removal, got %d", len(c1.tables))
	}
	if len(c2.tables) != 0 {
		t.Errorf("expected c2.tables to be empty after batch removal, got %d", len(c2.tables))
	}
}
