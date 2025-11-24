package ecs

import (
	"testing"
)

func TestResource(t *testing.T) {
	w := NewWorld(1024)
	var get Resource[Grid]
	get = get.New(w)

	expectFalse(t, get.Has())
	gridResource := NewGrid(100, 200)
	get.Add(&gridResource)

	expectTrue(t, get.Has())
	grid := get.Get()

	expectEqual(t, Grid{100, 200}, *grid)

	get.Remove()
	expectFalse(t, get.Has())
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
	w := NewWorld()

	res := NewRes()
	AddResource(w, &res)

	resOut := *GetResource[ResInterface](w)

	expectNotNil(t, resOut)
	expectEqual(t, "test", resOut.MyMethod())
}
