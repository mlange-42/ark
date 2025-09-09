// Demonstrates that ECS can be mixed with non-ECS data structures, using a kdtree.
package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/mlange-42/ark/ecs"
	"gonum.org/v1/gonum/spatial/kdtree"
)

// Position component.
type Position struct {
	X float64
	Y float64
}

func main() {
	// Create a new World.
	world := ecs.NewWorld()

	// Create a mapper to build entities.
	builder := ecs.NewMap1[Position](&world)

	// Create a list of points for the kd-tree.
	points := kdPoints{}

	// Create 1000 entities with random positions and add them to the points list.
	builder.NewBatchFn(1000, func(entity ecs.Entity, a *Position) {
		a.X = rand.Float64() * 100.0
		a.Y = rand.Float64() * 100.0
		points = append(points, kdPoint{Position: *a, Entity: entity})
	})

	// Build a kd-tree and a keeper for the 5 nearest neighbors.
	tree := kdtree.New(points, false)
	keep := kdtree.NewNKeeper(5)

	// Find the 5 nearest neighbors to the point (50,50).
	center := kdPoint{Position: Position{X: 50, Y: 50}}
	tree.NearestSet(keep, center)

	// Print the results.
	fmt.Println("5 nearest neighbors to (50,50):")
	for _, c := range keep.Heap {
		p := c.Comparable.(kdPoint)
		fmt.Printf("Entity ID: %d, distance: %0.2f\n", p.Entity.ID(), math.Sqrt(p.Distance(center)))
	}
}

// ############ Implementations of kd-tree interfaces ############

// kdPoint is a kdtree.Comparable implementations.
type kdPoint struct {
	Position
	Entity ecs.Entity
}

// Compare satisfies the axis comparisons method of the kdtree.Comparable interface.
func (p kdPoint) Compare(c kdtree.Comparable, d kdtree.Dim) float64 {
	q := c.(kdPoint)
	switch d {
	case 0:
		return p.X - q.X
	case 1:
		return p.Y - q.Y
	default:
		panic("illegal dimension")
	}
}

// Dims returns the number of dimensions to be considered.
func (p kdPoint) Dims() int { return 2 }

// Distance returns the squared distance between the receiver and c.
func (p kdPoint) Distance(c kdtree.Comparable) float64 {
	q := c.(kdPoint)
	dx := p.X - q.X
	dy := p.Y - q.Y
	return dx*dx + dy*dy
}

// kdPoints is a collection of the place type that satisfies kdtree.Interface.
type kdPoints []kdPoint

func (p kdPoints) Index(i int) kdtree.Comparable         { return p[i] }
func (p kdPoints) Len() int                              { return len(p) }
func (p kdPoints) Pivot(d kdtree.Dim) int                { return plane{kdPoints: p, Dim: d}.Pivot() }
func (p kdPoints) Slice(start, end int) kdtree.Interface { return p[start:end] }

// plane is required to help kdPoint.
type plane struct {
	kdtree.Dim
	kdPoints
}

func (p plane) Less(i, j int) bool {
	switch p.Dim {
	case 0:
		return p.kdPoints[i].X < p.kdPoints[j].X
	case 1:
		return p.kdPoints[i].Y < p.kdPoints[j].Y
	default:
		panic("illegal dimension")
	}
}
func (p plane) Pivot() int { return kdtree.Partition(p, kdtree.MedianOfMedians(p)) }
func (p plane) Slice(start, end int) kdtree.SortSlicer {
	p.kdPoints = p.kdPoints[start:end]
	return p
}
func (p plane) Swap(i, j int) {
	p.kdPoints[i], p.kdPoints[j] = p.kdPoints[j], p.kdPoints[i]
}
