+++
title = 'World statistics'
type = "docs"
weight = 90
description = "Ark's world statistics feature for engine insights."
+++
Ark only exposes the API required for actual use.
Therefore, internals like the number of [archetypes](../architecture), memory used to store components etc. are not directly accessible.

However, it might sometimes be useful to have access to such metrics,
for example in order to judge effects of different ways of implementing something.
Otherwise, users would have to rely on logic reasoning and sufficient understanding of Ark to derive these numbers.

For that sake, Ark provides statistics about its internals, prepared in a compact and digestible form.

## Accessing statistics

All internal statistics can be accessed via {{< api ecs World.Stats >}},
which returns a {{< api "ecs/stats" World "*stats.World" >}}.
This, in turn, contains the other stats types described below.
All these types have a method `String()` to bring them into a compact, human-readable form. 

{{< code-func stats_test.go TestWorldStats >}}

Which prints:

```text
World     -- Components: 2, Archetypes: 2, Filters: 0, Memory: 56.0 kB, Locked: false
             Components: Position, Heading
Entities  -- Used: 100, Recycled: 0, Total: 100, Capacity: 1026
Archetype -- Tables:    1, Comps:  0, Entities:      0, Cap:   1024, Mem:     8.0 kB, Per entity:    8 B
             Components:
Archetype -- Tables:    1, Comps:  2, Entities:    100, Cap:   1024, Mem:    32.0 kB, Per entity:   32 B
             Components: Position, Heading
```

## World stats

{{< api "ecs/stats" World stats.World >}} provides world information like a list of all component types
and the total memory reserved for entities and components.
Further, it contains {{< api "ecs/stats" Entities stats.Entities >}} and
a {{< api "ecs/stats" Archetype stats.Archetype >}} for each archetype.

## Entity stats

{{< api "ecs/stats" Entities stats.Entities >}} contains information about the entity pool,
like capacity, alive entities and available entities for recycling.

## Archetype stats

{{< api "ecs/stats" Archetype stats.Archetype >}} provides information about an archetype, like its components,
memory in total and per entity, and more state information.

Further, it contains a {{< api "ecs/stats" Table stats.Table >}} for each table.

## Table stats

{{< api "ecs/stats" Table stats.Table >}} contains size, capacity and memory information for a table.
Tables are used to represent sub-archetypes with the same components, but a different combination
of [relationship](../relations) targets.
