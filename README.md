[![Ark (logo)](https://github.com/user-attachments/assets/4bbe57c6-2e16-43be-ad5e-0cf26c220f21)](https://github.com/mlange-42/ark)
[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/ark/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/ark/actions/workflows/tests.yml)
[![Coverage Status](https://img.shields.io/coverallsCoverage/github/mlange-42/ark?logo=coveralls)](https://badge.coveralls.io/github/mlange-42/ark?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/mlange-42/ark)](https://goreportcard.com/report/github.com/mlange-42/ark)
[![Go Reference](https://img.shields.io/badge/reference-%23007D9C?logo=go&logoColor=white&labelColor=gray)](https://pkg.go.dev/github.com/mlange-42/ark)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/ark)
[![MIT license](https://img.shields.io/badge/MIT-brightgreen?label=license)](https://github.com/mlange-42/ark/blob/main/LICENSE)

Ark is a work in progress Entity Component System for Go.

Ark implements the lessons learned from my other Go ECS: [Arche](https://github.com/mlange-42/arche).
If you are familiar with Arche, you will feel at home.
The primary aims are:

- More feature-complete entity relationships
- Making it even faster than Arche (already achieved)
- A more structured design due to better planning of features
- More focus on the generic API

## Installation

To use Ark in a Go project, run:

```shell
go get github.com/mlange-42/ark
```

## Usage

Below is the classical Position/Velocity example that every ECS shows in the docs.

For documentation besides the [API docs](https://pkg.go.dev/github.com/mlange-42/ark),
see Arche and its [user guide](https://mlange-42.github.io/arche/) for now.
Ark closely resembles Arche's generic API, and most information about Arche also applies to Ark.

```go
package main

import (
	"math/rand"

	"github.com/mlange-42/ark/ecs"
)

// Position component
type Position struct {
	X float64
	Y float64
}

// Velocity component
type Velocity struct {
	X float64
	Y float64
}

func main() {
	// Create a World with given initial capacity.
	world := ecs.NewWorld(1024)

	// Create a component mapper.
	mapper := ecs.NewMap2[Position, Velocity](&world)

	// Create entities.
	for i := 0; i < 1000; i++ {
		// Create a new Entity with components.
		_ = mapper.NewEntity(
			&Position{X: rand.Float64() * 100, Y: rand.Float64() * 100},
			&Velocity{X: rand.NormFloat64(), Y: rand.NormFloat64()},
		)
	}

	// Create a generic filter.
	filter := ecs.NewFilter2[Position, Velocity](&world)

	// Time loop.
	for t := 0; t < 1000; t++ {
		// Get a fresh query.
		query := filter.Query()
		// Iterate it
		for query.Next() {
			// Component access through the Query.
			pos, vel := query.Get()
			// Update component fields.
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}
```

## Feature road map

At the moment, Ark supports all basic ECS functionality.
However, please be aware that the API is still unstable.

- [x] Create and remove entities
- [x] Create entities with components
- [x] Add and remove components
- [x] Exchange components (add/remove in one operation)
- [x] Queries with basic component filters
- [x] Advanced filters like `With`, `Without` and `Exclusive`
- [x] World component access for specific entities
- [x] ECS resources
- [x] Entity relationships
- [ ] Batch operations
- [ ] Unsafe API for runtime types
- [ ] (De)-serialization (`ark-serde`)
- [ ] Event system
- [ ] Comprehensive user guide

## License

This project is distributed under the [MIT license](./LICENSE).
