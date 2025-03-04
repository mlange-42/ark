+++
title = 'Quickstart'
type = 'docs'
weight = 10
description = 'Quickstart guide to install and use Ark.'
[params]
prev = "/"
+++
This page shows how to install Ark, and gives a minimal usage example.

Finally, it points into possible directions to continue.

## Installation

To use Ark in a Go project, run:

```bash
go get github.com/mlange-42/ark
```

## Usage example

Here is the classical Position/Velocity example that every ECS shows in the docs.

{{< code example_test.go >}}

## What's next?

If you ask **"What is ECS?"**, take a look at the great [**ECS FAQ**](https://github.com/SanderMertens/ecs-faq) by Sander Mertens, the author of the [Flecs](http://flecs.dev) ECS.

To learn how to use Ark, read the following sections,
browse the [API documentation](https://pkg.go.dev/github.com/mlange-42/ark),
or take a look at the {{< repo "tree/main/_examples" examples >}} in the {{< repo "" "GitHub repository" >}}.
