package level

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/adm87/deepdown/scripts/assets"
	"github.com/adm87/deepdown/scripts/camera"
	"github.com/adm87/deepdown/scripts/debug"
	"github.com/adm87/deepdown/scripts/deepdown"
	"github.com/adm87/deepdown/scripts/input"
	"github.com/adm87/deepdown/scripts/input/actions"
	"github.com/adm87/deepdown/scripts/physics"
	"github.com/adm87/tiled"
	"github.com/adm87/tiled/tilemap"
	"github.com/adm87/utilities/hash"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Player struct {
	physics.BoxCollider

	Data tilemap.Data
}

type Level struct {
	ctx deepdown.Context

	tilemap *tilemap.Map
	camera  *camera.Camera
	player  *Player

	world *physics.World

	op ebiten.DrawImageOptions
}

func NewLevel(ctx deepdown.Context, targetWidth, targetHeight float32) *Level {
	world := physics.NewWorld(ctx)
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

	if err := l.BuildStaticCollision(tiled.ObjectGroupByName(tmx, "Static")); err != nil {
		return err
	}

	if err := l.BuildPlayer(tiled.ObjectGroupByName(tmx, "Player"), tmx); err != nil {
		return err
	}

	l.clampCamera()
	return nil
}

func (l *Level) Update(dt float64) {
	if input.IsActive(actions.MoveLeft) {
		l.player.Velocity[0] -= actions.MovementSpeed
	}
	if input.IsActive(actions.MoveRight) {
		l.player.Velocity[0] += actions.MovementSpeed
	}
	if jump := input.GetBinding[*input.KeyPressDurationBinding](actions.Jump); jump != nil {
		if l.player.OnGround && jump.JustReleased() {
			pressure := jump.Pressure()
			l.player.Velocity[1] = actions.JumpVelocity * float32(pressure)
			l.player.OnGround = false
		}
	}
}

func (l *Level) FixedUpdate(dt float64) {
	minX, minY, maxX, maxY := l.camera.Viewport()
	l.world.Update(dt, minX, minY, maxX, maxY)
}

func (l *Level) LateUpdate(dt float64) {
	l.camera.X = l.player.X + l.player.Width/2
	l.camera.Y = l.player.Y + l.player.Height/2
	l.player.Data.X = l.player.X
	l.player.Data.Y = l.player.Y
	l.clampCamera()
}

func (l *Level) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 30, G: 30, B: 30, A: 255})

	l.tilemap.Frame().Set(l.camera.Viewport())
	l.tilemap.BufferFrame()

	mat := l.camera.Matrix()
	itr := l.tilemap.Itr()

	for tiles := itr.Next(); tiles != nil; tiles = itr.Next() {
		l.DrawTileBatch(screen, tiles, mat)
	}

	l.DrawTile(&l.player.Data, screen, mat)

	if debug.DrawCollisionCells {
		l.DrawCollisionCells(screen, mat, l.world.QueryStaticCells(l.camera.Viewport()), color.RGBA{R: 255, A: 255})
		l.DrawCollisionCells(screen, mat, l.world.QueryBodyCells(l.camera.Viewport()), color.RGBA{G: 255, A: 255})
	}

	if debug.DrawPotentialCollisions {
		l.DrawPotentialCollisions(screen, mat, l.world.QueryStatic(l.camera.Viewport()), color.RGBA{B: 255, A: 255})
		l.DrawPotentialCollisions(screen, mat, l.world.QueryBody(l.player.AABB()), color.RGBA{R: 255, G: 255, A: 255})
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Vel: %.2f, %.2f\nOnGround: %v", l.player.Velocity[0], l.player.Velocity[1], l.player.OnGround))
}

func (l *Level) DrawTileBatch(screen *ebiten.Image, tiles []tilemap.Data, mat ebiten.GeoM) {
	for i := range tiles {
		l.DrawTile(&tiles[i], screen, mat)
	}
}

