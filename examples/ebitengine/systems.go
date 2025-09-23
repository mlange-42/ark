package main

import (
	"image/color"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mlange-42/ark/ecs"
)

// CreateStars system initializes the stars.
type CreateStars struct{}

// Initialize the system.
func (s *CreateStars) Initialize(w *ecs.World) {
	set := ecs.GetResource[Settings](w)
	builder := ecs.NewMap3[From, To, Brightness](w)

	builder.NewBatchFn(set.StarsCount, func(entity ecs.Entity, from *From, to *To, br *Brightness) {
		resetStar(from, to, br, set)
	})
}

// Update the system.
func (s *CreateStars) Update(w *ecs.World) {}

// MoveStars system moves stars around.
type MoveStars struct {
	filter *ecs.Filter2[From, To]
}

// Initialize the system.
func (s *MoveStars) Initialize(w *ecs.World) {
	s.filter = s.filter.New(w)
}

// Update the system.
func (s *MoveStars) Update(w *ecs.World) {
	scale := ecs.GetResource[Settings](w).Scale
	mouseX, mouseY := ebiten.CursorPosition()
	x, y := float32(mouseX)*scale, float32(mouseY)*scale

	query := s.filter.Query()
	for query.Next() {
		from, to := query.Get()
		from.X = to.X
		from.Y = to.Y
		to.X += (to.X - x) / 32
		to.Y += (to.Y - y) / 32
	}
}

// ResetStars system resets stars when they reach the edge of the screen.
type ResetStars struct {
	filter *ecs.Filter3[From, To, Brightness]
}

// Initialize the system.
func (s *ResetStars) Initialize(w *ecs.World) {
	s.filter = s.filter.New(w)
}

// Update the system.
func (s *ResetStars) Update(w *ecs.World) {
	set := ecs.GetResource[Settings](w)

	query := s.filter.Query()
	for query.Next() {
		from, to, br := query.Get()

		if from.X < 0 || set.ScreenWidth*set.Scale < from.X ||
			from.Y < 0 || set.ScreenHeight*set.Scale < from.Y {
			resetStar(from, to, br, set)
		}
	}
}

// BrightnessStars system adjusts star brightness.
type BrightnessStars struct {
	filter *ecs.Filter1[Brightness]
}

// Initialize the system.
func (s *BrightnessStars) Initialize(w *ecs.World) {
	s.filter = s.filter.New(w)
}

// Update the system.
func (s *BrightnessStars) Update(w *ecs.World) {
	query := s.filter.Query()
	for query.Next() {
		br := query.Get()
		br.V++
		if 0xff < br.V {
			br.V = 0xff
		}
	}
}

// RenderStars system renders the stars.
type RenderStars struct {
	filter *ecs.Filter3[From, To, Brightness]
}

// Initialize the system.
func (s *RenderStars) Initialize(w *ecs.World) {
	s.filter = s.filter.New(w)
}

// Draw the system.
func (s *RenderStars) Draw(w *ecs.World, screen *ebiten.Image) {
	scale := ecs.GetResource[Settings](w).Scale

	query := s.filter.Query()
	for query.Next() {
		from, to, br := query.Get()
		c := color.RGBA{
			R: uint8(0xbb * br.V / 0xff),
			G: uint8(0xdd * br.V / 0xff),
			B: uint8(0xff * br.V / 0xff),
			A: 0xff}
		vector.StrokeLine(screen, from.X/scale, from.Y/scale, to.X/scale, to.Y/scale, 1, c, true)
	}
}

func resetStar(from *From, to *To, br *Brightness, set *Settings) {
	to.X = rand.Float32() * set.ScreenWidth * set.Scale
	to.Y = rand.Float32() * set.ScreenHeight * set.Scale

	from.X = to.X
	from.Y = to.Y

	br.V = rand.Float32() * 0xff
}
