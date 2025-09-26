+++
title = 'Event system'
type = "docs"
weight = 80
description = "Ark's event system and observers."
+++
Ark provides an an event system with observers that allow an application to react on events,
such as adding and removing components and entities.

Observers can filter for the events they are interested in, in several ways.
A callback function with is executed for the affected entity whenever an observer's filter matches.

## Example

{{< code-func events_test.go TestEventsBasic >}}

## Event types

Observers are specific for different event types, and each observer can react only to one type.
See [below](#combining-multiple-types) for how to react on multiple different types.

- **OnCreateEntity** &mdash; Fires after a new entity is created.  
- **OnRemoveEntity** &mdash; Fires before an entity is removed.
- **OnAddComponents** &mdash; Fires after components are added to an existing entity.
- **OnRemoveComponents** &mdash; Fires before components are removed from an entity.
- **OnSetComponents** &mdash; Fires after existing components are set from an entity.

If multiple components are added/removed/set for an entity,
one event is emitted for the entire operation.

## Combining multiple types

Observers can be combined to react to multiple event types in a single callback function.
Below is a combination of observers to react on component addition as well as removal.
The callback is set up to be able to distinguish between these event types (if needed).

{{< code-func events_test.go TestCombineObservers >}}

## Filters

**Component events** (`OnAddComponents`, `OnRemoveComponents`, `OnSetComponents`)
can be filtered by the affected components using {{< api ecs Observer.For >}}.
If this filter is not used, the observer triggers on any events of its type.
If a single component is specified, the observer triggers only on add/remove/set of this component.
If multiple components are specified, all must be added/removed/set at the same time for the observer to trigger.

**All events** can be filtered by the composition of the affected entity via
{{< api ecs Observer.With >}}, {{< api ecs Observer.Without >}} and {{< api ecs Observer.Exclusive >}}, just like [queries](../queries/).

## Event timing

The time an event is fired relative to the operation it is related to depends on the event's type.
The observer callbacks are executed immediately by any fired event.

Events for entity creation and for adding or setting components are fired after the operation.
Hence, the new or changed components can be inspected in the observer's callback.
In this case, the world is in an [unlocked](../queries#world-lock) state when the callback is executed.

Events for entity or component removal are fired before the operation.
This way, the entity or component to be removed can be inspected in the observer's callback.
In this case, the world is [locked](../queries#world-lock) when the callback is executed.

For [batch operations](../batch), all events are fired before or after the entire batch, respectively.
For batch creation or addition, events are fired after the potential batch callback
is executed for all entities, allowing to inspect the result.
