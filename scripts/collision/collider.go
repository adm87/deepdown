package collision

import (
	"math"

	"github.com/adm87/deepdown/scripts/geom"
)

const (
	MinVelocityThreshold float64 = 0.01
)

type ColliderFlags uint8

const (
	FlagNone     ColliderFlags = 0
	FlagOnGround ColliderFlags = 1 << iota
	FlagOnSlope
)

type AABB interface {
	Bounds() (minX, minY, maxX, maxY float32)
}

type Collider interface {
	AABB

	Info() ColliderInfo
	Equals(other Collider) bool

	applyVelocity(dt float64) bool
}

type ColliderInfo struct {
	id       uint32
	priority uint8

	Shape     Shape
	Layer     Layer
	Type      Type
	Detection Detection
	Flags     ColliderFlags

	Offset   [2]float32
	Velocity [2]float32

	previous [2]float32
}

var nextColliderID uint32

func (ci ColliderInfo) ID() uint32 {
	return ci.id
}

// =========== TriangleCollider ==========

type TriangleCollider struct {
	geom.Triangle
	ColliderInfo
}

func NewTriangleCollider(x, y float32, points []float32) *TriangleCollider {
	nextColliderID++
	return &TriangleCollider{
		Triangle: geom.NewTriangle(x, y, points),
		ColliderInfo: ColliderInfo{
			id:        nextColliderID,
			priority:  75,
			Layer:     DefaultLayer,
			Type:      Static,
			Detection: DiscreteDetection,
			Shape:     Triangle,
			Offset:    [2]float32{0, 0},
		},
	}
}

func (tc *TriangleCollider) applyVelocity(dt float64) bool {
	tc.previous[0] = tc.Triangle.X
	tc.previous[1] = tc.Triangle.Y

	if math.Abs(float64(tc.ColliderInfo.Velocity[0])) < MinVelocityThreshold {
		tc.ColliderInfo.Velocity[0] = 0
	}
	if math.Abs(float64(tc.ColliderInfo.Velocity[1])) < MinVelocityThreshold {
		tc.ColliderInfo.Velocity[1] = 0
	}
	if tc.ColliderInfo.Velocity[0] == 0 && tc.ColliderInfo.Velocity[1] == 0 {
		return false
	}

	tc.Triangle.X += tc.ColliderInfo.Velocity[0] * float32(dt)
	tc.Triangle.Y += tc.ColliderInfo.Velocity[1] * float32(dt)
	return true
}

func (tc *TriangleCollider) Info() ColliderInfo {
	return tc.ColliderInfo
}

func (tc *TriangleCollider) Equals(other Collider) bool {
	return tc.ColliderInfo.id == other.Info().id
}

func (tc *TriangleCollider) SetLayer(layer Layer) {
	tc.ColliderInfo.Layer = layer
}

func (tc *TriangleCollider) SetType(colliderType Type) {
	tc.ColliderInfo.Type = colliderType
}

func (tc *TriangleCollider) SetDetection(detection Detection) {
	tc.ColliderInfo.Detection = detection
}

func (tc *TriangleCollider) Bounds() (minX, minY, maxX, maxY float32) {
	minX, minY = tc.Triangle.Min()
	maxX, maxY = tc.Triangle.Max()
	return minX + tc.Offset[0], minY + tc.Offset[1], maxX + tc.Offset[0], maxY + tc.Offset[1]
}

// =========== BoxCollider ==========

type BoxCollider struct {
	geom.Rectangle
	ColliderInfo
}

func NewBoxCollider(x, y, width, height float32) *BoxCollider {
	nextColliderID++
	return &BoxCollider{
		Rectangle: geom.NewRectangle(x, y, width, height),
		ColliderInfo: ColliderInfo{
			id:        nextColliderID,
			priority:  50,
			Layer:     DefaultLayer,
			Type:      Static,
			Detection: DiscreteDetection,
			Shape:     Box,
			Offset:    [2]float32{0, 0},
		},
	}
}

