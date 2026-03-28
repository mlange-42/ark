package main

import (
	"testing"
)

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

	loop()

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

	loop()

	for b.Loop() {
		loop()
	}
}

func aos128Byte(b *testing.B, n int) {
	entities := make([]Aos128Byte, n)
	for i := range n {
		entities[i] = Aos128Byte{
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

	loop()

	for b.Loop() {
		loop()
	}
}

func aos256Byte(b *testing.B, n int) {
	entities := make([]Aos256Byte, n)
	for i := range n {
		entities[i] = Aos256Byte{
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

	loop()

	for b.Loop() {
		loop()
	}
}
