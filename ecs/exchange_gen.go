package ecs

// Code generated by go generate; DO NOT EDIT.

import "unsafe"

// Exchange1 allows to exchange components of entities.
// It adds the given components. Use [Exchange1.Removes]
// to set components to be removed.
type Exchange1[A any] struct {
	world  *World
	ids    []ID
	remove []ID
}

// NewExchange1 creates an [Exchange1].
func NewExchange1[A any](world *World) *Exchange1[A] {
	ids := []ID{
		ComponentID[A](world),
	}
	return &Exchange1[A]{
		world: world,
		ids:   ids,
	}
}

// Removes sets the components that this [Exchange1] removes.
func (ex *Exchange1[A]) Removes(components ...Comp) *Exchange1[A] {
	ids := make([]ID, len(components))
	for i, c := range components {
		ids[i] = ex.world.componentID(c.tp)
	}
	ex.remove = ids
	return ex
}

// Exchange performs the exchange on the given entity, adding the provided components
// and removing those previously specified with [Exchange1.Removes].
func (ex *Exchange1[A]) Exchange(entity Entity, a *A) {
	ex.world.exchange(entity, ex.ids, ex.remove, []unsafe.Pointer{
		unsafe.Pointer(a),
	})
}

// Exchange2 allows to exchange components of entities.
// It adds the given components. Use [Exchange2.Removes]
// to set components to be removed.
type Exchange2[A any, B any] struct {
	world  *World
	ids    []ID
	remove []ID
}

// NewExchange2 creates an [Exchange2].
func NewExchange2[A any, B any](world *World) *Exchange2[A, B] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
	}
	return &Exchange2[A, B]{
		world: world,
		ids:   ids,
	}
}

// Removes sets the components that this [Exchange2] removes.
func (ex *Exchange2[A, B]) Removes(components ...Comp) *Exchange2[A, B] {
	ids := make([]ID, len(components))
	for i, c := range components {
		ids[i] = ex.world.componentID(c.tp)
	}
	ex.remove = ids
	return ex
}

// Exchange performs the exchange on the given entity, adding the provided components
// and removing those previously specified with [Exchange2.Removes].
func (ex *Exchange2[A, B]) Exchange(entity Entity, a *A, b *B) {
	ex.world.exchange(entity, ex.ids, ex.remove, []unsafe.Pointer{
		unsafe.Pointer(a),
		unsafe.Pointer(b),
	})
}

// Exchange3 allows to exchange components of entities.
// It adds the given components. Use [Exchange3.Removes]
// to set components to be removed.
type Exchange3[A any, B any, C any] struct {
	world  *World
	ids    []ID
	remove []ID
}

// NewExchange3 creates an [Exchange3].
func NewExchange3[A any, B any, C any](world *World) *Exchange3[A, B, C] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
		ComponentID[C](world),
	}
	return &Exchange3[A, B, C]{
		world: world,
		ids:   ids,
	}
}

// Removes sets the components that this [Exchange3] removes.
func (ex *Exchange3[A, B, C]) Removes(components ...Comp) *Exchange3[A, B, C] {
	ids := make([]ID, len(components))
	for i, c := range components {
		ids[i] = ex.world.componentID(c.tp)
	}
	ex.remove = ids
	return ex
}

// Exchange performs the exchange on the given entity, adding the provided components
// and removing those previously specified with [Exchange3.Removes].
func (ex *Exchange3[A, B, C]) Exchange(entity Entity, a *A, b *B, c *C) {
	ex.world.exchange(entity, ex.ids, ex.remove, []unsafe.Pointer{
		unsafe.Pointer(a),
		unsafe.Pointer(b),
		unsafe.Pointer(c),
	})
}

// Exchange4 allows to exchange components of entities.
// It adds the given components. Use [Exchange4.Removes]
// to set components to be removed.
type Exchange4[A any, B any, C any, D any] struct {
	world  *World
	ids    []ID
	remove []ID
}

