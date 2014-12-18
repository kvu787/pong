package lib

type Vector_s struct {
	X,Y float64
}

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

func DegreesToRadians(deg float64) float64 {
	return deg / 360.0 * (2 * math.Pi)
