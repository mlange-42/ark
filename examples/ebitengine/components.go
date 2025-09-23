package main

// Settings is a global resource.
type Settings struct {
	ScreenWidth  float32
	ScreenHeight float32
	Scale        float32
	StarsCount   int
}

// From is a component holding current star coordinates.
type From struct {
	X, Y float32
}

// To is a component holding target star coordinates.
type To struct {
	X, Y float32
}

// Brightness is a component holding star brightness.
type Brightness struct {
	V float32
}
