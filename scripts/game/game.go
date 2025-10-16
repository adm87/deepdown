package game

import (
	"github.com/adm87/deepdown/data"
	"github.com/adm87/deepdown/scripts/assets"
	"github.com/adm87/deepdown/scripts/camera"
	"github.com/adm87/deepdown/scripts/debug"
	"github.com/adm87/deepdown/scripts/deepdown"
	"github.com/adm87/deepdown/scripts/level"
	"github.com/adm87/tiled"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WindowTitle = "Deepdown"

	TargetWidth  = 1280
	TargetHeight = 720

	Scale = 0.2
)

type Game struct {
	ctx deepdown.Context

	lvl    *level.Level
	camera *camera.Camera
}

func NewGame(ctx deepdown.Context) *Game {
	assets.MustLoad(
		data.Img10x10,
		data.SampleMap,
		data.SampleSheet,
		data.TilemapPacked,
	)

	ebiten.SetWindowTitle(WindowTitle)
	ebiten.SetWindowSize(TargetWidth, TargetHeight)

	tilemap := tiled.NewTilemap()
	tilemap.SetTmx(assets.MustGet[*tiled.Tmx](data.SampleMap))

	minX, minY, maxX, maxY := tilemap.Bounds()

	width := TargetWidth * Scale
	height := TargetHeight * Scale

	cam := camera.NewCamera(float32((minX+maxX)/2), float32((minY+maxY)/2), float32(width), float32(height))

	lvl := level.NewLevel(ctx)
	lvl.SetTmx(assets.MustGet[*tiled.Tmx](data.SampleMap))

	return &Game{
		ctx:    ctx,
		lvl:    lvl,
		camera: cam,
	}
}

func (g *Game) Update() error {
	if err := debug.PollInput(); err != nil {
		return err
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.camera.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.camera.Y += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.camera.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.camera.X += 2
	}

	minX, minY, maxX, maxY := g.lvl.Bounds()
	if g.camera.X < float32(minX)+g.camera.Width/2 {
		g.camera.X = float32(minX) + g.camera.Width/2
	}
	if g.camera.Y < float32(minY)+g.camera.Height/2 {
		g.camera.Y = float32(minY) + g.camera.Height/2
	}
	if g.camera.X > float32(maxX)-g.camera.Width/2 {
		g.camera.X = float32(maxX) - g.camera.Width/2
	}
	if g.camera.Y > float32(maxY)-g.camera.Height/2 {
		g.camera.Y = float32(maxY) - g.camera.Height/2
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.lvl.Draw(screen, g.camera)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return TargetWidth * Scale, TargetHeight * Scale
}
