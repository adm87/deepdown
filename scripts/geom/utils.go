package geom

import (
	"math"
)

func abs(value float32) float32 {
	if value < 0 {
		return -value
	}
	return value
}

func ComputeSlopeNormal(p1, p2 [2]float32) [2]float32 {
	edgeX := float64(p2[0] - p1[0])
	edgeY := float64(p2[1] - p1[1])

	nx := -edgeY
	ny := edgeX

	length := math.Sqrt(nx*nx + ny*ny)
	if length == 0 {
		return [2]float32{0, 0}
	}

	return [2]float32{-float32(nx / length), -float32(ny / length)}
}

func FindSlope(points [6]float32) ([2][2]float32, [2]float32, SlopeType) {
	// Find the longest edge among the three points
	var slope [2][2]float32
	longEdgeIdx := 0
	longestDist := float32(0)

	for i := range 3 {
		j := (i + 1) % 3
		p1 := [2]float32{points[i*2], points[i*2+1]}
		p2 := [2]float32{points[j*2], points[j*2+1]}
		dist := Distance(p1, p2)
		if dist > longestDist {
			longestDist = dist
			longEdgeIdx = i
		}
	}

	i := longEdgeIdx
	j := (i + 1) % 3
	k := (i + 2) % 3

	p1 := [2]float32{points[i*2], points[i*2+1]}
	p2 := [2]float32{points[j*2], points[j*2+1]}
	corner := [2]float32{points[k*2], points[k*2+1]}

	slope[0], slope[1] = p1, p2

	var slopeType SlopeType
	switch {
	case p1[1] == p2[1]:
		slopeType = SlopeFlat
	case (p1[0] < p2[0] && p1[1] < p2[1]) || (p1[0] > p2[0] && p1[1] > p2[1]):
		slopeType = SlopeAcending
	default:
		slopeType = SlopeDecending
	}

	return slope, corner, slopeType
}

func FindTriangleSurfaceAt(x float32, triangle *Triangle) (surfaceY float32, found bool) {
	points := triangle.Points()

	// Test each edge
	for i := range 3 {
		x1, y1 := triangle.X+points[i*2], triangle.Y+points[i*2+1]
		x2, y2 := triangle.X+points[((i+1)%3)*2], triangle.Y+points[((i+1)%3)*2+1]

		// Skip vertical or horizontal edges
		if abs(x2-x1) < 0.1 || abs(y2-y1) < 0.1 {
			continue
		}

		// Check if edge spans x coordinate
		if x >= min(x1, x2) && x <= max(x1, x2) {
			t := (x - x1) / (x2 - x1)
			return y1 + t*(y2-y1), true
		}
	}

	return 0, false
}

func Distance(p1, p2 [2]float32) float32 {
	dx := p2[0] - p1[0]
	dy := p2[1] - p1[1]
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

func LineIntersects(aLine [2][2]float32, bLine [2][2]float32) ([2]float32, bool) {
	x1, y1 := aLine[0][0], aLine[0][1]
	x2, y2 := aLine[1][0], aLine[1][1]
	x3, y3 := bLine[0][0], bLine[0][1]
	x4, y4 := bLine[1][0], bLine[1][1]

	denominator := (y4-y3)*(x2-x1) - (x4-x3)*(y2-y1)

	// Check if lines are parallel
	if denominator == 0 {
		return [2]float32{}, false
	}

	uA := ((x4-x3)*(y1-y3) - (y4-y3)*(x1-x3)) / denominator
	uB := ((x2-x1)*(y1-y3) - (y2-y1)*(x1-x3)) / denominator

	if uA >= 0 && uA <= 1 && uB >= 0 && uB <= 1 {
		return [2]float32{(x1 + uA*(x2-x1)), (y1 + uA*(y2-y1))}, true
	}

	return [2]float32{}, false
}
