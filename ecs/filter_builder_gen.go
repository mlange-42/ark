package ecs

// Code generated by go generate; DO NOT EDIT.

// FilterBuilder0 builds a [Filter0].
type FilterBuilder0 struct {
	world *World
	ids   []ID
}

// NewFilter0 creates a new [FilterBuilder0].
//
// Use [FilterBuilder0.Build] to obtain a [Filter0].
func NewFilter0(world *World) *FilterBuilder0 {
	ids := []ID{}

	return &FilterBuilder0{
		world: world,
		ids:   ids,
	}
}

// Build creates a [Filter0] from this builder.
func (q *FilterBuilder0) Build() *Filter0 {
	return &Filter0{
		world: q.world,
		ids:   q.ids,
		mask:  All(q.ids[:]...),
	}
}

// FilterBuilder1 builds a [Filter1].
type FilterBuilder1[A any] struct {
	world *World
	ids   []ID
}

// NewFilter1 creates a new [FilterBuilder1].
//
// Use [FilterBuilder1.Build] to obtain a [Filter1].
func NewFilter1[A any](world *World) *FilterBuilder1[A] {
	ids := []ID{
		ComponentID[A](world),
	}

	return &FilterBuilder1[A]{
		world: world,
		ids:   ids,
	}
}

// Build creates a [Filter1] from this builder.
func (q *FilterBuilder1[A]) Build() *Filter1[A] {
	return &Filter1[A]{
		world: q.world,
		ids:   q.ids,
		mask:  All(q.ids[:]...),
	}
}

// FilterBuilder2 builds a [Filter2].
type FilterBuilder2[A any, B any] struct {
	world *World
	ids   []ID
}

// NewFilter2 creates a new [FilterBuilder2].
//
// Use [FilterBuilder2.Build] to obtain a [Filter2].
func NewFilter2[A any, B any](world *World) *FilterBuilder2[A, B] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
	}

	return &FilterBuilder2[A, B]{
		world: world,
		ids:   ids,
	}
}

// Build creates a [Filter2] from this builder.
func (q *FilterBuilder2[A, B]) Build() *Filter2[A, B] {
	return &Filter2[A, B]{
		world: q.world,
		ids:   q.ids,
		mask:  All(q.ids[:]...),
	}
}

// FilterBuilder3 builds a [Filter3].
type FilterBuilder3[A any, B any, C any] struct {
	world *World
	ids   []ID
}

// NewFilter3 creates a new [FilterBuilder3].
//
// Use [FilterBuilder3.Build] to obtain a [Filter3].
func NewFilter3[A any, B any, C any](world *World) *FilterBuilder3[A, B, C] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
		ComponentID[C](world),
	}

	return &FilterBuilder3[A, B, C]{
		world: world,
		ids:   ids,
	}
}

// Build creates a [Filter3] from this builder.
func (q *FilterBuilder3[A, B, C]) Build() *Filter3[A, B, C] {
	return &Filter3[A, B, C]{
		world: q.world,
		ids:   q.ids,
		mask:  All(q.ids[:]...),
	}
}

// FilterBuilder4 builds a [Filter4].
type FilterBuilder4[A any, B any, C any, D any] struct {
	world *World
	ids   []ID
}

// NewFilter4 creates a new [FilterBuilder4].
//
// Use [FilterBuilder4.Build] to obtain a [Filter4].
func NewFilter4[A any, B any, C any, D any](world *World) *FilterBuilder4[A, B, C, D] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
		ComponentID[C](world),
		ComponentID[D](world),
	}

	return &FilterBuilder4[A, B, C, D]{
		world: world,
		ids:   ids,
	}
}

// Build creates a [Filter4] from this builder.
func (q *FilterBuilder4[A, B, C, D]) Build() *Filter4[A, B, C, D] {
	return &Filter4[A, B, C, D]{
		world: q.world,
		ids:   q.ids,
		mask:  All(q.ids[:]...),
	}
}

