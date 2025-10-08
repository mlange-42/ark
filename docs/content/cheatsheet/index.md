+++
title = 'Cheat sheet'
type = "docs"
weight = 1000
description = "Cheat sheet for frequent use cases."
+++
Frequently used Ark operations for quick lookup.

{{< details-buttons paragraph="true" >}}

## World creation

{{< details-buttons group="world" >}}

The world is the central ECS data storage.
Most applications will use exactly one world.

{{% details closed="true" group="world" title="🌍 Create a World with default initial capacity" %}}
{{< code-func cheatsheet_test.go TestCreateWorld >}}
API: {{< api ecs World.New >}}
{{% /details %}}

{{% details closed="true" group="world" title="🌍 Create a World with a specific initial capacity" %}}
{{< code-func cheatsheet_test.go TestCreateWorldConfig >}}
API: {{< api ecs World.New >}}
{{% /details %}}

## Create entities

{{< details-buttons group="create" >}}

{{% details closed="true" group="create" title="✨ Create an **entity without components**" %}}
Create an **entity without components**:
{{< code-func cheatsheet_test.go TestCreateEmpty >}}
API: {{< api ecs World.NewEntity >}}
{{% /details %}}

{{% details closed="true" group="create" title="✨ A **component mapper** is required for creating **entities with components**" %}}
Component mappers should be stored and re-used for best performance.
{{< code-func cheatsheet_test.go TestCreateMapper >}}
API: {{< api ecs Map1 >}},  {{< api ecs Map2 >}}, ...
{{% /details %}}

{{% details closed="true" group="create" title="✨ Create a **single entity**, given some components" %}}
{{< code-func cheatsheet_test.go TestCreateEntity >}}
API: {{< api ecs Map2.NewEntity >}}
{{% /details %}}

{{% details closed="true" group="create" title="✨ Create a **single entity** using a callback" %}}
{{< code-func cheatsheet_test.go TestCreateEntityFn >}}
API: {{< api ecs Map2.NewEntityFn >}}
{{% /details %}}

{{% details closed="true" group="create" title="✨ Create **many entities** more efficiently, all with the same component values" %}}
{{< code-func cheatsheet_test.go TestCreateBatch >}}
API: {{< api ecs Map2.NewBatch >}}
{{% /details %}}

{{% details closed="true" group="create" title="✨ Create **many entities**, using a callback for individual initialization" %}}
{{< code-func cheatsheet_test.go TestCreateBatchFn >}}
API: {{< api ecs Map2.NewBatchFn >}}
{{% /details %}}

## Remove entities

{{< details-buttons group="delete" >}}

{{% details closed="true" group="delete" title="❌ Remove a **single entity**" %}}
{{< code-func cheatsheet_test.go TestRemoveEntity >}}
API: {{< api ecs World.RemoveEntity >}}
{{% /details %}}


{{% details closed="true" group="delete" title="❌ Remove **all entities** that match a filter" %}}
{{< code-func cheatsheet_test.go TestRemoveEntities >}}
API: {{< api ecs World.RemoveEntities >}}
{{% /details %}}

{{% details closed="true" group="delete" title="❌ With a **callback** to do something with entities before their removal" %}}
{{< code-func cheatsheet_test.go TestRemoveEntitiesFn >}}
API: {{< api ecs World.RemoveEntities >}}
{{% /details %}}

## Add/remove components

{{< details-buttons group="components" >}}

{{% details closed="true" group="components" title="🧩 A **component mapper** is required for adding and removing components" %}}
It adds or removes the given components to/from entities.
Component mappers should be stored and re-used for best performance.
{{< code-func cheatsheet_test.go TestCreateMapper >}}
API: {{< api ecs Map1 >}},  {{< api ecs Map2 >}}, ...
{{% /details %}}

{{% details closed="true" group="components" title="🧩 Add and remove components to/from a **single entity**" %}}
{{< code-func cheatsheet_test.go TestAddRemoveComponents >}}
API: {{< api ecs Map2.Add >}}, {{< api ecs Map2.Remove >}}
{{% /details %}}

{{% details closed="true" group="components" title="🧩 Add components to **all entities** matching a filter" %}}
{{< code-func cheatsheet_test.go TestAddBatch >}}
API: {{< api ecs Map2.AddBatch >}}
{{% /details %}}

{{% details closed="true" group="components" title="🧩 Add components to **all entities** matching a filter, with individual initialization" %}}
{{< code-func cheatsheet_test.go TestAddBatchFn >}}
API: {{< api ecs Map2.AddBatchFn >}}
{{% /details %}}

{{% details closed="true" group="components" title="🧩 Remove components from **all entities** matching a filter." %}}
The callback can be used to do something with entities before component removal.
{{< code-func cheatsheet_test.go TestRemoveBatch >}}
API: {{< api ecs Map2.RemoveBatch >}}
{{% /details %}}

## Filters and queries

{{< details-buttons group="queries" >}}

