package cheatsheet

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/mlange-42/ark/ecs"
)

type Position struct {
	X float64
	Y float64
}

type Velocity struct {
	X float64
	Y float64
}

type Altitude struct {
	Z float64
}

type Health struct{}

type Grid struct{}

func NewGrid(sx, sy int) Grid {
	return Grid{}
}

var world = ecs.NewWorld()
var mapper = ecs.NewMap2[Position, Velocity](&world)
var filter = ecs.NewFilter2[Position, Velocity](&world).Exclusive()
var entity = world.NewEntity()

var registry = ecs.EventRegistry{}
var OnCollisionDetected = registry.NewEventType()

func TestCreateWorld(t *testing.T) {
	world := ecs.NewWorld()
	_ = &world
}

func TestCreateWorldConfig(t *testing.T) {
	world := ecs.NewWorld(1024)
	_ = &world
}

func TestCreateEmpty(t *testing.T) {
	e := world.NewEntity()
	_ = e
}

func TestCreateMapper(t *testing.T) {
	mapper := ecs.NewMap2[Position, Velocity](&world)
	_ = mapper
}

func TestCreateEntity(t *testing.T) {
	e := mapper.NewEntity(
		&Position{X: 100, Y: 100},
		&Velocity{X: 1, Y: -1},
	)
	_ = e
}

func TestCreateEntityFn(t *testing.T) {
	e := mapper.NewEntityFn(func(pos *Position, vel *Velocity) {
		pos.X, pos.Y = 100, 100
		vel.X, vel.Y = 1, -1
	})
	_ = e
}

func TestCreateBatch(t *testing.T) {
	mapper.NewBatch(
		100,
		&Position{X: 100, Y: 100},
		&Velocity{X: 1, Y: -1},
	)
}

func TestCreateBatchFn(t *testing.T) {
	mapper.NewBatchFn(
		100,
		func(e ecs.Entity, pos *Position, vel *Velocity) {
			pos.X, pos.Y = rand.Float64()*100, rand.Float64()*100
			vel.X, vel.Y = rand.NormFloat64(), rand.NormFloat64()
		})
}

func TestRemoveEntity(t *testing.T) {
	world.RemoveEntity(entity)
}

func TestRemoveEntities(t *testing.T) {
	filter := ecs.NewFilter2[Position, Velocity](&world).Exclusive()

	world.RemoveEntities(filter.Batch(), nil)
}

func TestRemoveEntitiesFn(t *testing.T) {
	world.RemoveEntities(filter.Batch(), func(entity ecs.Entity) {
		// Do something before removal
	})
}

func TestAddRemoveComponents(t *testing.T) {
	entity := world.NewEntity()

	mapper.Add(
		entity,
		&Position{X: 100, Y: 100},
		&Velocity{X: 1, Y: -1},
	)

	mapper.Remove(entity)
}

func TestAddBatch(t *testing.T) {
	filter := ecs.NewFilter1[Altitude](&world).Exclusive()

	mapper.AddBatch(
		filter.Batch(),
		&Position{X: 100, Y: 100},
		&Velocity{X: 1, Y: -1},
	)
}

func TestAddBatchFn(t *testing.T) {
	filter := ecs.NewFilter1[Altitude](&world).Exclusive()

	mapper.AddBatchFn(
		filter.Batch(),
		func(e ecs.Entity, pos *Position, vel *Velocity) {
			pos.X, pos.Y = rand.Float64()*100, rand.Float64()*100
			vel.X, vel.Y = rand.NormFloat64(), rand.NormFloat64()
		})
}

func TestRemoveBatch(t *testing.T) {
	filter := ecs.NewFilter2[Position, Velocity](&world)

	mapper.RemoveBatch(
		filter.Batch(),
		func(entity ecs.Entity) { /* ... */ })
}

func TestFilterQuery(t *testing.T) {
	filter := ecs.NewFilter2[Position, Velocity](&world)

	query := filter.Query()
	for query.Next() {
		pos, vel := query.Get()
		pos.X += vel.X
		pos.Y += vel.Y
	}
}

func TestFilterWith(t *testing.T) {
	filter := ecs.NewFilter2[Position, Velocity](&world).
		With(ecs.C[Altitude]())
	_ = filter
}

func TestFilterWithout(t *testing.T) {
	filter := ecs.NewFilter2[Position, Velocity](&world).
		Without(ecs.C[Altitude]())
	_ = filter
}

func TestFilterExclusive(t *testing.T) {
	filter := ecs.NewFilter2[Position, Velocity](&world).
		Exclusive()
	_ = filter
}

func TestFilterWithWithout(t *testing.T) {
	filter := ecs.NewFilter1[Position](&world).
		With(ecs.C[Velocity]()).
		With(ecs.C[Altitude]()).
		Without(ecs.C[Health]())
	_ = filter
}

func TestQueryCount(t *testing.T) {
	query := filter.Query()
	fmt.Println(query.Count())

	query.Close()
}

func TestResourcesQuick(t *testing.T) {
	grid := NewGrid(100, 100)

	ecs.AddResource(&world, &grid)
	_ = ecs.GetResource[Grid](&world)
}

func TestResources(t *testing.T) {
	gridResource := ecs.NewResource[Grid](&world)

	grid := gridResource.Get()
	_ = grid
}

func TestObserver(t *testing.T) {
	ecs.Observe(ecs.OnCreateEntity).
		Do(func(e ecs.Entity) {
			// Do something with the newly created entity.
		}).
		Register(&world)
}

func TestObserverFilter(t *testing.T) {
	ecs.Observe(ecs.OnAddComponents).
		For(ecs.C[Position]()).
		For(ecs.C[Velocity]()).
		Do(func(e ecs.Entity) {
			// Do something with the entity.
		}).
		Register(&world)
}

func TestObserverFilterGeneric(t *testing.T) {
	ecs.Observe2[Position, Velocity](ecs.OnAddComponents).
		Do(func(e ecs.Entity, pos *Position, vel *Velocity) {
			// Do something with the entity and the components
		}).
		Register(&world)
}

func TestObserverWithWithout(t *testing.T) {
	ecs.Observe1[Position](ecs.OnAddComponents).
		With(ecs.C[Velocity]()).
		Without(ecs.C[Altitude]()).
		Do(func(e ecs.Entity, pos *Position) {
			// Do something with the entity and the component
		}).
		Register(&world)
}

func TestCustomEventType(t *testing.T) {
	var registry = ecs.EventRegistry{}

	var OnCollisionDetected = registry.NewEventType()
	var OnInputReceived = registry.NewEventType()

	_, _ = OnCollisionDetected, OnInputReceived
}

func TestCustomEvent(t *testing.T) {
	event := world.Event(OnCollisionDetected).
		For(ecs.C[Position]())

	entity := mapper.NewEntity(&Position{}, &Velocity{})
	event.Emit(entity)
}

func TestCustomEventObserver(t *testing.T) {
	ecs.Observe1[Position](OnCollisionDetected).
		Do(func(e ecs.Entity, p *Position) {
			// Do something with the collision entity and the component
		}).
		Register(&world)
}
