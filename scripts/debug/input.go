package debug

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	DrawCollisionCells      = false
	DrawPotentialCollisions = false
)

func PollInput() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF9) {
		DrawCollisionCells = !DrawCollisionCells
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF10) {
		DrawPotentialCollisions = !DrawPotentialCollisions
	}

	return nil
}
