package physics

import (
	"github.com/adm87/deepdown/scripts/components"
	"github.com/adm87/deepdown/scripts/ecs/entity"
	"github.com/adm87/utilities/hash"
)

const (
	GridCellSize float32 = 8.0
	Gravity      float32 = 400.0
	Epsilon      float32 = 0.0001

	MinimumVelocityThreshold float64 = 0.01
	MaxVelocityRiseSpeed     float32 = -150.0
	MaxVelocityFallSpeed     float32 = 200.0

	GroundCheckDistance  float32 = 1.0
	GroundCheckTolerance float32 = 0.5

	VelocityDamping float32 = 0.75
)

func CalculateAABB(transfor *components.Transform, bounds *components.Bounds) (minX, minY, maxX, maxY float32) {
	x, y := transfor.Position()
	w, h := bounds.Size()
	ox, oy := bounds.Offset()

	minX = x + ox
	minY = y + oy
	maxX = minX + w
	maxY = minY + h
	return
}

type Contact struct {
}

type World struct {
	staticGrid *hash.Grid[entity.Entity] // Static world colliders
	bodyGrid   *hash.Grid[entity.Entity] // Dynamic and trigger body colliders

	contacts []Contact
}

func NewWorld() *World {
	return &World{
		staticGrid: hash.NewGrid[entity.Entity](GridCellSize, GridCellSize),
		bodyGrid:   hash.NewGrid[entity.Entity](GridCellSize, GridCellSize),
		contacts:   make([]Contact, 0),
	}
}

func (w *World) Add(entity entity.Entity) {
	switch components.GetCollision(entity).Type {
	case components.CollisionTypeStatic:
		w.insert(entity, w.staticGrid)
	default:
		w.insert(entity, w.bodyGrid)
	}
}

func (w *World) RemoveCollider(entity entity.Entity) {
	switch components.GetCollision(entity).Type {
	case components.CollisionTypeStatic:
		w.staticGrid.Remove(entity)
	default:
		w.bodyGrid.Remove(entity)
	}
}

func (w *World) Update(dt float64, minX, minY, maxX, maxY float32) {
	// activeBodies := w.bodyGrid.Query(minX, minY, maxX, maxY)
}

func (w *World) QueryStatic(minX, minY, maxX, maxY float32) []entity.Entity {
	return w.staticGrid.Query(minX, minY, maxX, maxY)
}

func (w *World) QueryBody(minX, minY, maxX, maxY float32) []entity.Entity {
	return w.bodyGrid.Query(minX, minY, maxX, maxY)
}

func (w *World) QueryStaticCells(minX, minY, maxX, maxY float32) []uint64 {
	return w.staticGrid.QueryCells(minX, minY, maxX, maxY)
}

func (w *World) QueryBodyCells(minX, minY, maxX, maxY float32) []uint64 {
	return w.bodyGrid.QueryCells(minX, minY, maxX, maxY)
}

func (w *World) insert(entity entity.Entity, grid *hash.Grid[entity.Entity]) {
	minX, minY, maxX, maxY := CalculateAABB(
		components.GetTransform(entity),
		components.GetBounds(entity),
	)
	grid.Insert(entity, minX, minY, maxX, maxY, hash.NoGridPadding)
}

