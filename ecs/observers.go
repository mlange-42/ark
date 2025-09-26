package ecs

// Observer2 is a generic observer for two components.
type Observer2[A any, B any] struct {
	observer Observer
	callback func(Entity, *A, *B)
}

// Observe2 creates a new Observer2.
func Observe2[A any, B any](evt EventType) *Observer2[A, B] {
	return &Observer2[A, B]{
		observer: Observer{
			event: evt,
			id:    maxObserverID,
		},
	}
}

// With adds components that entities must have to trigger the observer.
// If multiple components are provided, the entity must have all of them.
//
// Method calls can be chained, which has the same effect as calling with multiple arguments.
func (o *Observer2[A, B]) With(comps ...Comp) *Observer2[A, B] {
	if o.observer.id != maxObserverID {
		panic("can't modify a registered observer")
	}
	if len(comps) == 0 {
		return o
	}
	o.observer.with = append(o.observer.with, comps...)
	o.observer.hasWith = true
	return o
}

// Without adds components that entities must not have to trigger the observer.
// If multiple components are provided, the entity must not have any of them.
//
// Method calls can be chained, which has the same effect as calling with multiple arguments.
func (o *Observer2[A, B]) Without(comps ...Comp) *Observer2[A, B] {
	if o.observer.id != maxObserverID {
		panic("can't modify a registered observer")
	}
	if len(comps) == 0 {
		return o
	}
	o.observer.without = append(o.observer.without, comps...)
	o.observer.hasWithout = true
	return o
}

// Exclusive makes the observer exclusive in the sense that the components given by [Observer.With]
// are matched exactly, and no other components are allowed.
//
// Overwrites components set via [Observer.Without].
func (o *Observer2[A, B]) Exclusive() *Observer2[A, B] {
	if o.observer.id != maxObserverID {
		panic("can't modify a registered observer")
	}
	o.observer.exclusive = true
	o.observer.hasWithout = true
	return o
}

// Do sets the observer's callback. Must be called exactly once before registration.
func (o *Observer2[A, B]) Do(fn func(Entity, *A, *B)) *Observer2[A, B] {
	if o.callback != nil {
		panic("observer already has a callback")
	}
	o.callback = fn
	return o
}

// Register this observer. This is mandatory for the observer to take effect.
func (o *Observer2[A, B]) Register(w *World) *Observer2[A, B] {
	if o.callback == nil {
		panic("observer callback must be set via Do before registering")
	}
	if o.observer.id != maxObserverID {
		panic("observer is already registered")
	}

	idA := ComponentID[A](w)
	idB := ComponentID[A](w)
	storageA := &w.storage.components[idA.id]
	storageB := &w.storage.components[idB.id]

	o.observer.callback = func(e Entity) {
		index := &w.storage.entities[e.id]
		row := uintptr(index.row)
		o.callback(
			e,
			(*A)(storageA.columns[index.table].Get(row)),
			(*B)(storageB.columns[index.table].Get(row)),
		)
	}

	// TODO: optimize to use IDs.
	o.observer.For(
		C[A](),
		C[B](),
	)
	w.registerObserver(&o.observer)
	return o
}

// Unregister this observer.
func (o *Observer2[A, B]) Unregister(w *World) *Observer2[A, B] {
	if o.observer.id == maxObserverID {
		panic("observer is not registered")
	}
	w.unregisterObserver(&o.observer)
	return o
}
