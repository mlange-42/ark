# Ark examples

This folder contains examples for the [Ark](https://github.com/mlange-42/ark) Entity Component System.

- [readme](./readme/main.go): The minimal example from the README, for automatic testing by the GitHub CI.
- [systems](./systems/main.go): Demonstrates how to implement systems and a scheduler.
- [world_lock](./world_lock/main.go): Demonstrates how to manipulate entities despite the world being locked during query iteration.
- [relations](./relations/main.go): Demonstrates the use of Ark's entity relationships feature.
- [ebitengine](./ebitengine/): Demonstrates how to use Ark with the [Ebiten](https://ebitengine.org/) game engine.
- [parallel](./parallel/main.go): Demonstrates how to run multiple simulations in parallel.
- [entity_grid](./entity_grid/main.go): Demonstrates that ECS can be mixed with non-ECS data structures, using a grid of entities.
- [kdtree](./kdtree/main.go): Demonstrates that ECS can be mixed with non-ECS data structures, using a kdtree.

## Running the examples

1. Clone the repository

```
git clone https://github.com/mlange-42/ark
cd ark
```

2. Most examples can be run from directory `examples` like this:

```
cd examples
go run ./<example>
```

An exception if the [ebitengine](./ebitengine/) example.
Due to its relatively heavy dependencies, it is a separate module and must be run from within `examples/ebitengine`:

```
cd examples/ebitengine
go run .
```
