+++
title = 'Design philosophy'
type = 'docs'
weight = 130
description = "Ark's design philosophy and limitations."
+++
Unlike most other ECS implementations, Ark is primarily designed for the development of scientific,
individual-based models rather than for game development.
This motivates some design decisions, with an emphasis on **simplicity**, **safety** and **performance**.
Still, Ark remains perfectly usable for game development.

## Simple and type-safe API

The {{< api ecs World >}} object provides a minimal and pure ECS core — a data store for entities and components with robust querying and iteration capabilities.

Ark does not include built-in systems or an update loop.
These are left to the user, offering flexibility and control.
For a more complete, ready-to-use setup, see the [ark-tools](https://github.com/mlange-42/ark-tools) module.

Ark leverages Go’s generics to provide a fully type-safe API for queries, component access and manipulation. This design offers several advantages:

- Compile-time safety: Component types are enforced at compile time, eliminating runtime type assertions and reducing the risk of subtle bugs.
- Zero reflection: Unlike many ECS frameworks that rely on reflection, Ark’s generic API avoids it entirely, leading to superb performance.
- Optimized queries: Generic filters and mappers are statically typed and reusable, minimizing allocations and maximizing throughput.

For scenarios where component types are only known at runtime,
such as serialization, deserialization, or dynamic inspection,
Ark offers an additional [Unsafe API](../unsafe).
While it sacrifices compile-time type safety, this API provides maximum flexibility for advanced use cases.
It complements the generic API by enabling dynamic access to components,
making Ark suitable for tooling, debugging, and data-driven workflows where static typing isn't feasible.

## Performance-driven design

Ark is engineered for high performance, especially in large-scale simulations. Key optimizations include:

- [Archetype](../architecture)-based storage: Enables cache-friendly memory layout and fast iteration.
- Batch operations: Mass manipulation of entities and components is highly efficient.
- Reusable filters and mappers: Designed to minimize allocations and maximize throughput.

For hard numbers on performance, see the [Benchmarks](../benchmarks) chapter.
For a comparison that shows Ark outperforming many other Go ECS libraries, see the [go-ecs-benchmarks](https://github.com/mlange-42/go-ecs-benchmarks) repository.

## Determinism

Ark guarantees deterministic and reproducible iteration order.
While it doesn’t preserve insertion order or maintain consistent order across successive iterations. Both are impossible in archetype-based ECS. However, identical operations on the same World will always produce the same iteration sequence.

This deterministic behavior ensures that simulations yield consistent results across runs, platforms, and environments.

To reinforce this reliability, Ark is built with zero external dependencies. This eliminates variability introduced by third-party libraries and ensures that performance and behavior remain predictable and stable over time.

## Safety first: panic on violations

Ark puts an emphasis on safety and on avoiding undefined behavior.
It panics on unexpected operations, like removing a dead entity,
adding a component that is already present, or attempting to change a locked world.

While panics may seem unidiomatic in Go, Ark’s scientific context demands strict behavior.
Explicit error handling in performance-critical paths is impractical,
and silent failures are unacceptable.
For details, see the [Error handling](../errors) chapter.

## Limitations

The **number of component types** per World is capped at 256, a deliberate performance-oriented decision. This constraint enables extremely fast component lookups by using compact, array-based internal representations.

The **number of entities** alive at any one time is limited to just under 5 billion (`uint32` ID).

Ark is **not thread-safe**. This design choice avoids internal locking mechanisms, which would introduce overhead and complexity. In scientific modeling, where large numbers of simulations are often executed in parallel, this approach is more efficient and scalable.

> [!NOTE]
> The following feature is not yet released and is planned for Ark v0.6.0.
> You can try it out on the `main` branch.

**Concurrent query execution** is yet possible if the queries don't access the same entities concurrently.
For example, [entity relationships](../relations/) can be used to split up entities
of the same archetype to process them in parallel.
See the [parallel queries](https://github.com/mlange-42/ark/blob/main/examples/parallel_queries/main.go) example.