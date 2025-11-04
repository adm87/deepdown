package physics

import (
	"math"

	"github.com/adm87/deepdown/scripts/deepdown"
	"github.com/adm87/utilities/hash"
)

const (
	GridCellSize float32 = 8.0
	Gravity      float32 = 400.0
	Epsilon      float32 = 0.0001

	MinimumVelocityThreshold float64 = 0.01
	MaxVelocityRiseSpeed     float32 = -150.0
	MaxVelocityFallSpeed     float32 = 200.0
)

func clamp[T float32 | float64](value, min, max T) T {
	return T(math.Max(float64(min), math.Min(float64(max), float64(value))))
}

type CollisionPair uint64

func EncodePair(id1, id2 uint32) CollisionPair {
	if id1 > id2 {
		id1, id2 = id2, id1
	}
	return CollisionPair(uint64(id1)<<32 | uint64(id2))
}

func DecodePair(pair CollisionPair) (id1, id2 uint32) {
	id1 = uint32(pair >> 32)
	id2 = uint32(pair & 0xFFFFFFFF)
	return id1, id2
}

type World struct {
	ctx deepdown.Context

	staticGrid *hash.Grid[Collider] // Static world colliders
	bodyGrid   *hash.Grid[Collider] // Dynamic and trigger body colliders
}

func NewWorld(ctx deepdown.Context) *World {
	return &World{
		ctx:        ctx,
		staticGrid: hash.NewGrid[Collider](GridCellSize, GridCellSize),
		bodyGrid:   hash.NewGrid[Collider](GridCellSize, GridCellSize),
	}
}

func (w *World) AddCollider(collider Collider) {
	switch collider.Info().State {
	case ColliderStateStatic:
		w.insert(collider, w.staticGrid)
	default:
		w.insert(collider, w.bodyGrid)
	}
}

func (w *World) RemoveCollider(collider Collider) {
	switch collider.Info().State {
	case ColliderStateStatic:
		w.staticGrid.Remove(collider)
	default:
		w.bodyGrid.Remove(collider)
	}
}

func (w *World) Update(dt float64, minX, minY, maxX, maxY float32) {
	activeBodies := w.bodyGrid.Query(minX, minY, maxX, maxY)

	w.preupdate(dt, activeBodies)

	w.handleCollisions(activeBodies)

	w.postupdate(activeBodies)
}

func (w *World) QueryStatic(minX, minY, maxX, maxY float32) []Collider {
	return w.staticGrid.Query(minX, minY, maxX, maxY)
}

func (w *World) QueryBody(minX, minY, maxX, maxY float32) []Collider {
	return w.bodyGrid.Query(minX, minY, maxX, maxY)
}

func (w *World) QueryStaticCells(minX, minY, maxX, maxY float32) []uint64 {
	return w.staticGrid.QueryCells(minX, minY, maxX, maxY)
}

func (w *World) QueryBodyCells(minX, minY, maxX, maxY float32) []uint64 {
	return w.bodyGrid.QueryCells(minX, minY, maxX, maxY)
}

func (w *World) insert(collider Collider, grid *hash.Grid[Collider]) {
	minX, minY, maxX, maxY := collider.AABB()
	switch c := collider.(type) {
	// case *TriangleCollider:
	// 	grid.InsertFunc(c, minX, minY, maxX, maxY, hash.NoGridPadding, func(cMinX, cMinY, cMaxX, cMaxY float32) bool {
	// 		return c.Triangle.IntersectsAABB(cMinX, cMinY, cMaxX, cMaxY)
	// 	})
	default:
		grid.Insert(c, minX, minY, maxX, maxY, hash.NoGridPadding)
	}
}

func (w *World) preupdate(dt float64, activeBodies []Collider) {
	for i := range activeBodies {
		info := activeBodies[i].Info()

		info.OnGround = w.isGrounded(activeBodies[i], info)
		if !info.OnGround {
			info.Velocity[1] = clamp(info.Velocity[1]+Gravity*float32(dt), MaxVelocityRiseSpeed, MaxVelocityFallSpeed)
		}

		if math.Abs(float64(info.Velocity[0])) < MinimumVelocityThreshold {
			info.Velocity[0] = 0
		}
		if math.Abs(float64(info.Velocity[1])) < MinimumVelocityThreshold {
			info.Velocity[1] = 0
		}

		x, y := activeBodies[i].Position()

		if info.Velocity[0] == 0 && info.Velocity[1] == 0 {
			info.nextPosition[0] = x
			info.nextPosition[1] = y
			continue
		}

		info.nextPosition[0] = x + info.Velocity[0]*float32(dt)
		info.nextPosition[1] = y + info.Velocity[1]*float32(dt)

		info.Velocity[0] *= 0.8
	}
}

