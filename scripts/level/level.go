package level

import (
	"image"
	"log/slog"
	"math"
	"sync"

	"github.com/adm87/deepdown/scripts/assets"
	"github.com/adm87/deepdown/scripts/camera"
	"github.com/adm87/deepdown/scripts/collision"
	"github.com/adm87/deepdown/scripts/deepdown"
	"github.com/adm87/tiled"
	"github.com/adm87/utilities/hashgrid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Level struct {
	ctx deepdown.Context

	tilemap *tiled.Tilemap
	world   *collision.World

	op ebiten.DrawImageOptions
}

func NewLevel(ctx deepdown.Context) *Level {
	return &Level{
		ctx:     ctx,
		tilemap: tiled.NewTilemap(),
		world:   collision.NewWorld(ctx),
		op:      ebiten.DrawImageOptions{},
	}
}

func (l *Level) Bounds() (minX, minY, maxX, maxY int32) {
	return l.tilemap.Bounds()
}

func (l *Level) SetTmx(tmx *tiled.Tmx) error {
	l.tilemap.SetTmx(tmx)
	return l.buildLevel()
}

func (l *Level) Draw(screen *ebiten.Image, camera *camera.Camera) {
	minX, minY, maxX, maxY := camera.Viewport()

	itr, err := l.tilemap.GetTiles(int32(minX), int32(minY), int32(maxX), int32(maxY))
	if err != nil {
		l.ctx.Logger().Error("error", slog.Any("err", err))
		return
	}

	// TASK: Batch tiles to reduce draw calls

	for tiles := itr.Next(); tiles != nil; tiles = itr.Next() {
		l.drawTiles(screen, camera, tiles)
	}

	keys := l.world.Grid.Keys()
	for _, k := range keys {
		x, y := hashgrid.DecodeGridKey(k)
		l.drawWorldCell(screen, camera, x, y, collision.GridCellSize)
	}
}

func (l *Level) drawWorldCell(screen *ebiten.Image, camera *camera.Camera, x, y int32, cellSize float32) {
	path := vector.Path{}
	view := camera.Matrix()

	minX, minY := view.Apply(float64(x)*float64(cellSize), float64(y)*float64(cellSize))
	maxX, maxY := view.Apply(float64(x+1)*float64(cellSize), float64(y+1)*float64(cellSize))

	path.MoveTo(float32(minX), float32(minY))
	path.LineTo(float32(maxX), float32(minY))
	path.LineTo(float32(maxX), float32(maxY))
	path.LineTo(float32(minX), float32(maxY))
	path.Close()

	cs := ebiten.ColorScale{}
	cs.Scale(1, 0, 0, 1)
	cs.ScaleAlpha(0.5)
	vector.StrokePath(screen, &path, &vector.StrokeOptions{
		Width: 1.0,
	}, &vector.DrawPathOptions{
		ColorScale: cs,
	})

}

func (l *Level) drawTiles(screen *ebiten.Image, camera *camera.Camera, tiles []tiled.TileData) {
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
		l.op.GeoM.Concat(camera.Matrix())

		screen.DrawImage(img.SubImage(srcRect).(*ebiten.Image), &l.op)
	}
}

func (l *Level) buildLevel() error {
	errch := make(chan error)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := BuildCollision(l.ctx.Logger(), l.world, tiled.ObjectGroupByName(l.tilemap.Tmx, "WorldCollision")); err != nil {
			l.ctx.Logger().Error("error", slog.Any("err", err))
			errch <- err
			return
		}
	}()

	// TASK: Build other systems (e.g. spawn points, triggers, etc.)

	wg.Wait()

	close(errch)
	return <-errch
}