{{% details closed="true" group="queries" title="🔍 Use **filters and queries** to iterate entities" %}}
Filters should be stored and re-used for best performance.  
Always create a new query before iterating.
{{< code-func cheatsheet_test.go TestFilterQuery >}}
API: {{< api ecs Filter1 >}}, {{< api ecs Filter2 >}}, ..., {{< api ecs Filter2.Query >}}
{{% /details %}}

{{% details closed="true" group="queries" title="🔍 Filters can match **additional components**" %}}
For components the entities should have, but that are not accessed in the query.
{{< code-func cheatsheet_test.go TestFilterWith >}}
API: {{< api ecs Filter2.With >}}
{{% /details %}}

{{% details closed="true" group="queries" title="🔍 Filters can **exclude components**" %}}
{{< code-func cheatsheet_test.go TestFilterWithout >}}
API: {{< api ecs Filter2.Without >}}
{{% /details %}}

{{% details closed="true" group="queries" title="🔍 Filters can be **exclusive** on the given components" %}}
This filter matches only entities with exactly the given components.
{{< code-func cheatsheet_test.go TestFilterExclusive >}}
API: {{< api ecs Filter2.Exclusive >}}
{{% /details %}}

{{% details closed="true" group="queries" title="🔍 Filters can **combine** multiple conditions" %}}
{{< code-func cheatsheet_test.go TestFilterWithWithout >}}
API: {{< api ecs Filter2.With >}}, {{< api ecs Filter2.Without >}}
{{% /details %}}

{{% details closed="true" group="queries" title="🔍 Queries can **count entities** without iterating" %}}
Note that a query that is not iterated must be closed explicitly.
{{< code-func cheatsheet_test.go TestQueryCount >}}
API: {{< api ecs Query.Count >}}, {{< api ecs Query.Close >}}
{{% /details %}}

## Resources

{{< details-buttons group="resources" >}}

Resources can be used for "global", singleton-like data structures that are not associated to particular entities.

{{% details closed="true" group="resources" title="📦 Adding and getting resources, the simple but slower way (&approx;20ns)" %}}
{{< code-func cheatsheet_test.go TestResourcesQuick >}}
API: {{< api ecs AddResource >}}, {{< api ecs GetResource >}}
{{% /details %}}

{{% details closed="true" group="resources" title="📦 For repeated access, better use a resource accessor (`Get()` &approx;1ns)" %}}
{{< code-func cheatsheet_test.go TestResources >}}
(Creating the accessor does not add the actual `Grid` resource!)

API: {{< api ecs Resource >}}, {{< api ecs Resource.Get >}}
{{% /details %}}

## Events and observers

{{< details-buttons group="events" >}}

> [!NOTE]
> This feature is not yet released and is planned for Ark v0.6.0.
> You can try it out on the `main` branch.

{{% details closed="true" group="events" title="👀 **Create and register** observers for ECS lifecycle events" %}}
Gets notified on any creation of an entity.
{{< code-func cheatsheet_test.go TestObserver >}}
API: {{< api ecs Observer >}}
{{% /details %}}

{{% details closed="true" group="events" title="👀 Observers can **filter** for certain components" %}}
Gets notified when a `Position` and a `Velocity` component are added to an entity.
{{< code-func cheatsheet_test.go TestObserverFilter >}}
API: {{< api ecs Observer.For >}}
{{% /details %}}

{{% details closed="true" group="events" title="👀 Observers can **process matched components**" %}}
Gets notified when a `Position` and a `Velocity` component are added to an entity,
with both available in the callback.
{{< code-func cheatsheet_test.go TestObserverFilterGeneric >}}
API: {{< api ecs Observer1 >}}, {{< api ecs Observer2 >}}, ...
{{% /details %}}

{{% details closed="true" group="events" title="👀 Observers can also filter for the **entity composition**" %}}
Gets notified when a `Position` component is added to an entity
that also has `Velocity` but not `Altitude`.
{{< code-func cheatsheet_test.go TestObserverWithWithout >}}
API: {{< api ecs Observer1.With >}}, {{< api ecs Observer1.Without >}}
{{% /details %}}

{{% details closed="true" group="events" title="📣 **Custom event types** can be created using a registry" %}}
{{< code-func cheatsheet_test.go TestCustomEventType >}}
API: {{< api ecs EventRegistry >}}, {{< api ecs EventRegistry.NewEventType >}}
{{% /details %}}

{{% details closed="true" group="events" title="📣 **Custom events** can emitted by the user" %}}
{{< code-func cheatsheet_test.go TestCustomEvent >}}
API: {{< api ecs World.Event >}}, {{< api ecs Event >}}
{{% /details %}}

{{% details closed="true" group="events" title="📣 Custom **events can be observed** like pre-defined events" %}}
{{< code-func cheatsheet_test.go TestCustomEventObserver >}}
API: {{< api ecs Observer >}}, {{< api ecs Observer1 >}}, {{< api ecs Observer2 >}}, ...
{{% /details %}}

<br />
<br />

{{< details-buttons paragraph="true" >}}
