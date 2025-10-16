package geom

type Rectangle struct {
	X, Y          float32
	Width, Height float32
}

func NewRectangle(x, y, width, height float32) Rectangle {
	return Rectangle{X: x, Y: y, Width: width, Height: height}
}

func (r *Rectangle) Center() (x, y float32) {
	return r.X + r.Width/2, r.Y + r.Height/2
}

func (r *Rectangle) Contains(x, y float32) bool {
	return x >= r.X && x <= r.X+r.Width && y >= r.Y && y <= r.Y+r.Height
}

func (r *Rectangle) Intersects(other *Rectangle) bool {
	return r.X < other.X+other.Width && r.X+r.Width > other.X &&
		r.Y < other.Y+other.Height && r.Y+r.Height > other.Y
}

// ========== AABB interface ==========

func (r *Rectangle) Min() (x, y float32) {
	return r.X, r.Y
}

func (r *Rectangle) Max() (x, y float32) {
	return r.X + r.Width, r.Y + r.Height
}

// ========== AABB interface ==========
