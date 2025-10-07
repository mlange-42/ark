+++
title = 'Cheat sheet'
type = "docs"
weight = 1000
description = "Cheat sheet for frequent use cases."
+++
Frequently used Ark operations for quick lookup.

## World creation

The world is the central ECS data storage.
Most applications will use exactly one world.

{{% details closed="true" title="🌍 Create a World with default initial capacity" %}}
{{< code-func cheatsheet_test.go TestCreateWorld >}}
{{% /details %}}

{{% details closed="true" title="🌍 Create a World with a specific initial capacity" %}}
{{< code-func cheatsheet_test.go TestCreateWorldConfig >}}
{{% /details %}}

## Create entities

{{% details closed="true" title="✨ Create an **entity without components**" %}}
Create an **entity without components**:
{{< code-func cheatsheet_test.go TestCreateEmpty >}}
{{% /details %}}

{{% details closed="true" title="✨ A **component mapper** is required for creating **entities with components**" %}}
{{< code-func cheatsheet_test.go TestCreateMapper >}}
{{% /details %}}

{{% details closed="true" title="✨ Create a **single entity**, given some components" %}}
{{< code-func cheatsheet_test.go TestCreateEntity >}}
{{% /details %}}

{{% details closed="true" title="✨ Create a **single entity** using a callback" %}}
{{< code-func cheatsheet_test.go TestCreateEntityFn >}}
{{% /details %}}

{{% details closed="true" title="✨ Create **many entities** more efficiently, all with the same component values" %}}
{{< code-func cheatsheet_test.go TestCreateBatch >}}
{{% /details %}}

{{% details closed="true" title="✨ Create **many entities**, using a callback for individual initialization" %}}
{{< code-func cheatsheet_test.go TestCreateBatchFn >}}
{{% /details %}}

## Remove entities

{{% details closed="true" title="❌ Remove a **single entity**" %}}
{{< code-func cheatsheet_test.go TestRemoveEntity >}}
{{% /details %}}


{{% details closed="true" title="❌ Remove **all entities** that match a filter" %}}
{{< code-func cheatsheet_test.go TestRemoveEntities >}}
{{% /details %}}

{{% details closed="true" title="❌ With a **callback** to do something with entities before their removal" %}}
{{< code-func cheatsheet_test.go TestRemoveEntitiesFn >}}
{{% /details %}}

## Add/remove components

{{% details closed="true" title="🧩 A **component mapper** is required for adding and removing components" %}}
It adds or removes the given components to/from entities.
{{< code-func cheatsheet_test.go TestCreateMapper >}}
{{% /details %}}

{{% details closed="true" title="🧩 Add and remove components to/from a **single entity**" %}}
{{< code-func cheatsheet_test.go TestAddRemoveComponents >}}
{{% /details %}}

{{% details closed="true" title="🧩 Add components to **all entities** matching a filter" %}}
{{< code-func cheatsheet_test.go TestAddBatch >}}
{{% /details %}}

{{% details closed="true" title="🧩 Add components to **all entities** matching a filter, with individual initialization" %}}
{{< code-func cheatsheet_test.go TestAddBatchFn >}}
{{% /details %}}

{{% details closed="true" title="🧩 Remove components from **all entities** matching a filter." %}}
The callback can be used to do something with entities before component removal.
{{< code-func cheatsheet_test.go TestRemoveBatch >}}
{{% /details %}}

## Resources

Resources can be used for "global", singleton-like data structures that are not associated to particular entities.

{{% details closed="true" title="📦 Adding and getting resources, the simple but slower way (&approx;20ns)" %}}
{{< code-func cheatsheet_test.go TestResourcesQuick >}}
{{% /details %}}

{{% details closed="true" title="📦 For repeated access, better use a resource accessor (`Get()` &approx;1ns)" %}}
{{< code-func cheatsheet_test.go TestResources >}}
(Creating the accessor does not add the actual `Grid` resource!)
{{% /details %}}

## Events and observers

Create and register observers for ECS lifecycle events:
