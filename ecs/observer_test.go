package ecs

import (
	"testing"
)

func TestObserve(t *testing.T) {
	w := NewWorld()
	obs := Observe(OnAddComponents).Do(func(e Entity) {}).Register(&w)

	expectPanicsWithValue(t, "can't modify a registered observer",
		func() {
			obs.For(C[Position]())
		})
	expectPanicsWithValue(t, "can't modify a registered observer",
		func() {
			obs.With(C[Position]())
		})
	expectPanicsWithValue(t, "can't modify a registered observer",
		func() {
			obs.Without(C[Position]())
		})
	expectPanicsWithValue(t, "can't modify a registered observer",
		func() {
			obs.Exclusive()
		})

	obs = Observe(OnCreateEntity).For(C[Position]())
	expectTrue(t, obs.hasWith)
	obs = Observe(OnRemoveEntity).For(C[Position]())
	expectTrue(t, obs.hasWith)

	expectPanicsWithValue(t, "observer callback must be set via Do before registering",
		func() {
			Observe(OnCreateEntity).Register(&w)
		})
	expectPanicsWithValue(t, "observer already has a callback",
		func() {
			Observe(OnCreateEntity).Do(func(e Entity) {}).Do(func(e Entity) {})
		})

	obs = Observe(OnAddComponents).Do(func(e Entity) {})

	obs = obs.For()
	expectEqual(t, 0, len(obs.comps))
	expectFalse(t, obs.hasComps)

	obs = obs.For(C[Position]())
	expectEqual(t, 1, len(obs.comps))
	expectTrue(t, obs.hasComps)

	obs = obs.With()
	expectEqual(t, 0, len(obs.with))
	expectFalse(t, obs.hasWith)

	obs = obs.With(C[Position]())
	expectEqual(t, 1, len(obs.with))
	expectTrue(t, obs.hasWith)

	obs = obs.Without()
	expectEqual(t, 0, len(obs.without))
	expectFalse(t, obs.hasWithout)

	obs = obs.Without(C[Position]())
	expectEqual(t, 1, len(obs.without))
	expectTrue(t, obs.hasWithout)
}

func TestObserverRegister(t *testing.T) {
	w := NewWorld()

	obs1 := Observe(OnCreateEntity).
		With(C[Position]()).
		With(C[Velocity]()).
		Without(C[Heading]()).
		Do(func(e Entity) {}).
		Register(&w)
	expectTrue(t, w.storage.observers.HasObservers(OnCreateEntity))

	obs2 := Observe(OnCreateEntity).
		With(C[Position]()).
		Do(func(e Entity) {}).
		Register(&w)

	expectPanicsWithValue(t, "observer is already registered",
		func() {
			obs1.Register(&w)
		})

	obs1.Unregister(&w)
	expectTrue(t, w.storage.observers.HasObservers(OnCreateEntity))
	obs2.Unregister(&w)
	expectFalse(t, w.storage.observers.HasObservers(OnCreateEntity))

	expectPanicsWithValue(t, "observer is not registered",
		func() {
			obs1.Unregister(&w)
		})

	obs1 = Observe(OnCreateEntity).Do(func(e Entity) {}).Register(&w)
	_ = Observe(OnCreateEntity).With(C[Position]()).Do(func(e Entity) {}).Register(&w)
	obs3 := Observe(OnCreateEntity).With(C[Velocity]()).Do(func(e Entity) {}).Register(&w)
	obs3.Unregister(&w)
	expectTrue(t, w.storage.observers.anyNoWith[OnCreateEntity])
	obs1.Unregister(&w)
	expectFalse(t, w.storage.observers.anyNoWith[OnCreateEntity])

	anyNoComps := &w.storage.observers.anyNoComps[OnAddComponents]
	obs1 = Observe(OnAddComponents).For(C[Velocity]()).Do(func(e Entity) {}).Register(&w)
	obs2 = Observe(OnAddComponents).For(C[Position]()).Do(func(e Entity) {}).Register(&w)
	expectFalse(t, *anyNoComps)
	obs3 = Observe(OnAddComponents).Do(func(e Entity) {}).Register(&w)
	expectTrue(t, *anyNoComps)
	obs2.Unregister(&w)
	expectTrue(t, *anyNoComps)
	obs3.Unregister(&w)
	expectFalse(t, *anyNoComps)
}