func (w *World) preupdate(dt float64, activeBodies []entity.Entity) {
	// for i := range activeBodies {
	// 	info := activeBodies[i].Info()

	// 	// Apply gravity and clamp vertical velocity
	// 	velY := clamp(info.Velocity[1]+Gravity*float32(dt), MaxVelocityRiseSpeed, MaxVelocityFallSpeed)

	// 	// Check ground state
	// 	info.OnGround = w.isGrounded(activeBodies[i], info, velY*float32(dt))

	// 	// Update coyote time tracking
	// 	if info.OnGround {
	// 		info.timeSinceLeftGround = 0
	// 	} else {
	// 		info.timeSinceLeftGround += float32(dt)
	// 	}

	// 	// Apply vertical velocity only when airborne
	// 	if !info.OnGround {
	// 		info.Velocity[1] = velY
	// 	}

	// 	// Zero out negligible velocities
	// 	if math.Abs(float64(info.Velocity[0])) < MinimumVelocityThreshold {
	// 		info.Velocity[0] = 0
	// 	}
	// 	if math.Abs(float64(info.Velocity[1])) < MinimumVelocityThreshold {
	// 		info.Velocity[1] = 0
	// 	}

	// 	// Calculate next position
	// 	x, y := activeBodies[i].Position()
	// 	if info.Velocity[0] == 0 && info.Velocity[1] == 0 {
	// 		info.nextPosition[0] = x
	// 		info.nextPosition[1] = y
	// 		continue
	// 	}

	// 	info.nextPosition[0] = x + info.Velocity[0]*float32(dt)
	// 	info.nextPosition[1] = y + info.Velocity[1]*float32(dt)

	// 	// Apply horizontal damping
	// 	info.Velocity[0] *= VelocityDamping
	// }
}

func (w *World) postupdate(activeBodies []entity.Entity) {
	// for i := range activeBodies {
	// 	info := activeBodies[i].Info()

	// 	x, y := activeBodies[i].Position()
	// 	nX, nY := info.nextPosition[0], info.nextPosition[1]

	// 	if nX == x && nY == y {
	// 		continue
	// 	}

	// 	info.prevPosition[0] = x
	// 	info.prevPosition[1] = y

	// 	activeBodies[i].SetPosition(nX, nY)

	// 	w.bodyGrid.Remove(activeBodies[i])
	// 	w.insert(activeBodies[i], w.bodyGrid)
	// }
}

func (w *World) handleCollisions(activeBodies []entity.Entity) {
	w.contacts = w.contacts[:0]
	for i := range activeBodies {
		id := activeBodies[i]

		bounds := components.GetBounds(id)
		collision := components.GetCollision(id)
		transform := components.GetTransform(id)
		if collision.Mode == components.CollisionModeIgnore {
			continue
		}

		// ========== Check static collisions ==========

		others := w.staticGrid.Query(CalculateAABB(transform, bounds))
		for j := range others {
			otherId := others[j]

			otherCollision := components.GetCollision(otherId)

			if otherCollision.Mode == components.CollisionModeIgnore || others[j].Equals(activeBodies[i]) {
				continue
			}

			if !components.ShouldCollide(collision.Layer, otherCollision.Layer) {
				continue
			}

			// if contact, overlaps := CheckOverlap(activeBodies[i], others[j]); overlaps {
			// 	info.collisions = append(info.collisions, contact)
			// }
		}

		// w.resolveStaticCollisions(info)

		// ========== Check body collisions ==========

		// TASK: Implement dynamic body collision detection and event dispatching
	}
}

// func (w *World) isGrounded(collider Collider, info *ColliderInfo, travelled float32) bool {
// 	minX, minY, maxX, maxY := collider.AABB()
// 	centerX := (minX + maxX) / 2
// 	queryDistance := max(travelled, GroundCheckDistance)

// 	others := w.staticGrid.Query(minX, minY, maxX, maxY+queryDistance)
// 	for _, other := range others {
// 		otherInfo := other.Info()
// 		if otherInfo.Mode == CollisionModeIgnore || info.id == otherInfo.id {
// 			continue
// 		}

// 		if !ShouldCollide(info.Layer, otherInfo.Layer) || !otherInfo.IsFloor() {
// 			continue
// 		}

// 		var surfaceY float32
// 		switch o := other.(type) {
// 		case *BoxCollider:
// 			// Check horizontal overlap for box colliders
// 			oMinX, oMinY, oMaxX, _ := o.AABB()
// 			if maxX <= oMinX || minX >= oMaxX {
// 				continue
// 			}
// 			surfaceY = oMinY

