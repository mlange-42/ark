+++
title = 'Entity relationships'
type = "docs"
weight = 50
description = "Entity relationships usage and details."
+++

In a basic ECS, relations between entities, like hierarchies, can be represented
by storing entities in components.
E.g., we could have a child component like this:

```go
type ChildOf struct {
    Parent ecs.Entity
}
```

Or, alternatively, a parent component with many children:

```go
type Parent struct {
    Children []ecs.Entity
}
```

In conjunction with [component mappers](../operations#component-mappers), this is often sufficient.
However, we are not able to leverage the power of queries to e.g. get all children of a particular parent in an efficient way.

To make entity relations even more useful and efficient, Ark supports them as a first class feature.
Relations are added to and removed from entities just like components,
and hence can be queried like components, with the usual efficiency.
This is achieved by creating separate archetypes
for relations with different target entities.

## Relation components

To use entity relations, create components that have *embedded* an {{< api ecs RelationMarker >}} as their first member:

```go
type ChildOf struct {
    ecs.RelationMarker
}
```

That's all to make a component be treated as an entity relationship by Ark.
The component can contain further variables, but the marker must be the first one.

## Creating relations

Most methods of `MapX` (e.g. {{< api ecs Map2 >}}) provide var-args for specifying relationship targets.
These are of type {{< api Relation >}}, which is an interface with multiple implementations:

{{< api Rel >}} is safe, but has some run-time overhead for component ID lookup during creation.

{{< api RelIdx >}} is fast but less safe.

See the examples below for their usage.

### On new entities

When creating entities, we can use a `MapX` (e.g. {{< api ecs Map2 >}}):

{{< code-func relations_test.go TestNewEntity >}}

For the faster variant {{< api RelIdx >}}, note that the first argument
is the zero-based index of the relation component in the {{< api ecs Map2 >}}'s generic parameters.

If there are multiple relation components, multiple {{< api Rel >}}/{{< api RelIdx >}} arguments can (and must) be used.

### When adding components

Relation target must also be given when adding relation components to an entity:

{{< code-func relations_test.go TestAdd >}}

## Set and get relations

We can also change the target entity of an already assigned relation component.
This is done via {{< api ecs Map2.SetRelations >}} et al.:

{{< code-func relations_test.go TestSetRelations >}}

Note that multiple relation targets can be changed in the same call.

Similarly, relation targets can be obtained with {{< api ecs Map2.GetRelation >}} et al.:

{{< code-func relations_test.go TestGetRelation >}}

Note that, due to Go's limitations on generics, the slow generic way is not possible here.

For a simpler syntax and when only a single relation component is accessed,
{{< api ecs Map >}} can be used alternatively:

{{< code-func relations_test.go TestMap >}}

## Batch operations

All [batch operation](../batch) methods of `MapX` (e.g. {{< api ecs Map2.NewBatch >}}) can be used with relation targets just like the normal component operations shown above.

## Filters and queries

[Filters](../queries) support entity relationships using the same syntax as shown in the examples above.

There are two ways to specify target entities to filter for: when building the filter, and when getting the query.
Both ways can be combined.

Relation targets given via {{< api ecs Map2.Relations >}} when building a filter are best used for permanent or long-lived targets.

{{< code-func relations_test.go TestFilter1 >}}

With [cached filters](../queries#filter-caching), the targets specified this way are included in the cache.
For short-lived targets, it is better to pass them when building a query with {{< api ecs Map2.Query >}}

{{< code-func relations_test.go TestFilter2 >}}

These targets are not cached, but the same filter can be used for different targets.

Filters also support both {{< api Rel >}} and {{< api RelIdx >}}.
In the filter examples above, we used the slow but safe {{< api Rel >}} when building the filter.
When getting the query, we use the faster {{< api RelIdx >}},
because in real-world use cases this is called more frequently than the one-time filter construction.

Relation targets not specified by the filter are treated as wildcard.
This means that the filter matches entities with any target.

## Dead target entities

Entities that are the target of any relationships can be removed from the world like any other entity.
When this happens, all entities that have this target in a relation get assigned to the zero entity as target.
The respective archetype is de-activated and marked for potential re-use for another target entity.

## When to use, and when not

When using Ark's entity relations, an archetype is created for each target entity of a relation.
Thus, entity relations are not efficient if the number of target entities is high (tens of thousands),
while only a low number of entities has a relation to each particular target (less than a few dozens).
Particularly in the extreme case of 1:1 relations, storing entities in components
as explained in the introduction of this chapter is more efficient.

However, with a moderate number of relation targets, particularly with many entities per target,
entity relations are very efficient.

Beyond use cases where the relation target is a "physical" entity that appears
in a simulation or game, targets can also be more abstract, like categories.
Examples:

 - Different tree species in a forest model.
 - Behavioral states in a finite state machine.
 - The opposing factions in a strategy game.
 - Render layers in a game or other graphical application.

This concept is particularly useful for things that would best be expressed by components,
but the possible components (or categories) are only known at runtime.
Thus, it is not possible to create ordinary components for them.
However, these categories can be represented by entities, which are used as relation targets.

## Longer example

To conclude this chapter, here is a longer example that uses Ark's entity relationships feature
to represent animals of different species in multiple farms.

{{< code relations_example_test.go >}}

Note that this examples uses the safe and clear, but slower generic variant to specify relationship targets.
As an optimization, {{< api RelIdx >}} could be used instead of {{< api Rel >}}, particularly for queries.