func TestObserverManager(t *testing.T) {
	m := newObserverManager()
	expectPanicsWithValue(t, "can't unregister observer, not found",
		func() {
			m.RemoveObserver(&Observer{id: 13})
		})
}

func TestObserverOnCreateEntity(t *testing.T) {
	w := NewWorld()

	posID := ComponentID[Position](&w)
	velID := ComponentID[Velocity](&w)
	builder1 := NewMap1[Position](&w)
	builder2 := NewMap1[Velocity](&w)

	callCount := 0
	isBatch := false

	Observe(OnCreateEntity).
		With(C[Position]()).
		Do(func(e Entity) {
			expectEqual(t, isBatch, w.IsLocked())
			callCount++
		}).
		Register(&w)

	builder1.NewEntity(&Position{})
	expectEqual(t, 1, callCount)
	builder2.NewEntity(&Velocity{})
	expectEqual(t, 1, callCount)

	w.Unsafe().NewEntity(posID)
	expectEqual(t, 2, callCount)
	w.Unsafe().NewEntity(velID)
	expectEqual(t, 2, callCount)

	isBatch = true
	builder1.NewBatch(10, &Position{})
	expectEqual(t, 12, callCount)
	builder2.NewBatch(10, &Velocity{})
	expectEqual(t, 12, callCount)
}

func TestObserverOnCreateEntities(t *testing.T) {
	w := NewWorld()

	callCount := 0

	Observe(OnCreateEntity).
		Do(func(e Entity) {
			expectTrue(t, w.IsLocked())
			callCount++
		}).
		Register(&w)

	w.NewEntities(10, nil)
	expectEqual(t, 10, callCount)
}

func TestObserverOnRemoveEntity(t *testing.T) {
	w := NewWorld()

	builder1 := NewMap1[Position](&w)
	builder2 := NewMap1[Velocity](&w)
	filter1 := NewFilter1[Position](&w)
	filter2 := NewFilter1[Velocity](&w)

	callCount := 0

	Observe(OnRemoveEntity).
		With(C[Position]()).
		Do(func(e Entity) {
			expectTrue(t, w.IsLocked())
			callCount++
		}).
		Register(&w)

	e := builder1.NewEntity(&Position{})
	w.RemoveEntity(e)
	expectEqual(t, 1, callCount)
	e = builder2.NewEntity(&Velocity{})
	w.RemoveEntity(e)
	expectEqual(t, 1, callCount)

	builder1.NewBatch(10, &Position{})
	w.RemoveEntities(filter1.Batch(), nil)
	expectEqual(t, 11, callCount)
	builder2.NewBatch(10, &Velocity{})
	w.RemoveEntities(filter2.Batch(), nil)
	expectEqual(t, 11, callCount)
}

func TestObserverOnAddRemove(t *testing.T) {
	w := NewWorld()

	builder1 := NewMap1[Position](&w)
	builder2 := NewMap1[Velocity](&w)
	builder3 := NewMap1[Heading](&w)
	filter1 := NewFilter1[Position](&w)
	filter2 := NewFilter1[Velocity](&w)

	callAdd := 0
	callRemove := 0
	isBatch := false

	Observe(OnAddComponents).
		For(C[Position]()).
		Do(func(e Entity) {
			expectEqual(t, isBatch, w.IsLocked())
			callAdd++
		}).
		Register(&w)

	Observe(OnRemoveComponents).
		For(C[Position]()).
		Do(func(e Entity) {
			expectTrue(t, w.IsLocked())
			callRemove++
		}).
		Register(&w)

	e := w.NewEntity()
	builder3.Add(e, &Heading{})
	expectEqual(t, 0, callAdd)
	builder1.Add(e, &Position{})
	expectEqual(t, 1, callAdd)
	builder2.Add(e, &Velocity{})
	expectEqual(t, 1, callAdd)

	builder3.Remove(e)
	expectEqual(t, 0, callRemove)
	builder2.Remove(e)
	expectEqual(t, 0, callRemove)
	builder1.Remove(e)
	expectEqual(t, 1, callRemove)

	isBatch = true

	builder2.NewBatch(10, &Velocity{})
	builder1.AddBatch(filter2.Batch(), &Position{})
	expectEqual(t, 11, callAdd)

	builder1.RemoveBatch(filter1.Batch(), nil)
	expectEqual(t, 11, callRemove)
}

