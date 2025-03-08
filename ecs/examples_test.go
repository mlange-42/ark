package ecs_test

import "github.com/mlange-42/ark/ecs"

type Position struct {
	X, Y float64
}

type Velocity struct {
	X, Y float64
}

type Grid struct {
	Width  int
	Height int
}

func NewGrid(width, height int) Grid {
	return Grid{
		Width:  width,
		Height: height,
	}
}

func (g *Grid) Get(x, y int) ecs.Entity {
	return ecs.Entity{}
}
