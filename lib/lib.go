package lib

func IsPointInsideRectangle1(px, py, rx, ry, rw, rh float64) bool {
	var (
		minX float64 = rx
		maxX float64 = rx + rw
		minY float64 = ry
		maxY float64 = ry + rh
	)

	return minX <= px && px <= maxX && minY <= py && py <= maxY
}

func IsPointInsideCircle(px, py, cx, cy, rad float64) bool {
	point := Vector_s{px, py}
	center := Vector_s{cx, cy}
	return VectorMagnitude(VectorSub(point, center)) <= rad
}
