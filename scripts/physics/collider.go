package physics

import (
	"github.com/adm87/deepdown/scripts/geom"
)

type Collider interface {
	AABB() (minX, minY, maxX, maxY float32) // Axis-Aligned Bounding Box
	Equals(other Collider) bool             // Equality check
	Info() *ColliderInfo                    // Returns collider info
	Position() (x, y float32)               // Current position
	SetPosition(x, y float32)               // Sets current position
}

type ColliderInfo struct {
	Movement

	id    uint32
	Layer Layer
	State State
	Type  Type
	Mode  Mode

	OnGround bool

	contacts []Contact
}

type Movement struct {
	Velocity [2]float32 // Velocity

	nextPosition [2]float32 // Next position
	prevPosition [2]float32 // Previous position
}

type Contact struct {
	Normal [2]float32
	Depth  float32
}

// =========== Collider Types ==========

type Type uint8

const (
	ColliderTypeBox Type = iota
)

func (ct Type) String() string {
	switch ct {
	case ColliderTypeBox:
		return "Box"
	default:
		return "Unknown"
	}
}

func (ct Type) IsValid() bool {
	return ct <= ColliderTypeBox
}

// =========== Collider State ==========

type State uint8

const (
	ColliderStateStatic State = iota
	ColliderStateDynamic
	ColliderStateTrigger
)

func (cs State) String() string {
	switch cs {
	case ColliderStateStatic:
		return "Solid"
	case ColliderStateDynamic:
		return "Dynamic"
	case ColliderStateTrigger:
		return "Trigger"
	default:
		return "Unknown"
	}
}

func (cs State) IsValid() bool {
	return cs <= ColliderStateTrigger
}

// =========== Collision Mode ==========

type Mode uint8

const (
	CollisionModeIgnore Mode = iota
	CollisionModeDiscrete
	CollisionModeContinuous
)

func (cm Mode) String() string {
	switch cm {
	case CollisionModeDiscrete:
		return "Discrete"
	case CollisionModeContinuous:
		return "Continuous"
	default:
		return "Unknown"
	}
}

func (cm Mode) IsValid() bool {
	return cm <= CollisionModeContinuous
}

// =========== Box Collider ==========

type BoxCollider struct {
	ColliderInfo
	geom.Rectangle
}

func (bc *BoxCollider) AABB() (minX, minY, maxX, maxY float32) {
	minX = bc.nextPosition[0]
	minY = bc.nextPosition[1]
	maxX = minX + bc.Width
	maxY = minY + bc.Height
	return
}

func (bc *BoxCollider) Equals(other Collider) bool {
	return bc.ColliderInfo.id == other.Info().id
}

func (bc *BoxCollider) Info() *ColliderInfo {
	return &bc.ColliderInfo
}

func (bc *BoxCollider) Position() (x, y float32) {
	return bc.X, bc.Y
}

func (bc *BoxCollider) SetPosition(x, y float32) {
	bc.prevPosition[0], bc.prevPosition[1] = bc.X, bc.Y
	bc.X, bc.Y = x, y
	bc.nextPosition[0], bc.nextPosition[1] = x, y
}

// =========== Triangle Collider ==========

type TriangleCollider struct {
	ColliderInfo
	geom.Triangle
}

func (tc *TriangleCollider) AABB() (minX, minY, maxX, maxY float32) {
	return
}

func (tc *TriangleCollider) Equals(other Collider) bool {
	return tc.ColliderInfo.id == other.Info().id
}

func (tc *TriangleCollider) Info() *ColliderInfo {
	return &tc.ColliderInfo
}

func (tc *TriangleCollider) Position() (x, y float32) {
	return tc.X, tc.Y
}

func (tc *TriangleCollider) SetPosition(x, y float32) {
	tc.X, tc.Y = x, y
}
