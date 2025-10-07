+++
title = 'Cheat sheet'
type = "docs"
weight = 1000
description = "Cheat sheet for frequent use cases."
+++
## Create entities

Create an **entity without components**:

{{< code-func cheatsheet_test.go TestCreateEmpty >}}

A **component mapper** is required for creating **entities with components**:

{{< code-func cheatsheet_test.go TestCreateMapper >}}

Create a **single entity**, given some components:

{{< code-func cheatsheet_test.go TestCreateEntity >}}

Create a **single entity** using a callback:

{{< code-func cheatsheet_test.go TestCreateEntityFn >}}

Create **many entities** more efficiently, all with the same component values:

{{< code-func cheatsheet_test.go TestCreateBatch >}}

Create **many entities**, using a callback for individual initialization:

{{< code-func cheatsheet_test.go TestCreateBatchFn >}}

## Remove entities

Remove a **single entity**:

{{< code-func cheatsheet_test.go TestRemoveEntity >}}

Remove **all entities** that match a filter:

{{< code-func cheatsheet_test.go TestRemoveEntities >}}

You can use a callback to do something with entities before their removal:

{{< code-func cheatsheet_test.go TestRemoveEntitiesFn >}}

## Add/remove components

A **component mapper** is required for adding and removing components.
It adds or removes the given components from entities:

{{< code-func cheatsheet_test.go TestCreateMapper >}}

Add and remove components to/from a **single entity**:

{{< code-func cheatsheet_test.go TestAddRemoveComponents >}}

Add components to **all entities** matching a filter:

{{< code-func cheatsheet_test.go TestAddBatch >}}

Add components to **all entities** matching a filter, with individual initialization:

{{< code-func cheatsheet_test.go TestAddBatchFn >}}

Remove components from **all entities** matching a filter.
The callback can be used to do something with entities before component removal:

{{< code-func cheatsheet_test.go TestRemoveBatch >}}
