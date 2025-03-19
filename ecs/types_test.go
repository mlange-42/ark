package ecs

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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

type CompI struct {
	X, Y float64
}

type CompJ struct {
	X, Y float64
}

type CompK struct {
	X, Y float64
}

type CompL struct {
	X, Y float64
}

func TestID(t *testing.T) {
	id := ID{12}
	assert.EqualValues(t, 12, id.Index())
	rid := ResID{13}
	assert.EqualValues(t, 13, rid.Index())
}

func TestTypeSizes(t *testing.T) {
	printTypeSize[Entity]()
	printTypeSize[entityIndex]()
	printTypeSize[entityPool]()
	printTypeSize[bitPool]()
	printTypeSize[lock]()
	printTypeSize[bitMask64]()
	printTypeSize[bitMask256]()
	printTypeSize[World]()
	printTypeSize[storage]()
	printTypeSize[node]()
	printTypeSize[archetype]()
	printTypeSize[table]()
	printTypeSize[column]()
	printTypeSize[entityColumn]()
	printTypeSize[cacheEntry]()
	printTypeSize[cursor]()
	printTypeSizeName[Filter2[Position, Velocity]]("Filter2")
	printTypeSizeName[Query2[Position, Velocity]]("Query2")
	printTypeSizeName[Query4[CompA, CompB, CompC, CompD]]("Query4")
	printTypeSizeName[Query8[CompA, CompB, CompC, CompD, CompE, CompF, CompG, CompH]]("Query8")
	printTypeSize[Resources]()
	printTypeSizeName[Resource[Position]]("Resource")
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
