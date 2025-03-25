+++
title = 'Resources'
type = "docs"
weight = 50
description = "Resources for singleton objects."
+++

Resources are data structures that are unique to an ECS world.
Examples could be the current game/simulation tick, a grid that your entities live on,
or an acceleration structure for spatial indexing.
As such, they can be thought of as components that exist only once and are not associated to an entity.

## Adding resources

Resources are Go structs that can contain any types of variables, just like [components](../concepts/#components).
Simply instantiate your resource and add a pointer to it to the world using {{< api ecs AddResource >}},
typically during world initialization:

{{< code-func resources_test.go TestAddResource >}}

The original resource struct can be stored and modified,
and changes are reflected in code that retrieves the resource from the world (see the following sections).

## Direct access

Resources can be retrieved from the world by their type, with {{< api ecs GetResource >}}:

{{< code-func resources_test.go TestResourceWorld >}}

However, this method has an overhead of approx. 20ns for the type lookup.
It is sufficient for one-time use of a resource.
When accessing a resource regularly, [Resource mappers](#resource-mappers) should be used.

## Resource mappers

Resource mappers are a more efficient way for retrieving a resource repeatedly.
To use them, create an {{< api ecs Resource >}}, store it, and use it for retrieval:

{{< code-func resources_test.go TestResourceMapper >}}

This way, resource access takes less than 1ns.

Resource mappers can also be used the add and remove resources, and to check for their existence:

{{< code-func resources_test.go TestResourceMapperAddRemove >}}
