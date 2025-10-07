+++
title = 'Cheat sheet'
type = "docs"
weight = 1000
description = "Cheat sheet for frequent use cases."
+++
## Create entities

Create an **entity without components**:

{{< code-func cheatsheet_test.go TestCreateEmpty >}}

For creating **entities with components**, you need a **component mapper**:

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
