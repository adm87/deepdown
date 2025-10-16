package debug

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawPolygon(screen *ebiten.Image, startX, startY float32, points []float32, col color.RGBA, view ebiten.GeoM) {
	if len(points) < 4 {
		return
	}

	path := vector.Path{}
	for i := 0; i < len(points); i += 2 {
		x, y := view.Apply(float64(startX+points[i]), float64(startY+points[i+1]))
		if i == 0 {
			path.MoveTo(float32(x), float32(y))
		} else {
			path.LineTo(float32(x), float32(y))
		}
	}
	path.Close()

	do := &vector.DrawPathOptions{}
	do.ColorScale.ScaleWithColor(col)

	// vector.StrokePath(screen, &path, &vector.StrokeOptions{
	// 	Width: 1,
	// }, do)
	vector.FillPath(screen, &path, nil, do)
}
