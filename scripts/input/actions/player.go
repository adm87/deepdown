package actions

import (
	"github.com/adm87/deepdown/scripts/input"
)

const (
	MoveLeft = input.InputAction(iota)
	MoveRight
	Jump
)

const (
	MovementHoldThresh = 0.15
	MovementSpeed      = 10
	MovementDampening  = 0.8
	JumpVelocity       = -75.0
	JumpThresh         = 0.1
	Gravity            = 400.0
)