func (bc *BoxCollider) applyVelocity(dt float64) bool {
	bc.previous[0] = bc.Rectangle.X
	bc.previous[1] = bc.Rectangle.Y

	if math.Abs(float64(bc.ColliderInfo.Velocity[0])) < MinVelocityThreshold {
		bc.ColliderInfo.Velocity[0] = 0
	}
	if math.Abs(float64(bc.ColliderInfo.Velocity[1])) < MinVelocityThreshold {
		bc.ColliderInfo.Velocity[1] = 0
	}
	if bc.ColliderInfo.Velocity[0] == 0 && bc.ColliderInfo.Velocity[1] == 0 {
		return false
	}

	bc.Rectangle.X += bc.ColliderInfo.Velocity[0] * float32(dt)
	bc.Rectangle.Y += bc.ColliderInfo.Velocity[1] * float32(dt)
	return true
}

func (bc *BoxCollider) Info() ColliderInfo {
	return bc.ColliderInfo
}

func (bc *BoxCollider) Equals(other Collider) bool {
	return bc.ColliderInfo.id == other.Info().id
}

func (bc *BoxCollider) SetLayer(layer Layer) {
	bc.ColliderInfo.Layer = layer
}

func (bc *BoxCollider) SetType(colliderType Type) {
	bc.ColliderInfo.Type = colliderType
}

func (bc *BoxCollider) SetDetection(detection Detection) {
	bc.ColliderInfo.Detection = detection
}

func (bc *BoxCollider) Bounds() (minX, minY, maxX, maxY float32) {
	minX, minY = bc.Rectangle.Min()
	maxX, maxY = bc.Rectangle.Max()
	return minX + bc.Offset[0], minY + bc.Offset[1], maxX + bc.Offset[0], maxY + bc.Offset[1]
}

// =========== PolygonCollider ==========

type PolygonCollider struct {
	geom.Polygon
	ColliderInfo
}

func NewPolygonCollider(x, y float32, points []float32) *PolygonCollider {
	nextColliderID++
	return &PolygonCollider{
		Polygon: geom.NewPolygon(x, y, points),
		ColliderInfo: ColliderInfo{
			id:        nextColliderID,
			Layer:     DefaultLayer,
			Type:      Static,
			Detection: DiscreteDetection,
			Shape:     Polygon,
			Offset:    [2]float32{0, 0},
		},
	}
}

func (pc *PolygonCollider) applyVelocity(dt float64) bool {
	pc.previous[0] = pc.Polygon.X
	pc.previous[1] = pc.Polygon.Y

	if math.Abs(float64(pc.ColliderInfo.Velocity[0])) < MinVelocityThreshold {
		pc.ColliderInfo.Velocity[0] = 0
	}
	if math.Abs(float64(pc.ColliderInfo.Velocity[1])) < MinVelocityThreshold {
		pc.ColliderInfo.Velocity[1] = 0
	}
	if pc.ColliderInfo.Velocity[0] == 0 && pc.ColliderInfo.Velocity[1] == 0 {
		return false
	}

	pc.Polygon.X += pc.ColliderInfo.Velocity[0] * float32(dt)
	pc.Polygon.Y += pc.ColliderInfo.Velocity[1] * float32(dt)
	return true
}

func (pc *PolygonCollider) Info() ColliderInfo {
	return pc.ColliderInfo
}

func (pc *PolygonCollider) Equals(other Collider) bool {
	return pc.ColliderInfo.id == other.Info().id
}

func (pc *PolygonCollider) SetLayer(layer Layer) {
	pc.ColliderInfo.Layer = layer
}

func (pc *PolygonCollider) SetType(colliderType Type) {
	pc.ColliderInfo.Type = colliderType
}

func (pc *PolygonCollider) SetDetection(detection Detection) {
	pc.ColliderInfo.Detection = detection
}

func (pc *PolygonCollider) Bounds() (minX, minY, maxX, maxY float32) {
	minX, minY = pc.Polygon.Min()
	maxX, maxY = pc.Polygon.Max()
	return minX + pc.Offset[0], minY + pc.Offset[1], maxX + pc.Offset[0], maxY + pc.Offset[1]
}
