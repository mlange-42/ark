+++
title = 'Event system'
type = "docs"
weight = 80
description = "Ark's event system and observers."
+++
Ark provides an event system with observers that allow an application to react on events,
such as adding and removing components and entities.

Observers can [filter](#filters) for the events they are interested in, in several ways.
A callback function is executed for the affected entity whenever an observer's filter matches.

In addition to built-in lifecycle events like `OnCreateEntity` or `OnAddComponents`,
Ark supports [custom event types](#custom-events) that enable domain-specific triggers.
These events can be emitted manually and observed with the same filtering and callback mechanisms,
making them ideal for modeling interactions such as user input, synchronization, or game logic.

Observers are lightweight, composable, and follow the same declarative patterns as Ark’s [query](../queries/) system.
They provide fine-grained control over when and how logic is executed.
This design encourages a declarative, data-driven approach while maintaining performance and flexibility.

## Example

{{< code-func events_test.go TestEventsBasic >}}

## Event types

Observers are specific for different event types, and each observer can react only to one type.
See [below](#combining-multiple-types) for how to react on multiple different types.

- **OnCreateEntity** &mdash; Emitted after a new entity is created.  
- **OnRemoveEntity** &mdash; Emitted before an entity is removed.
- **OnAddComponents** &mdash; Emitted after components are added to an existing entity.
- **OnRemoveComponents** &mdash; Emitted before components are removed from an entity.
- **OnSetComponents** &mdash; Emitted after existing components are set from an entity.
- **OnAddRelations** &mdash; Emitted after relation targets are added to an entity.*
- **OnRemoveRelations** &mdash; Emitted before relation targets are removed from an entity.*

If multiple components are added/removed/set for an entity,
one event is emitted for the entire operation.

\* *Relation events are emitted when entities with relations are created or removed, when relation components are added or removed, as well as when targets are set without changing components.*

## Combining multiple types

Observers can be combined to react to multiple event types in a single callback function.
Below is a combination of observers to react on component addition as well as removal.
The callback is set up to be able to distinguish between these event types (if needed).

{{< code-func events_test.go TestCombineObservers >}}

## Filters

Observers filter for the components specified by their generic parameters.
Additional components can be specified using {{< api ecs Observer.For >}},
but these are not directly accessible in the callback.

Observers only trigger when all specified components (in parameters and in `For`)
are affected in a single operation.
For example, if an observer watches `Position` and `Velocity`,
both must be added or removed together for the observer to activate

Further, events can be filtered by the composition of the affected entity via
{{< api ecs Observer.With >}}, {{< api ecs Observer.Without >}} and {{< api ecs Observer.Exclusive >}}, just like [queries](../queries/).

**Examples** (leaving out observer registration):

Both observers are triggered when an entity with `Position` is created.
The first one has direct access to the component in the callback while the second does not:

{{< code-func events_test.go TestObserveCreate >}}

Both observers are triggered when an entity with `Position` as well as `Velocity` is created:

{{< code-func events_test.go TestObserve2Create >}}

An observer that is triggered when any entity is created, irrespective of its components:

{{< code-func events_test.go TestObserveCreateEmpty >}}

An observer that is triggered when a `Position` component is added to an existing entity:

{{< code-func events_test.go TestObserveAdd >}}

An observer that is triggered when a `Position` component is added to an entity
that has `Velocity`, but not `Altitude` (or rather, had before the operation):

{{< code-func events_test.go TestObserveAddWith >}}

## Event timing

The time an event is emitted relative to the operation it is related to depends on the event's type.
The observer callbacks are executed immediately by any emitted event.

Events for entity creation and for adding or setting components are emitted after the operation.
Hence, the new or changed components can be inspected in the observer's callback.
If emitted from individual operations, the world is in an [unlocked](../queries#world-lock) state when the callback is executed. Contrary, when emitted from a batch operation, the world is [locked](../queries#world-lock).

Events for entity or component removal are emitted before the operation.
This way, the entity or component to be removed can be inspected in the observer's callback.
In this case, the world is [locked](../queries#world-lock) when the callback is executed.

For [batch operations](../batch), all events are emitted before or after the entire batch, respectively.
For batch creation or addition, events are emitted after the potential batch callback
is executed for all entities, allowing to inspect the result.

Note that observer order is undefined. Observers are not necessarily triggered
in the same order as they were registered.

## Custom events

Custom events in Ark allow developers to define and emit their own event types,
enabling application-specific logic such as UI interactions, game state changes,
or other domain-specific triggers.
These events support the same filtering and observer mechanisms as built-in events.

Define custom event types using {{< api ecs EventRegistry.NewEventType >}}:

{{< code-func events_test.go TestNewEventType 0 8 >}}

Ideally, custom event types are stored as global variables of the applications.

Alteratively, if all custom events are defined in one place, constants can be used like this:

{{< code-func events_test.go TestNewEventTypeIota 0 6 >}}

Use custom events like this:

{{< code-func events_test.go TestEventEmit >}}

Observers might not be interested in components, or in more than one component.
This is also supported by custom events:

{{< code-func events_test.go TestEventClick >}}

Here, the event is created and emitted in a single expression.
However, it is recommended to store events after construction and to reuse them for `Emit`.
Reusing event instances is especially important for events with components,
as it avoids repeated lookups and improves runtime efficiency.
The overhead for component ID lookup is &approx;20ns per component.

For custom events, observer [filters](#filters) work exactly the same as for predefined events.
The components in the generic parameters of the observer, as well as those defined by `For`,
are matched against the components of the event.
`With`, `Without` and `Exclusive` are matched against the entity for which the event is emitted.

Note that custom events can also be emitted for the zero entity:

{{< code-func events_test.go TestEventZeroEntity >}}
