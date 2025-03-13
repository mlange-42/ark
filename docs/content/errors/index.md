+++
title = 'Error handling'
type = "docs"
weight = 90
description = "Ark's way to handle errors."
+++

Ark puts an emphasis on safety, and on avoiding undefined and unexpected behavior.
It panics on unexpected operations, like removing a dead entity,
adding a component that is already present, or attempting to change a locked world.

This may not seem idiomatic for Go.
However, explicit error handling in performance hot spots is not an option.
Neither is silent failure or ignoring invalid operations, given the scientific background of Ark.
Therefore, Ark panics.

## Debug build

Ark tries to give informative error messages on invalid operations or other misuse.
In performance hot spots like [queries](../queries) or [component mappers](../operations#component-mappers),
however, this is not always possible without degrading performance.
As an example, {{< api ecs Query2.Get >}} panics (deliberately) with `invalid memory address or nil pointer dereference`
when called after query iteration finished.

In case of uninformative errors in queries or mappers, try to run your project using the build tag `debug`:

```
go run -tags debug .
```

This enables additional checks for more helpful error messages, at the cost of a performance penalty.
Note that this does not change when Ark panics, it only improves the error messages.

If you still get uninformative error messages from inside Ark, please [create an issue!](https://github.com/mlange-42/ark/issues/new)
Either there are missing debug checks, or there is a bug in Ark.
