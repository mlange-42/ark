# Ark

[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/ark/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/ark/actions/workflows/tests.yml)
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

## Feature road map

At the moment, Ark supports all basic ECS functionality.
However, please be aware that the API is still unstable.

- [x] Create and remove entities
- [x] Create entities with components
- [x] Add and remove components
- [x] Exchange components (add/remove in one operation)
- [x] Queries with basic component filters
- [x] Advanced filters like `With` and `Without`
- [x] World component access for specific entities
- [x] ECS resources
- [ ] Entity relationships
- [ ] Batch operations
- [ ] Unsafe API for runtime types
- [ ] (De)-serialization (`ark-serde`)
- [ ] Event system

## License

This project is distributed under the [MIT license](./LICENSE).
