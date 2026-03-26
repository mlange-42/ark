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

// Payload1 component.
type Payload1 struct {
	X float64
	Y float64
}

// Payload2 component.
type Payload2 struct {
	X float64
	Y float64
}

// Payload3 component.
type Payload3 struct {
	X float64
	Y float64
}

// Payload4 component.
type Payload4 struct {
	X float64
	Y float64
}

// Payload5 component.
type Payload5 struct {
	X float64
	Y float64
}

// Payload6 component.
type Payload6 struct {
	X float64
	Y float64
}

// Aos16Byte entity.
type Aos16Byte struct {
	Pos Position
	Vel Velocity
}

// Aos32Byte entity.
type Aos32Byte struct {
	Pos Position
	Vel Velocity
	P1  Payload1
	P2  Payload2
}

// Aos64Byte entity.
type Aos64Byte struct {
	Pos Position
	Vel Velocity
	P1  Payload1
	P2  Payload2
	P3  Payload3
	P4  Payload4
	P5  Payload5
	P6  Payload6
}
