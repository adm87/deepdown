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

// ========== AABB interface ==========

func (p *Polygon) Min() (x, y float32) {
	if len(p.points) < 2 {
		return p.X, p.Y
	}
	return p.X + p.minX, p.Y + p.minY
}

func (p *Polygon) Max() (x, y float32) {
	if len(p.points) < 2 {
		return p.X, p.Y
	}
	return p.X + p.maxX, p.Y + p.maxY
}

func computeAABB(points []float32) (minX, minY, maxX, maxY float32) {
	if len(points)%2 != 0 {
		panic("points length must be even")
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

// ========== AABB interface ==========
