+++
title = 'Batch operations'
type = "docs"
weight = 60
description = "Ark's queries and filters."
+++

In an archetype-based ECS, creation and removal of entities or components are relatively costly operations.
For these operations, Ark provides batched versions.
This allows to create or manipulate a large number of entities much faster than one by one.
Most batch methods come in two flavors. A "normal" one, and one that runs a callback function on each affected entity.

## Creating entities

Entity creation is probably the most common use case for batching.
When the number of similar entities that are to be created is known,
creation can be batched with {{< api ecs Map2.NewBatch >}} et al.:

{{< code-func batch_test.go TestNewBatch >}}

{{< api ecs Map2.NewBatchFn >}} can be used for more flexible initialization:

{{< code-func batch_test.go TestNewBatchFn >}}

## Components

Components can be added, removed or exchanged in batch operations.
For these operations, {{< api ecs Map2 >}}, {{< api ecs Exchange2 >}} etc.
provide batch versions of the respective methods.
Component batch operations take an {{< api ecs Batch >}} filter as an argument to determine the affected entities:

{{< code-func batch_test.go TestBatchComponents >}}

## Removing entities

Entities can be removed in batches using {{< api ecs World.RemoveEntities >}}:

{{< code-func batch_test.go TestRemoveEntities >}}
