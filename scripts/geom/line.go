package geom

type Line struct {
	X1, Y1 float32
	X2, Y2 float32
}

func NewLine(x1, y1, x2, y2 float32) Line {
	return Line{X1: x1, Y1: y1, X2: x2, Y2: y2}
}

// ========== AABB interface ==========

func (l *Line) Min() (x, y float32) {
	if l.X1 < l.X2 {
		x = l.X1
	} else {
		x = l.X2
	}
	if l.Y1 < l.Y2 {
		y = l.Y1
	} else {
		y = l.Y2
	}
	return
}

func (l *Line) Max() (x, y float32) {
	if l.X1 > l.X2 {
		x = l.X1
	} else {
		x = l.X2
	}
	if l.Y1 > l.Y2 {
		y = l.Y1
	} else {
		y = l.Y2
	}
	return
}
