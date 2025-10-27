package level

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/adm87/deepdown/scripts/assets"
	"github.com/adm87/deepdown/scripts/camera"
	"github.com/adm87/deepdown/scripts/collision"
	"github.com/adm87/deepdown/scripts/deepdown"
	"github.com/adm87/tiled"
	"github.com/adm87/tiled/tilemap"
	"github.com/adm87/utilities/hash"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Level struct {
	ctx deepdown.Context

	tilemap *tilemap.Map
	camera  *camera.Camera

	world  *collision.World
	player *collision.BoxCollider

	op ebiten.DrawImageOptions
}

func NewLevel(ctx deepdown.Context, targetWidth, targetHeight float32) *Level {
	world := collision.NewWorld(ctx)
	return &Level{
		ctx:     ctx,
		tilemap: tilemap.NewMap(),
		camera:  camera.NewCamera(0, 0, targetWidth, targetHeight),
		op:      ebiten.DrawImageOptions{},
		world:   world,
	}
}

func (l *Level) Camera() *camera.Camera {
	return l.camera
}

func (l *Level) Bounds() (minX, minY, maxX, maxY int32) {
	return 0, 0, 0, 0
}

func (l *Level) SetTmx(tmx *tiled.Tmx) error {
	l.tilemap.SetTmx(tmx)
	l.tilemap.Frame().Set(l.camera.Viewport())

	c, err := BuildLevel(l.ctx.Logger(), l.world, tmx)
	if err != nil {
		return err
	}

	l.player = c.(*collision.BoxCollider)

	l.camera.X = l.player.X + l.player.Width/2
	l.camera.Y = l.player.Y + l.player.Height/2

	l.world.OnEnter = l.OnCollision
	l.world.OnStay = l.OnCollision

	return nil
}

func (l *Level) Update() error {
	l.player.Velocity[1] += 0.5 // Gravity

	// Early Update Phase
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		l.player.Velocity[0] -= 0.2
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		l.player.Velocity[0] += 0.2
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		l.player.Velocity[1] -= 1
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		l.player.X = l.camera.X - l.player.Width/2
		l.player.Y = l.camera.Y - l.player.Height/2
		l.player.Velocity[0] = 0
		l.player.Velocity[1] = 0
	}

	// Fixed Update Phase - TODO implement fixed timestep
	l.world.UpdateColliders()
	l.world.CheckCollisions()

	// Late Update Phase
	l.camera.X = l.player.X + l.player.Width/2
	l.camera.Y = l.player.Y + l.player.Height/2

	maxX := l.tilemap.Tmx.Width * l.tilemap.Tmx.TileWidth
	maxY := l.tilemap.Tmx.Height * l.tilemap.Tmx.TileHeight

	l.camera.X = float32(math.Max(float64(l.camera.Width)/2, float64(l.camera.X)))
	l.camera.Y = float32(math.Max(float64(l.camera.Height)/2, float64(l.camera.Y)))
	l.camera.X = float32(math.Min(float64(maxX)-float64(l.camera.Width)/2, float64(l.camera.X)))
	l.camera.Y = float32(math.Min(float64(maxY)-float64(l.camera.Height)/2, float64(l.camera.Y)))

	l.player.Velocity[0] *= 0.8
	l.player.Velocity[1] *= 0.9
	return nil
}

func (l *Level) Draw(screen *ebiten.Image) {
	l.tilemap.Frame().Set(l.camera.Viewport())
	l.tilemap.BufferFrame()

	mat := l.camera.Matrix()
	itr := l.tilemap.Itr()

	for tiles := itr.Next(); tiles != nil; tiles = itr.Next() {
		l.DrawTileBatch(screen, tiles, mat)
	}

	// l.DrawCollisionCells(screen, mat)
	l.DrawPotentialCollisions(screen, mat)

	msg := fmt.Sprintf("Velocity X: %.2f Y: %.2f", l.player.Velocity[0], l.player.Velocity[1])
	ebitenutil.DebugPrintAt(screen, msg, 10, 10)
}

