package ecs

import (
	"math"
)

type observerID uint32

const maxObserverID = math.MaxUint32

// EventType is the type for event identifiers.
// See below for predefined event types.
//
// Use [NewEventType] to create custom events types.
// See [Event] and [World.Event] for using custom events.
//
// See [Observer] for details on events and observers.
type EventType uint8

// Predefined event types.
const (

	// OnCreateEntity event.
	// Emitted after an entity is created.
	OnCreateEntity EventType = iota

	// OnRemoveEntity event.
	// Emitted before an entity is removed.
	OnRemoveEntity

	// OnAddComponents event.
	// Emitted after components are added to an entity.
	OnAddComponents

	// OnRemoveComponents event.
	// Emitted before components are removed from an entity.
	OnRemoveComponents

	// OnSetComponents event.
	// Emitted after components are set from an entity.
	OnSetComponents

	// OnAddRelations event.
	// Emitted after relation targets are added to an entity.
	OnAddRelations

	// OnRemoveRelations event.
	// Emitted before relation targets are removed from an entity.
	OnRemoveRelations

	// Marker for number of event types.
	eventsEnd
)

var nextUserEvent = eventsEnd - 1

// NewEventType creates a new EventType for custom events.
// Custom event types should be stored in global variables.
//
// The maximum number of event types is 255, with 5 predefined and 250 potential custom types.
//
// See [Event] and [World.Event] for using custom events.
func NewEventType() EventType {
	if nextUserEvent == math.MaxUint8 {
		panic("reached maximum number of custom event types")
	}
	nextUserEvent++
	return EventType(nextUserEvent)
}

// Event is a custom event.
//
// Create events using [World.Event].
type Event struct {
	world     *World
	eventType EventType
	mask      bitMask
}

// For sets the event's component types. Optional.
//
// For best performance, store the event after setting component types,
// and re-use afterwards.
func (e Event) For(comps ...Comp) Event {
	for i := range comps {
		id := TypeID(e.world, comps[i].tp)
		e.mask.Set(id.id, true)
	}
	return e
}

// Emit the event for the given entity.
func (e Event) Emit(entity Entity) {
	e.world.emitEvent(&e, entity)
}

type observerManager struct {
	observers    [][]*Observer
	hasObservers []bool
	allComps     []bitMask
	allWith      []bitMask
	anyNoComps   []bool
	anyNoWith    []bool
	pool         intPool[observerID]
	indices      map[observerID]int
}

func newObserverManager() observerManager {
	return observerManager{
		observers:    make([][]*Observer, math.MaxUint8),
		hasObservers: make([]bool, math.MaxUint8),
		anyNoComps:   make([]bool, math.MaxUint8),
		anyNoWith:    make([]bool, math.MaxUint8),
		allComps:     make([]bitMask, math.MaxUint8),
		allWith:      make([]bitMask, math.MaxUint8),
		pool:         newIntPool[observerID](32),
		indices:      map[observerID]int{},
	}
}

func (m *observerManager) AddObserver(o *Observer, w *World) {
	if o.callback == nil {
		panic("observer callback must be set via Do before registering")
	}
	if o.id != maxObserverID {
		panic("observer is already registered")
	}

	o.id = m.pool.Get()

	for _, c := range o.comps {
		id := TypeID(w, c.tp)
		o.compsMask.Set(id.id, true)
	}
	for _, c := range o.with {
		id := TypeID(w, c.tp)
		o.withMask.Set(id.id, true)
	}
	if o.exclusive {
		o.withoutMask = o.withMask.Not()
	} else {
		for _, c := range o.without {
			id := TypeID(w, c.tp)
			o.withoutMask.Set(id.id, true)
		}
	}

	m.observers[o.event] = append(m.observers[o.event], o)
	m.indices[o.id] = len(m.observers[o.event]) - 1
	m.hasObservers[o.event] = true

	if o.hasWith {
		m.allWith[o.event].OrI(&o.withMask)
	} else {
		m.anyNoWith[o.event] = true
	}

	if o.event == OnCreateEntity || o.event == OnRemoveEntity {
		return
	}

	if o.hasComps {
		m.allComps[o.event].OrI(&o.compsMask)
	} else {
		m.anyNoComps[o.event] = true
	}
}

func (m *observerManager) RemoveObserver(o *Observer) {
	if o.id == maxObserverID {
		panic("observer is not registered")
	}

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

	var allWith bitMask
	m.anyNoWith[o.event] = false
	for _, o := range m.observers[o.event] {
		if !o.hasWith {
			m.anyNoWith[o.event] = true
			break
		}
		allWith.OrI(&o.withMask)
	}
	m.allWith[o.event] = allWith

	if o.event == OnCreateEntity || o.event == OnRemoveEntity {
		return
	}

	var allComps bitMask
	m.anyNoComps[o.event] = false
	for _, o := range m.observers[o.event] {
		if !o.hasComps {
			m.anyNoComps[o.event] = true
			break
		}
		allComps.OrI(&o.compsMask)
	}
	m.allComps[o.event] = allComps
}

