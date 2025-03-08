// Package ecs provides the core API of Ark, an Entity Component System (ECS) for Go.
//
// See the top-level module [github.com/mlange-42/ark] for an overview.
//
// 🕮 Also read Ark's [User Guide]!
//
// # Outline
//
//   - [World] provides basic functionality that does not require generics,
//     like [World.NewEntity], [World.Alive], [World.RemoveEntity], etc.
//   - [Filter1], [Filter2], etc. provide filters for querying entities by their components.
//   - [Query1], [Query2], etc. are the actual query iterators, and provide functionality like
//     [Query1.Next], [Query1.Get] and [Query1.Entity].
//   - [Map1], [Map2], etc. provide access to and manipulation of multiple components.
//     like [Map1.Get], [Map1.Add], [Map1.Remove], etc.
//   - [Map] provides access to a single component, like [Map.Get] and [Map.Add], [Map.Remove].
//   - [Exchange1], [Exchange2] etc. allows to add, remove and exchange components.
//   - [Resource] provides access the world's [Resources].
//   - See the separate module [ark-serde] for serialization.
//
// # ECS Manipulations
//
// This section gives an overview on how to achieve typical ECS operations in Ark.
//
// Access data:
//   - Create a filter: [NewFilter2].
//   - Create a [Query2]: [Filter2.Query].
//   - Iterate a Query: [Query2.Next], [Query2.Get].
//   - Access components of entities: [Map.Get], [Map2.Get].
//   - Access relationship targets: [Map.GetRelation], [Map2.GetRelation].
//
// Manipulate a single [Entity]:
//   - Create an entity: [World.NewEntity]
//   - Create an entity with components: [Map2.NewEntity], [Map2.NewEntityFn].
//   - Add components to an entity: [Map.Add], [Map.AddFn], [Map2.Add], [Map2.AddFn], [Exchange2.Add], [Exchange2.AddFn].
//   - Remove components from an entity: [Map.Remove], [Map2.Remove], [Exchange2.Remove].
//   - Exchange components of an entity: [Exchange2.Exchange], [Exchange2.ExchangeFn].
//   - Change relationship targets: [Map.SetRelation], [Map2.SetRelations].
//   - Remove an entity from the world: [World.RemoveEntity].
//
// Manipulate entities in batches:
//   - Create entities: [World.NewEntities]
//   - Create entities with components: [Map2.NewBatch], [Map2.NewBatchFn].
//   - Add components to entities: [Map2.AddBatch], [Map2.AddBatchFn], [Exchange2.AddBatch], [Exchange2.AddBatchFn].
//   - Remove components from entities: [Map2.RemoveBatch], [Exchange2.RemoveBatch].
//   - Exchange components of entities: [Exchange2.ExchangeBatch], [Exchange2.ExchangeBatchFn].
//   - Change relationship targets: [Map2.SetRelationsBatch].
//   - Remove entities from the world: [World.RemoveEntities].
//
// # Build tags
//
// Ark provides two build tags:
//   - tiny: Reduces the maximum number of components to 64, for faster mask-related operations and smaller archetype memory footprint.
//   - debug: Improves error messages on incorrect use, at the cost of performance. Use this if you get panics from queries or maps.
//
// When building your application, use them like this:
//
//	go build -tags tiny .
//	go build -tags debug .
//	go build -tags tiny,debug .
//
// [ark-serde]: https://github.com/mlange-42/ark-serde/
// [User Guide]: https://mlange-42.github.io/ark/
package ecs
