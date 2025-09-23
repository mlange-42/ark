// Demonstrates how to use Ark with the [Ebiten] game engine.
// Uses the [Stars] example from the Ebiten docs.
//
// On Linux and MacOS, you nee a C compiler to run this example.
// See the [Ebiten] docs for details.
//
// [Ebiten]: https://ebitengine.org/
// [Stars]: https://ebitengine.org/en/examples/stars.html
package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/ark/ecs"
)

func main() {
	// Create a game instance.
	// See the initializer for how the game logic is plugged together
	// in a modular way.
	game := NewGame()

	// Create a window.
	s := ecs.GetResource[Settings](&game.World)
	ebiten.SetWindowSize(int(s.ScreenWidth), int(s.ScreenHeight))
	ebiten.SetWindowTitle("Stars!")

	// Run the game.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
