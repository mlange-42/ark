package ecs

// Observer for ECS events.
//
// Observers react to structural changes, such as entity creation, removal, and component addition/removal.
// Use the methods Observe, With, Without, and Do to configure the observer before registering it.
//
// See [EventType] for available events.
// See also [Observer1], [Observer2], etc.
type Observer struct {
	compsMask   bitMask
	withMask    bitMask
	withoutMask bitMask
	callback    func(Entity)
	comps       []Comp
	with        []Comp
	without     []Comp
	id          observerID
	event       EventType
	hasComps    bool
	hasWithout  bool
	hasWith     bool
	exclusive   bool
}

// Observe creates a new ECS event observer for the specified event type.
//
// Observers react to structural changes, such as entity creation, removal, and component addition/removal.
// Use the methods For, With, Without, Exclusive, and Do to configure the observer before registering it.
//
// See also [Observe1], [Observe2], etc.
func Observe(evt EventType) *Observer {
	return &Observer{
		event: evt,
		id:    maxObserverID,
	}
}

// For adds components that the observer observes.
// The component events, the observer triggers if these components are added to or removed from an entity.
//
// If not specified, the observer triggers on any component addition or removal.
// If multiple components are provided, all must be added/removed at the same time to trigger the observer.
//
// For entity events (OnCreateEntity, OnRemoveEntity), is has the same effect as With.
//
// Method calls can be chained, which has the same effect as calling with multiple arguments.
func (o *Observer) For(comps ...Comp) *Observer {
	if o.id != maxObserverID {
		panic("can't modify a registered observer")
	}
	if len(comps) == 0 {
		return o
	}
	if o.event == OnCreateEntity || o.event == OnRemoveEntity {
		return o.With(comps...)
	}
	o.comps = append(o.comps, comps...)
	o.hasComps = true
	return o
}

// With adds components that entities must have to trigger the observer.
// If multiple components are provided, the entity must have all of them.
//
// Method calls can be chained, which has the same effect as calling with multiple arguments.
func (o *Observer) With(comps ...Comp) *Observer {
	if o.id != maxObserverID {
		panic("can't modify a registered observer")
	}
	if len(comps) == 0 {
		return o
	}
	o.with = append(o.with, comps...)
	o.hasWith = true
	return o
}

// Without adds components that entities must not have to trigger the observer.
// If multiple components are provided, the entity must not have any of them.
//
// Method calls can be chained, which has the same effect as calling with multiple arguments.
func (o *Observer) Without(comps ...Comp) *Observer {
	if o.id != maxObserverID {
		panic("can't modify a registered observer")
	}
	if len(comps) == 0 {
		return o
	}
	o.without = append(o.without, comps...)
	o.hasWithout = true
	return o
}

// Exclusive makes the observer exclusive in the sense that the components given by [Observer.With]
// are matched exactly, and no other components are allowed.
//
// Overwrites components set via [Observer.Without].
func (o *Observer) Exclusive() *Observer {
	if o.id != maxObserverID {
		panic("can't modify a registered observer")
	}
	o.exclusive = true
	o.hasWithout = true
	return o
}

// Do sets the observer's callback. Must be called exactly once before registration.
func (o *Observer) Do(fn func(Entity)) *Observer {
	if o.callback != nil {
		panic("observer already has a callback")
	}
	o.callback = fn
	return o
}

// Register this observer. This is mandatory for the observer to take effect.
func (o *Observer) Register(w *World) *Observer {
	w.registerObserver(o)
	return o
}

// Unregister this observer.
func (o *Observer) Unregister(w *World) *Observer {
	w.unregisterObserver(o)
	return o
}
