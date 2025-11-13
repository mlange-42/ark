package ecs

import (
	"encoding"
	"encoding/json"
	"runtime"
	"testing"
)

func TestEntity(t *testing.T) {
	e := newEntity(100)
	expectEqual(t, 100, e.ID())
	expectEqual(t, 0, e.Gen())
}

func TestEntityIndex(t *testing.T) {
	index := entityIndex{}
	expectEqual(t, 0, index.table)
	expectEqual(t, 0, index.row)
}

func TestReservedEntities(t *testing.T) {
	w := NewWorld(1024)

	zero := Entity{}
	wildcard := Entity{1, 0}

	expectFalse(t, w.Alive(zero))
	expectFalse(t, w.Alive(wildcard))

	expectTrue(t, zero.IsZero())
	expectFalse(t, wildcard.IsZero())
	expectTrue(t, wildcard.isWildcard())
}

func TestEntityMarshal(t *testing.T) {
	e := Entity{2, 3}

	var _ json.Marshaler = &e
	var _ json.Unmarshaler = &e

	jsonData, err := json.Marshal(&e)
	expectNil(t, err)

	e2 := Entity{}
	err = json.Unmarshal(jsonData, &e2)
	expectNil(t, err)

	expectEqual(t, e2, e)

	err = e2.UnmarshalJSON([]byte("pft"))
	expectNotNil(t, err)
}

func TestEntityMarshalBinary(t *testing.T) {
	e := Entity{2, 3}

	var _ encoding.BinaryMarshaler = &e
	var _ encoding.BinaryUnmarshaler = &e
	var _ encoding.BinaryAppender = &e

	binData, err := e.MarshalBinary()
	expectNil(t, err)

	e2 := Entity{}
	err = e2.UnmarshalBinary(binData)
	expectNil(t, err)
	expectEqual(t, 8, len(binData))

	expectEqual(t, e2, e)

	err = e2.UnmarshalBinary(make([]byte, 9))
	expectNotNil(t, err)

	e = Entity{4, 5}
	binData, err = e.AppendBinary(binData)
	expectNil(t, err)
	expectEqual(t, 16, len(binData))

	err = e2.UnmarshalBinary(binData[8:])
	expectNil(t, err)
	expectEqual(t, e2, e)

	err = e2.UnmarshalBinary(binData[:8])
	expectNil(t, err)
	expectEqual(t, e2, Entity{2, 3})
}

func BenchmarkEntityMarshalBinary_1000(b *testing.B) {
	w := NewWorld()

	entities := make([]Entity, 0, 1000)
	w.NewEntities(1000, func(e Entity) {
		entities = append(entities, e)
	})

	var binData []byte
	loop := func() {
		for _, e := range entities {
			binData, _ = e.MarshalBinary()
		}
	}
	for b.Loop() {
		loop()
	}

	runtime.KeepAlive(binData)
}

func BenchmarkEntityAppendBinary_1000(b *testing.B) {
	w := NewWorld()

	entities := make([]Entity, 0, 1000)
	w.NewEntities(1000, func(e Entity) {
		entities = append(entities, e)
	})

	binData := make([]byte, 0, 8000)
	loop := func() {
		binData = binData[:0]
		for _, e := range entities {
			binData, _ = e.AppendBinary(binData)
		}
	}
	for b.Loop() {
		loop()
	}

	runtime.KeepAlive(binData)
}

func BenchmarkEntityUnmarshalBinary_1000(b *testing.B) {
	w := NewWorld()

	entities := make([][]byte, 0, 1000)
	w.NewEntities(1000, func(e Entity) {
		binData, _ := e.MarshalBinary()
		entities = append(entities, binData)
	})

	var entity Entity
	loop := func() {
		for _, e := range entities {
			_ = entity.UnmarshalBinary(e)
		}
	}
	for b.Loop() {
		loop()
	}

	runtime.KeepAlive(entity)
}
