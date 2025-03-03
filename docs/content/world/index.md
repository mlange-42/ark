+++
title = 'The World'
type = "docs"
weight = 30
description = "The World as Ark's central data storage."
+++

The *World* ({{< api ecs World >}}) is the central data storage in Ark.
It manages and stores entities ({{< api ecs Entity >}}), their components, as well as [Resources](../resources).

Here, we only deal with world creation.
Most world functionality is covered in chapters [Entities & Components](../entities) and [World Entity Access](../world-access).

## World creation

To create a world, use {{< api ecs NewWorld >}} with an initial capacity:

{{< code-func world_test.go TestWorldSimple >}}

The initial capacity is used to initialize archetypes, the entity list, etc.

## Reset the world

For systematic simulations, it is possible to reset a populated world for reuse:

{{< code-func world_test.go TestWorldReset >}}
