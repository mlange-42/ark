# Changelog

## [[v0.2.0]](https://github.com/mlange-42/ark/compare/v0.1.0...v0.2.0)

### Features

- Adds `MapX.NewBatch` and `MapX.NewBatchFn` for fast batch-creation of entities (#65)
- Adds `ExchangeX.Add` and `ExchangeX.Remove` (#70)
- Adds `Map.GetRelationUnchecked` (#75)
- Adds build tag `tiny` for 64 bit masks (#77)
- Adds `QueryX.GetRelation`, `MapX.GerRelation` and `MapX.GetRelationUnchecked` (#79)
- Adds unsafe, ID-based API (#80, #81, #82)
- Adds several functions and methods required for [ark-serde](https://github.com/mlange-42/ark-serde) (#83, #85)
- Adds resource shortcut functions `GetResource` and `AddResource` (#86)
- Adds `World.Reset` (#89)

## [[v0.1.0]](https://github.com/mlange-42/ark/tree/v0.1.0)

Initial release of Ark.

Basic ECS implementation, as well as entity relationships.