// FilterBuilder5 builds a [Filter5].
type FilterBuilder5[A any, B any, C any, D any, E any] struct {
	world *World
	ids   []ID
}

// NewFilter5 creates a new [FilterBuilder5].
//
// Use [FilterBuilder5.Build] to obtain a [Filter5].
func NewFilter5[A any, B any, C any, D any, E any](world *World) *FilterBuilder5[A, B, C, D, E] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
		ComponentID[C](world),
		ComponentID[D](world),
		ComponentID[E](world),
	}

	return &FilterBuilder5[A, B, C, D, E]{
		world: world,
		ids:   ids,
	}
}

// Build creates a [Filter5] from this builder.
func (q *FilterBuilder5[A, B, C, D, E]) Build() *Filter5[A, B, C, D, E] {
	return &Filter5[A, B, C, D, E]{
		world: q.world,
		ids:   q.ids,
		mask:  All(q.ids[:]...),
	}
}

// FilterBuilder6 builds a [Filter6].
type FilterBuilder6[A any, B any, C any, D any, E any, F any] struct {
	world *World
	ids   []ID
}

// NewFilter6 creates a new [FilterBuilder6].
//
// Use [FilterBuilder6.Build] to obtain a [Filter6].
func NewFilter6[A any, B any, C any, D any, E any, F any](world *World) *FilterBuilder6[A, B, C, D, E, F] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
		ComponentID[C](world),
		ComponentID[D](world),
		ComponentID[E](world),
		ComponentID[F](world),
	}

	return &FilterBuilder6[A, B, C, D, E, F]{
		world: world,
		ids:   ids,
	}
}

// Build creates a [Filter6] from this builder.
func (q *FilterBuilder6[A, B, C, D, E, F]) Build() *Filter6[A, B, C, D, E, F] {
	return &Filter6[A, B, C, D, E, F]{
		world: q.world,
		ids:   q.ids,
		mask:  All(q.ids[:]...),
	}
}

// FilterBuilder7 builds a [Filter7].
type FilterBuilder7[A any, B any, C any, D any, E any, F any, G any] struct {
	world *World
	ids   []ID
}

// NewFilter7 creates a new [FilterBuilder7].
//
// Use [FilterBuilder7.Build] to obtain a [Filter7].
func NewFilter7[A any, B any, C any, D any, E any, F any, G any](world *World) *FilterBuilder7[A, B, C, D, E, F, G] {
	ids := []ID{
		ComponentID[A](world),
		ComponentID[B](world),
		ComponentID[C](world),
		ComponentID[D](world),
		ComponentID[E](world),
		ComponentID[F](world),
		ComponentID[G](world),
	}

	return &FilterBuilder7[A, B, C, D, E, F, G]{
		world: world,
		ids:   ids,
	}
}

// Build creates a [Filter7] from this builder.
func (q *FilterBuilder7[A, B, C, D, E, F, G]) Build() *Filter7[A, B, C, D, E, F, G] {
	return &Filter7[A, B, C, D, E, F, G]{
		world: q.world,
		ids:   q.ids,
		mask:  All(q.ids[:]...),
	}
}

// FilterBuilder8 builds a [Filter8].
type FilterBuilder8[A any, B any, C any, D any, E any, F any, G any, H any] struct {
	world *World
	ids   []ID
}

// NewFilter8 creates a new [FilterBuilder8].
//
// Use [FilterBuilder8.Build] to obtain a [Filter8].
func NewFilter8[A any, B any, C any, D any, E any, F any, G any, H any](world *World) *FilterBuilder8[A, B, C, D, E, F, G, H] {
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

	return &FilterBuilder8[A, B, C, D, E, F, G, H]{
		world: world,
		ids:   ids,
	}
}

// Build creates a [Filter8] from this builder.
func (q *FilterBuilder8[A, B, C, D, E, F, G, H]) Build() *Filter8[A, B, C, D, E, F, G, H] {
	return &Filter8[A, B, C, D, E, F, G, H]{
		world: q.world,
		ids:   q.ids,
		mask:  All(q.ids[:]...),
	}
}