func (l *Level) DrawTile(data *tilemap.Data, screen *ebiten.Image, mat ebiten.GeoM) {
	tileset, err := l.tilemap.GetTileset(data.TsIdx)
	if err != nil {
		println(err.Error())
		return
	}

	tsx := assets.MustGet[*tiled.Tsx](assets.AssetHandle(tileset.Source))
	img := assets.MustGet[*ebiten.Image](assets.AssetHandle(tsx.Image.Source))

	srcX := (int32(data.TileID) % tsx.Columns) * tsx.TileWidth
	srcY := (int32(data.TileID) / tsx.Columns) * tsx.TileHeight
	srcRect := image.Rect(int(srcX), int(srcY), int(srcX+tsx.TileWidth), int(srcY+tsx.TileHeight))

	distX := float64(data.X) + float64(tsx.TileOffset.X)
	distY := float64(data.Y) + float64(tsx.TileOffset.Y)
	distY -= float64(tsx.TileHeight) - float64(l.tilemap.Tmx.TileHeight) // Align to bottom of tile

	l.op.GeoM.Reset()

	if data.FlipFlag.Diagonal() {
		l.op.GeoM.Rotate(math.Pi * 0.5)
		l.op.GeoM.Scale(-1, 1)
		l.op.GeoM.Translate(float64(tsx.TileHeight-tsx.TileWidth), 0)
	}

	if data.FlipFlag.Horizontal() {
		l.op.GeoM.Scale(-1, 1)
		l.op.GeoM.Translate(float64(tsx.TileWidth), 0)
	}

	if data.FlipFlag.Vertical() {
		l.op.GeoM.Scale(1, -1)
		l.op.GeoM.Translate(0, float64(tsx.TileHeight))
	}

	l.op.GeoM.Translate(distX, distY)
	l.op.GeoM.Concat(mat)

	screen.DrawImage(img.SubImage(srcRect).(*ebiten.Image), &l.op)
}

func (l *Level) DrawCollisionCells(screen *ebiten.Image, mat ebiten.GeoM, cells []uint64, col color.RGBA) {
	width, height := physics.GridCellSize, physics.GridCellSize
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
	op.ColorScale.ScaleWithColor(col)

	vector.StrokePath(screen, &path, &vector.StrokeOptions{
		Width: 1,
	}, op)
}

func (l *Level) DrawPotentialCollisions(screen *ebiten.Image, mat ebiten.GeoM, colliders []physics.Collider, col color.RGBA) {
	path := vector.Path{}

	for i := range colliders {
		cMinX, cMinY, cMaxX, cMaxY := colliders[i].AABB()

		minX, minY := mat.Apply(float64(cMinX), float64(cMinY))
		maxX, maxY := mat.Apply(float64(cMaxX), float64(cMaxY))

		path.MoveTo(float32(minX), float32(minY))
		path.LineTo(float32(maxX), float32(minY))
		path.LineTo(float32(maxX), float32(maxY))
		path.LineTo(float32(minX), float32(maxY))
		path.LineTo(float32(minX), float32(minY))
	}

	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(col)

	vector.StrokePath(screen, &path, &vector.StrokeOptions{
		Width: 1,
	}, op)
}

func (l *Level) clampCamera() {
	maxX := l.tilemap.Tmx.Width * l.tilemap.Tmx.TileWidth
	maxY := l.tilemap.Tmx.Height * l.tilemap.Tmx.TileHeight

	l.camera.X = float32(math.Max(float64(l.camera.Width)/2, float64(l.camera.X)))
	l.camera.Y = float32(math.Max(float64(l.camera.Height)/2, float64(l.camera.Y)))
	l.camera.X = float32(math.Min(float64(maxX)-float64(l.camera.Width)/2, float64(l.camera.X)))
	l.camera.Y = float32(math.Min(float64(maxY)-float64(l.camera.Height)/2, float64(l.camera.Y)))
}
