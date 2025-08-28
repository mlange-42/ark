+++
title = 'Architecture'
type = 'docs'
weight = 130
description = "Ark's internal ECS architecture."
+++

Ark uses an archetype-based architecture.
This chapter explains the concept and Ark's take on it.

## Archetypes

The ASCII graph below illustrates the approach of archetypes.
Each archetype stores the component data for all entities with exactly the same components.
Archetype storage can be imagined as a table, with entities in rows and components in columns.
Additionally, an archetype has a column for the actual entity (1st column in the figure).

In the illustration below, the first archetype holds all entities with (only/exactly) the components A, B and C,
as well as their components.
Similarly, the second archetype contains all entities with A and C, and their components.

```text
 Entities   Archetypes   Bit masks    Queries

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

Each archetype contains a bit-mask, encoding its component composition for fast comparison.
This way, queries can easily identify their relevant archetypes, and then simply iterate entities linearly,
which is very fast and cache-friendly.
Components can be accessed through a query very efficiently (&approx;1ns per component).

## World entity access

For getting components for an entity outside of queries, the world contains a list that is indexed by the entity ID (left-most in the figure above). For each entity, it references its current archetype and the position of the entity in the archetype. This way, getting components for entities (i.e. random access) is fast, although not as fast as in queries (≈2ns vs. 1ns).

Note that the entity list also contains entities that are currently not alive,
because they were removed from the {{< api ecs World >}}.
These entities are recycled when new entities are requested from the world.
Therefore, besides the ID shown in the illustration, each entity also has a generation
variable. It is incremented on each "reincarnation",
which allows to distinguish recycled from dead entities, as well as from previous or later "incarnations".

## Performance

Obviously, archetypes are an optimization for iteration speed.
But they also come with a downside. Adding or removing components to/from an entity requires moving all the components of the entity to another archetype.
This takes roughly 10-20ns per involved component.
To reduce the number of archetype changes, it is recommended to add/remove/exchange multiple components at the same time rather than one after the other (see chapter [Performance tips](../performance) for more details).

For more numbers on performance, see chapter [Benchmarks](../benchmarks). 

## Details

Actually, the explanation above is quite simplified.
Particularly it leaves out [Entity relationships](../relations) and the *archetypes graph*.

### Archetype graph

When components are added to or removed from an entity, it is necessary to find its new archetype.
To accelerate the search, a graph of *archetype nodes* (or just *nodes*) is used.
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

Nodes and connections are created as needed. When searching for an archetype, the algorithm proceeds transition by transition.
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

At the beginning of this chapter, archetypes were described as being tables.
However, with Ark's [Entity relationships](../relations) feature, archetypes can contain multiple tables,
one for each unique combination of relationship targets.
As an example, we have components `A`, `B` and `R`, where `R` is a relation.
Further, we have two parent entities `E1` and `E2`.
When you create some entities with components `A B R(E1)` and `A B R(E2)`, i.e. with relation targets `E1` and `E2`,
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

Normal archetype tables without a relation are never removed, because they are not a temporary thing.
For relation archetypes, however, things are different.
Once a target entity dies, it will never appear again (actually it could, after dying another 4294967294 times).

In Ark, empty tables with a dead target are recycled.
They are deactivated, but their allocated memory for entities and components is retained.
When a table in the same archetype, but for another target entity is requested, a recycled table is reused if available.
To be able to efficiently detect whether a table can be removed,
a bitset is used to keep track of entities that are the target of a relation.
