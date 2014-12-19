package main

import (
	"math"
	"math/rand"
	"time"

	sf "bitbucket.org/kvu787/gosfml2"
)

/**************/
/*** VECTOR ***/
/**************/

func NewPolar(magnitude float64, radians float64) Vector_s {
	var unitVector Vector_s = Vector_s{X: math.Cos(radians), Y: math.Sin(radians)}
	return VectorMul(magnitude, unitVector)
}

func VectorAdd(v1, v2 Vector_s) Vector_s {
	return Vector_s{v1.X + v2.X, v1.Y + v2.Y}
}

func VectorMul(k float64, v Vector_s) Vector_s {
	return Vector_s{k * v.X, k * v.Y}
}

func VectorSub(v1, v2 Vector_s) Vector_s {
	return VectorAdd(v1, VectorMul(-1, v2))
}

func VectorDiv(k float64, v Vector_s) Vector_s {
	return VectorMul(1.0/k, v)
}

func VectorAngle(v Vector_s) float64 {
	return math.Atan2(v.Y, v.X)
}

func VectorMagnitude(v Vector_s) float64 {
	return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2))
}

func VectorDotProduct(v1, v2 Vector_s) float64 {
	return v1.X*v2.X + v1.Y*v2.Y
}

func VectorProjection(v1, v2 Vector_s) Vector_s {
	coefficient := VectorDotProduct(v1, v2) / math.Pow(VectorMagnitude(v2), 2)
	return VectorMul(coefficient, v2)
}

func VectorRejection(v1, v2 Vector_s) Vector_s {
	return VectorSub(v1, VectorProjection(v1, v2))
}

func VectorRotate(v Vector_s, radians float64) Vector_s {
	var currentRadians float64 = VectorAngle(v)
	var newRadians float64 = currentRadians + radians
	var newUnitVector Vector_s = Vector_s{X: math.Cos(newRadians), Y: math.Sin(newRadians)}
	return VectorMul(VectorMagnitude(v), newUnitVector)
}

func VectorScale(magnitude float64, v Vector_s) Vector_s {
	return VectorMul(magnitude/VectorMagnitude(v), v)
}

func VectorCrossProduct(v1 Vector_s, v2 Vector_s) float64 {
	return v1.X*v2.Y - v1.Y*v2.X
}

func VectorNormalize(v Vector_s) Vector_s {
	return VectorScale(1, v)
}

func DegreesToRadians(deg float64) float64 {
	return deg / 360.0 * (2 * math.Pi)
}

/*********************/
/*** INTERSECTIONS ***/
/*********************/

func AreRectangleSegmentIntersecting(r Rectangle_s, s Segment_s) bool {
	var segments [4]Segment_s = RectangleSegments(r)
	var i int

	for i = 0; i < 4; i++ {
		if AreSegmentsIntersecting(s, segments[i]) {
			return true
		}
	}
	return false
}

func AreSegmentsIntersecting(s1 Segment_s, s2 Segment_s) bool {
	// check if each segment's endpoints are on opposite sides
	// use cross product to check 'sidedness'

	var a Vector_s // segment vector
	var b Vector_s // endpoint vector 1
	var c Vector_s // endpoint vector 2
	var xp1 float64
	var xp2 float64

	// use s1 as line
	a = VectorSub(s1.End, s1.Start)
	b = VectorSub(s2.Start, s1.Start)
	c = VectorSub(s2.End, s1.Start)
	xp1 = VectorCrossProduct(a, b)
	xp2 = VectorCrossProduct(a, c)

	if (xp1 < 0) == (xp2 < 0) {
		return false
	}

	// use s2 as line
	a = VectorSub(s2.End, s2.Start)
	b = VectorSub(s1.Start, s2.Start)
	c = VectorSub(s1.End, s2.Start)
	xp1 = VectorCrossProduct(a, b)
	xp2 = VectorCrossProduct(a, c)

	return (xp1 < 0) != (xp2 < 0)
}

func AreRectanglesIntersecting(r1 Rectangle_s, r2 Rectangle_s) bool {
	// assume position is rectangle center

	if IsPointInsideRectangle(r1.Position, r2) || IsPointInsideRectangle(r2.Position, r1) {
		return true
	}

	var segments1 [4]Segment_s = RectangleSegments(r1)
	var segments2 [4]Segment_s = RectangleSegments(r2)
	var i int
	var j int
	for i = 0; i < 4; i++ {
		for j = 0; j < 4; j++ {
			if AreSegmentsIntersecting(segments1[i], segments2[j]) {
				return true
			}
		}
	}

	return false
}

func AreRectangleCircleIntersecting(r Rectangle_s, c Circle_s) bool {
	var segments [4]Segment_s = RectangleSegments(r)
	var i int

	if IsPointInsideRectangle(c.Position, r) {
		return true
	} else {
		for i = 0; i < 4; i++ {
			if AreCircleSegmentIntersecting(c, segments[i]) {
				return true
			}
		}
		return false
	}
}

func AreCircleLineIntersecting(c Circle_s, s Segment_s) bool {
	var a Vector_s = VectorSub(c.Position, s.Start)
	var b Vector_s = VectorSub(s.End, s.Start)
	var rejection Vector_s = VectorRejection(a, b)
	return VectorMagnitude(rejection) < c.Radius
}

func AreCircleSegmentIntersecting(c Circle_s, s Segment_s) bool {
	if !AreCircleLineIntersecting(c, s) {
		return false
	}
	var a Vector_s = VectorSub(c.Position, s.Start)
	var b Vector_s = VectorSub(s.End, s.Start)
	var proj Vector_s = VectorProjection(a, b)
	var m Vector_s = VectorAdd(s.Start, proj)
	var isPointOnSegment bool = ((m.X < s.Start.X) != (m.X < s.End.X)) || ((m.Y < s.Start.Y) != (m.Y < s.End.Y))
	return isPointOnSegment
}

