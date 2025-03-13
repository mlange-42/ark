# Changelog

## [[unpublished]](https://github.com/mlange-42/ark/compare/v0.3.0...main)

### Features

- Adds `QueryX.Count` (#175)
- World stats contain separate reserved/used memory fields (#177)
- Provides `MapX` for up to 12 components (#182)
- Adds `MapX.Set` and `Map.Set` (#183)
- Adds function `ResourceTypeID` (#184)
- `MapX.Get` etc. return `nil` for missing components instead of panic (#189)

### Documentation

- Adds user guide section on Ark's error handling philosophy (#170)
- Adds information on uninformative `Map.Get` errors to API docs (#171)
- Adds benchmarks for build tag `tiny` (#180)

### Performance

- Optimizes table creation and world stats by re-use of item sizes from archetypes (#177)
- Optimizes table column lookup, speeding up component operations and unsafe queries (#178)
- Batch operations use the filter cache for cached filters (#191)

### Bugfixes

- Fails with a more informative error when creating entities with missing relation targets (#169)
- Fixes false-positive debug checks on registered queries and unsafe get (#179)

### Other

- Reduces the maximum number of world locks to 64 to fail earlier on unclosed queries (#185)

## [[v0.3.0]](https://github.com/mlange-42/ark/compare/v0.2.0...v0.3.0)

### Breaking changes

- Relation component marker renamed from `ecs.Relation` to `ecs.RelationMarker` (#120)
- Relation target syntax changed, `ecs.Rel` is now `ecs.RelIdx` (#120)
- World argument in unsafe/ID-based filter moved from method `Query` to filter constructor (#137)
- Removes `Mask` and `MaskTotalBits` (makes it private/internal) (#158, #160)

### Features

- Adds `MapX.AddBatch` and `MapX.AddBatchFn` batch operations (#100)
- Adds `MapX.RemoveBatch` batch operation (#100)
- Adds `World.RemoveEntities` batch operation (#100)
- Adds `ExchangeX.AddBatch` and `ExchangeX.AddBatchFn` batch operations (#101)
- Adds `ExchangeX.RemoveBatch` batch operation (#100)
- Adds `ExchangeX.ExchangeBatch` and `ExchangeX.ExchangeBatchFn` batch operations (#101)
- Adds `Map.SetRelationBatch` and `MapX.SetRelationsBatch` (#102)
- Adds `World.NewEntities` (#136)
- Adds all batch operations for `Map` (164)
- Initial world capacity is optional, 2nd value for relation archetype initial capacity (#109)
- Filters have permanent and ad-hoc relation targets (#113)
- Filters can be registered to the cache to speed up query iteration (#114, #116, #122)
- Relation targets can be specified by either component type or component index (#120)
- `ExchangeX.Removes` can be called multiple times in chains (#124)
- `FilterX.Relations` can be called multiple times in chains (#129)
- Adds methods `Entity.ID` and `Entity.Gen` for debugging purposes (#134)
- Adds `MapX.NewEntityFn`, `MapX.AddFn`, `Map.NewEntityFn` and `Map.AddFn` (#145)
- Adds `ExchangeX.AddFn` and `ExchangeX.ExchangeFn` (#145)
- Adds `World.Stats` for extracting world statistics (#147)
- Adds `Unsafe.IDs` to get all component IDs of an entity (#149)
- Adds methods `ID.Index` and `ResID.Index` (#159)

### Documentation

- Adds a section on tools to the README (#94, #95)
- Adds a benchmark for unsafe queries and world access (#98)
- Adds a comprehensive user guide (#106, #107, #121, #123, #125, #126, #127, #133, #135)
- Adds tables with benchmarks to the user guide (#141)

### Performance

- Implements archetype graph for faster lookup on transitions (#104)
- Optimizes query creation with 50% speedup (#144)
- Optimizes component operations with average 20% speedup (#146)

### Bugfixes

- Locks the world during batch entity creation callback (#97)
- Checks the validity of relations when creating tables (#131)
- Checks the validity of relations in filters and component mappers (#132)

## [[v0.2.0]](https://github.com/mlange-42/ark/compare/v0.1.0...v0.2.0)

### Features

- Adds `MapX.NewBatch` and `MapX.NewBatchFn` for fast batch-creation of entities (#65)
- Adds `ExchangeX.Add` and `ExchangeX.Remove` (#70)
- Adds `Map.GetRelationUnchecked` (#75)
- Adds build tag `tiny` for 64 bit masks (#77)
- Adds `QueryX.GetRelation`, `MapX.GetRelation` and `MapX.GetRelationUnchecked` (#79)
- Adds unsafe, ID-based API (#80, #81, #82)
- Adds several functions and methods required for [ark-serde](https://github.com/mlange-42/ark-serde) (#83, #85)
- Adds resource shortcut functions `GetResource` and `AddResource` (#86)
- Adds `World.Reset` (#89)

## [[v0.1.0]](https://github.com/mlange-42/ark/tree/v0.1.0)

Initial release of Ark.

Basic ECS implementation, as well as entity relationships.