func (l *Level) DrawTileBatch(screen *ebiten.Image, tiles []tilemap.Data, mat ebiten.GeoM) {
	for i := range tiles {
		tileset, err := l.tilemap.GetTileset(tiles[i].TsIdx)
		if err != nil {
			println(err.Error())
			return
		}

		tsx := assets.MustGet[*tiled.Tsx](assets.AssetHandle(tileset.Source))
		img := assets.MustGet[*ebiten.Image](assets.AssetHandle(tsx.Image.Source))

		srcX := (int32(tiles[i].TileID) % tsx.Columns) * tsx.TileWidth
		srcY := (int32(tiles[i].TileID) / tsx.Columns) * tsx.TileHeight
		srcRect := image.Rect(int(srcX), int(srcY), int(srcX+tsx.TileWidth), int(srcY+tsx.TileHeight))

		distX := float64(tiles[i].X) + float64(tsx.TileOffset.X)
		distY := float64(tiles[i].Y) + float64(tsx.TileOffset.Y)
		distY -= float64(tsx.TileHeight) - float64(l.tilemap.Tmx.TileHeight) // Align to bottom of tile

		l.op.GeoM.Reset()

		if tiles[i].FlipFlag.Diagonal() {
			l.op.GeoM.Rotate(math.Pi * 0.5)
			l.op.GeoM.Scale(-1, 1)
			l.op.GeoM.Translate(float64(tsx.TileHeight-tsx.TileWidth), 0)
		}

		if tiles[i].FlipFlag.Horizontal() {
			l.op.GeoM.Scale(-1, 1)
			l.op.GeoM.Translate(float64(tsx.TileWidth), 0)
		}

		if tiles[i].FlipFlag.Vertical() {
			l.op.GeoM.Scale(1, -1)
			l.op.GeoM.Translate(0, float64(tsx.TileHeight))
		}

		l.op.GeoM.Translate(distX, distY)
		l.op.GeoM.Concat(mat)

		screen.DrawImage(img.SubImage(srcRect).(*ebiten.Image), &l.op)
	}
}

func (l *Level) DrawCollisionCells(screen *ebiten.Image, mat ebiten.GeoM) {
	cells := l.world.QueryCells(l.player.Bounds())
	width, height := l.world.GetCellSize()
	path := vector.Path{}

	for i := range cells {
		x, y := hash.DecodeGridKey(cells[i])

		minX, minY := mat.Apply(float64(x)*float64(width), float64(y)*float64(height))
		maxX, maxY := mat.Apply(float64(x)*float64(width)+float64(width), float64(y)*float64(height)+float64(height))

		path.MoveTo(float32(minX), float32(minY))
		path.LineTo(float32(maxX), float32(minY))
		path.LineTo(float32(maxX), float32(maxY))
		path.LineTo(float32(minX), float32(maxY))
		path.LineTo(float32(minX), float32(minY))
	}

	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(color.RGBA{R: 255, A: 255})

	vector.StrokePath(screen, &path, &vector.StrokeOptions{
		Width: 1,
	}, op)
}

func (l *Level) DrawPotentialCollisions(screen *ebiten.Image, mat ebiten.GeoM) {
	minX, minY, maxX, maxY := l.player.Bounds()

	colliders := l.world.Query(minX, minY, maxX, maxY)
	path := vector.Path{}

	for i := range colliders {
		cMinX, cMinY, cMaxX, cMaxY := colliders[i].Bounds()

		minX, minY := mat.Apply(float64(cMinX), float64(cMinY))
		maxX, maxY := mat.Apply(float64(cMaxX), float64(cMaxY))

		path.MoveTo(float32(minX), float32(minY))
		path.LineTo(float32(maxX), float32(minY))
		path.LineTo(float32(maxX), float32(maxY))
		path.LineTo(float32(minX), float32(maxY))
		path.LineTo(float32(minX), float32(minY))
	}

	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(color.RGBA{G: 255, A: 255})

	vector.StrokePath(screen, &path, &vector.StrokeOptions{
		Width: 1,
	}, op)
}
