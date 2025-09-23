package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/ark/ecs"
)

// LogicSystem interface.
type LogicSystem interface {
	Initialize(w *ecs.World)
	Update(w *ecs.World)
}

// RenderSystem interface.
type RenderSystem interface {
	Initialize(w *ecs.World)
	Draw(w *ecs.World, screen *ebiten.Image)
}

// Game is the main game type.
type Game struct {
	World         ecs.World      // The ECS world as data store.
	LogicSystems  []LogicSystem  // Systems for game logic.
	RenderSystems []RenderSystem // Systems for rendering.
}

// NewGame creates a new game.
func NewGame() *Game {
	// Create the game object.
	g := Game{
		World: ecs.NewWorld(),
	}

	// Add "global" resources (like game settings, a grid, ...).
	ecs.AddResource(&g.World, &Settings{
		ScreenWidth:  640,
		ScreenHeight: 480,
		Scale:        64,
		StarsCount:   1024,
	})

	// Add logic systems.
	g.LogicSystems = append(g.LogicSystems,
		&CreateStars{},
		&MoveStars{},
		&BrightnessStars{},
		&ResetStars{},
	)

	// Add render systems.
	g.RenderSystems = append(g.RenderSystems,
		&RenderStars{},
	)

	// Initialize logic systems.
	for _, s := range g.LogicSystems {
		s.Initialize(&g.World)
	}
	// Initialize render systems.
	for _, s := range g.RenderSystems {
		s.Initialize(&g.World)
	}

	return &g
}

// Update the game.
func (g *Game) Update() error {
	// Update all logic systems.
	for _, s := range g.LogicSystems {
		s.Update(&g.World)
	}
	return nil
}

// Draw the game's graphics.
func (g *Game) Draw(screen *ebiten.Image) {
	// Update/draw all render systems.
	for _, s := range g.RenderSystems {
		s.Draw(&g.World, screen)
	}
}

// Layout the game's screen.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	set := ecs.GetResource[Settings](&g.World)
	return int(set.ScreenWidth), int(set.ScreenHeight)
}