func (w *World) postupdate(activeBodies []Collider) {
	for i := range activeBodies {
		info := activeBodies[i].Info()

		x, y := activeBodies[i].Position()
		nX, nY := info.nextPosition[0], info.nextPosition[1]

		if nX == x && nY == y {
			continue
		}

		info.prevPosition[0] = x
		info.prevPosition[1] = y

		activeBodies[i].SetPosition(nX, nY)

		w.bodyGrid.Remove(activeBodies[i])
		w.insert(activeBodies[i], w.bodyGrid)
	}
}

func (w *World) handleCollisions(activeBodies []Collider) {
	for i := range activeBodies {
		info := activeBodies[i].Info()
		info.collisions = info.collisions[:0]

		if info.Mode == CollisionModeIgnore {
			continue
		}

		others := w.staticGrid.Query(activeBodies[i].AABB())
		for j := range others {
			otherInfo := others[j].Info()

			if otherInfo.Mode == CollisionModeIgnore || info.id == otherInfo.id {
				continue
			}

			if !ShouldCollide(info.Layer, otherInfo.Layer) {
				continue
			}

			if contact, overlaps := CheckOverlap(activeBodies[i], others[j]); overlaps {
				info.collisions = append(info.collisions, contact)
			}
		}

		// w.resolveStaticContacts(info)
	}
}

func (w *World) isGrounded(collider Collider, info *ColliderInfo) bool {
	return false // or true if ground is close enough
}

func (w *World) resolveStaticContacts(info *ColliderInfo) {
	if len(info.collisions) == 0 {
		return
	}

	var horizontal *Collision
	var vertical *Collision
	var slope *Collision

	for i := range info.collisions {
		contact := &info.collisions[i]

		switch contact.other.(type) {
		case *TriangleCollider:
			if slope == nil || contact.Depth > slope.Depth {
				slope = contact
			}
		default:
			if math.Abs(float64(contact.Normal[1])) > math.Abs(float64(contact.Normal[0])) {
				if vertical == nil || contact.Depth > vertical.Depth {
					vertical = contact
				}
			} else {
				if horizontal == nil || contact.Depth > horizontal.Depth {
					horizontal = contact
				}
			}
		}
	}

	if slope != nil {
		resolveStaticSlopeCollision(info, slope)
		info.OnGround = true
	} else if vertical != nil {
		if resolveStaticVerticalCollision(info, vertical) {
			info.OnGround = true
		}
	}

	if horizontal != nil {
		resolveStaticHorizontalCollision(info, horizontal)
	}
}

func resolveStaticSlopeCollision(info *ColliderInfo, contact *Collision) {
	info.nextPosition[0] += contact.Normal[0] * contact.Depth
	info.nextPosition[1] += contact.Normal[1] * contact.Depth

	if contact.Normal[1] < 0 {
		info.Velocity[1] = 0
	} else {
		dotProduct := info.Velocity[0]*contact.Normal[0] + info.Velocity[1]*contact.Normal[1]
		if dotProduct < 0 {
			info.Velocity[0] -= contact.Normal[0] * dotProduct
			info.Velocity[1] -= contact.Normal[1] * dotProduct
		}
	}
}

func resolveStaticVerticalCollision(info *ColliderInfo, contact *Collision) bool {
	info.nextPosition[1] += contact.Normal[1] * contact.Depth

	if (contact.Normal[1] < 0 && info.Velocity[1] > 0) ||
		(contact.Normal[1] > 0 && info.Velocity[1] < 0) {
		info.Velocity[1] = 0
	}

	return contact.Normal[1] < 0
}

func resolveStaticHorizontalCollision(info *ColliderInfo, contact *Collision) {
	info.nextPosition[0] += contact.Normal[0] * contact.Depth

	if (contact.Normal[0] < 0 && info.Velocity[0] > 0) ||
		(contact.Normal[0] > 0 && info.Velocity[0] < 0) {
		info.Velocity[0] = 0
	}
}