// NewExchange4 creates an [Exchange4].
func NewExchange4[A any, B any, C any, D any](world *World) *Exchange4[A, B, C, D] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
		ComponentID[C](world),
		ComponentID[D](world),
	}
	return &Exchange4[A, B, C, D]{
		world: world,
		ids:   ids,
	}
}

// Removes sets the components that this [Exchange4] removes.
func (ex *Exchange4[A, B, C, D]) Removes(components ...Comp) *Exchange4[A, B, C, D] {
	ids := make([]ID, len(components))
	for i, c := range components {
		ids[i] = ex.world.componentID(c.tp)
	}
	ex.remove = ids
	return ex
}

// Exchange performs the exchange on the given entity, adding the provided components
// and removing those previously specified with [Exchange4.Removes].
func (ex *Exchange4[A, B, C, D]) Exchange(entity Entity, a *A, b *B, c *C, d *D) {
	ex.world.exchange(entity, ex.ids, ex.remove, []unsafe.Pointer{
		unsafe.Pointer(a),
		unsafe.Pointer(b),
		unsafe.Pointer(c),
		unsafe.Pointer(d),
	})
}

// Exchange5 allows to exchange components of entities.
// It adds the given components. Use [Exchange5.Removes]
// to set components to be removed.
type Exchange5[A any, B any, C any, D any, E any] struct {
	world  *World
	ids    []ID
	remove []ID
}

// NewExchange5 creates an [Exchange5].
func NewExchange5[A any, B any, C any, D any, E any](world *World) *Exchange5[A, B, C, D, E] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
		ComponentID[C](world),
		ComponentID[D](world),
		ComponentID[E](world),
	}
	return &Exchange5[A, B, C, D, E]{
		world: world,
		ids:   ids,
	}
}

// Removes sets the components that this [Exchange5] removes.
func (ex *Exchange5[A, B, C, D, E]) Removes(components ...Comp) *Exchange5[A, B, C, D, E] {
	ids := make([]ID, len(components))
	for i, c := range components {
		ids[i] = ex.world.componentID(c.tp)
	}
	ex.remove = ids
	return ex
}

// Exchange performs the exchange on the given entity, adding the provided components
// and removing those previously specified with [Exchange5.Removes].
func (ex *Exchange5[A, B, C, D, E]) Exchange(entity Entity, a *A, b *B, c *C, d *D, e *E) {
	ex.world.exchange(entity, ex.ids, ex.remove, []unsafe.Pointer{
		unsafe.Pointer(a),
		unsafe.Pointer(b),
		unsafe.Pointer(c),
		unsafe.Pointer(d),
		unsafe.Pointer(e),
	})
}

// Exchange6 allows to exchange components of entities.
// It adds the given components. Use [Exchange6.Removes]
// to set components to be removed.
type Exchange6[A any, B any, C any, D any, E any, F any] struct {
	world  *World
	ids    []ID
	remove []ID
}

// NewExchange6 creates an [Exchange6].
func NewExchange6[A any, B any, C any, D any, E any, F any](world *World) *Exchange6[A, B, C, D, E, F] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
		ComponentID[C](world),
		ComponentID[D](world),
		ComponentID[E](world),
		ComponentID[F](world),
	}
	return &Exchange6[A, B, C, D, E, F]{
		world: world,
		ids:   ids,
	}
}

// Removes sets the components that this [Exchange6] removes.
func (ex *Exchange6[A, B, C, D, E, F]) Removes(components ...Comp) *Exchange6[A, B, C, D, E, F] {
	ids := make([]ID, len(components))
	for i, c := range components {
		ids[i] = ex.world.componentID(c.tp)
	}
	ex.remove = ids
	return ex
}

