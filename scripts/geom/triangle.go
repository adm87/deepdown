package geom

import (
	"fmt"
	"math"
)

const Sqrt1 = 0.707

type SlopeType uint8

const (
	SlopeNone      SlopeType = iota
	SlopeDownLeft            // \ - descending from left to right
	SlopeUpLeft              // / - ascending from left to right
	SlopeDownRight           // / - descending from right to left (visually same as UpLeft)
	SlopeUpRight             // \ - ascending from right to left (visually same as DownLeft)
)

type Triangle struct {
	X, Y       float32
	Nx, Ny     float32
	minX, minY float32
	maxX, maxY float32
	points     [6]float32
	slopeType  SlopeType
}

func NewTriangle(x, y float32, points [6]float32) Triangle {
	t := Triangle{
		X: x,
		Y: y,
	}
	t.SetPoints(points)
	return t
}

func (t *Triangle) SetPoints(points [6]float32) {
	if err := validateTrianglePoints(points); err != nil {
		panic(err)
	}

	t.minX, t.minY, t.maxX, t.maxY = computeAABB(points[:])
	t.slopeType = classifySlope45Degree(points)
	nx, ny := computeSlopeNormal(t.slopeType)
	t.Nx, t.Ny = float32(nx), float32(ny)
	t.points = points
}

func (t *Triangle) SlopeType() SlopeType {
	return t.slopeType
}

func (t *Triangle) Points() []float32 {
	return t.points[:]
}

func (t *Triangle) GetVertex(i int) (x, y float32) {
	i *= 2
	return t.X + t.points[i], t.Y + t.points[i+1]
}

func (t *Triangle) ContainsPoint(px, py float32) bool {
	x1, y1 := t.GetVertex(0)
	x2, y2 := t.GetVertex(1)
	x3, y3 := t.GetVertex(2)

	// Calculate cross products for each edge
	d1 := (px-x2)*(y1-y2) - (x1-x2)*(py-y2)
	d2 := (px-x3)*(y2-y3) - (x2-x3)*(py-y3)
	d3 := (px-x1)*(y3-y1) - (x3-x1)*(py-y1)

	// Check if all have same sign
	hasNeg := (d1 < 0) || (d2 < 0) || (d3 < 0)
	hasPos := (d1 > 0) || (d2 > 0) || (d3 > 0)

	return !(hasNeg && hasPos)
}

// ========== AABB interface ==========

func (t *Triangle) Min() (x, y float32) {
	return t.X + t.minX, t.Y + t.minY
}

func (t *Triangle) Max() (x, y float32) {
	return t.X + t.maxX, t.Y + t.maxY
}

func (t *Triangle) IntersectsAABB(minX, minY, maxX, maxY float32) bool {
	if t.X+t.maxX < minX || t.X+t.minX > maxX || t.Y+t.maxY < minY || t.Y+t.minY > maxY {
		return false
	}

	for i := range 3 {
		vx, vy := t.GetVertex(i)
		if vx >= minX && vx <= maxX && vy >= minY && vy <= maxY {
			return true
		}
	}

	corners := []float32{minX, minY, maxX, minY, maxX, maxY, minX, maxY}
	for i := 0; i < 8; i += 2 {
		if t.ContainsPoint(corners[i], corners[i+1]) {
			return true
		}
	}

	return false
}

// ========== AABB interface ==========

func validateTrianglePoints(points [6]float32) error {
	if len(points) != 6 {
		return fmt.Errorf("triangle requires exactly 6 values (3 points with x,y), got %d", len(points))
	}

	// Extract the three points
	p1 := [2]float32{points[0], points[1]}
	p2 := [2]float32{points[2], points[3]}
	p3 := [2]float32{points[4], points[5]}

	if pointsEqual(p1, p2) || pointsEqual(p2, p3) || pointsEqual(p1, p3) {
		return fmt.Errorf("triangle has duplicate points, cannot form valid triangle")
	}

	if areCollinear(p1, p2, p3) {
		return fmt.Errorf("triangle points are collinear, cannot form valid triangle")
	}

	area := calculateTriangleArea(p1, p2, p3)
	if area < 0.001 {
		return fmt.Errorf("triangle area too small (%f), likely degenerate", area)
	}

	const maxTriangleSize = 10000.0
	if area > maxTriangleSize {
		return fmt.Errorf("triangle area too large (%f), exceeds maximum (%f)", area, maxTriangleSize)
	}

	return nil
}

func classifySlope45Degree(points [6]float32) SlopeType {
	for i := 0; i < 6; i += 2 {
		j := (i + 2) % 6

		dx := points[j] - points[i]
		dy := points[j+1] - points[i+1]

		absDx := float32(math.Abs(float64(dx)))
		absDy := float32(math.Abs(float64(dy)))

		if math.Abs(float64(absDx-absDy)) < 0.1 {
			if dx > 0 && dy < 0 {
				return SlopeDownLeft
			} else if dx > 0 && dy > 0 {
				return SlopeUpLeft
			} else if dx < 0 && dy < 0 {
				return SlopeUpRight
			} else if dx < 0 && dy > 0 {
				return SlopeDownRight
			}
		}
	}

	return SlopeNone
}

func computeSlopeNormal(slopeType SlopeType) (nx, ny float64) {
	switch slopeType {
	case SlopeDownLeft:
		return Sqrt1, Sqrt1 // \ slope
	case SlopeUpLeft:
		return -Sqrt1, Sqrt1 // / slope
	case SlopeDownRight:
		return -Sqrt1, -Sqrt1 // / slope
	case SlopeUpRight:
		return Sqrt1, -Sqrt1 // \ slope
	default:
		return 0.0, 0.0
	}
}

func pointsEqual(p1, p2 [2]float32) bool {
	const tolerance = 0.001
	return math.Abs(float64(p1[0]-p2[0])) < tolerance &&
		math.Abs(float64(p1[1]-p2[1])) < tolerance
}

func areCollinear(p1, p2, p3 [2]float32) bool {
	v1x := p2[0] - p1[0]
	v1y := p2[1] - p1[1]
	v2x := p3[0] - p1[0]
	v2y := p3[1] - p1[1]

	crossProduct := v1x*v2y - v1y*v2x

	const tolerance = 0.001
	return math.Abs(float64(crossProduct)) < tolerance
}

func calculateTriangleArea(p1, p2, p3 [2]float32) float32 {
	v1x := p2[0] - p1[0]
	v1y := p2[1] - p1[1]
	v2x := p3[0] - p1[0]
	v2y := p3[1] - p1[1]

	crossProduct := v1x*v2y - v1y*v2x
	return float32(math.Abs(float64(crossProduct))) * 0.5
}
