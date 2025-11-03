package geom

const Sqrt1 = 0.707

type SlopeType uint8

const (
	SlopeNone SlopeType = iota
	SlopeFlat
	SlopeAcending
	SlopeDecending
)

// Triangle represents a right-angled triangle defined by three points.
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
	t.minX, t.minY, t.maxX, t.maxY = computeAABB(points[:])
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
