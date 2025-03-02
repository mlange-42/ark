package ecs

// Code generated by go generate; DO NOT EDIT.

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExchange1(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap1[CompA](&w)
	ex := NewExchange1[CompA](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Exchange(e, &CompA{})
	assert.False(t, posMap.HasAll(e))
	assert.True(t, mapper.HasAll(e))
}

func TestExchange1Add(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap1[CompA](&w)
	ex := NewExchange1[CompA](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Add(e, &CompA{})
	assert.True(t, posMap.HasAll(e))
	assert.True(t, mapper.HasAll(e))
}

func TestExchange1Remove(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap1[CompA](&w)
	ex := NewExchange1[CompA](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Remove(e)
	assert.False(t, posMap.HasAll(e))
	assert.False(t, mapper.HasAll(e))
}

func TestExchange1AddBatch(t *testing.T) {
	n := 12
	w := NewWorld(8)

	exchange := NewExchange1[CompA](&w).Removes(C[CompA]())
	posMap := NewMap1[Position](&w)
	posVelMap := NewMap2[Position, Velocity](&w)

	cnt := 1
	posMap.NewBatchFn(n, func(entity Entity, pos *Position) {
		pos.X = float64(cnt)
		cnt++
	})
	posVelMap.NewBatchFn(n, func(entity Entity, pos *Position, _ *Velocity) {
		pos.X = float64(cnt)
		cnt++
	})
	assert.Equal(t, 2*n+1, cnt)

	filter := NewFilter1[Position](&w)
	exchange.AddBatch(filter.Batch(), &CompA{})

	filter2 := NewFilter1[CompA](&w)
	query := filter2.Query()
	cnt = 0
	for query.Next() {
		pos := posMap.Get(query.Entity())
		assert.Greater(t, pos.X, 0.0)
		cnt++
	}
	assert.Equal(t, 2*n, cnt)

	exchange.RemoveBatch(filter2.Batch(), nil)

	query = filter2.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)
}

func TestExchange1AddBatchFn(t *testing.T) {
	n := 12
	w := NewWorld(8)

	exchange := NewExchange1[CompA](&w).Removes(C[CompA]())
	posMap := NewMap1[Position](&w)
	posVelMap := NewMap2[Position, Velocity](&w)

	cnt := 1
	posMap.NewBatchFn(n, func(entity Entity, pos *Position) {
		pos.X = float64(cnt)
		cnt++
	})
	posVelMap.NewBatchFn(n, func(entity Entity, pos *Position, _ *Velocity) {
		pos.X = float64(cnt)
		cnt++
	})
	assert.Equal(t, 2*n+1, cnt)

	filter := NewFilter1[Position](&w)
	cnt = 0
	exchange.AddBatchFn(filter.Batch(), func(entity Entity, a *CompA) {
		a.X = float64(cnt)
		cnt++
	})

	filter2 := NewFilter1[CompA](&w)
	query := filter2.Query()
	cnt = 0
	for query.Next() {
		a := query.Get()
		assert.EqualValues(t, cnt, a.X)
		pos := posMap.Get(query.Entity())
		assert.Greater(t, pos.X, 0.0)
		cnt++
	}
	assert.Equal(t, 2*n, cnt)

	cnt = 0
	exchange.RemoveBatch(filter2.Batch(), func(entity Entity) {
		cnt++
	})
	assert.Equal(t, 2*n, cnt)

	query = filter2.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)
}

func TestExchange2(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap2[CompA, CompB](&w)
	ex := NewExchange2[CompA, CompB](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Exchange(e, &CompA{}, &CompB{})
	assert.False(t, posMap.HasAll(e))
	assert.True(t, mapper.HasAll(e))
}

func TestExchange2Add(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap2[CompA, CompB](&w)
	ex := NewExchange2[CompA, CompB](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Add(e, &CompA{}, &CompB{})
	assert.True(t, posMap.HasAll(e))
	assert.True(t, mapper.HasAll(e))
}

func TestExchange2Remove(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap2[CompA, CompB](&w)
	ex := NewExchange2[CompA, CompB](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Remove(e)
	assert.False(t, posMap.HasAll(e))
	assert.False(t, mapper.HasAll(e))
}

func TestExchange2AddBatch(t *testing.T) {
	n := 12
	w := NewWorld(8)

	exchange := NewExchange2[CompA, CompB](&w).Removes(C[CompA]())
	posMap := NewMap1[Position](&w)
	posVelMap := NewMap2[Position, Velocity](&w)

	cnt := 1
	posMap.NewBatchFn(n, func(entity Entity, pos *Position) {
		pos.X = float64(cnt)
		cnt++
	})
	posVelMap.NewBatchFn(n, func(entity Entity, pos *Position, _ *Velocity) {
		pos.X = float64(cnt)
		cnt++
	})
	assert.Equal(t, 2*n+1, cnt)

	filter := NewFilter1[Position](&w)
	exchange.AddBatch(filter.Batch(), &CompA{}, &CompB{})

	filter2 := NewFilter1[CompA](&w)
	query := filter2.Query()
	cnt = 0
	for query.Next() {
		pos := posMap.Get(query.Entity())
		assert.Greater(t, pos.X, 0.0)
		cnt++
	}
	assert.Equal(t, 2*n, cnt)

	exchange.RemoveBatch(filter2.Batch(), nil)

	query = filter2.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)
}

func TestExchange2AddBatchFn(t *testing.T) {
	n := 12
	w := NewWorld(8)

	exchange := NewExchange2[CompA, CompB](&w).Removes(C[CompA]())
	posMap := NewMap1[Position](&w)
	posVelMap := NewMap2[Position, Velocity](&w)

	cnt := 1
	posMap.NewBatchFn(n, func(entity Entity, pos *Position) {
		pos.X = float64(cnt)
		cnt++
	})
	posVelMap.NewBatchFn(n, func(entity Entity, pos *Position, _ *Velocity) {
		pos.X = float64(cnt)
		cnt++
	})
	assert.Equal(t, 2*n+1, cnt)

	filter := NewFilter1[Position](&w)
	cnt = 0
	exchange.AddBatchFn(filter.Batch(), func(entity Entity, a *CompA, b *CompB) {
		a.X = float64(cnt)
		cnt++
	})

	filter2 := NewFilter1[CompA](&w)
	query := filter2.Query()
	cnt = 0
	for query.Next() {
		a := query.Get()
		assert.EqualValues(t, cnt, a.X)
		pos := posMap.Get(query.Entity())
		assert.Greater(t, pos.X, 0.0)
		cnt++
	}
	assert.Equal(t, 2*n, cnt)

	cnt = 0
	exchange.RemoveBatch(filter2.Batch(), func(entity Entity) {
		cnt++
	})
	assert.Equal(t, 2*n, cnt)

	query = filter2.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)
}

func TestExchange3(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap3[CompA, CompB, CompC](&w)
	ex := NewExchange3[CompA, CompB, CompC](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Exchange(e, &CompA{}, &CompB{}, &CompC{})
	assert.False(t, posMap.HasAll(e))
	assert.True(t, mapper.HasAll(e))
}

func TestExchange3Add(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap3[CompA, CompB, CompC](&w)
	ex := NewExchange3[CompA, CompB, CompC](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Add(e, &CompA{}, &CompB{}, &CompC{})
	assert.True(t, posMap.HasAll(e))
	assert.True(t, mapper.HasAll(e))
}

func TestExchange3Remove(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap3[CompA, CompB, CompC](&w)
	ex := NewExchange3[CompA, CompB, CompC](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Remove(e)
	assert.False(t, posMap.HasAll(e))
	assert.False(t, mapper.HasAll(e))
}

func TestExchange3AddBatch(t *testing.T) {
	n := 12
	w := NewWorld(8)

	exchange := NewExchange3[CompA, CompB, CompC](&w).Removes(C[CompA]())
	posMap := NewMap1[Position](&w)
	posVelMap := NewMap2[Position, Velocity](&w)

	cnt := 1
	posMap.NewBatchFn(n, func(entity Entity, pos *Position) {
		pos.X = float64(cnt)
		cnt++
	})
	posVelMap.NewBatchFn(n, func(entity Entity, pos *Position, _ *Velocity) {
		pos.X = float64(cnt)
		cnt++
	})
	assert.Equal(t, 2*n+1, cnt)

	filter := NewFilter1[Position](&w)
	exchange.AddBatch(filter.Batch(), &CompA{}, &CompB{}, &CompC{})

	filter2 := NewFilter1[CompA](&w)
	query := filter2.Query()
	cnt = 0
	for query.Next() {
		pos := posMap.Get(query.Entity())
		assert.Greater(t, pos.X, 0.0)
		cnt++
	}
	assert.Equal(t, 2*n, cnt)

	exchange.RemoveBatch(filter2.Batch(), nil)

	query = filter2.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)
}

func TestExchange3AddBatchFn(t *testing.T) {
	n := 12
	w := NewWorld(8)

	exchange := NewExchange3[CompA, CompB, CompC](&w).Removes(C[CompA]())
	posMap := NewMap1[Position](&w)
	posVelMap := NewMap2[Position, Velocity](&w)

	cnt := 1
	posMap.NewBatchFn(n, func(entity Entity, pos *Position) {
		pos.X = float64(cnt)
		cnt++
	})
	posVelMap.NewBatchFn(n, func(entity Entity, pos *Position, _ *Velocity) {
		pos.X = float64(cnt)
		cnt++
	})
	assert.Equal(t, 2*n+1, cnt)

	filter := NewFilter1[Position](&w)
	cnt = 0
	exchange.AddBatchFn(filter.Batch(), func(entity Entity, a *CompA, b *CompB, c *CompC) {
		a.X = float64(cnt)
		cnt++
	})

	filter2 := NewFilter1[CompA](&w)
	query := filter2.Query()
	cnt = 0
	for query.Next() {
		a := query.Get()
		assert.EqualValues(t, cnt, a.X)
		pos := posMap.Get(query.Entity())
		assert.Greater(t, pos.X, 0.0)
		cnt++
	}
	assert.Equal(t, 2*n, cnt)

	cnt = 0
	exchange.RemoveBatch(filter2.Batch(), func(entity Entity) {
		cnt++
	})
	assert.Equal(t, 2*n, cnt)

	query = filter2.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)
}

func TestExchange4(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap4[CompA, CompB, CompC, CompD](&w)
	ex := NewExchange4[CompA, CompB, CompC, CompD](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Exchange(e, &CompA{}, &CompB{}, &CompC{}, &CompD{})
	assert.False(t, posMap.HasAll(e))
	assert.True(t, mapper.HasAll(e))
}

func TestExchange4Add(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap4[CompA, CompB, CompC, CompD](&w)
	ex := NewExchange4[CompA, CompB, CompC, CompD](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Add(e, &CompA{}, &CompB{}, &CompC{}, &CompD{})
	assert.True(t, posMap.HasAll(e))
	assert.True(t, mapper.HasAll(e))
}

func TestExchange4Remove(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap4[CompA, CompB, CompC, CompD](&w)
	ex := NewExchange4[CompA, CompB, CompC, CompD](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Remove(e)
	assert.False(t, posMap.HasAll(e))
	assert.False(t, mapper.HasAll(e))
}

func TestExchange4AddBatch(t *testing.T) {
	n := 12
	w := NewWorld(8)

	exchange := NewExchange4[CompA, CompB, CompC, CompD](&w).Removes(C[CompA]())
	posMap := NewMap1[Position](&w)
	posVelMap := NewMap2[Position, Velocity](&w)

	cnt := 1
	posMap.NewBatchFn(n, func(entity Entity, pos *Position) {
		pos.X = float64(cnt)
		cnt++
	})
	posVelMap.NewBatchFn(n, func(entity Entity, pos *Position, _ *Velocity) {
		pos.X = float64(cnt)
		cnt++
	})
	assert.Equal(t, 2*n+1, cnt)

	filter := NewFilter1[Position](&w)
	exchange.AddBatch(filter.Batch(), &CompA{}, &CompB{}, &CompC{}, &CompD{})

	filter2 := NewFilter1[CompA](&w)
	query := filter2.Query()
	cnt = 0
	for query.Next() {
		pos := posMap.Get(query.Entity())
		assert.Greater(t, pos.X, 0.0)
		cnt++
	}
	assert.Equal(t, 2*n, cnt)

	exchange.RemoveBatch(filter2.Batch(), nil)

	query = filter2.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)
}

func TestExchange4AddBatchFn(t *testing.T) {
	n := 12
	w := NewWorld(8)

	exchange := NewExchange4[CompA, CompB, CompC, CompD](&w).Removes(C[CompA]())
	posMap := NewMap1[Position](&w)
	posVelMap := NewMap2[Position, Velocity](&w)

	cnt := 1
	posMap.NewBatchFn(n, func(entity Entity, pos *Position) {
		pos.X = float64(cnt)
		cnt++
	})
	posVelMap.NewBatchFn(n, func(entity Entity, pos *Position, _ *Velocity) {
		pos.X = float64(cnt)
		cnt++
	})
	assert.Equal(t, 2*n+1, cnt)

	filter := NewFilter1[Position](&w)
	cnt = 0
	exchange.AddBatchFn(filter.Batch(), func(entity Entity, a *CompA, b *CompB, c *CompC, d *CompD) {
		a.X = float64(cnt)
		cnt++
	})

	filter2 := NewFilter1[CompA](&w)
	query := filter2.Query()
	cnt = 0
	for query.Next() {
		a := query.Get()
		assert.EqualValues(t, cnt, a.X)
		pos := posMap.Get(query.Entity())
		assert.Greater(t, pos.X, 0.0)
		cnt++
	}
	assert.Equal(t, 2*n, cnt)

	cnt = 0
	exchange.RemoveBatch(filter2.Batch(), func(entity Entity) {
		cnt++
	})
	assert.Equal(t, 2*n, cnt)

	query = filter2.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)
}

func TestExchange5(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap5[CompA, CompB, CompC, CompD, CompE](&w)
	ex := NewExchange5[CompA, CompB, CompC, CompD, CompE](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Exchange(e, &CompA{}, &CompB{}, &CompC{}, &CompD{}, &CompE{})
	assert.False(t, posMap.HasAll(e))
	assert.True(t, mapper.HasAll(e))
}

func TestExchange5Add(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap5[CompA, CompB, CompC, CompD, CompE](&w)
	ex := NewExchange5[CompA, CompB, CompC, CompD, CompE](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Add(e, &CompA{}, &CompB{}, &CompC{}, &CompD{}, &CompE{})
	assert.True(t, posMap.HasAll(e))
	assert.True(t, mapper.HasAll(e))
}

func TestExchange5Remove(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap5[CompA, CompB, CompC, CompD, CompE](&w)
	ex := NewExchange5[CompA, CompB, CompC, CompD, CompE](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Remove(e)
	assert.False(t, posMap.HasAll(e))
	assert.False(t, mapper.HasAll(e))
}

func TestExchange5AddBatch(t *testing.T) {
	n := 12
	w := NewWorld(8)

	exchange := NewExchange5[CompA, CompB, CompC, CompD, CompE](&w).Removes(C[CompA]())
	posMap := NewMap1[Position](&w)
	posVelMap := NewMap2[Position, Velocity](&w)

	cnt := 1
	posMap.NewBatchFn(n, func(entity Entity, pos *Position) {
		pos.X = float64(cnt)
		cnt++
	})
	posVelMap.NewBatchFn(n, func(entity Entity, pos *Position, _ *Velocity) {
		pos.X = float64(cnt)
		cnt++
	})
	assert.Equal(t, 2*n+1, cnt)

	filter := NewFilter1[Position](&w)
	exchange.AddBatch(filter.Batch(), &CompA{}, &CompB{}, &CompC{}, &CompD{}, &CompE{})

	filter2 := NewFilter1[CompA](&w)
	query := filter2.Query()
	cnt = 0
	for query.Next() {
		pos := posMap.Get(query.Entity())
		assert.Greater(t, pos.X, 0.0)
		cnt++
	}
	assert.Equal(t, 2*n, cnt)

	exchange.RemoveBatch(filter2.Batch(), nil)

	query = filter2.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)
}

func TestExchange5AddBatchFn(t *testing.T) {
	n := 12
	w := NewWorld(8)

	exchange := NewExchange5[CompA, CompB, CompC, CompD, CompE](&w).Removes(C[CompA]())
	posMap := NewMap1[Position](&w)
	posVelMap := NewMap2[Position, Velocity](&w)

	cnt := 1
	posMap.NewBatchFn(n, func(entity Entity, pos *Position) {
		pos.X = float64(cnt)
		cnt++
	})
	posVelMap.NewBatchFn(n, func(entity Entity, pos *Position, _ *Velocity) {
		pos.X = float64(cnt)
		cnt++
	})
	assert.Equal(t, 2*n+1, cnt)

	filter := NewFilter1[Position](&w)
	cnt = 0
	exchange.AddBatchFn(filter.Batch(), func(entity Entity, a *CompA, b *CompB, c *CompC, d *CompD, e *CompE) {
		a.X = float64(cnt)
		cnt++
	})

	filter2 := NewFilter1[CompA](&w)
	query := filter2.Query()
	cnt = 0
	for query.Next() {
		a := query.Get()
		assert.EqualValues(t, cnt, a.X)
		pos := posMap.Get(query.Entity())
		assert.Greater(t, pos.X, 0.0)
		cnt++
	}
	assert.Equal(t, 2*n, cnt)

	cnt = 0
	exchange.RemoveBatch(filter2.Batch(), func(entity Entity) {
		cnt++
	})
	assert.Equal(t, 2*n, cnt)

	query = filter2.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)
}

func TestExchange6(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap6[CompA, CompB, CompC, CompD, CompE, CompF](&w)
	ex := NewExchange6[CompA, CompB, CompC, CompD, CompE, CompF](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Exchange(e, &CompA{}, &CompB{}, &CompC{}, &CompD{}, &CompE{}, &CompF{})
	assert.False(t, posMap.HasAll(e))
	assert.True(t, mapper.HasAll(e))
}

func TestExchange6Add(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap6[CompA, CompB, CompC, CompD, CompE, CompF](&w)
	ex := NewExchange6[CompA, CompB, CompC, CompD, CompE, CompF](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Add(e, &CompA{}, &CompB{}, &CompC{}, &CompD{}, &CompE{}, &CompF{})
	assert.True(t, posMap.HasAll(e))
	assert.True(t, mapper.HasAll(e))
}

func TestExchange6Remove(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap6[CompA, CompB, CompC, CompD, CompE, CompF](&w)
	ex := NewExchange6[CompA, CompB, CompC, CompD, CompE, CompF](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Remove(e)
	assert.False(t, posMap.HasAll(e))
	assert.False(t, mapper.HasAll(e))
}

func TestExchange6AddBatch(t *testing.T) {
	n := 12
	w := NewWorld(8)

	exchange := NewExchange6[CompA, CompB, CompC, CompD, CompE, CompF](&w).Removes(C[CompA]())
	posMap := NewMap1[Position](&w)
	posVelMap := NewMap2[Position, Velocity](&w)

	cnt := 1
	posMap.NewBatchFn(n, func(entity Entity, pos *Position) {
		pos.X = float64(cnt)
		cnt++
	})
	posVelMap.NewBatchFn(n, func(entity Entity, pos *Position, _ *Velocity) {
		pos.X = float64(cnt)
		cnt++
	})
	assert.Equal(t, 2*n+1, cnt)

	filter := NewFilter1[Position](&w)
	exchange.AddBatch(filter.Batch(), &CompA{}, &CompB{}, &CompC{}, &CompD{}, &CompE{}, &CompF{})

	filter2 := NewFilter1[CompA](&w)
	query := filter2.Query()
	cnt = 0
	for query.Next() {
		pos := posMap.Get(query.Entity())
		assert.Greater(t, pos.X, 0.0)
		cnt++
	}
	assert.Equal(t, 2*n, cnt)

	exchange.RemoveBatch(filter2.Batch(), nil)

	query = filter2.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)
}

func TestExchange6AddBatchFn(t *testing.T) {
	n := 12
	w := NewWorld(8)

	exchange := NewExchange6[CompA, CompB, CompC, CompD, CompE, CompF](&w).Removes(C[CompA]())
	posMap := NewMap1[Position](&w)
	posVelMap := NewMap2[Position, Velocity](&w)

	cnt := 1
	posMap.NewBatchFn(n, func(entity Entity, pos *Position) {
		pos.X = float64(cnt)
		cnt++
	})
	posVelMap.NewBatchFn(n, func(entity Entity, pos *Position, _ *Velocity) {
		pos.X = float64(cnt)
		cnt++
	})
	assert.Equal(t, 2*n+1, cnt)

	filter := NewFilter1[Position](&w)
	cnt = 0
	exchange.AddBatchFn(filter.Batch(), func(entity Entity, a *CompA, b *CompB, c *CompC, d *CompD, e *CompE, f *CompF) {
		a.X = float64(cnt)
		cnt++
	})

	filter2 := NewFilter1[CompA](&w)
	query := filter2.Query()
	cnt = 0
	for query.Next() {
		a := query.Get()
		assert.EqualValues(t, cnt, a.X)
		pos := posMap.Get(query.Entity())
		assert.Greater(t, pos.X, 0.0)
		cnt++
	}
	assert.Equal(t, 2*n, cnt)

	cnt = 0
	exchange.RemoveBatch(filter2.Batch(), func(entity Entity) {
		cnt++
	})
	assert.Equal(t, 2*n, cnt)

	query = filter2.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)
}

func TestExchange7(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap7[CompA, CompB, CompC, CompD, CompE, CompF, CompG](&w)
	ex := NewExchange7[CompA, CompB, CompC, CompD, CompE, CompF, CompG](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Exchange(e, &CompA{}, &CompB{}, &CompC{}, &CompD{}, &CompE{}, &CompF{}, &CompG{})
	assert.False(t, posMap.HasAll(e))
	assert.True(t, mapper.HasAll(e))
}

func TestExchange7Add(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap7[CompA, CompB, CompC, CompD, CompE, CompF, CompG](&w)
	ex := NewExchange7[CompA, CompB, CompC, CompD, CompE, CompF, CompG](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Add(e, &CompA{}, &CompB{}, &CompC{}, &CompD{}, &CompE{}, &CompF{}, &CompG{})
	assert.True(t, posMap.HasAll(e))
	assert.True(t, mapper.HasAll(e))
}

func TestExchange7Remove(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap7[CompA, CompB, CompC, CompD, CompE, CompF, CompG](&w)
	ex := NewExchange7[CompA, CompB, CompC, CompD, CompE, CompF, CompG](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Remove(e)
	assert.False(t, posMap.HasAll(e))
	assert.False(t, mapper.HasAll(e))
}

func TestExchange7AddBatch(t *testing.T) {
	n := 12
	w := NewWorld(8)

	exchange := NewExchange7[CompA, CompB, CompC, CompD, CompE, CompF, CompG](&w).Removes(C[CompA]())
	posMap := NewMap1[Position](&w)
	posVelMap := NewMap2[Position, Velocity](&w)

	cnt := 1
	posMap.NewBatchFn(n, func(entity Entity, pos *Position) {
		pos.X = float64(cnt)
		cnt++
	})
	posVelMap.NewBatchFn(n, func(entity Entity, pos *Position, _ *Velocity) {
		pos.X = float64(cnt)
		cnt++
	})
	assert.Equal(t, 2*n+1, cnt)

	filter := NewFilter1[Position](&w)
	exchange.AddBatch(filter.Batch(), &CompA{}, &CompB{}, &CompC{}, &CompD{}, &CompE{}, &CompF{}, &CompG{})

	filter2 := NewFilter1[CompA](&w)
	query := filter2.Query()
	cnt = 0
	for query.Next() {
		pos := posMap.Get(query.Entity())
		assert.Greater(t, pos.X, 0.0)
		cnt++
	}
	assert.Equal(t, 2*n, cnt)

	exchange.RemoveBatch(filter2.Batch(), nil)

	query = filter2.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)
}

func TestExchange7AddBatchFn(t *testing.T) {
	n := 12
	w := NewWorld(8)

	exchange := NewExchange7[CompA, CompB, CompC, CompD, CompE, CompF, CompG](&w).Removes(C[CompA]())
	posMap := NewMap1[Position](&w)
	posVelMap := NewMap2[Position, Velocity](&w)

	cnt := 1
	posMap.NewBatchFn(n, func(entity Entity, pos *Position) {
		pos.X = float64(cnt)
		cnt++
	})
	posVelMap.NewBatchFn(n, func(entity Entity, pos *Position, _ *Velocity) {
		pos.X = float64(cnt)
		cnt++
	})
	assert.Equal(t, 2*n+1, cnt)

	filter := NewFilter1[Position](&w)
	cnt = 0
	exchange.AddBatchFn(filter.Batch(), func(entity Entity, a *CompA, b *CompB, c *CompC, d *CompD, e *CompE, f *CompF, g *CompG) {
		a.X = float64(cnt)
		cnt++
	})

	filter2 := NewFilter1[CompA](&w)
	query := filter2.Query()
	cnt = 0
	for query.Next() {
		a := query.Get()
		assert.EqualValues(t, cnt, a.X)
		pos := posMap.Get(query.Entity())
		assert.Greater(t, pos.X, 0.0)
		cnt++
	}
	assert.Equal(t, 2*n, cnt)

	cnt = 0
	exchange.RemoveBatch(filter2.Batch(), func(entity Entity) {
		cnt++
	})
	assert.Equal(t, 2*n, cnt)

	query = filter2.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)
}

func TestExchange8(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap8[CompA, CompB, CompC, CompD, CompE, CompF, CompG, CompH](&w)
	ex := NewExchange8[CompA, CompB, CompC, CompD, CompE, CompF, CompG, CompH](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Exchange(e, &CompA{}, &CompB{}, &CompC{}, &CompD{}, &CompE{}, &CompF{}, &CompG{}, &CompH{})
	assert.False(t, posMap.HasAll(e))
	assert.True(t, mapper.HasAll(e))
}

func TestExchange8Add(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap8[CompA, CompB, CompC, CompD, CompE, CompF, CompG, CompH](&w)
	ex := NewExchange8[CompA, CompB, CompC, CompD, CompE, CompF, CompG, CompH](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Add(e, &CompA{}, &CompB{}, &CompC{}, &CompD{}, &CompE{}, &CompF{}, &CompG{}, &CompH{})
	assert.True(t, posMap.HasAll(e))
	assert.True(t, mapper.HasAll(e))
}

func TestExchange8Remove(t *testing.T) {
	w := NewWorld(16)

	posMap := NewMap2[Position, Velocity](&w)
	mapper := NewMap8[CompA, CompB, CompC, CompD, CompE, CompF, CompG, CompH](&w)
	ex := NewExchange8[CompA, CompB, CompC, CompD, CompE, CompF, CompG, CompH](&w).Removes(C[Velocity](), C[Position]())

	e := posMap.NewEntity(&Position{}, &Velocity{})

	ex.Remove(e)
	assert.False(t, posMap.HasAll(e))
	assert.False(t, mapper.HasAll(e))
}

func TestExchange8AddBatch(t *testing.T) {
	n := 12
	w := NewWorld(8)

	exchange := NewExchange8[CompA, CompB, CompC, CompD, CompE, CompF, CompG, CompH](&w).Removes(C[CompA]())
	posMap := NewMap1[Position](&w)
	posVelMap := NewMap2[Position, Velocity](&w)

	cnt := 1
	posMap.NewBatchFn(n, func(entity Entity, pos *Position) {
		pos.X = float64(cnt)
		cnt++
	})
	posVelMap.NewBatchFn(n, func(entity Entity, pos *Position, _ *Velocity) {
		pos.X = float64(cnt)
		cnt++
	})
	assert.Equal(t, 2*n+1, cnt)

	filter := NewFilter1[Position](&w)
	exchange.AddBatch(filter.Batch(), &CompA{}, &CompB{}, &CompC{}, &CompD{}, &CompE{}, &CompF{}, &CompG{}, &CompH{})

	filter2 := NewFilter1[CompA](&w)
	query := filter2.Query()
	cnt = 0
	for query.Next() {
		pos := posMap.Get(query.Entity())
		assert.Greater(t, pos.X, 0.0)
		cnt++
	}
	assert.Equal(t, 2*n, cnt)

	exchange.RemoveBatch(filter2.Batch(), nil)

	query = filter2.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)
}

func TestExchange8AddBatchFn(t *testing.T) {
	n := 12
	w := NewWorld(8)

	exchange := NewExchange8[CompA, CompB, CompC, CompD, CompE, CompF, CompG, CompH](&w).Removes(C[CompA]())
	posMap := NewMap1[Position](&w)
	posVelMap := NewMap2[Position, Velocity](&w)

	cnt := 1
	posMap.NewBatchFn(n, func(entity Entity, pos *Position) {
		pos.X = float64(cnt)
		cnt++
	})
	posVelMap.NewBatchFn(n, func(entity Entity, pos *Position, _ *Velocity) {
		pos.X = float64(cnt)
		cnt++
	})
	assert.Equal(t, 2*n+1, cnt)

	filter := NewFilter1[Position](&w)
	cnt = 0
	exchange.AddBatchFn(filter.Batch(), func(entity Entity, a *CompA, b *CompB, c *CompC, d *CompD, e *CompE, f *CompF, g *CompG, h *CompH) {
		a.X = float64(cnt)
		cnt++
	})

	filter2 := NewFilter1[CompA](&w)
	query := filter2.Query()
	cnt = 0
	for query.Next() {
		a := query.Get()
		assert.EqualValues(t, cnt, a.X)
		pos := posMap.Get(query.Entity())
		assert.Greater(t, pos.X, 0.0)
		cnt++
	}
	assert.Equal(t, 2*n, cnt)

	cnt = 0
	exchange.RemoveBatch(filter2.Batch(), func(entity Entity) {
		cnt++
	})
	assert.Equal(t, 2*n, cnt)

	query = filter2.Query()
	cnt = 0
	for query.Next() {
		cnt++
	}
	assert.Equal(t, 0, cnt)
}
