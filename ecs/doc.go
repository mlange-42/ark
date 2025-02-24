// Package ecs provides the core API of Ark, an Entity Component System (ECS) for Go.
//
// See the top-level module [github.com/mlange-42/ark] for an overview.
//
// # Outline
//
//   - [World] provides basic functionality that does not require generics,
//     like [World.NewEntity], [World.Alive], [World.RemoveEntity], etc.
//   - [Filter0], [Filter1], etc. provide generic filters and query generation using [Filter0.Query] and friends.
//   - [Query0], [Query1], are the actual query iterators, and provide functionality like
//     [Query1.Next], [Query1.Get] and [Query1.Entity].
//   - [Map] provides generic access to a single component using world access, like [Map.Get] and [Map.Add], [Map.Remove].
//   - [Map1], [Map2], etc. provide generic access to multiple components using world access,
//     like [Map1.Get], [Map1.Add], [Map1.Remove], etc.
//     They an also be used to create entities with components, with [Map1.NewEntity] etc.
//   - [Exchange1], [Exchange2] etc allows to add, remove and exchange components, incl. as batch operations.
//   - [Resource] provides generic access to a resource from [ecs.Resources].
//
// # Build tags
//
// Arche provides the build tag `debug`, which improves error messages on misuse, at the cost of performance.
// Use this if you get panics from queries or maps.
//
// When building your application, use it like this:
//
//	go build -tags debug .
package ecs
