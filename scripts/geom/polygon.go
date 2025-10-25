package geom

type Polygon struct {
	X, Y       float32   // Positional offset
	minX, minY float32   // Minimum extent (relative to X,Y)
	maxX, maxY float32   // Maximum extent (relative to X,Y)
	points     []float32 // Flat array of x,y pairs: [x0, y0, x1, y1, x2, y2, ...]
}

func NewPolygon(x, y float32, points []float32) Polygon {
	if len(points)%2 != 0 {
		panic("points length must be even")
	}
	minX, minY, maxX, maxY := computeAABB(points)
	return Polygon{
		X:      x,
		Y:      y,
		minX:   minX,
		minY:   minY,
		maxX:   maxX,
		maxY:   maxY,
		points: points,
	}
}

func (p *Polygon) VertexCount() int {
	return len(p.points) / 2
}

func (p *Polygon) GetVertex(i int) (x, y float32) {
	i *= 2
	return p.X + p.points[i], p.Y + p.points[i+1]
}

func (p *Polygon) IntersectsAABB(minX, minY, maxX, maxY float32) bool {
	// AABB
	if p.X+p.maxX < minX || p.X+p.minX > maxX ||
		p.Y+p.maxY < minY || p.Y+p.minY > maxY {
		return false
	}

	// Any polygon vertex inside rectangle?
	for i := 0; i < len(p.points); i += 2 {
		x, y := p.X+p.points[i], p.Y+p.points[i+1]
		if x >= minX && x <= maxX && y >= minY && y <= maxY {
			return true
		}
	}

	// Any rectangle corner inside polygon?
	corners := [4][2]float32{
		{minX, minY},
		{maxX, minY},
		{maxX, maxY},
		{minX, maxY},
	}
	for i := range corners {
		if p.ContainsPoint(corners[i][0], corners[i][1]) {
			return true
		}
	}

	// Check for edge intersections
	rectEdges := [4][4]float32{
		{minX, minY, maxX, minY}, // bottom
		{maxX, minY, maxX, maxY}, // right
		{maxX, maxY, minX, maxY}, // top
		{minX, maxY, minX, minY}, // left
	}

	// Check each polygon edge against rectangle edges
	for i := 0; i < len(p.points); i += 2 {
		j := (i + 2) % len(p.points)
		px1, py1 := p.X+p.points[i], p.Y+p.points[i+1]
		px2, py2 := p.X+p.points[j], p.Y+p.points[j+1]

		for k := range rectEdges {
			if lineSegmentsIntersect(px1, py1, px2, py2, rectEdges[k][0], rectEdges[k][1], rectEdges[k][2], rectEdges[k][3]) {
				return true
			}
		}
	}

	return false
}

func (p *Polygon) ContainsPoint(x, y float32) bool {
	// Quick AABB rejection
	if x < p.X+p.minX || x > p.X+p.maxX || y < p.Y+p.minY || y > p.Y+p.maxY {
		return false
	}

	inside := false
	j := len(p.points) - 2
	for i := 0; i < len(p.points); i += 2 {
		xi, yi := p.X+p.points[i], p.Y+p.points[i+1]
		xj, yj := p.X+p.points[j], p.Y+p.points[j+1]
		if ((yi > y) != (yj > y)) && (x < (xj-xi)*(y-yi)/(yj-yi)+xi) {
			inside = !inside
		}
		j = i
	}
	return inside
}

// ========== AABB interface ==========

func (p *Polygon) Min() (x, y float32) {
	if len(p.points) < 6 {
		return p.X, p.Y
	}
	return p.X + p.minX, p.Y + p.minY
}

func (p *Polygon) Max() (x, y float32) {
	if len(p.points) < 6 {
		return p.X, p.Y
	}
	return p.X + p.maxX, p.Y + p.maxY
}

func computeAABB(points []float32) (minX, minY, maxX, maxY float32) {
	if len(points)%2 != 0 {
		panic("points length must be even")
	}
	if len(points) == 0 {
		return 0, 0, 0, 0
	}

	minX, minY = points[0], points[1]
	maxX, maxY = points[0], points[1]

	for i := 2; i < len(points); i += 2 {
		x, y := points[i], points[i+1]
		if x < minX {
			minX = x
		} else if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		} else if y > maxY {
			maxY = y
		}
	}

	return
}

// Helper function for line segment intersection
func lineSegmentsIntersect(x1, y1, x2, y2, x3, y3, x4, y4 float32) bool {
	// Orientation function
	orientation := func(px, py, qx, qy, rx, ry float32) int {
		val := (qy-py)*(rx-qx) - (qx-px)*(ry-qy)
		if val == 0 {
			return 0 // collinear
		}
		if val > 0 {
			return 1 // clockwise
		}
		return 2 // counterclockwise
	}

	// Check if point q lies on segment pr
	onSegment := func(px, py, qx, qy, rx, ry float32) bool {
		return qx <= max(px, rx) && qx >= min(px, rx) && qy <= max(py, ry) && qy >= min(py, ry)
	}

	o1 := orientation(x1, y1, x2, y2, x3, y3)
	o2 := orientation(x1, y1, x2, y2, x4, y4)
	o3 := orientation(x3, y3, x4, y4, x1, y1)
	o4 := orientation(x3, y3, x4, y4, x2, y2)

	// General case
	if o1 != o2 && o3 != o4 {
		return true
	}

	// Special cases for collinear points
	if o1 == 0 && onSegment(x1, y1, x3, y3, x2, y2) {
		return true
	}
	if o2 == 0 && onSegment(x1, y1, x4, y4, x2, y2) {
		return true
	}
	if o3 == 0 && onSegment(x3, y3, x1, y1, x4, y4) {
		return true
	}
	if o4 == 0 && onSegment(x3, y3, x2, y2, x4, y4) {
		return true
	}

	return false
}

// Helper functions for min/max
func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
