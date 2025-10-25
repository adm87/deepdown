package level

import (
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
	return &Level{
		ctx:     ctx,
		tilemap: tilemap.NewMap(),
		camera:  camera.NewCamera(0, 0, targetWidth, targetHeight),
		world:   collision.NewWorld(ctx),
		op:      ebiten.DrawImageOptions{},
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
	return nil
}

func (l *Level) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		l.player.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		l.player.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		l.player.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		l.player.X += 1
	}

	l.world.UpdateCollider(l.player, hash.NoGridPadding)

	l.camera.X = l.player.X + l.player.Width/2
	l.camera.Y = l.player.Y + l.player.Height/2

	l.tilemap.Frame().Set(l.camera.Viewport())

	l.world.CheckCollisions()
	return nil
}

func (l *Level) Draw(screen *ebiten.Image) {
	l.tilemap.BufferFrame()

	mat := l.camera.Matrix()
	itr := l.tilemap.Itr()

	for tiles := itr.Next(); tiles != nil; tiles = itr.Next() {
		l.DrawTileBatch(screen, tiles, mat)
	}

	l.DrawCollisionCells(screen, mat)
	l.DrawPotentialCollisions(screen, mat)
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
	cells := l.world.GetCells()
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
