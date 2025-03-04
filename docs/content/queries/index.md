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

{{< api ecs Query2.Get >}} returns all queries components of the current entity.
The current entity can be obtained with {{< api ecs Query2.Entity >}}.

## World lock

The world gets locked for [component operations](../operations/) when a query is created.
The lock is automatically released when query iteration has finished.
When breaking out of the iteration, the query must be closed manually with {{< api ecs Query2.Close >}}.

The lock prevents entity creation and removal as well as adding and removing components.
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

{{< api ecs Filter2.Without >}} (and related methods) make the filter exclusive on the given components,
i.e. is excludes all other components:

{{< code-func queries_test.go TestQueriesExclusive >}}

### Optional

There is no `Optional` provided, as it would require an additional check in {{< api ecs Query2.Get >}} et al.
Instead, use {{< api ecs Map.Has >}}, {{< api ecs Map.Get >}} or similar methods in {{< api ecs Map2 >}} et al.

{{< code-func queries_test.go TestQueriesOptional >}}

## Filter caching
