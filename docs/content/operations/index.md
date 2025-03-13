+++
title = 'Component operations'
type = "docs"
weight = 40
description = "Manipulating components."
+++

Components contain the data associated to your game or simulation entities.
This chapter explains how to to manipulate them.
For general information on components, see section [Components](../concepts#components) in chapter [Concepts](../concepts).

## Component mappers

Component mappers {{< api ecs Map1 >}}, {{< api ecs Map2 >}} etc.
are helpers that allow to create entities with components,
to add components to entities, and to remove components from entities.
They are parametrized by the component types they handle.

{{< code-func operations_test.go TestComponentMapper >}}

In this example, the `2` in {{< api ecs NewMap2 >}} denotes the number of mapped components.
Unfortunately, this is required due to the limitations of Go's generics.

In addition to {{< api ecs Map1 >}}, {{< api ecs Map2 >}}, etc., there is {{< api ecs Map >}}.
It is a dedicated mapper for a single component and provides a few additional methods.

Component mappers are also used to access components for specific entities:

{{< code-func operations_test.go TestComponentMapperGet >}}

> [!IMPORTANT]
> The component pointers obtained should never be stored
> outside of the current context, as they are not persistent inside the world.

## Component exchange

Adding and removing components are relatively costly operations,
as entities and their components must be moved between archetypes.
It is most efficient to perform component additions and removals in a single operation,
instead of using multiple operations.

For that sake, Ark provides {{< api ecs Exchange1 >}}, {{< api ecs Exchange2 >}} etc.,
to do additions and removals in one go.
It adds the components given by the generics, and removes those specified with {{< api ecs Exchange2.Removes >}}.

{{< code-func operations_test.go TestExchange >}}
