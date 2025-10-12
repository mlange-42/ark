+++
title = 'Concepts'
type = "docs"
weight = 20
description = "ECS concepts in Ark."
+++
This chapter gives a brief explanation of ECS concepts and how they are represented in Ark.

## The World

The World ({{< api ecs World >}}) is the central data storage in Ark.
It manages and stores entities ({{< api ecs Entity >}}), their [components](#components), as well as [Resources](../resources).

To create a world with default settings, use {{< api ecs NewWorld >}}:

{{< code-func concepts_test.go TestWorldSimple 0 1 >}}

A world can also be configured with an initial capacity for [archetypes](../architecture), the entity list, etc:

{{< code-func concepts_test.go TestWorldConfig 0 1 >}}

For systematic simulations, it is possible to reset a populated world for reuse:

{{< code-func concepts_test.go TestWorldReset >}}

## Entities

Entities ({{< api ecs Entity >}}) are used to represent the "objects" in a game or simulation.
In Ark, an entity is just an opaque ID (with an associated current generation) that allows access
to [components](#components) associated to it.

Entities without any components can be created through the ({{< api ecs World >}}):

{{< code-func concepts_test.go TestCreateEntitySimple 0 1 >}}

For creating entities with components, [component mappers](../components#component-mappers) are used.
Entities can be removed or deleted like this:

{{< code-func concepts_test.go TestRemoveEntity >}}

Entities can be stored safely, in components or elsewhere.
However, always store them by value, never by pointer.
When dealing with stored entities, it may be required to check whether they are still alive:

{{< code-func concepts_test.go TestEntityAlive >}}

In Ark, entities are returned to a pool when they are removed from the world.
These entities can be recycled, with the same ID ({{< api ecs Entity.ID Entity.ID >}}), but an incremented generation ({{< api ecs Entity.Gen Entity.Gen >}}).
This allows to determine whether an entity held by the user is still alive, despite it was potentially recycled.

## Components

Components contain the data, or state variables, associated with an entity.
Each entity can have an arbitrary combination of components,
but can only have one instance of each.
Components can be added to and removed from entities at runtime.

Components are identified by their type. So e.g., all instances of `Position`
are of the same component type. Each entity can have only one `Position`.

Components are simple Go structs and can contain variables of any type,
including slices and pointers.
Typically, components don't have any functions.
Particularly, contrary to object oriented programming (OOP),
components don't contain game or simulation logic.
In ECS, logic is performed by [systems](#systems), using [queries](#queries).

For optimal performance and modularity, components should be small
and only contain closely related state variables that are typically used together.
"Good" components are e.g. `Position`, `Velocity`, `Age`, `Sex`, etc.
"Bad" components are large monolithic things with many state variables like `Player` or `Animal`.

Components can also be labels or tags, which means that they don't contain any data
but are just used to tag entities, like `Female` and `Male`.

A world can contain up to 256 different component types (64 with build tag `ark_tiny`).

See chapter [Component operations](../operations) for how to create entities with components,
adding and removing components, and other details.

## Queries

Queries are the core feature for writing logic in an ECS.
A query iterates over all entities that possess all the component types specified by the query.
Note that these entities may contain further components, which are ignored.

For best performance, filters are used to create queries:

{{< code-func concepts_test.go TestQuery >}}

Filters are relatively costly to create, as they require lookup of component IDs.
This takes around 20ns per component.
Thus, make sure to create filters only once and store them, e.g. in [systems](#systems).
Then, create a new query from the filter each time before the iteration loop.

> [!IMPORTANT]
> The component pointers obtained from a query should never be stored
> outside of the current context (i.e. the query loop), as they are not persistent inside the world.

For advanced filters, caching and other details, see chapter [Filters & queries](../queries).

## Systems

Systems perform the logic of your game or simulation, using queries.
Ark does not provide systems or a scheduler for them.
You can create your own interface for system, matching your game engine if you are using one.
Alternatively, [ark-tools](https://github.com/mlange-42/ark-tools) provides systems,
a scheduler, and other useful stuff for Ark. See there for a full example.

## Resources

Resources are singular data structures in an ECS world.
As such, they can be thought of as components that exist only once and are not associated to an entity.
Examples could be the current game/simulation tick, a grid that your entities live on,
or an acceleration structure for spatial indexing.

As with [components](#components), resources are Go structs that can contain any types of variables.

{{< code-func concepts_test.go TestResource 0 7 >}}

See chapter [Resources](../resources) for more details.

## Relationships

Entity relationships are a powerful, advanced ECS feature that was first introduced by [Flecs](https://www.flecs.dev/flecs/).
They serve the efficient representation of entity hierarchies, groupings, or other relationships.

Relationships can also be realized by storing entities (or lists of entities) in components.
However, Ark's relations feature allows for more efficiency and more comprehensible logic, as relationships can be used in [queries](#queries).

Some illustrative examples:

- Iterate/get all game objects in a grid cell.
- Iterate/get all animals in a herd.
- Iterate/get all plants of a certain species.
- Build hierarchies, like a scene graph.

Compared to Flecs, entity relations in Ark are more limited.
Each entity can have an arbitrary number of relationships to other entities,
but for each relation type (i.e. relation component), there can be only one target entity.
This is primarily a performance consideration.

For usage and more details, see chapter [Entity relationships](../relations).
