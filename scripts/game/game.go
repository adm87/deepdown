package game

import (
	"github.com/adm87/deepdown/data"
	"github.com/adm87/deepdown/scripts/assets"
	"github.com/adm87/deepdown/scripts/debug"
	"github.com/adm87/deepdown/scripts/deepdown"
	"github.com/adm87/deepdown/scripts/input"
	"github.com/adm87/deepdown/scripts/input/actions"
	"github.com/adm87/deepdown/scripts/level"
	"github.com/adm87/tiled"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WindowTitle = "Deepdown"

	TargetWidth  = 1280
	TargetHeight = 720

	Scale         = 0.15
	MaxFixedSteps = 5
)

type Game struct {
	ctx deepdown.Context

	lvl *level.Level

	dt              float64
	fixDt           float64
	accumulatedTime float64
}

func NewGame(ctx deepdown.Context) *Game {
	assets.MustLoad(
		data.GymCollision,
		data.SampleSheet,
		data.TilemapPacked,
	)

	ebiten.SetWindowTitle(WindowTitle)
	ebiten.SetWindowSize(int(TargetWidth), int(TargetHeight))

	width := float32(TargetWidth) * float32(Scale)
	height := float32(TargetHeight) * float32(Scale)

	lvl := level.NewLevel(ctx, width, height)
	lvl.SetTmx(assets.MustGet[*tiled.Tmx](data.GymCollision))

	input.Register(
		input.NewKeyHoldBinding(
			actions.MoveLeft,
			actions.MovementHoldThresh,
			[2]ebiten.Key{ebiten.KeyA, ebiten.KeyLeft},
		),
		input.NewKeyHoldBinding(
			actions.MoveRight,
			actions.MovementHoldThresh,
			[2]ebiten.Key{ebiten.KeyD, ebiten.KeyRight},
		),
		input.NewKeyPressDurationBinding(
			actions.Jump,
			actions.JumpThresh,
			[2]ebiten.Key{ebiten.KeySpace, ebiten.KeyUp},
		),
	)

	return &Game{
		ctx:   ctx,
		lvl:   lvl,
		fixDt: 1.0 / 60.0,
	}
}

func (g *Game) Update() error {
	if err := debug.PollInput(); err != nil {
		return err
	}

	g.dt = 1.0 / float64(ebiten.TPS())
	g.accumulatedTime += g.dt

	if g.accumulatedTime > MaxFixedSteps*g.fixDt {
		g.accumulatedTime = MaxFixedSteps * g.fixDt
	}

	input.Update(g.dt)

	g.lvl.Update(g.dt)
	for g.accumulatedTime >= g.fixDt {
		g.lvl.FixedUpdate(g.fixDt)
		g.accumulatedTime -= g.fixDt
	}
	g.lvl.LateUpdate(g.dt)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.lvl.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	width := int(float64(TargetWidth) * Scale)
	height := int(float64(TargetHeight) * Scale)
	return width, height
}
