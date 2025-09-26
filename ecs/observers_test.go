package ecs

import (
	"fmt"
	"testing"
)

func TestObserver2(t *testing.T) {
	w := NewWorld()
	builder := NewMap2[Position, Velocity](&w)

	Observe2[Position, Velocity](OnCreateEntity).
		Do(func(e Entity, p *Position, v *Velocity) {
			fmt.Printf("%#v\n", p)
			fmt.Printf("%#v\n", v)
		}).
		Register(&w)

	//e := w.NewEntity()
	//builder.Add(e, &Position{}, &Velocity{})

	builder.NewEntity(&Position{1, 2}, &Velocity{3, 4})
}