// 		case *TriangleCollider:
// 			// Check if center point is within triangle bounds
// 			oMinX, _, oMaxX, _ := o.AABB()
// 			if centerX < oMinX || centerX > oMaxX {
// 				continue
// 			}
// 			y, found := geom.FindTriangleSurfaceAt(centerX, &o.Triangle)
// 			if !found {
// 				continue
// 			}
// 			surfaceY = y

// 		default:
// 			continue
// 		}

// 		// Check if within ground detection range
// 		distance := surfaceY - maxY
// 		if distance >= -GroundCheckTolerance && distance <= queryDistance {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (w *World) resolveStaticCollisions(info *ColliderInfo) {
// 	if len(info.collisions) == 0 {
// 		return
// 	}

// 	var horizontal *Collision
// 	var vertical *Collision
// 	var slope *Collision

// 	for i := range info.collisions {
// 		contact := &info.collisions[i]

// 		switch contact.other.(type) {
// 		case *TriangleCollider:
// 			if slope == nil || contact.Depth > slope.Depth {
// 				slope = contact
// 			}
// 		default:
// 			if math.Abs(float64(contact.Normal[1])) > math.Abs(float64(contact.Normal[0])) {
// 				if vertical == nil || contact.Depth > vertical.Depth {
// 					vertical = contact
// 				}
// 			} else {
// 				if horizontal == nil || contact.Depth > horizontal.Depth {
// 					horizontal = contact
// 				}
// 			}
// 		}
// 	}

// 	if slope != nil {
// 		w.resolveStaticSlopeCollision(info, slope)
// 	} else if vertical != nil {
// 		w.resolveStaticVerticalCollision(info, vertical)
// 	}

// 	if horizontal != nil {
// 		otherInfo := horizontal.other.Info()
// 		if slope != nil {
// 			if otherInfo.IsWall() {
// 				w.resolveStaticHorizontalCollision(info, horizontal)
// 			}
// 		} else {
// 			w.resolveStaticHorizontalCollision(info, horizontal)
// 		}
// 	}
// }

// func (w *World) resolveStaticSlopeCollision(info *ColliderInfo, col *Collision) {
// 	// Move out of collision along slope normal
// 	info.nextPosition[0] += col.Normal[0] * col.Depth
// 	info.nextPosition[1] += col.Normal[1] * col.Depth

// 	if col.Normal[1] < 0 {
// 		// Colliding with ceiling - stop upward velocity
// 		info.Velocity[1] = 0
// 	} else {
// 		// Redirect velocity along slope surface
// 		dotProduct := info.Velocity[0]*col.Normal[0] + info.Velocity[1]*col.Normal[1]
// 		if dotProduct < 0 {
// 			info.Velocity[0] -= col.Normal[0] * dotProduct
// 			info.Velocity[1] -= col.Normal[1] * dotProduct
// 		}
// 	}
// }

// func (w *World) resolveStaticVerticalCollision(info *ColliderInfo, col *Collision) {
// 	// Move out of collision vertically
// 	info.nextPosition[1] += col.Normal[1] * col.Depth

// 	// Stop velocity if moving into the collision
// 	if (col.Normal[1] < 0 && info.Velocity[1] > 0) || (col.Normal[1] > 0 && info.Velocity[1] < 0) {
// 		info.Velocity[1] = 0
// 	}
// }

// func (w *World) resolveStaticHorizontalCollision(info *ColliderInfo, col *Collision) {
// 	// Move out of collision horizontally
// 	info.nextPosition[0] += col.Normal[0] * col.Depth

// 	// Stop velocity if moving into the collision
// 	if (col.Normal[0] < 0 && info.Velocity[0] > 0) || (col.Normal[0] > 0 && info.Velocity[0] < 0) {
// 		info.Velocity[0] = 0
// 	}
// }
