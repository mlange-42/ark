<div align="center" width="100%">

[![Ark (logo)](https://github.com/user-attachments/assets/4bbe57c6-2e16-43be-ad5e-0cf26c220f21)](https://github.com/mlange-42/ark)
[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/ark/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/ark/actions/workflows/tests.yml)
[![Coverage Status](https://img.shields.io/coverallsCoverage/github/mlange-42/ark?logo=coveralls)](https://badge.coveralls.io/github/mlange-42/ark?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/mlange-42/ark)](https://goreportcard.com/report/github.com/mlange-42/ark)
[![User Guide](https://img.shields.io/badge/user_guide-%23007D9C?logo=go&logoColor=white&labelColor=gray)](https://mlange-42.github.io/ark/)
[![Go Reference](https://img.shields.io/badge/reference-%23007D9C?logo=go&logoColor=white&labelColor=gray)](https://pkg.go.dev/github.com/mlange-42/ark)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/ark)
[![MIT license](https://img.shields.io/badge/MIT-brightgreen?label=license)](https://github.com/mlange-42/ark/blob/main/LICENSE)

Ark is an archetype-based [Entity Component System](https://en.wikipedia.org/wiki/Entity_component_system) (ECS) for [Go](https://go.dev/).
It is the successor of [Arche](https://github.com/mlange-42/arche).

&mdash;&mdash;

[Features](#features) &nbsp; &bull; &nbsp; [Installation](#installation) &nbsp; &bull; &nbsp; [Usage](#usage) &nbsp; &bull; &nbsp; [Tools](#tools)
</div>

## Features

- Designed for performance and highly optimized.
- Well-documented, type-safe [API](https://pkg.go.dev/github.com/mlange-42/ark).
- Entity relationships as first-class feature.
- Fast batch operations for mass manipulation.
- No systems. Just queries. Use your own structure (or the [Tools](#tools)).
- World serialization and deserialization with [ark-serde](https://github.com/mlange-42/ark-serde).

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
	"math/rand/v2"

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
	for range 1000 {
		// Create a new Entity with components.
		_ = mapper.NewEntity(
			&Position{X: rand.Float64() * 100, Y: rand.Float64() * 100},
			&Velocity{X: rand.NormFloat64(), Y: rand.NormFloat64()},
		)
	}

	// Create a generic filter.
	filter := ecs.NewFilter2[Position, Velocity](&world)

	// Time loop.
	for range 5000 {
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

## Tools

- [ark-serde](https://github.com/mlange-42/ark-serde) provides JSON serialization and deserialization for Ark's World.
- [ark-tools](https://github.com/mlange-42/ark-tools) provides systems, a scheduler, and other useful stuff for Ark.

## License

This project is distributed under the [MIT license](./LICENSE).
