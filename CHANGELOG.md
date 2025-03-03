# Changelog

## [[unpublished]](https://github.com/mlange-42/ark/compare/v0.2.0...main)

### Features

- Adds `MapX.AddBatch` and `MapX.AddBatchFn` batch operations (#100)
- Adds `MapX.RemoveBatch` batch operation (#100)
- Adds `World.RemoveEntities` batch operation (#100)
- Adds `ExchangeX.AddBatch` and `ExchangeX.AddBatchFn` batch operations (#101)
- Adds `ExchangeX.RemoveBatch` batch operation (#100)
- Adds `ExchangeX.ExchangeBatch` and `ExchangeX.ExchangeBatchFn` batch operations (#101)
- Adds `Map.SetRelationBatch` and `MapX.SetRelationsBatch` (#102)
- Initial world capacity is optional, 2nd value for relation archetype initial capacity (#109)

### Documentation

- Adds a section on tool to the README (#94, #95)
- Adds a benchmark for unsafe queries and world access (#98)

### Performance

- Implements archetype graph for faster lookup on transitions (#104)

### Bugfixes

- Locks the world during batch entity creation callback (#97)

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
