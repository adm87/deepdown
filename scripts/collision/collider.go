package collision

import (
	"github.com/adm87/deepdown/scripts/geom"
)

type AABB interface {
	Bounds() (minX, minY, maxX, maxY float32)
}

type Collider interface {
	AABB

	Info() ColliderInfo
}

type ColliderInfo struct {
	id uint32

	Layer     Layer
	Type      Type
	Detection Detection

	Offset [2]float32
}

var nextColliderID uint32

func (ci ColliderInfo) ID() uint32 {
	return ci.id
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
			Layer:     DefaultLayer,
			Type:      Static,
			Detection: DiscreteDetection,
			Offset:    [2]float32{0, 0},
		},
	}
}

func (bc *BoxCollider) Info() ColliderInfo {
	return bc.ColliderInfo
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
			Offset:    [2]float32{0, 0},
		},
	}
}

func (pc *PolygonCollider) Info() ColliderInfo {
	return pc.ColliderInfo
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
