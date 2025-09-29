# Changelog

## [[unpublished]](https://github.com/mlange-42/ark/compare/v0.5.2...main)

### Breaking changed

- Converts interface `Relation` to a concrete struct type (#306)
- Removes `RelationID`, use `Relation` instead (#306)

### Features

- Adds an event system for ECS operations entity creation/removal, component addition/removal and relation target changes (#330, #331, #333, #342, #344, #352)
- Adds custom event support for the event system (#340, #348)
- Adds method `World.Shrink` for freeing memory that exceeds current requirements (#323)

### Performance

- Speeds up relation handling by making `Relation` interface a concrete type (#306)
- Skips `Filter.Without` and `FilterX.Without` when called with zero arguments (#309)
- Resets archetype tables in bulk, avoiding table ID lookup (#313)
- Speeds up freeing tables when removing relation targets, by adding an index lookup (#314)
- Optimizes dumping and loading entities for serialization by pre-allocating memory (#315)
- Speeds up archetype switching by 20% by inlining the mask (#317)
- Speeds up query creation by optimized world locking (#318)
- Optimizes bit mask methods by using less math and more bit-wise operations (#319)
- Speeds up memory allocation for archetype tables and bulk copying by avoiding reflection (#321)

### Documentation

- More explicit API reference docs for `Filter.Exclusive` and `FilterX.Exclusive` (#322)
- Adds a stand-alone example for using Ark with the [Ebiten](https://ebitengine.org/) game engine (#329)
- Adds a user guide chapter on the new event system (#334, #345, #346, #347, #353)
- Adds stand-alone examples for built-in and custom events (#354)

### Bugfixes

- Zeroes memory of non-trivial component types using reflection to inform GC about invalidated pointers (#324)

### Other

- Adds benchmarks for `World.Stats` (#310)
- World stats are now JSON (de)serializable (#311)
- Improves error message on attempt to modify a locked world (#312)
- Changes callback signatures for better IDE autocomplete (#343)

## [[v0.5.2]](https://github.com/mlange-42/ark/compare/v0.5.1...v0.5.2)

### Performance

- Simplifies calculation of item sizes by considering Go's memory alignment rules (#301)

### Documentation

- Excludes code generation from the API reference by moving it to an internal package (#300)
- Clarifies the use of entity relations with filters and maps in the API reference (#305)

### Bugfixes

- Checks that queried/used relation components are actually part of filters and maps (#305, fixes #304)

## [[v0.5.1]](https://github.com/mlange-42/ark/compare/v0.5.0...v0.5.1)

### Bugfixes

- Fixes crash on table creation when removing a relation component (#297, fixes #295)

## [[v0.5.0]](https://github.com/mlange-42/ark/compare/v0.4.8...v0.5.0)

### Breaking changes

- Renames build tags to `ark_tiny` and `ark_debug` to avoid conflicts (#284)

### Features

- Adds methods `QueryX.EntityAt` and `Query.EntityAt` for drawing random entities (#280)

### Performance

- Uses component index to speed up `QueryX.Count` and `QueryX.EntityAt` (#278)
- Includes `FilterX.With` args in component index optimization (#282)

### Documentation

- Adds more examples for queries to API reference (#283)
- Adds a user guide chapter on design philosophy and limitations (#285)
- Adds stand-alone examples for world lock, kd-tree, entity relations and parallel simulations (#288)
- Removes line numbers in user guide code blocks to fix copy button functionality (#291)

### Bugfixes

- Prevents relations specified when querying from getting registered in cache later (#292)

### Other

- Adds debug and non-debug versions for unsafe queries (#279)

## [[v0.4.8]](https://github.com/mlange-42/ark/compare/v0.4.7...v0.4.8)

### Bugfixes

- Fixes component index for pointer persistence (#277)

## [[v0.4.7]](https://github.com/mlange-42/ark/compare/v0.4.6...v0.4.7)

### Performance

- Adds a component index (mapping from components to archetypes) to speed up queries in case of many archetypes (#269, #271, #272)

### Documentation

- Adds a user guide chapter on Ark's archetype-based architecture (#262, #270)

### Other

- Gets rid of the `stretchr/testify` dependency (#275)

## [[v0.4.6]](https://github.com/mlange-42/ark/compare/v0.4.5...v0.4.6)

### Bugfixes

- Fixes duplicate entities in queries by removing tables from internal lookup for relations (#265, co-authored by [ganyyy](https://github.com/ganyyy), fixes #263)

## [[v0.4.5]](https://github.com/mlange-42/ark/compare/v0.4.4...v0.4.5)

### Bugfixes

- Fixes panic on removing a relation target when there are multiple relations on an entity (#260, fixes #259)

## [[v0.4.4]](https://github.com/mlange-42/ark/compare/v0.4.3...v0.4.4)

### Bugfixes

- Fixes bug where a table of relation entities was accidentally iterated twice (#258, fixes #256)

## [[v0.4.3]](https://github.com/mlange-42/ark/compare/v0.4.2...v0.4.3)

### Bugfixes

- Fixes bug in moving entities between archetype tables, caused by non-persistent pointers (#244, fixes #243)
- Fixes premature garbage collection of slices and pointers by copying using reflect (#245)

### Performance

- For trivial component types, get rid of the performance degradation caused by #245 (#249)

## [[v0.4.2]](https://github.com/mlange-42/ark/compare/v0.4.1...v0.4.2)

### Performance

- Reduces the number of initial archetypes to 16 for faster world creation (#234)
- Uses `reflect.TypeFor`, speeding up component ID lookup by 20% (#239 by [LucDrenth](https://github.com/LucDrenth))

### Documentation

- Fixes typos in API docs (#235 by [LucDrenth](https://github.com/LucDrenth))

## [[v0.4.1]](https://github.com/mlange-42/ark/compare/v0.4.0...v0.4.1)

### Documentation

- Adds an example on how to implement systems (#221, #223)
- Adds an example on how to use non-ECS data structures with entities (#222)
- Adds chapter on resources to the user guide (#226)
- Adds a section on limitations of entity relationships to the user guide (#227)
- Tweaks and fixes entity relationships documentation (#230)

## [[v0.4.0]](https://github.com/mlange-42/ark/compare/v0.3.0...v0.4.0)

### Breaking changes

- Removes redundant information from `stats`: `stats.World.ComponentCount` and `stats.Archetype.Components` (#192)
- Unsafe `Filter` and `Query` renamed to `UnsafeFilter` and `UnsafeQuery` (#206)
- Constructors for `MapX` and `Map` return pointers instead of structs, for consistency (#208)

### Features

- Adds `QueryX.Count` (#175)
- World stats contain separate reserved/used memory fields (#177)
- Provides `MapX` for up to 12 components (#182)
- Adds `MapX.Set` and `Map.Set` (#183)
- Adds function `ResourceTypeID` (#184)
- `MapX.Get` etc. return `nil` for missing components instead of panic (#189)
- Adds `FilterX.New`, `MapX.New` and `ExchangeX.New` to avoid repetition of generics (#207 by [Tener](https://github.com/Tener))

### Documentation

- Adds user guide chapter on Ark's error handling philosophy (#170)
- Adds information on uninformative `Map.Get` errors to API docs (#171)
- Adds benchmarks for build tag `tiny` (#180)
- Adds user guide chapter on world statistics (#193)

### Performance

- Optimizes table creation and world stats by re-use of item sizes from archetypes (#177)
- Optimizes table column lookup, speeding up component operations and unsafe queries (#178)
- Batch operations use the filter cache for cached filters (#191)
- Optimizes query creation with 30% speedup (#197)
- More reliable inlining of query methods (#198)
- Optimizes memory sizes of internal types (#199, #200)
- Optimizes handling of relations in filters (#203)
- Rework of component assignment for a bugfix, accidentally speeding up assigning operations (#205)
- Uses a dedicated column type for entities in tables/archetypes, with a small speedup for component operations (216)

### Bugfixes

- Fails with a more informative error when creating entities with missing relation targets (#169)
- Fixes false-positive debug checks on registered queries and unsafe get (#179)
- Fixes a bookkeeping bug when recycling relationship tables (#204)
- Fixes garbage collected pointers when assigning components using non-`...Fn` methods (#205)
- Fixes error on removing relationship components (#212)
- Checks that entities actually have relationship components to set target for (#212)

### Other

- Reduces the maximum number of world locks to 64 to fail earlier on unclosed queries (#185)
- Ark is dual-licensed under MIT and Apache 2.0 (#202)

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
