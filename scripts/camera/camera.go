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

func (c *Camera) Viewport() [4]float32 {
	return [4]float32{c.X - c.Width/2, c.Y - c.Height/2, c.X + c.Width/2, c.Y + c.Height/2}
}

func (c *Camera) Matrix() ebiten.GeoM {
	c.matrix.Reset()
	c.matrix.Translate(float64(c.Width/2-c.X), float64(c.Height/2-c.Y))
	return c.matrix
}
