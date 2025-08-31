package ecs_test

import (
	"testing"

	"github.com/mlange-42/ark/ecs"
)

func TestResource(t *testing.T) {
	w := ecs.NewWorld(1024)
	get := ecs.NewResource[Grid](&w)
	if get.Has() {
		t.Errorf("expected resource to not exist")
	}
	gridResource := NewGrid(100, 200)
	get.Add(&gridResource)
	if !get.Has() {
		t.Errorf("expected resource to exist")
	}
	grid := get.Get()
	if *grid != (Grid{100, 200}) {
		t.Errorf("expected grid to be %v, got %v", Grid{100, 200}, *grid)
	}
	get.Remove()
	if get.Has() {
		t.Errorf("expected resource to not exist")
	}
}

type ResInterface interface {
	MyMethod() string
}

type ResImpl struct{}

func NewRes() ResInterface {
	return &ResImpl{}
}

func (r *ResImpl) MyMethod() string {
	return "test"
}

func TestResourceInterface(t *testing.T) {
	w := ecs.NewWorld()
	res := NewRes()
	ecs.AddResource(&w, &res)
	resOut := *ecs.GetResource[ResInterface](&w)
	if resOut == nil {
		t.Errorf("expected resource to not be nil")
	}
	if resOut.MyMethod() != "test" {
		t.Errorf("expected MyMethod to return 'test', got '%s'", resOut.MyMethod())
	}
}

func ExampleResource() {
	// Create a world.
	world := ecs.NewWorld()

	// Create a resource.
	gridResource := NewGrid(100, 100)
	// Add it to the world.
	ecs.AddResource(&world, &gridResource)

	// Resource access in systems.
	// Create and store a resource accessor.
	gridAccess := ecs.NewResource[Grid](&world)

	// Use the resource.
	grid := gridAccess.Get()
	entity := grid.Get(13, 42)
	_ = entity
	// Output:
}
