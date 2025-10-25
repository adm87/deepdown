package game

import (
	"github.com/adm87/deepdown/data"
	"github.com/adm87/deepdown/scripts/assets"
	"github.com/adm87/deepdown/scripts/debug"
	"github.com/adm87/deepdown/scripts/deepdown"
	"github.com/adm87/deepdown/scripts/level"
	"github.com/adm87/tiled"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WindowTitle = "Deepdown"

	TargetWidth  int32 = 1280
	TargetHeight int32 = 720

	Scale float32 = 0.2
)

type Game struct {
	ctx deepdown.Context

	lvl *level.Level
}

func NewGame(ctx deepdown.Context) *Game {
	assets.MustLoad(
		data.Img10x10,
		data.SampleMap,
		data.SampleSheet,
		data.TilemapPacked,
	)

	ebiten.SetWindowTitle(WindowTitle)
	ebiten.SetWindowSize(int(TargetWidth), int(TargetHeight))

	width := float32(TargetWidth) * Scale
	height := float32(TargetHeight) * Scale

	lvl := level.NewLevel(ctx, width, height)
	lvl.SetTmx(assets.MustGet[*tiled.Tmx](data.SampleMap))

	return &Game{
		ctx: ctx,
		lvl: lvl,
	}
}

func (g *Game) Update() error {
	if err := debug.PollInput(); err != nil {
		return err
	}

	return g.lvl.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.lvl.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(float32(TargetWidth) * Scale), int(float32(TargetHeight) * Scale)
}
