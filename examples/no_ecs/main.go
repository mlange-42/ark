// Demonstrates that ECS can be mixed with non-ECS data structures, as long as they store entities.
package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/mlange-42/ark/ecs"
)

// CellCoord component.
type CellCoord struct {
	X int
	Y int
}

func main() {
	// Create a new World.
	world := ecs.NewWorld()

	// Create a non-ECS grid data structure,
	// and add is as a resource.
	grid := NewGrid(15, 20)
	ecs.AddResource(&world, &grid)

	// Create enities on the grid.
	createGridEntities(&world, 100)

	// Run a simulation
	run(&world)
}

func run(world *ecs.World) {
	// Get the grid resource.
	gridRes := ecs.NewResource[Grid](world)
	grid := gridRes.Get()

	// Create a mapper for component access.
	mapper := ecs.NewMap1[CellCoord](world)

	// Print coords of entities in the grid.
	b := strings.Builder{}
	for y := range grid.Height {
		for x := range grid.Width {
			// Get an entity by coordinate.
			entity := grid.Data[y][x]
			// If no entity there...
			if entity.IsZero() {
				b.WriteString("    -    ")
				continue
			}
			// Otherwise, get the CellCoord component.
			coord := mapper.Get(entity)
			// Print it.
			b.WriteString(fmt.Sprintf(" (%2d,%2d) ", coord.X, coord.Y))
		}
		b.WriteRune('\n')
	}
	fmt.Println(b.String())
}

func createGridEntities(world *ecs.World, count int) {
	// Get the grid resource.
	gridRes := ecs.NewResource[Grid](world)
	grid := gridRes.Get()

	// Create a mapper as a builder.
	builder := ecs.NewMap1[CellCoord](world)

	// Put some entities into the grid.
	cnt := 0
	for cnt < count {
		// Draw random coordinates.
		x, y := rand.Intn(grid.Width), rand.Intn(grid.Height)
		// Skip if there is already an entity.
		if !grid.Data[y][x].IsZero() {
			continue
		}
		// Create an entity.
		entity := builder.NewEntity(&CellCoord{X: x, Y: y})
		// Place the entity in the grid.
		grid.Data[y][x] = entity
		cnt++
	}
}

// Grid resource / data structure.
type Grid struct {
	Data   [][]ecs.Entity
	Width  int
	Height int
}

// NewGrid creates a ner Grid of the given size.
func NewGrid(w, h int) Grid {
	grid := make([][]ecs.Entity, h)
	for i := 0; i < h; i++ {
		grid[i] = make([]ecs.Entity, w)
	}
	return Grid{
		Data:   grid,
		Width:  w,
		Height: h,
	}
}
