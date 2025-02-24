package ecs

import (
	"fmt"
	"reflect"
	"testing"
)

type Position struct {
	X, Y float64
}

type Velocity struct {
	X, Y float64
}

type Heading struct {
	H float64
}

type Label struct{}

type ChildOf struct {
	Relation
}

type ChildOf2 struct {
	Relation
}

type SliceComp struct {
	Slice []int
}

type PointerComp struct {
	Ptr   *PointerType
	Value int
}

type PointerType struct {
	Pos *Position
}

type CompA struct {
	X, Y float64
}

type CompB struct {
	X, Y float64
}

func TestTypeSizes(t *testing.T) {
	printTypeSize[Entity]()
	printTypeSize[entityIndex]()
	printTypeSize[Mask]()
	printTypeSize[World]()
	printTypeSize[column]()
	printTypeSize[table]()
	printTypeSize[archetype]()
	printTypeSizeName[[]column]("slice")
	printTypeSizeName[*Position]("pointer")
	printTypeSizeName[reflect.Value]("reflect.Value")
}

func printTypeSize[T any]() {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	fmt.Printf("%18s: %5d B\n", tp.Name(), tp.Size())
}

func printTypeSizeName[T any](name string) {
	tp := reflect.TypeOf((*T)(nil)).Elem()
	fmt.Printf("%18s: %5d B\n", name, tp.Size())
}
