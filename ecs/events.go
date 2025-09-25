package ecs

import "math"

type observerID uint16

const maxObserverID = math.MaxUint16

// EventType is the type for event identifiers.
// See [Observer] for details on events.
type EventType uint8

const (
	// OnCreateEntity event.
	OnCreateEntity EventType = iota
	// OnRemoveEntity event.
	OnRemoveEntity
	// OnAddComponents event.
	OnAddComponents
	// OnRemoveComponents event.
	OnRemoveComponents
	// OnSetComponents event.
	OnSetComponents
	// OnChangeTarget event.
	//OnChangeTarget

	eventsEnd
)

// Observer for ECS events.
//
// Observers react to structural changes, such as entity creation, removal, and component addition/removal.
// Use the methods NewObserver, With, Without, and Do to configure the observer before registering it.
type Observer struct {
	event       EventType
	with        []Comp
	without     []Comp
	withMask    bitMask
	withoutMask bitMask
	hasWithout  bool
	callback    func(Entity)
	id          observerID
}

// NewObserver creates a new ECS event observer for the specified event type.
//
// Observers react to structural changes, such as entity creation, removal, and component addition/removal.
// Use the methods With, Without, and Do to configure the observer before registering it.
func NewObserver(evt EventType) *Observer {
	return &Observer{
		event: evt,
		id:    maxObserverID,
	}
}

// With adds components the observer observes.
func (o *Observer) With(comps ...Comp) *Observer {
	if o.id != maxObserverID {
		panic("can't modify a registered observer")
	}
	o.with = append(o.with, comps...)
	return o
}

// Without adds components the observer excludes.
// Only valid for [OnCreateEntity] and [OnRemoveEntity].
func (o *Observer) Without(comps ...Comp) *Observer {
	if o.id != maxObserverID {
		panic("can't modify a registered observer")
	}
	if o.event != OnCreateEntity && o.event != OnRemoveEntity {
		panic("can use Observer.Without only for OnCreateEntity and OnRemoveEntity events")
	}
	if len(comps) == 0 {
		return o
	}
	o.without = append(o.without, comps...)
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

// Register this observer.
func (o *Observer) Register(w *World) *Observer {
	if o.callback == nil {
		panic("observer callback must be set via Do before registering")
	}
	if o.id != maxObserverID {
		panic("observer is already registered")
	}
	w.registerObserver(o)
	return o
}

// Unregister this observer.
func (o *Observer) Unregister(w *World) *Observer {
	if o.id == maxObserverID {
		panic("observer is not registered")
	}
	w.unregisterObserver(o)
	return o
}

type observerManager struct {
	observers    [][]*Observer
	hasObservers []bool
	pool         intPool[observerID]
	indices      map[observerID]int
}

func newObserverManager() observerManager {
	return observerManager{
		observers:    make([][]*Observer, eventsEnd),
		hasObservers: make([]bool, eventsEnd),
		pool:         newIntPool[observerID](32),
		indices:      map[observerID]int{},
	}
}

func (m *observerManager) AddObserver(o *Observer, reg *componentRegistry) {
	o.id = m.pool.Get()

	for _, c := range o.with {
		id, _ := reg.ComponentID(c.tp)
		o.withMask.Set(id, true)
	}
	for _, c := range o.without {
		id, _ := reg.ComponentID(c.tp)
		o.withoutMask.Set(id, true)
	}

	m.observers[o.event] = append(m.observers[o.event], o)
	m.indices[o.id] = len(m.observers[o.event]) - 1
	m.hasObservers[o.event] = true
}

func (m *observerManager) RemoveObserver(o *Observer) {
	idx, ok := m.indices[o.id]
	if !ok {
		panic("can't unregister observer, not found")
	}
	delete(m.indices, o.id)

	observers := m.observers[o.event]
	observers[idx].id = maxObserverID

	last := len(observers) - 1
	if idx != last {
		observers[idx], observers[last] = observers[last], observers[idx]
		m.indices[observers[idx].id] = idx
	}
	observers[last] = nil
	m.observers[o.event] = observers[:last]
	m.hasObservers[o.event] = last > 0
}

func (m *observerManager) HasObservers(evt EventType) bool {
	return m.hasObservers[evt]
}

func (m *observerManager) FireCreateEntity(e Entity, mask *bitMask) {
	if !m.hasObservers[OnCreateEntity] {
		return
	}
	m.doFireCreateEntity(e, mask)
}

func (m *observerManager) doFireCreateEntity(e Entity, mask *bitMask) {
	observers := m.observers[OnCreateEntity]
	for _, o := range observers {
		if mask.Contains(&o.withMask) && (!o.hasWithout || !mask.ContainsAny(&o.withoutMask)) {
			o.callback(e)
		}
	}
}

func (m *observerManager) doFireRemoveEntity(e Entity, mask *bitMask) {
	observers := m.observers[OnRemoveEntity]
	for _, o := range observers {
		if mask.Contains(&o.withMask) && (!o.hasWithout || !mask.ContainsAny(&o.withoutMask)) {
			o.callback(e)
		}
	}
}

func (m *observerManager) FireAdd(e Entity, oldMask *bitMask, newMask *bitMask) {
	if !m.hasObservers[OnAddComponents] {
		return
	}
	m.doFireAdd(e, oldMask, newMask)
}

func (m *observerManager) doFireAdd(e Entity, oldMask *bitMask, newMask *bitMask) {
	observers := m.observers[OnAddComponents]
	for _, o := range observers {
		if newMask.Contains(&o.withMask) && !oldMask.ContainsAny(&o.withMask) {
			o.callback(e)
		}
	}
}

func (m *observerManager) FireRemove(e Entity, oldMask *bitMask, newMask *bitMask) {
	if !m.hasObservers[OnRemoveComponents] {
		return
	}
	m.doFireRemove(e, oldMask, newMask)
}

func (m *observerManager) doFireRemove(e Entity, oldMask *bitMask, newMask *bitMask) {
	observers := m.observers[OnRemoveComponents]
	for _, o := range observers {
		if oldMask.Contains(&o.withMask) && !newMask.ContainsAny(&o.withMask) {
			o.callback(e)
		}
	}
}

func (m *observerManager) FireSet(e Entity, mask *bitMask) {
	if !m.hasObservers[OnSetComponents] {
		return
	}
	m.doFireSet(e, mask)
}

func (m *observerManager) doFireSet(e Entity, mask *bitMask) {
	observers := m.observers[OnSetComponents]
	for i := range observers {
		o := observers[i]
		if mask.Contains(&o.withMask) {
			o.callback(e)
		}
	}
}
