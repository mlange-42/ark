package ecs

import (
	"fmt"
	"math"
)

type observerID uint32

const maxObserverID = math.MaxUint32

// EventType is the type for event identifiers.
// See below for predefined event types.
//
// Use an [EventRegistry] to create custom events types.
// See [Event] and [World.Event] for using custom events.
//
// See [Observer] for details on events and observers.
type EventType uint8

// Predefined event types.
const (
	customEvent EventType = iota + 248

	// OnCreateEntity event.
	// Emitted after an entity is created.
	OnCreateEntity

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
	// Includes creating entities, added components as well as setting relation targets.
	OnAddRelations

	// OnRemoveRelations event.
	// Emitted before relation targets are removed from an entity.
	// Includes removing entities, removing components as well as setting relation targets.
	OnRemoveRelations
)

// EventRegistry for creating new custom event types (see [EventType]).
//
// Your application should have a single, global EventRegistry.
// Using event types from multiple registries in the same [World]
// leads to conflicts.
type EventRegistry struct {
	nextID EventType
}

// NewEventType creates a new EventType for custom events.
// Custom event types should be stored in global variables.
//
// The maximum number of event types is limited to 255, with 7 predefined and 248 potential custom types.
//
// See [Event] and [World.Event] for using custom events.
func (r *EventRegistry) NewEventType() EventType {
	if r.nextID > customEvent {
		panic(fmt.Sprintf("reached maximum number of %d custom event types", customEvent))
	}
	id := r.nextID
	r.nextID++
	return id
}

// Event is a custom event.
//
// Create events using [World.Event].
// Create custom event types using an [EventRegistry].
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
		e.mask.Set(id.id)
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
	indices      map[observerID]observerIndex
	totalCount   uint32
	maxEventType EventType
}

type observerIndex struct {
	idx   uint32
	event EventType
}

func newObserverManager() observerManager {
	maxEvents := math.MaxUint8 + 1
	return observerManager{
		observers:    make([][]*Observer, maxEvents),
		hasObservers: make([]bool, maxEvents),
		anyNoComps:   make([]bool, maxEvents),
		anyNoWith:    make([]bool, maxEvents),
		allComps:     make([]bitMask, maxEvents),
		allWith:      make([]bitMask, maxEvents),
		pool:         newIntPool[observerID](32),
		indices:      map[observerID]observerIndex{},
	}
}

