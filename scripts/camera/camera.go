package camera

import (
	"github.com/adm87/deepdown/scripts/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	*geom.Rectangle

	matrix ebiten.GeoM
}

func NewCamera(x, y, width, height float32) *Camera {
	return &Camera{
		Rectangle: &geom.Rectangle{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
		},
	}
}

func (c *Camera) Viewport() (minX, minY, maxX, maxY float32) {
	minX = c.X - c.Width/2
	minY = c.Y - c.Height/2
	maxX = c.X + c.Width/2
	maxY = c.Y + c.Height/2
	return
}

func (c *Camera) Matrix() ebiten.GeoM {
	c.matrix.Reset()
	c.matrix.Translate(float64(c.Width/2-c.X), float64(c.Height/2-c.Y))
	return c.matrix
}
