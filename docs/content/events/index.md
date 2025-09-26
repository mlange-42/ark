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

## Event types

### OnCreateEntity

### OnRemoveEntity

### OnAddComponents

### OnRemoveComponents