func (m *observerManager) HasObservers(evt EventType) bool {
	return m.hasObservers[evt]
}

func (m *observerManager) FireCreateEntity(e Entity, mask *bitMask) {
	if !m.hasObservers[OnCreateEntity] {
		return
	}
	m.doFireCreateEntity(e, mask, true)
}

func (m *observerManager) doFireCreateEntity(e Entity, mask *bitMask, earlyOut bool) bool {
	if earlyOut && !m.anyNoWith[OnCreateEntity] && !m.allWith[OnCreateEntity].ContainsAny(mask) {
		return false
	}
	observers := m.observers[OnCreateEntity]
	found := false
	for _, o := range observers {
		if o.hasWith && !mask.Contains(&o.withMask) {
			continue
		}
		if o.hasWithout && mask.ContainsAny(&o.withoutMask) {
			continue
		}
		o.callback(e)
		found = true
	}
	return found
}

func (m *observerManager) doFireRemoveEntity(e Entity, mask *bitMask, earlyOut bool) bool {
	if earlyOut && !m.anyNoWith[OnRemoveEntity] && !m.allWith[OnRemoveEntity].ContainsAny(mask) {
		return false
	}
	observers := m.observers[OnRemoveEntity]
	found := false
	for _, o := range observers {
		if o.hasWith && !mask.Contains(&o.withMask) {
			continue
		}
		if o.hasWithout && mask.ContainsAny(&o.withoutMask) {
			continue
		}
		o.callback(e)
		found = true
	}
	return found
}

func (m *observerManager) FireAdd(e Entity, oldMask *bitMask, newMask *bitMask) {
	if !m.hasObservers[OnAddComponents] {
		return
	}
	m.doFireAdd(e, oldMask, newMask, true)
}

func (m *observerManager) doFireAdd(e Entity, oldMask *bitMask, newMask *bitMask, earlyOut bool) bool {
	if earlyOut {
		if !m.anyNoComps[OnAddComponents] &&
			(!m.allComps[OnAddComponents].ContainsAny(newMask) || oldMask.Contains(&m.allComps[OnAddComponents])) {
			return false
		}
		if !m.anyNoWith[OnAddComponents] && !m.allWith[OnAddComponents].ContainsAny(oldMask) {
			return false
		}
	}
	observers := m.observers[OnAddComponents]
	found := false
	for _, o := range observers {
		if o.hasComps && (!newMask.Contains(&o.compsMask) || oldMask.ContainsAny(&o.compsMask)) {
			continue
		}
		if o.hasWith && !oldMask.Contains(&o.withMask) {
			continue
		}
		if o.hasWithout && oldMask.ContainsAny(&o.withoutMask) {
			continue
		}
		o.callback(e)
		found = true
	}
	return found
}

func (m *observerManager) doFireRemove(e Entity, oldMask *bitMask, newMask *bitMask, earlyOut bool) bool {
	if earlyOut {
		if !m.anyNoComps[OnRemoveComponents] &&
			(!m.allComps[OnRemoveComponents].ContainsAny(oldMask) || newMask.Contains(&m.allComps[OnRemoveComponents])) {
			return false
		}
		if !m.anyNoWith[OnRemoveComponents] && !m.allWith[OnRemoveComponents].ContainsAny(oldMask) {
			return false
		}
	}
	observers := m.observers[OnRemoveComponents]
	found := false
	for _, o := range observers {
		if o.hasComps && (newMask.Contains(&o.compsMask) || !oldMask.ContainsAny(&o.compsMask)) {
			continue
		}
		if o.hasWith && !oldMask.Contains(&o.withMask) {
			continue
		}
		if o.hasWithout && oldMask.ContainsAny(&o.withoutMask) {
			continue
		}
		o.callback(e)
		found = true
	}
	return found
}

func (m *observerManager) doFireSet(e Entity, mask *bitMask, newMask *bitMask) {
	if !m.anyNoComps[OnSetComponents] && !m.allComps[OnSetComponents].ContainsAny(mask) {
		return
	}
	if !m.anyNoWith[OnSetComponents] && !m.allWith[OnSetComponents].ContainsAny(newMask) {
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

func (m *observerManager) FireCustom(evt EventType, e Entity, mask, entityMask *bitMask) {
	if !m.hasObservers[evt] {
		return
	}
	m.doFireCustom(evt, e, mask, entityMask)
}

func (m *observerManager) doFireCustom(evt EventType, e Entity, mask, entityMask *bitMask) {
	if !m.anyNoComps[evt] && !m.allComps[evt].ContainsAny(mask) {
		return
	}
	if !m.anyNoWith[evt] && !m.allWith[evt].ContainsAny(entityMask) {
		return
	}
	observers := m.observers[evt]
	for _, o := range observers {
		if o.hasComps && !mask.Contains(&o.compsMask) {
			continue
		}
		if o.hasWith && !entityMask.Contains(&o.withMask) {
			continue
		}
		if o.hasWithout && entityMask.ContainsAny(&o.withoutMask) {
			continue
		}
		o.callback(e)
	}
}
