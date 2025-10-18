+++
title = 'Architecture'
type = 'docs'
weight = 140
description = "Ark's internal ECS architecture."
+++

Ark uses an archetype-based architecture.
This chapter explains the concept and Ark's take on it.

## Archetypes

The ASCII graph below illustrates the concept of archetypes in an entity-component system.

Each archetype represents a unique combination of components and stores data for all entities that share exactly that combination. You can think of an archetype as a table, where rows correspond to entities
and columns represent components. The first column always contains the entity identifiers themselves.

In the illustration below, the first archetype stores all entities that have components A, B and C,
as well as their associated component data.
The second archetype contains all entities with A and C, and their corresponding data.

```text
 Entities   Archetypes   Bitmasks     Queries

   E         E Comps
  |0|       |2|A|B|C|    111...   <-.      <---.
  |1|---.   |8|A|B|C|               |          |
  |2|   '-->|1|A|B|C|               |          |
  |3|       |3|A|B|C|               |--Q(A,C)  |
  |4|                               |  101...  |
  |6|   .-->|7|A|C|      101...   <-'          |--Q(B)
  |7|---'   |6|A|C|                            |  010...
  |8|       |4|A|C|                            |
  |9|---.                                      |
  |.|   |   |5|B|C|      011...            <---'
  |.|   '-->|9|B|C|
  |.|
  |.| <===> [Entity pool]
```
*Illustration of Ark's archetype-based architecture.*

Each archetype maintains a bitmask that encodes its unique component composition. 
This compact representation enables fast bitwise comparisons,
allowing queries to quickly determine which archetypes are relevant.

Once matching archetypes are identified, queries can linearly iterate over their entities, a process that is both highly efficient and cache-friendly, thanks to the tight memory layout of archetype tables.

Component access through queries is extremely fast (below 1ns per component) due to predictable memory access patterns and minimal indirection.

## World entity access

To retrieve components for a specific entity outside of query execution,
the World maintains a list indexed by entity ID (as shown leftmost in the diagram above).
Each entry in this list points to the entity's archetype and the position within the archetype's table.

This setup enables fast random access to component data, though slightly slower than query-based iteration (≈2ns vs. <1ns), due to the additional indirection.

Note that the entity list also contains entities that are currently not alive,
because they were removed from the {{< api ecs World >}}.
These entities are recycled when new entities are requested from the world.
Therefore, besides the ID shown in the illustration, each entity also has a generation
variable. It is incremented on each "reincarnation",
which allows to distinguish recycled from dead entities, as well as from previous or later "incarnations".

## Entity recycling and generations

The entity list also includes entities that are no longer alive because they have been removed from the World. These inactive entities are recycled when new entities are created, allowing efficient reuse of memory and IDs.

To safely distinguish between recycled and stale references, each entity carries a generation counter.
This counter is incremented every time an entity is "reincarnated".
It ensures that references to old or removed entities can be detected and invalidated,
and that different "incarnations" of same entity ID can be distinguished.

## Performance

Archetypes are primarily designed to maximize iteration speed by grouping entities with identical component sets into tightly packed memory layouts.
This structure enables blazing-fast traversal and component access during queries.

However, this optimization comes with a trade-off: Adding or removing components from an entity requires relocating it to a different archetype, essentially moving all of its component data. This operation typically costs &approx;20ns per involved component.

To reduce the number of archetype changes, it is recommended to add/remove/exchange multiple components at the same time rather than one after the other. Further, operations can be [batched](../batch) to manipulate many entities in a single command. See chapter [Performance tips](../performance) for more details.

For detailed benchmarks and performance metrics, refer to the [Benchmarks](../benchmarks) chapter.

## Details

The previous explanation offers a simplified view of archetypes. To fully understand the system, we need to consider two advanced concepts: the Archetype Graph and [Entity relationships](../relations).

### Archetype graph

When components are added to or removed from an entity, the system must locate the corresponding archetype that matches the new component composition. To accelerate this process, Ark uses a dynamic graph of archetype nodes (or just nodes).
The figure below illustrates the concept.

Each arrow represents the transition between two archetypes when a single component is added (solid arrow head)
or removed (empty arrow head).
Following these transitions, the archetype resulting from addition and/or removal of an arbitrary number
of components can be found easily.

{{< html >}}
<img alt="Archetype graph light" width="600" class="light" src="./images/archetype-graph.svg"></img>
<img alt="Archetype graph dark" width="600" class="dark" src="./images/archetype-graph-dark.svg"></img>
{{< /html >}}  
*Illustration of the archetype graph. Letters represent components. Boxes represent archetype nodes.
Arrows represent transitions when a single component is added or removed.*

Nodes and transitions are created on demand. When searching for an archetype, the algorithm proceeds transition by transition.
When looking for the next archetype, established transitions are checked first.
If this is not successful, the resulting component mask is used to search through all nodes.
On success, a new connection is established.
If the required node was still not found, a new node is created.
Then, the next transition it processed and so on, until the final node is found.
Only then, an archetype is created for the node.

As a result, the graph will usually not be fully connected.
There will also not be all possible nodes (combinations of components) present.
Nodes that are only traversed by the search but never receive entities contain no archetype and are called inactive.

During a game or simulation run, the graph stabilizes quickly.
Then, only the fast following of transitions is required to find an archetype when components are added or removed.
Transitions are stored in the nodes with lookup approx. 10 times faster than Go's `map`.

### Entity relations

Earlier, archetypes were described as flat tables.
However, with Ark’s [Entity relationships](../relations) feature,
archetypes can contain multiple sub-tables, each corresponding to a unique combination of relation targets.

As an example, we have components `A`, `B` and `R`, where `R` is a relation.
Further, we have two parent entities `E1` and `E2`.
When you create some entities with components `A B R(E1)` and `A B R(E2)`,
i.e. with relation targets `E1` and `E2`,
the following archetype is created:

```text

  Archetype [ A B R ]
    |
    |--- E1   E Comps
    |        |3|A|B|R|
    |        |6|A|B|R|
    |        |7|A|B|R|
    |
    '--- E2   E Comps
             |4|A|B|R|
             |5|A|B|R|
```

When querying without specifying a target, the archetype's tables are simply iterated if the archetype matches the filter.
When querying with a relation target (and the archetype matches), the table for the target entity is looked up in a standard Go `map`.

If the archetype contains multiple relation components, a `map` lookup is used to get all tables matching the target that is specified first. These tables are simply iterated if no further target is specified. If more than one target is specified, the selected tables are checked for these further targets and skipped if they don't match.

### Archetype removal

Normal archetype tables without a relation are never removed, because they are not considered temporary.
For relation archetypes, however, things are different.
Once a target entity dies, it will never appear again (actually it could, after dying another 4,294,967,294 times).

In Ark, empty tables with a dead target are recycled.
They are deactivated, but their allocated memory for entities and components is retained.
When a table in the same archetype, but for another target entity is requested, a recycled table is reused if available.
To be able to efficiently detect whether a table can be removed,
a bitset is used to keep track of entities that are the target of a relation.
