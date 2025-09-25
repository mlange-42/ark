package ecs

import (
	"math"
)

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
}

// Observe creates a new ECS event observer for the specified event type.
//
// Observers react to structural changes, such as entity creation, removal, and component addition/removal.
// Use the methods With, Without, and Do to configure the observer before registering it.
func Observe(evt EventType) *Observer {
	return &Observer{
		event: evt,
		id:    maxObserverID,
	}
}

// For adds components that the observer observes.
// Can only be used with OnAddComponents and OnRemoveComponents.
func (o *Observer) For(comps ...Comp) *Observer {
	if o.id != maxObserverID {
		panic("can't modify a registered observer")
	}
	if o.event != OnAddComponents && o.event != OnRemoveComponents && o.event != OnSetComponents {
		panic("can use Observer.For only for OnAddComponents, OnRemoveComponents and OnSetComponents events")
	}
	if len(comps) == 0 {
		return o
	}
	o.comps = append(o.comps, comps...)
	o.hasComps = true
	return o
}

// With adds components that entities must have to trigger.
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

// Without adds components the observer excludes.
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
	compsOr      []bitMask
	withOr       []bitMask
	anyNoComps   []bool
	anyNoWith    []bool
	pool         intPool[observerID]
	indices      map[observerID]int
}

func newObserverManager() observerManager {
	return observerManager{
		observers:    make([][]*Observer, eventsEnd),
		hasObservers: make([]bool, eventsEnd),
		anyNoComps:   make([]bool, eventsEnd),
		anyNoWith:    make([]bool, eventsEnd),
		compsOr:      make([]bitMask, eventsEnd),
		withOr:       make([]bitMask, eventsEnd),
		pool:         newIntPool[observerID](32),
		indices:      map[observerID]int{},
	}
}

func (m *observerManager) AddObserver(o *Observer, w *World) {
	o.id = m.pool.Get()

	for _, c := range o.comps {
		id := TypeID(w, c.tp)
		o.compsMask.Set(id.id, true)
	}
	for _, c := range o.with {
		id := TypeID(w, c.tp)
		o.withMask.Set(id.id, true)
	}
	for _, c := range o.without {
		id := TypeID(w, c.tp)
		o.withoutMask.Set(id.id, true)
	}

	m.observers[o.event] = append(m.observers[o.event], o)
	m.indices[o.id] = len(m.observers[o.event]) - 1
	m.hasObservers[o.event] = true

	if o.hasWith {
		m.withOr[o.event].OrI(&o.withMask)
	} else {
		m.anyNoWith[o.event] = true
	}

	if o.event == OnCreateEntity || o.event == OnRemoveEntity {
		return
	}

	if o.hasComps {
		m.compsOr[o.event].OrI(&o.compsMask)
	} else {
		m.anyNoComps[o.event] = true
	}
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

	var withOr bitMask
	m.anyNoWith[o.event] = false
	for _, o := range m.observers[o.event] {
		if !o.hasWith {
			m.anyNoWith[o.event] = true
			break
		}
		withOr.OrI(&o.withMask)
	}
	m.withOr[o.event] = withOr

	if o.event == OnCreateEntity || o.event == OnRemoveEntity {
		return
	}

	var compsOr bitMask
	m.anyNoComps[o.event] = false
	for _, o := range m.observers[o.event] {
		if !o.hasComps {
			m.anyNoComps[o.event] = true
			break
		}
		compsOr.OrI(&o.compsMask)
	}
	m.compsOr[o.event] = compsOr
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
	if !m.anyNoWith[OnCreateEntity] && !m.withOr[OnCreateEntity].ContainsAny(mask) {
		return
	}
	observers := m.observers[OnCreateEntity]
	for _, o := range observers {
		if o.hasWith && !mask.Contains(&o.withMask) {
			continue
		}
		if o.hasWithout && mask.ContainsAny(&o.withoutMask) {
			continue
		}
		o.callback(e)
	}
}

func (m *observerManager) doFireRemoveEntity(e Entity, mask *bitMask) {
	if !m.anyNoWith[OnRemoveEntity] && !m.withOr[OnRemoveEntity].ContainsAny(mask) {
		return
	}
	observers := m.observers[OnRemoveEntity]
	for _, o := range observers {
		if o.hasWith && !mask.Contains(&o.withMask) {
			continue
		}
		if o.hasWithout && mask.ContainsAny(&o.withoutMask) {
			continue
		}
		o.callback(e)
	}
}

func (m *observerManager) FireAdd(e Entity, oldMask *bitMask, newMask *bitMask) {
	if !m.hasObservers[OnAddComponents] {
		return
	}
	m.doFireAdd(e, oldMask, newMask)
}

func (m *observerManager) doFireAdd(e Entity, oldMask *bitMask, newMask *bitMask) {
	if !m.anyNoComps[OnAddComponents] &&
		(!m.compsOr[OnAddComponents].ContainsAny(newMask) || oldMask.Contains(&m.compsOr[OnAddComponents])) {
		return
	}
	if !m.anyNoWith[OnAddComponents] && !m.withOr[OnAddComponents].ContainsAny(newMask) {
		return
	}
	observers := m.observers[OnAddComponents]
	for _, o := range observers {
		if o.hasComps && (!newMask.Contains(&o.compsMask) || oldMask.ContainsAny(&o.compsMask)) {
			continue
		}
		if o.hasWith && !newMask.Contains(&o.withMask) {
			continue
		}
		if o.hasWithout && newMask.ContainsAny(&o.withoutMask) {
			continue
		}
		o.callback(e)
	}
}

func (m *observerManager) doFireRemove(e Entity, oldMask *bitMask, newMask *bitMask) {
	if !m.anyNoComps[OnRemoveComponents] &&
		(!m.compsOr[OnRemoveComponents].ContainsAny(oldMask) || newMask.Contains(&m.compsOr[OnRemoveComponents])) {
		return
	}
	if !m.anyNoWith[OnRemoveComponents] && !m.withOr[OnRemoveComponents].ContainsAny(newMask) {
		return
	}
	observers := m.observers[OnRemoveComponents]
	for _, o := range observers {
		if o.hasComps && (newMask.Contains(&o.compsMask) || !oldMask.ContainsAny(&o.compsMask)) {
			continue
		}
		if o.hasWith && !newMask.Contains(&o.withMask) {
			continue
		}
		if o.hasWithout && newMask.ContainsAny(&o.withoutMask) {
			continue
		}
		o.callback(e)
	}
}

func (m *observerManager) doFireSet(e Entity, mask *bitMask, newMask *bitMask) {
	if !m.anyNoComps[OnSetComponents] && !m.compsOr[OnSetComponents].ContainsAny(mask) {
		return
	}
	if !m.anyNoWith[OnSetComponents] && !m.withOr[OnSetComponents].ContainsAny(newMask) {
		return
	}
	observers := m.observers[OnSetComponents]
	for _, o := range observers {
		if o.hasComps && !mask.Contains(&o.compsMask) {
			continue
		}
		if o.hasWith && !newMask.Contains(&o.withMask) {
			continue
		}
		if o.hasWithout && newMask.ContainsAny(&o.withoutMask) {
			continue
		}
		o.callback(e)
	}
}
