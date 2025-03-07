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
	RelationMarker
}

type ChildOf2 struct {
	RelationMarker
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

type CompC struct {
	X, Y float64
}

type CompD struct {
	X, Y float64
}

type CompE struct {
	X, Y float64
}

type CompF struct {
	X, Y float64
}

type CompG struct {
	X, Y float64
}

type CompH struct {
	X, Y float64
}

func TestTypeSizes(t *testing.T) {
	printTypeSize[Entity]()
	printTypeSize[entityIndex]()
	printTypeSize[bitMask]()
	printTypeSize[World]()
	printTypeSize[storage]()
	printTypeSize[node]()
	printTypeSize[archetype]()
	printTypeSize[table]()
	printTypeSize[column]()
	printTypeSize[cacheEntry]()
	printTypeSizeName[Filter2[Position, Velocity]]("Filter2")
	printTypeSizeName[Query2[Position, Velocity]]("Query2")
	printTypeSizeName[Query4[CompA, CompB, CompC, CompD]]("Query4")
	printTypeSizeName[Query8[CompA, CompB, CompC, CompD, CompE, CompF, CompG, CompH]]("Query8")
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
