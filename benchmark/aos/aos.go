package main

import (
	"testing"
)

func aos16Byte(b *testing.B, n int) {
	entities := make([]Aos16Byte, n)
	for i := range n {
		entities[i] = Aos16Byte{
			Pos: Position{0, 0},
			Vel: Velocity{1, 1},
		}
	}

	loop := func() {
		for i := range entities {
			e := &entities[i]
			e.Pos.X += e.Vel.X
			e.Pos.Y += e.Vel.Y
		}
	}

	for b.Loop() {
		loop()
	}
}

func aos32Byte(b *testing.B, n int) {
	entities := make([]Aos32Byte, n)
	for i := range n {
		entities[i] = Aos32Byte{
			Pos: Position{0, 0},
			Vel: Velocity{1, 1},
		}
	}

	loop := func() {
		for i := range entities {
			e := &entities[i]
			e.Pos.X += e.Vel.X
			e.Pos.Y += e.Vel.Y
		}
	}

	for b.Loop() {
		loop()
	}
}

func aos64Byte(b *testing.B, n int) {
	entities := make([]Aos64Byte, n)
	for i := range n {
		entities[i] = Aos64Byte{
			Pos: Position{0, 0},
			Vel: Velocity{1, 1},
		}
	}

	loop := func() {
		for i := range entities {
			e := &entities[i]
			e.Pos.X += e.Vel.X
			e.Pos.Y += e.Vel.Y
		}
	}

	for b.Loop() {
		loop()
	}
}
