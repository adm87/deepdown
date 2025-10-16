package collision

import "github.com/adm87/deepdown/scripts/geom"

type AABB interface {
	Min() (x, y float32)
	Max() (x, y float32)
}

type Collider interface {
	AABB

	Info() ColliderInfo
}

type ColliderInfo struct {
	Layer     Layer
	Behaviour Behaviour
	Detection Detection
}

// =========== BoxCollider ==========

type BoxCollider struct {
	geom.Rectangle
	ColliderInfo
}

func NewBoxCollider(x, y, width, height float32) *BoxCollider {
	return &BoxCollider{
		Rectangle: geom.NewRectangle(x, y, width, height),
		ColliderInfo: ColliderInfo{
			Layer:     DefaultLayer,
			Behaviour: StaticBehaviour,
			Detection: DiscreteDetection,
		},
	}
}

func (bc *BoxCollider) Info() ColliderInfo {
	return bc.ColliderInfo
}

func (bc *BoxCollider) WithLayer(layer Layer) *BoxCollider {
	bc.ColliderInfo.Layer = layer
	return bc
}

func (bc *BoxCollider) WithBehaviour(behaviour Behaviour) *BoxCollider {
	bc.ColliderInfo.Behaviour = behaviour
	return bc
}

func (bc *BoxCollider) WithDetection(detection Detection) *BoxCollider {
	bc.ColliderInfo.Detection = detection
	return bc
}

// =========== PolygonCollider ==========

type PolygonCollider struct {
	geom.Polygon
	ColliderInfo
}

func NewPolygonCollider(x, y float32, points []float32) *PolygonCollider {
	return &PolygonCollider{
		Polygon: geom.NewPolygon(x, y, points),
		ColliderInfo: ColliderInfo{
			Layer:     DefaultLayer,
			Behaviour: StaticBehaviour,
			Detection: DiscreteDetection,
		},
	}
}

func (pc *PolygonCollider) Info() ColliderInfo {
	return pc.ColliderInfo
}

func (pc *PolygonCollider) WithLayer(layer Layer) *PolygonCollider {
	pc.ColliderInfo.Layer = layer
	return pc
}

func (pc *PolygonCollider) WithBehaviour(behaviour Behaviour) *PolygonCollider {
	pc.ColliderInfo.Behaviour = behaviour
	return pc
}

func (pc *PolygonCollider) WithDetection(detection Detection) *PolygonCollider {
	pc.ColliderInfo.Detection = detection
	return pc
}
