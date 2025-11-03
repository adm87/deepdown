package physics

import (
	"sync"

	"github.com/adm87/deepdown/scripts/geom"
)

// Note: colliderIDCounter and its associated used within the physics package is for internal use only.
// Do not modify or access it directly from outside this package.
// It is not recommended to use this ID for game logic or other purposes.
var colliderIDCounter uint32 = 1

func nextColliderID() (id uint32) {
	id = colliderIDCounter
	colliderIDCounter++
	return
}

// ClearColliderPools clears all pools causing all unused colliders to be garbage collected.
// Any colliders still in use will not be garbage collected.
func ClearColliderPools() {
	colliderIDCounter = 1

	boxColliderPool = newBoxColliderPool()
	triangleColliderPool = newTriangleColliderPool()
}

// ReleaseCollider returns a collider to the appropriate physics collider pool.
func ReleaseCollider(collider Collider) {
	switch c := collider.(type) {
	case *BoxCollider:
		ReleaseBoxCollider(c)
	case *TriangleCollider:
		ReleaseTriangleCollider(c)
	default:
		panic("unknown collider type")
	}
}

// =========== Box Colliders ==========

var boxColliderPool = newBoxColliderPool()

func newBoxColliderPool() *sync.Pool {
	return &sync.Pool{
		New: func() any {
			return &BoxCollider{
				ColliderInfo: ColliderInfo{
					Movement:   Movement{},
					id:         nextColliderID(),
					Layer:      CollisionLayerDefault,
					Mode:       CollisionModeDiscrete,
					State:      ColliderStateStatic,
					Type:       ColliderTypeBox,
					collisions: make([]Collision, 0, 4),
				},
				Rectangle: geom.Rectangle{},
			}
		},
	}
}

// GetBoxCollider retrieves a BoxCollider from the physics collider pool.
func GetBoxCollider(x, y, width, height float32) *BoxCollider {
	bc := boxColliderPool.Get().(*BoxCollider)
	bc.X = x
	bc.Y = y
	bc.Width = width
	bc.Height = height
	bc.Movement.nextPosition[0] = x
	bc.Movement.nextPosition[1] = y
	bc.Movement.prevPosition[0] = x
	bc.Movement.prevPosition[1] = y
	return bc
}

// ReleaseBoxCollider returns a BoxCollider to the physics collider pool.
func ReleaseBoxCollider(bc *BoxCollider) {
	if bc == nil {
		panic("cannot release a nil BoxCollider")
	}
	bc.X = 0
	bc.Y = 0
	bc.Width = 0
	bc.Height = 0
	bc.Movement.prevPosition[0] = 0
	bc.Movement.prevPosition[1] = 0
	bc.Movement.Velocity[0] = 0
	bc.Movement.Velocity[1] = 0
	bc.ColliderInfo.Layer = CollisionLayerDefault
	bc.ColliderInfo.State = ColliderStateStatic
	bc.ColliderInfo.Type = ColliderTypeBox
	bc.ColliderInfo.Mode = CollisionModeDiscrete
	bc.collisions = bc.collisions[:0]
	boxColliderPool.Put(bc)
}

// =========== Triangle Colliders ==========

var triangleColliderPool = newTriangleColliderPool()

func newTriangleColliderPool() *sync.Pool {
	return &sync.Pool{
		New: func() any {
			return &TriangleCollider{
				ColliderInfo: ColliderInfo{
					Movement:   Movement{},
					id:         nextColliderID(),
					Layer:      CollisionLayerDefault,
					Mode:       CollisionModeDiscrete,
					State:      ColliderStateStatic,
					Type:       ColliderTypeBox,
					collisions: make([]Collision, 0, 4),
				},
				Triangle: geom.Triangle{},
			}
		},
	}
}

// GetTriangleCollider retrieves a TriangleCollider from the physics collider pool.
func GetTriangleCollider(x, y float32, points [6]float32) *TriangleCollider {
	tc := triangleColliderPool.Get().(*TriangleCollider)
	tc.SetPoints(points)
	tc.X = x
	tc.Y = y
	tc.Movement.nextPosition[0] = x
	tc.Movement.nextPosition[1] = y
	tc.Movement.prevPosition[0] = x
	tc.Movement.prevPosition[1] = y
	return tc
}

func ReleaseTriangleCollider(tc *TriangleCollider) {
	if tc == nil {
		panic("cannot release a nil TriangleCollider")
	}
	tc.X = 0
	tc.Y = 0
	tc.SetPoints([6]float32{})
	tc.Movement.nextPosition[0] = 0
	tc.Movement.nextPosition[1] = 0
	tc.Movement.prevPosition[0] = 0
	tc.Movement.prevPosition[1] = 0
	tc.Movement.Velocity[0] = 0
	tc.Movement.Velocity[1] = 0
	tc.ColliderInfo.Layer = CollisionLayerDefault
	tc.ColliderInfo.State = ColliderStateStatic
	tc.ColliderInfo.Type = ColliderTypeBox
	tc.ColliderInfo.Mode = CollisionModeDiscrete
	tc.collisions = tc.collisions[:0]
	triangleColliderPool.Put(tc)
}
