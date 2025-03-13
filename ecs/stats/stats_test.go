package stats

import (
	"fmt"
	"reflect"
	"testing"
)

func TestStats(t *testing.T) {
	stats := World{
		Entities:       Entities{},
		ComponentTypes: []reflect.Type{reflect.TypeOf(1)},
		Locked:         false,
		Archetypes: []Archetype{
			{
				Size:           1,
				Capacity:       128,
				ComponentIDs:   []uint8{0},
				ComponentTypes: []reflect.Type{reflect.TypeOf(1)},
			},
			{
				Size:           1,
				Capacity:       128,
				ComponentIDs:   []uint8{0},
				ComponentTypes: []reflect.Type{reflect.TypeOf(1)},
			},
		},
	}
	fmt.Println(stats.String())
}
