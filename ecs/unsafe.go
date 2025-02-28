package ecs

// Unsafe provides access to Ark's unsafe ID-based API.
// Get an instance via [World.Unsafe].
type Unsafe struct {
	world *World
}