func (m *observerManager) AddObserver(o *Observer, w *World) {
	if o.id != maxObserverID {
		panic("observer is already registered")
	}
	if o.callback == nil {
		panic("observer callback must be set via Do before registering")
	}

	o.id = m.pool.Get()

	o.hasComps, o.hasWith, o.hasWithout = false, false, false

	switch o.event {
	case OnAddRelations, OnRemoveRelations:
		for _, c := range o.comps {
			id := TypeID(w, c.tp)
			if !w.storage.registry.IsRelation[id.id] {
				panic(fmt.Sprintf("non-relation component %d in relation observer", id.id))
			}
			o.compsMask.Set(id.id)
			o.hasComps = true
		}
	case OnCreateEntity, OnRemoveEntity:
		for _, c := range o.comps {
			id := TypeID(w, c.tp)
			o.withMask.Set(id.id)
			o.hasWith = true
		}
	default:
		for _, c := range o.comps {
			id := TypeID(w, c.tp)
			o.compsMask.Set(id.id)
			o.hasComps = true
		}
	}

	for _, c := range o.with {
		id := TypeID(w, c.tp)
		o.withMask.Set(id.id)
		o.hasWith = true
	}
	if o.exclusive {
		o.withoutMask = o.withMask.Not()
		o.hasWithout = true
	} else {
		for _, c := range o.without {
			id := TypeID(w, c.tp)
			o.withoutMask.Set(id.id)
			o.hasWithout = true
		}
	}

	m.indices[o.id] = observerIndex{idx: uint32(len(m.observers[o.event])), event: o.event}
	m.observers[o.event] = append(m.observers[o.event], o)
	m.hasObservers[o.event] = true
	if o.event > m.maxEventType {
		m.maxEventType = o.event
	}
	m.totalCount++

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
	observers[idx.idx].id = maxObserverID

	last := uint32(len(observers) - 1)
	if idx.idx != last {
		observers[idx.idx], observers[last] = observers[last], observers[idx.idx]
		m.indices[observers[idx.idx].id] = idx
	}
	observers[last] = nil
	m.observers[o.event] = observers[:last]
	m.hasObservers[o.event] = last > 0
	m.totalCount--

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

func (m *observerManager) FireCreateEntityRel(e Entity, mask *bitMask) {
	if !m.hasObservers[OnAddRelations] {
		return
	}
	m.doFireCreateEntityRel(e, mask, true)
}

func (m *observerManager) doFireCreateEntityRel(e Entity, mask *bitMask, earlyOut bool) bool {
	if earlyOut {
		if !m.anyNoComps[OnAddRelations] && !m.allComps[OnAddRelations].ContainsAny(mask) {
			return false
		}
		if !m.anyNoWith[OnAddRelations] && !m.allWith[OnAddRelations].ContainsAny(mask) {
			return false
		}
	}
	observers := m.observers[OnAddRelations]
	found := false
	for _, o := range observers {
		if o.hasComps && !mask.Contains(&o.compsMask) {
			continue
		}
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

func (m *observerManager) doFireRemoveEntityRel(e Entity, mask *bitMask, earlyOut bool) bool {
	if earlyOut {
		if !m.anyNoComps[OnRemoveRelations] && !m.allComps[OnRemoveRelations].ContainsAny(mask) {
			return false
		}
		if !m.anyNoWith[OnRemoveRelations] && !m.allWith[OnRemoveRelations].ContainsAny(mask) {
			return false
		}
	}
	observers := m.observers[OnRemoveRelations]
	found := false
	for _, o := range observers {
		if o.hasComps && !mask.Contains(&o.compsMask) {
			continue
		}
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

func (m *observerManager) FireAdd(evt EventType, e Entity, oldMask *bitMask, newMask *bitMask) {
	if !m.hasObservers[evt] {
		return
	}
	m.doFireAdd(evt, e, oldMask, newMask, true)
}

func (m *observerManager) doFireAdd(evt EventType, e Entity, oldMask *bitMask, newMask *bitMask, earlyOut bool) bool {
	if earlyOut {
		if !m.anyNoComps[evt] &&
			(!m.allComps[evt].ContainsAny(newMask) || oldMask.Contains(&m.allComps[evt])) {
			return false
		}
		if !m.anyNoWith[evt] && !m.allWith[evt].ContainsAny(oldMask) {
			return false
		}
	}
	observers := m.observers[evt]
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

func (m *observerManager) doFireRemove(evt EventType, e Entity, oldMask *bitMask, newMask *bitMask, earlyOut bool) bool {
	if earlyOut {
		if !m.anyNoComps[evt] &&
			(!m.allComps[evt].ContainsAny(oldMask) || newMask.Contains(&m.allComps[evt])) {
			return false
		}
		if !m.anyNoWith[evt] && !m.allWith[evt].ContainsAny(oldMask) {
			return false
		}
	}
	observers := m.observers[evt]
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

func (m *observerManager) doFireSetRelations(evt EventType, e Entity, mask *bitMask, newMask *bitMask, earlyOut bool) bool {
	if earlyOut {
		if !m.anyNoComps[evt] && !m.allComps[evt].ContainsAny(mask) {
			return false
		}
		if !m.anyNoWith[evt] && !m.allWith[evt].ContainsAny(newMask) {
			return false
		}
	}
	observers := m.observers[evt]
	found := false
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
		found = true
	}
	return found
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

func (m *observerManager) Reset() {
	if len(m.indices) == 0 {
		m.maxEventType = 0
		return
	}

	for _, idx := range m.indices {
		m.observers[idx.event][idx.idx].id = maxObserverID
	}

	for i := range m.maxEventType + 1 {
		m.observers[i] = m.observers[i][:0]
		m.hasObservers[i] = false
		m.allComps[i].Reset()
		m.allWith[i].Reset()
		m.anyNoComps[i] = false
		m.anyNoWith[i] = false
	}

	m.indices = map[observerID]observerIndex{}
	m.pool.Reset()
	m.totalCount = 0
	m.maxEventType = 0
}
