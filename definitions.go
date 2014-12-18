package main

type Vector_s struct {
	X float64
	Y float64
}

type Circle_s struct {
	Position Vector_s
	Radius   float64
}

type Rectangle_s struct {
	Position Vector_s
	Width    float64
	Height   float64
}

type Paddle_s struct {
	Rectangle Rectangle_s
	Speed     float64
}

type Ball_s struct {
	Circle   Circle_s
	Velocity Vector_s
}

type Input_s struct {
	IsUpArrowClicked   bool
	IsDownArrowClicked bool
	IsWClicked         bool
	IsSClicked         bool
	IsWindowClosed     bool
}

type Segment_s struct {
	Start Vector_s
	End   Vector_s
}

type Color_s struct {
	R uint8
	G uint8
	B uint8
	A uint8
}