func TestObserverOnSet(t *testing.T) {
	w := NewWorld()

	builder1 := NewMap1[Position](&w)
	builder2 := NewMap1[Velocity](&w)

	callCount := 0

	Observe(OnSetComponents).
		For(C[Position]()).
		Do(func(e Entity) {
			expectFalse(t, w.IsLocked())
			callCount++
		}).
		Register(&w)

	e := builder1.NewEntity(&Position{})
	builder2.Add(e, &Velocity{})
	expectEqual(t, 0, callCount)

	builder1.Set(e, &Position{})
	expectEqual(t, 1, callCount)
	builder2.Set(e, &Velocity{})
	expectEqual(t, 1, callCount)
}

func TestObserverComponentsWith(t *testing.T) {
	w := NewWorld()

	builder1 := NewMap1[Position](&w)
	builder2 := NewMap1[Velocity](&w)

	callAdd := 0
	callRemove := 0
	callSet := 0

	Observe(OnAddComponents).
		For(C[Position]()).
		With(C[Velocity]()).
		Do(func(e Entity) {
			expectFalse(t, w.IsLocked())
			callAdd++
		}).
		Register(&w)

	Observe(OnAddComponents).
		For(C[Velocity]()).Do(func(e Entity) {}).Register(&w)

	Observe(OnRemoveComponents).
		For(C[Position]()).
		With(C[Velocity]()).
		Do(func(e Entity) {
			expectTrue(t, w.IsLocked())
			callRemove++
		}).
		Register(&w)

	Observe(OnRemoveComponents).
		For(C[Velocity]()).Do(func(e Entity) {}).Register(&w)

	Observe(OnSetComponents).
		For(C[Position]()).
		With(C[Velocity]()).
		Do(func(e Entity) {
			expectFalse(t, w.IsLocked())
			callSet++
		}).
		Register(&w)

	Observe(OnSetComponents).
		For(C[Velocity]()).Do(func(e Entity) {}).Register(&w)

	e := builder2.NewEntity(&Velocity{})

	builder1.Add(e, &Position{})
	expectEqual(t, 1, callAdd)

	builder1.Set(e, &Position{})
	expectEqual(t, 1, callSet)

	builder1.Remove(e)
	expectEqual(t, 1, callRemove)

	e = w.NewEntity()

	builder1.Add(e, &Position{})
	expectEqual(t, 1, callAdd)

	builder1.Set(e, &Position{})
	expectEqual(t, 1, callSet)

	builder1.Remove(e)
	expectEqual(t, 1, callRemove)
}

func TestObserverComponentsWithSimple(t *testing.T) {
	w := NewWorld()

	builder1 := NewMap1[Position](&w)
	builder2 := NewMap1[Velocity](&w)

	callAdd := 0
	callRemove := 0
	callSet := 0

	Observe(OnAddComponents).
		For(C[Position]()).
		With(C[Velocity]()).
		Do(func(e Entity) {
			expectFalse(t, w.IsLocked())
			callAdd++
		}).
		Register(&w)

	Observe(OnRemoveComponents).
		For(C[Position]()).
		With(C[Velocity]()).
		Do(func(e Entity) {
			expectTrue(t, w.IsLocked())
			callRemove++
		}).
		Register(&w)

	Observe(OnSetComponents).
		For(C[Position]()).
		With(C[Velocity]()).
		Do(func(e Entity) {
			expectFalse(t, w.IsLocked())
			callSet++
		}).
		Register(&w)

	e := builder2.NewEntity(&Velocity{})

	builder1.Add(e, &Position{})
	expectEqual(t, 1, callAdd)

	builder1.Set(e, &Position{})
	expectEqual(t, 1, callSet)

	builder1.Remove(e)
	expectEqual(t, 1, callRemove)

	e = w.NewEntity()

	builder1.Add(e, &Position{})
	expectEqual(t, 1, callAdd)

	builder1.Set(e, &Position{})
	expectEqual(t, 1, callSet)

	builder1.Remove(e)
	expectEqual(t, 1, callRemove)
}

