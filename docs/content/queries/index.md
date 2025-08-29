+++
title = 'Filters & queries'
type = "docs"
weight = 30
description = "Ark's queries and filters."
+++

Queries are the core feature for writing logic in an ECS.
A query iterates over all entities that possess all the component types specified by the query.

Queries are constructed from filters.
While queries are one-time use iterators that are cheap to create,
filters are more costly to create and should be stored permanently, e.g. in your [systems](../concepts#systems).

## Filters and queries

With basic filters, queries iterate all entities that have the given components,
and any additional components that are not of interest.

In the example below, the filter would match any entities that have
`Position` and `Velocity`, and potentially further components like `Altitude`.

{{< code-func queries_test.go TestQueriesBasic >}}

{{< api ecs Query2.Get >}} returns all queried components of the current entity.
The current entity can be obtained with {{< api ecs Query2.Entity >}}.

## Query performance

Queries iteration is what an [archetype](../architecture)-based ECS is optimized for, and it is really fast.
This has two reasons.

Firstly, all entities with the same component composition are stored in the same archetype, or "table".
This means that filters only need to be checked against archetypes,
and the entities of a matching archetype can be iterated without any further checks.

Secondly, all components of the same type (like `Position`) are stored in a dedicated column of the archetype.
A query only accesses the required components (i.e. columns), although entities may possess many more components.
Memory access is therefore completely linear and contiguous, and the CPUs cache is used as efficiently as possible.

## World lock

The world gets locked for [component operations](../operations/) when a query is created.
The lock is automatically released when query iteration has finished.
When breaking out of the iteration, the query must be closed manually with {{< api ecs Query2.Close >}}.

The lock prevents entity creation and removal, as well as adding and removing components.
Thus, it may be necessary to collect entities during the iteration, and perform the operation afterwards:

{{< code-func queries_test.go TestQueriesLock >}}

## Advanced filters

Filters can be further specified using method chaining.

### With

{{< api ecs Filter2.With >}} (and related methods) allow to specify components that the queried entities should possess,
but that are not used inside the query iteration:

{{< code-func queries_test.go TestQueriesWith >}}

`With` can also be called multiple times instead of specifying multiple components in one call:

{{< code-func queries_test.go TestQueriesWith2 >}}

### Without

{{< api ecs Filter2.Without >}} (and related methods) allow to specify components that the queried entities should *not* possess:

{{< code-func queries_test.go TestQueriesWithout >}}

As with `With`, `Without` can be called multiple times:

{{< code-func queries_test.go TestQueriesWithout2 >}}

### Exclusive

{{< api ecs Filter2.Exclusive >}} (and related methods) make the filter exclusive on the given components,
i.e. it excludes all other components:

{{< code-func queries_test.go TestQueriesExclusive >}}

### Optional

There is no `Optional` provided, as it would require an additional check in {{< api ecs Query2.Get >}} et al.
Instead, use {{< api ecs Map.Has >}}, {{< api ecs Map.Get >}} or similar methods in {{< api ecs Map2 >}} et al.:

{{< code-func queries_test.go TestQueriesOptional >}}

## Filter caching

Although queries are highly performant, a huge number of [archetypes](../architecture) (like hundreds or thousands) may cause a slowdown.
To prevent this slowdown, filters can be registered to the world's filter cache via
{{< api ecs Filter2.Register >}}:

{{< code-func queries_test.go TestQueriesCached >}}

For registered filters, a list of matching archetypes is cached internally.
Thus, no filter evaluations are required during iteration.
Instead, filters are only evaluated when a new archetype is created.

When a registered filter is not required anymore, it can be unregistered with
{{< api ecs Filter2.Unregister >}}.
However, this is rarely required as (registered) filters are usually used over an entire game session or simulation run.
