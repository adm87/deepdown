package geom

type Point struct {
	X, Y float32
}

func NewPoint(x, y float32) Point {
	return Point{X: x, Y: y}
}
