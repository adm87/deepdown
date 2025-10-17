package geom

type Polygon struct {
	X, Y       float32   // Positional offset
	minX, minY float32   // Minimum extent (relative to X,Y)
	maxX, maxY float32   // Maximum extent (relative to X,Y)
	points     []float32 // Flat array of x,y pairs: [x0, y0, x1, y1, x2, y2, ...]
}

func NewPolygon(x, y float32, points []float32) Polygon {
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

func (p *Polygon) IsEmpty() bool {
	return len(p.points) == 0
}

func (p *Polygon) VertexCount() int {
	return len(p.points) / 2
}

func (p *Polygon) GetVertex(i int) (x, y float32) {
	idx := i * 2
	return p.points[idx], p.points[idx+1]
}

// forEachEdge iterates over polygon edges without allocation.
// The callback function receives edge coordinates and returns true to continue iteration.
// Return false from the callback to break early.
func (p *Polygon) forEachEdge(fn func(x1, y1, x2, y2 float32) bool) {
	vertCount := len(p.points) / 2
	for i := 0; i < vertCount; i++ {
		x1 := p.points[i*2]
		y1 := p.points[i*2+1]
		x2 := p.points[(i+1)%vertCount*2]
		y2 := p.points[(i+1)%vertCount*2+1]
		if !fn(x1, y1, x2, y2) {
			break
		}
	}
}

func (p *Polygon) IntersectsRect(r *Rectangle) bool {
	if len(p.points) < 6 { // Need at least 3 vertices for a valid polygon
		return false
	}

	// AABB check
	if p.X+p.minX > r.X+r.Width || p.X+p.maxX < r.X ||
		p.Y+p.minY > r.Y+r.Height || p.Y+p.maxY < r.Y {
		return false
	}

	// Check if any of the polygon's vertices are inside the rectangle
	for i := 0; i < len(p.points); i += 2 {
		vx := p.X + p.points[i]
		vy := p.Y + p.points[i+1]
		if vx >= r.X && vx <= r.X+r.Width && vy >= r.Y && vy <= r.Y+r.Height {
			return true
		}
	}

	// Check if any of the rectangle's corners are inside the polygon
	corners := [8]float32{
		0, 0, // Top-left
		r.Width, 0, // Top-right
		r.Width, r.Height, // Bottom-right
		0, r.Height, // Bottom-left
	}
	for i := 0; i < 8; i += 2 {
		cornerX := r.X + corners[i]
		cornerY := r.Y + corners[i+1]
		if p.ContainsPoint(cornerX, cornerY) {
			return true
		}
	}

	// Check for edge intersections using zero-allocation iteration
	// Rectangle edges (relative to rectangle origin)
	rectEdges := [16]float32{
		0, 0, r.Width, 0, // top
		r.Width, 0, r.Width, r.Height, // right
		r.Width, r.Height, 0, r.Height, // bottom
		0, r.Height, 0, 0, // left
	}

	intersectionFound := false
	p.forEachEdge(func(x1, y1, x2, y2 float32) bool {
		// Transform polygon edge to world coordinates
		ex1, ey1 := p.X+x1, p.Y+y1
		ex2, ey2 := p.X+x2, p.Y+y2

		for i := 0; i < 16; i += 4 {
			// Transform rectangle edge to world coordinates
			rx1, ry1 := r.X+rectEdges[i], r.Y+rectEdges[i+1]
			rx2, ry2 := r.X+rectEdges[i+2], r.Y+rectEdges[i+3]

			if linesIntersect(ex1, ey1, ex2, ey2, rx1, ry1, rx2, ry2) {
				intersectionFound = true
				return false // Early exit
			}
		}
		return true // Continue iteration
	})

	return intersectionFound
}

func (p *Polygon) IntersectsPolygon(other *Polygon) bool {
	if len(p.points) < 6 || len(other.points) < 6 { // Need at least 3 vertices for a valid polygon
		return false
	}

	// AABB check
	if p.X+p.minX > other.X+other.maxX || p.X+p.maxX < other.X+other.minX ||
		p.Y+p.minY > other.Y+other.maxY || p.Y+p.maxY < other.Y+other.minY {
		return false
	}

	// Check if any of this polygon's vertices are inside the other polygon
	for i := 0; i < len(p.points); i += 2 {
		vx := p.X + p.points[i]
		vy := p.Y + p.points[i+1]
		if other.ContainsPoint(vx, vy) {
			return true
		}
	}

	// Check if any of the other polygon's vertices are inside this polygon
	for i := 0; i < len(other.points); i += 2 {
		vx := other.X + other.points[i]
		vy := other.Y + other.points[i+1]
		if p.ContainsPoint(vx, vy) {
			return true
		}
	}

	// Check for edge intersections using zero-allocation iteration
	intersectionFound := false
	p.forEachEdge(func(x1, y1, x2, y2 float32) bool {
		// Transform first polygon's edge to world coordinates
		ex1, ey1 := p.X+x1, p.Y+y1
		ex2, ey2 := p.X+x2, p.Y+y2

		other.forEachEdge(func(ox1, oy1, ox2, oy2 float32) bool {
			// Transform second polygon's edge to world coordinates
			otherX1, otherY1 := other.X+ox1, other.Y+oy1
			otherX2, otherY2 := other.X+ox2, other.Y+oy2

			if linesIntersect(ex1, ey1, ex2, ey2, otherX1, otherY1, otherX2, otherY2) {
				intersectionFound = true
				return false // Early exit from inner loop
			}
			return true // Continue inner iteration
		})

		return !intersectionFound // Continue outer loop only if no intersection found
	})

	return intersectionFound
}

func (p *Polygon) ContainsPoint(px, py float32) bool {
	if len(p.points) < 6 { // Need at least 3 vertices for a valid polygon
		return false
	}

	// AABB check
	if px < p.X+p.minX || px > p.X+p.maxX || py < p.Y+p.minY || py > p.Y+p.maxY {
		return false
	}

	// Check if point is exactly on a vertex
	vertCount := len(p.points) / 2
	for i := range vertCount {
		vx := p.X + p.points[i*2]
		vy := p.Y + p.points[i*2+1]
		if px == vx && py == vy {
			return true // Point exactly on vertex is considered inside
		}
	}

	// Check if point is exactly on an edge
	for i := range vertCount {
		x1 := p.X + p.points[i*2]
		y1 := p.Y + p.points[i*2+1]
		x2 := p.X + p.points[(i+1)%vertCount*2]
		y2 := p.Y + p.points[(i+1)%vertCount*2+1]

		if isPointOnLineSegment(px, py, x1, y1, x2, y2) {
			return true // Point exactly on edge is considered inside
		}
	}

	// Use the standard ray casting algorithm
	inside := false
	j := vertCount - 1
	for i := range vertCount {
		ix := p.X + p.points[i*2]
		iy := p.Y + p.points[i*2+1]
		jx := p.X + p.points[j*2]
		jy := p.Y + p.points[j*2+1]

		if ((iy > py) != (jy > py)) && (px < (jx-ix)*(py-iy)/(jy-iy)+ix) {
			inside = !inside
		}
		j = i
	}
	return inside
}

func linesIntersect(x1, y1, x2, y2, x3, y3, x4, y4 float32) bool {
	denom := (y4-y3)*(x2-x1) - (x4-x3)*(y2-y1)

	// Use absolute epsilon for better performance (avoid expensive max calculations)
	const epsilon = 1e-9
	if denom > -epsilon && denom < epsilon {
		return false // Parallel or nearly parallel lines
	}

	// Compute parametric values
	ua := ((x4-x3)*(y1-y3) - (y4-y3)*(x1-x3)) / denom
	if ua < 0 || ua > 1 {
		return false // Early exit if first parameter out of range
	}

	ub := ((x2-x1)*(y1-y3) - (y2-y1)*(x1-x3)) / denom
	return ub >= 0 && ub <= 1
}

func isPointOnLineSegment(px, py, x1, y1, x2, y2 float32) bool {
	const epsilon = 1e-9

	// Check if point is within the bounding box of the line segment
	minX, maxX := x1, x2
	if minX > maxX {
		minX, maxX = maxX, minX
	}
	minY, maxY := y1, y2
	if minY > maxY {
		minY, maxY = maxY, minY
	}

	if px < minX-epsilon || px > maxX+epsilon || py < minY-epsilon || py > maxY+epsilon {
		return false
	}

	// Check if point lies on the line using cross product
	dx1, dy1 := x2-x1, y2-y1
	dx2, dy2 := px-x1, py-y1

	// Cross product should be zero if points are collinear
	cross := dx1*dy2 - dy1*dx2
	return cross > -epsilon && cross < epsilon
}

// ========== AABB interface ==========

func (p *Polygon) Min() (x, y float32) {
	if len(p.points) < 6 { // Need at least 3 vertices
		return p.X, p.Y
	}
	return p.X + p.minX, p.Y + p.minY
}

func (p *Polygon) Max() (x, y float32) {
	if len(p.points) < 6 { // Need at least 3 vertices
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
		if points[i] < minX {
			minX = points[i]
		}
		if points[i+1] < minY {
			minY = points[i+1]
		}
		if points[i] > maxX {
			maxX = points[i]
		}
		if points[i+1] > maxY {
			maxY = points[i+1]
		}
	}

	return
}