func IsPointInsideRectangle(v Vector_s, r Rectangle_s) bool {
	var corners [4]Vector_s = RectangleCorners(r)

	var quadrant2Corner Vector_s = corners[1]
	var quadrant4Corner Vector_s = corners[3]

	var quadrant12Bound float64 = quadrant2Corner.X
	var quadrant23Bound float64 = quadrant2Corner.Y
	var quadrant34Bound float64 = quadrant4Corner.X
	var quadrant41Bound float64 = quadrant4Corner.Y

	return v.X < quadrant12Bound && v.Y > quadrant23Bound && v.X > quadrant34Bound && v.Y < quadrant41Bound
}

/************/
/*** MISC ***/
/************/

func GenerateRandomBallDirection() float64 {
	var spread float64 = DegreesToRadians(75)
	var a float64 = spread * (-0.5 + rand.Float64())
	if rand.Float64() > 0.5 {
		return a + math.Pi
	} else {
		return a
	}
}

func NormalizeScalar(k float64) float64 {
	return k / math.Abs(k)
}

func SegmentLength(s Segment_s) float64 {
	return VectorMagnitude(VectorSub(s.End, s.Start))
}

func Floor(x float64) int {
	return int(x)
}

func Clamp(min float64, val float64, max float64) float64 {
	if val < min {
		return min
	} else if val > max {
		return max
	} else {
		return val
	}
}

func RectangleSegments(r Rectangle_s) [4]Segment_s {
	var corners [4]Vector_s = RectangleCorners(r)

	return [4]Segment_s{
		Segment_s{Start: corners[0], End: corners[1]},
		Segment_s{Start: corners[1], End: corners[2]},
		Segment_s{Start: corners[2], End: corners[3]},
		Segment_s{Start: corners[3], End: corners[0]},
	}
}

// ordered by quadrant
func RectangleCorners(r Rectangle_s) [4]Vector_s {
	var pos Vector_s = r.Position
	var height float64 = r.Height
	var width float64 = r.Width

	var corners [4]Vector_s = [4]Vector_s{
		VectorAdd(pos, Vector_s{0.5 * width, 0.5 * height}),
		VectorAdd(pos, Vector_s{0.5 * width, -0.5 * height}),
		VectorAdd(pos, Vector_s{-0.5 * width, -0.5 * height}),
		VectorAdd(pos, Vector_s{-0.5 * width, 0.5 * height}),
	}

	return corners
}

func SegmentDirection(s Segment_s) Vector_s {
	return VectorSub(s.End, s.Start)
}

func PerpendicularVectorFromLineToPoint(segment Segment_s, point Vector_s) Vector_s {
	var a Vector_s = VectorSub(point, segment.Start)
	var b Vector_s = VectorSub(segment.End, segment.Start)
	return VectorRejection(a, b)
}

func SecondsToDuration(seconds float64) time.Duration {
	var nanoseconds float64 = seconds * 1e9
	return time.Duration(nanoseconds)
}

func Reflect(incident Vector_s, mirror Vector_s) Vector_s {
	var proj Vector_s = VectorProjection(incident, mirror)
	var rej Vector_s = VectorRejection(incident, mirror)
	return VectorAdd(proj, VectorMul(-1, rej))
}

/*********************/
/*** SFML SPECIFIC ***/
/*********************/

func SetupWindow(width int, height int) *sf.RenderWindow {
	return sf.NewRenderWindow(
		sf.VideoMode{Width: uint(width), Height: uint(height), BitsPerPixel: 32},
		"pong",
		sf.StyleDefault,
		sf.DefaultContextSettings())
}

func ClearWindow(color Color_s, window *sf.RenderWindow) {
	window.Clear(ToSFMLColor(color))
}

func DisplayWindow(window *sf.RenderWindow) {
	window.Display()
}

func RenderRectangle(rectangle Rectangle_s, window *sf.RenderWindow) {
	r, err := sf.NewRectangleShape()
	if err != nil {
		panic(err)
	}
	var width float32 = float32(rectangle.Width)
	var height float32 = float32(rectangle.Height)
	r.SetSize(sf.Vector2f{width, height})
	r.SetOrigin(sf.Vector2f{width / 2, height / 2})
	r.SetRotation(0)
	r.SetOutlineThickness(0)
	r.SetFillColor(sf.ColorWhite())
	r.SetPosition(ToSFML2Vector(rectangle.Position))
	window.Draw(r, sf.DefaultRenderStates())
}

func RenderCircle(circle Circle_s, window *sf.RenderWindow) {
	c, err := sf.NewCircleShape()
	if err != nil {
		panic(err)
	}
	r32 := float32(circle.Radius)
	c.SetRadius(r32)
	c.SetOrigin(sf.Vector2f{r32, r32})
	c.SetRotation(0)
	c.SetOutlineThickness(0)
	c.SetFillColor(sf.ColorWhite())
	c.SetPosition(ToSFML2Vector(circle.Position))
	window.Draw(c, sf.DefaultRenderStates())
}

func ToSFML2Vector(v Vector_s) sf.Vector2f {
	return sf.Vector2f{float32(v.X), float32(v.Y)}
}

func ToSFMLColor(c Color_s) sf.Color {
	return sf.Color{
		R: c.R,
		G: c.G,
		B: c.B,
		A: c.A,
	}
}
