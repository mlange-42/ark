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
	tp := typeOf[T]()
	fmt.Printf("%18s: %5d B\n", tp.Name(), tp.Size())
}

func printTypeSizeName[T any](name string) {
	tp := typeOf[T]()
	fmt.Printf("%18s: %5d B\n", name, tp.Size())
}