func TestObserverComponentsWithout(t *testing.T) {
	w := NewWorld()

	builder1 := NewMap1[Position](&w)
	builder2 := NewMap1[Velocity](&w)

	callAdd := 0
	callRemove := 0
	callSet := 0

	Observe(OnAddComponents).
		For(C[Position]()).
		Without(C[Velocity]()).
		Do(func(e Entity) {
			expectFalse(t, w.IsLocked())
			callAdd++
		}).
		Register(&w)

	Observe(OnRemoveComponents).
		For(C[Position]()).
		Without(C[Velocity]()).
		Do(func(e Entity) {
			expectTrue(t, w.IsLocked())
			callRemove++
		}).
		Register(&w)

	Observe(OnSetComponents).
		For(C[Position]()).
		Without(C[Velocity]()).
		Do(func(e Entity) {
			expectFalse(t, w.IsLocked())
			callSet++
		}).
		Register(&w)

	e := w.NewEntity()

	builder1.Add(e, &Position{})
	expectEqual(t, 1, callAdd)

	builder1.Set(e, &Position{})
	expectEqual(t, 1, callSet)

	builder1.Remove(e)
	expectEqual(t, 1, callRemove)

	e = builder2.NewEntity(&Velocity{})

	builder1.Add(e, &Position{})
	expectEqual(t, 1, callAdd)

	builder1.Set(e, &Position{})
	expectEqual(t, 1, callSet)

	builder1.Remove(e)
	expectEqual(t, 1, callRemove)
}

func TestObserverEntitiesWith(t *testing.T) {
	w := NewWorld()

	builder1 := NewMap1[Position](&w)
	builder2 := NewMap1[Velocity](&w)

	callAdd := 0
	callRemove := 0

	Observe(OnCreateEntity).Do(func(e Entity) {}).Register(&w)
	Observe(OnRemoveEntity).Do(func(e Entity) {}).Register(&w)

	Observe(OnCreateEntity).
		With(C[Velocity]()).
		Do(func(e Entity) { callAdd++ }).
		Register(&w)
	Observe(OnRemoveEntity).
		With(C[Velocity]()).
		Do(func(e Entity) { callRemove++ }).
		Register(&w)

	e1 := builder1.NewEntity(&Position{})
	expectEqual(t, 0, callAdd)
	e2 := builder2.NewEntity(&Velocity{})
	expectEqual(t, 1, callAdd)

	w.RemoveEntity(e1)
	expectEqual(t, 0, callRemove)
	w.RemoveEntity(e2)
	expectEqual(t, 1, callRemove)
}

func TestObserverExclusive(t *testing.T) {
	w := NewWorld()

	obs := Observe(OnCreateEntity).
		With(C[Position]()).
		Exclusive().
		Do(func(e Entity) {}).Register(&w)

	expectTrue(t, obs.withMask.Get(0))
	expectFalse(t, obs.withMask.Get(1))

	expectFalse(t, obs.withoutMask.Get(0))
	expectTrue(t, obs.withoutMask.Get(1))

	expectTrue(t, obs.exclusive)
	expectTrue(t, obs.hasWithout)
}

func TestObserverEntitiesWithout(t *testing.T) {
	w := NewWorld()

	builder1 := NewMap1[Position](&w)
	builder2 := NewMap1[Velocity](&w)

	callAdd := 0
	callRemove := 0

	Observe(OnCreateEntity).Do(func(e Entity) {}).Register(&w)
	Observe(OnRemoveEntity).Do(func(e Entity) {}).Register(&w)

	Observe(OnCreateEntity).
		Without(C[Velocity]()).
		Do(func(e Entity) { callAdd++ }).
		Register(&w)
	Observe(OnRemoveEntity).
		Without(C[Velocity]()).
		Do(func(e Entity) { callRemove++ }).
		Register(&w)

	e1 := builder2.NewEntity(&Velocity{})
	expectEqual(t, 0, callAdd)
	e2 := builder1.NewEntity(&Position{})
	expectEqual(t, 1, callAdd)

	w.RemoveEntity(e1)
	expectEqual(t, 0, callRemove)
	w.RemoveEntity(e2)
	expectEqual(t, 1, callRemove)
}

func TestObserverWildcardComponents(t *testing.T) {
	w := NewWorld()

	builder1 := NewMap1[Position](&w)

	callAdd := 0
	callRemove := 0
	callSet := 0

	Observe(OnAddComponents).
		Do(func(e Entity) {
			expectFalse(t, w.IsLocked())
			callAdd++
		}).
		Register(&w)

	Observe(OnRemoveComponents).
		Do(func(e Entity) {
			expectTrue(t, w.IsLocked())
			callRemove++
		}).
		Register(&w)

	Observe(OnSetComponents).
		Do(func(e Entity) {
			expectFalse(t, w.IsLocked())
			callSet++
		}).
		Register(&w)

	e := w.NewEntity()

	builder1.Add(e, &Position{})
	expectEqual(t, 1, callAdd)

	builder1.Set(e, &Position{})
	expectEqual(t, 1, callSet)

	builder1.Remove(e)
	expectEqual(t, 1, callRemove)
}

