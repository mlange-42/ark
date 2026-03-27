package main

// Position component.
type Position struct {
	X float64
	Y float64
}

// Velocity component.
type Velocity struct {
	X float64
	Y float64
}

// Payload32B component.
type Payload32B struct {
	X float64
	Y float64
	U float64
	V float64
}

// Payload64B component.
type Payload64B struct {
	X float64
	Y float64
	U float64
	V float64
	A float64
	B float64
	C float64
	D float64
}

// Payload128B component.
type Payload128B struct {
	X float64
	Y float64
	U float64
	V float64
	A float64
	B float64
	C float64
	D float64
	E float64
	F float64
	G float64
	H float64
	I float64
	J float64
	K float64
	L float64
}

// Aos32Byte entity.
type Aos32Byte struct {
	Pos Position
	Vel Velocity
}

// Aos64Byte entity.
type Aos64Byte struct {
	Pos Position
	Vel Velocity
	P32 Payload32B
}

// Aos128Byte entity.
type Aos128Byte struct {
	Pos Position
	Vel Velocity
	P32 Payload32B
	P64 Payload64B
}

// Aos256Byte entity.
type Aos256Byte struct {
	Pos  Position
	Vel  Velocity
	P32  Payload32B
	P64  Payload64B
	P128 Payload128B
}
