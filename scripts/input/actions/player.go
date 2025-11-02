package actions

import (
	"github.com/adm87/deepdown/scripts/input"
)

const (
	MoveLeft = input.InputAction(iota)
	MoveRight
	MoveUp
	MoveDown
	Jump
)

const (
	MovementHoldThresh = 0.1
	MovementSpeed      = 10
	MovementDampening  = 0.8
	JumpVelocity       = -85.0
	JumpThresh         = 0.1
	Gravity            = 400.0
)