func TestObserverWildcardEntities(t *testing.T) {
	w := NewWorld()

	builder1 := NewMap1[Position](&w)

	callAdd := 0
	callRemove := 0

	Observe(OnCreateEntity).
		Do(func(e Entity) {
			expectFalse(t, w.IsLocked())
			callAdd++
		}).
		Register(&w)

	Observe(OnRemoveEntity).
		Do(func(e Entity) {
			expectTrue(t, w.IsLocked())
			callRemove++
		}).
		Register(&w)

	e1 := w.NewEntity()
	expectEqual(t, 1, callAdd)
	e2 := builder1.NewEntity(&Position{})
	expectEqual(t, 2, callAdd)

	w.RemoveEntity(e1)
	expectEqual(t, 1, callRemove)
	w.RemoveEntity(e2)
	expectEqual(t, 2, callRemove)
}

func TestObserverCustomEvents(t *testing.T) {
	w := NewWorld()
	builder := NewMap1[Position](&w)

	e1 := w.NewEntity()
	e2 := builder.NewEntity(&Position{})

	NewEvent(customEvent, &w).Emit(e1)

	callCount := 0

	obs := Observe(customEvent).With(comp[Position]()).Do(func(e Entity) { callCount++ }).Register(&w)
	NewEvent(customEvent, &w).Emit(e1)
	obs.Unregister(&w)

	Observe(customEvent).Do(func(e Entity) { callCount++ }).Register(&w)
	Observe(customEvent).For(C[Position]()).Do(func(e Entity) { callCount++ }).Register(&w)
	Observe(customEvent).With(C[Position]()).Do(func(e Entity) { callCount++ }).Register(&w)
	Observe(customEvent).Without(C[Position]()).Do(func(e Entity) { callCount++ }).Register(&w)
	NewEvent(customEvent, &w).Emit(e1)
	NewEvent(customEvent, &w).Emit(e2)

	expectEqual(t, 4, callCount)
}

func benchmarkEventsPos(b *testing.B, n int) {
	w := NewWorld()

	if n > 0 {
		Observe(OnAddComponents).
			For(C[CompA]()).
			Do(func(e Entity) {}).
			Register(&w)
	}

	for i := 1; i < n; i++ {
		Observe(OnAddComponents).
			For(C[CompB]()).
			Do(func(e Entity) {}).
			Register(&w)
	}

	oldMask := newMask()
	newMask := newMask(ComponentID[CompA](&w))

	for b.Loop() {
		w.storage.observers.FireAdd(Entity{}, &oldMask, &newMask)
	}
}

func benchmarkEventsNeg(b *testing.B, n int) {
	w := NewWorld()

	for range n {
		Observe(OnAddComponents).
			For(C[CompA]()).
			Do(func(e Entity) {}).
			Register(&w)
	}

	oldMask := newMask()
	newMask := newMask(ComponentID[CompB](&w))

	for b.Loop() {
		w.storage.observers.FireAdd(Entity{}, &oldMask, &newMask)
	}
}

func BenchmarkEventsPos_0Obs(b *testing.B) {
	benchmarkEventsPos(b, 0)
}

func BenchmarkEventsPos_1Obs(b *testing.B) {
	benchmarkEventsPos(b, 1)
}

func BenchmarkEventsPos_2Obs(b *testing.B) {
	benchmarkEventsPos(b, 2)
}

func BenchmarkEventsPos_5Obs(b *testing.B) {
	benchmarkEventsPos(b, 5)
}

func BenchmarkEventsPos_10Obs(b *testing.B) {
	benchmarkEventsPos(b, 10)
}

func BenchmarkEventsNeg_0Obs(b *testing.B) {
	benchmarkEventsNeg(b, 0)
}

func BenchmarkEventsNeg_1Obs(b *testing.B) {
	benchmarkEventsNeg(b, 1)
}

func BenchmarkEventsNeg_2Obs(b *testing.B) {
	benchmarkEventsNeg(b, 2)
}

func BenchmarkEventsNeg_5Obs(b *testing.B) {
	benchmarkEventsNeg(b, 5)
}

func BenchmarkEventsNeg_10Obs(b *testing.B) {
	benchmarkEventsNeg(b, 10)
}
