package main

type Input_s struct {
	IsKeyUpPressed   bool
	IsKeyDownPressed bool
	IsWindowClosed   bool
}

type Paddle_s struct {
	Rectangle Rectangle_s
	Speed     float64
}

type Rectangle_s struct {
	Position Vector_s
	Width    float64
	Height   float64
}

type Vector_s struct {
	X float64
	Y float64
}