// Exchange performs the exchange on the given entity, adding the provided components
// and removing those previously specified with [Exchange6.Removes].
func (ex *Exchange6[A, B, C, D, E, F]) Exchange(entity Entity, a *A, b *B, c *C, d *D, e *E, f *F) {
	ex.world.exchange(entity, ex.ids, ex.remove, []unsafe.Pointer{
		unsafe.Pointer(a),
		unsafe.Pointer(b),
		unsafe.Pointer(c),
		unsafe.Pointer(d),
		unsafe.Pointer(e),
		unsafe.Pointer(f),
	})
}

// Exchange7 allows to exchange components of entities.
// It adds the given components. Use [Exchange7.Removes]
// to set components to be removed.
type Exchange7[A any, B any, C any, D any, E any, F any, G any] struct {
	world  *World
	ids    []ID
	remove []ID
}

// NewExchange7 creates an [Exchange7].
func NewExchange7[A any, B any, C any, D any, E any, F any, G any](world *World) *Exchange7[A, B, C, D, E, F, G] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
		ComponentID[C](world),
		ComponentID[D](world),
		ComponentID[E](world),
		ComponentID[F](world),
		ComponentID[G](world),
	}
	return &Exchange7[A, B, C, D, E, F, G]{
		world: world,
		ids:   ids,
	}
}

// Removes sets the components that this [Exchange7] removes.
func (ex *Exchange7[A, B, C, D, E, F, G]) Removes(components ...Comp) *Exchange7[A, B, C, D, E, F, G] {
	ids := make([]ID, len(components))
	for i, c := range components {
		ids[i] = ex.world.componentID(c.tp)
	}
	ex.remove = ids
	return ex
}

// Exchange performs the exchange on the given entity, adding the provided components
// and removing those previously specified with [Exchange7.Removes].
func (ex *Exchange7[A, B, C, D, E, F, G]) Exchange(entity Entity, a *A, b *B, c *C, d *D, e *E, f *F, g *G) {
	ex.world.exchange(entity, ex.ids, ex.remove, []unsafe.Pointer{
		unsafe.Pointer(a),
		unsafe.Pointer(b),
		unsafe.Pointer(c),
		unsafe.Pointer(d),
		unsafe.Pointer(e),
		unsafe.Pointer(f),
		unsafe.Pointer(g),
	})
}

// Exchange8 allows to exchange components of entities.
// It adds the given components. Use [Exchange8.Removes]
// to set components to be removed.
type Exchange8[A any, B any, C any, D any, E any, F any, G any, H any] struct {
	world  *World
	ids    []ID
	remove []ID
}

// NewExchange8 creates an [Exchange8].
func NewExchange8[A any, B any, C any, D any, E any, F any, G any, H any](world *World) *Exchange8[A, B, C, D, E, F, G, H] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
		ComponentID[C](world),
		ComponentID[D](world),
		ComponentID[E](world),
		ComponentID[F](world),
		ComponentID[G](world),
		ComponentID[H](world),
	}
	return &Exchange8[A, B, C, D, E, F, G, H]{
		world: world,
		ids:   ids,
	}
}

// Removes sets the components that this [Exchange8] removes.
func (ex *Exchange8[A, B, C, D, E, F, G, H]) Removes(components ...Comp) *Exchange8[A, B, C, D, E, F, G, H] {
	ids := make([]ID, len(components))
	for i, c := range components {
		ids[i] = ex.world.componentID(c.tp)
	}
	ex.remove = ids
	return ex
}

// Exchange performs the exchange on the given entity, adding the provided components
// and removing those previously specified with [Exchange8.Removes].
func (ex *Exchange8[A, B, C, D, E, F, G, H]) Exchange(entity Entity, a *A, b *B, c *C, d *D, e *E, f *F, g *G, h *H) {
	ex.world.exchange(entity, ex.ids, ex.remove, []unsafe.Pointer{
		unsafe.Pointer(a),
		unsafe.Pointer(b),
		unsafe.Pointer(c),
		unsafe.Pointer(d),
		unsafe.Pointer(e),
		unsafe.Pointer(f),
		unsafe.Pointer(g),
		unsafe.Pointer(h),
	})
}
