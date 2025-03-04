+++
title = 'Unsafe API'
type = "docs"
weight = 70
description = "Ark's unsafe, ID-based API."
+++
So far, we used the type-safe, generic API of Ark throughout.
However, there may be use cases where component types are not known at compile-time,
like serialization and de-serializations.

For these cases, Ark offers an unsafe, ID-based API.
It is accessible via {{< api ecs World.Unsafe >}}.

## Component IDs

Internally, each component type is mapped to an {{< api ecs ID >}}.
We don't see it in the generic API, but we can use it for more flexibility in the ID-based API.
IDs can be obtained by the function {{< api ecs ComponentID >}}.
It a component is not yet registered, it gets registered upon first use.

{{< code-func unsafe_test.go TestUnsafeIDs >}}

## Creating entities

Entities are created with {{< api ecs Unsafe.NewEntity >}}, giving the desired component IDs:

{{< code-func unsafe_test.go TestUnsafeNewEntity >}}

## Filters and queries

Filters and queries work similar to the generic API, but also use component IDs:

{{< code-func unsafe_test.go TestUnsafeQuery >}}

## Component access

Components of entities can be accessed outside queries using {{< api ecs Unsafe.Get >}}/{{< api ecs Unsafe.Has >}}:

{{< code-func unsafe_test.go TestUnsafeGet >}}

## Component operations

Components can be added and removed using methods of {{< api ecs Unsafe >}}:

{{< code-func unsafe_test.go TestUnsafeComponents >}}

## Limitations

Besides not being type-safe (i.e. you can cast to anything without an error), the unsafe API has a few limitations:

- It is slower than the type-safe API.
- There are currently no batch operations provided for it.
