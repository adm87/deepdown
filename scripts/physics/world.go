package physics

import (
	"math"

	"github.com/adm87/deepdown/scripts/components"
	"github.com/adm87/deepdown/scripts/ecs/entity"
	"github.com/adm87/deepdown/scripts/geom"
	"github.com/adm87/utilities/hash"
	"github.com/adm87/utilities/sparse"
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

func CalculateAABB(transform *components.Transform, bounds *components.Bounds) [4]float32 {
	x, y := transform.Position()
	return calculateAABB([2]float32{x, y}, bounds)
}

func calculateAABB(position [2]float32, bounds *components.Bounds) [4]float32 {
	w, h := bounds.Size()
	ox, oy := bounds.Offset()
	x, y := position[0], position[1]
	return [4]float32{x + ox, y + oy, x + ox + w, y + oy + h}
}

func clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func abs(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

type Contact struct {
	shape  components.CollisionShape
	isWall bool
	depth  float32
	normal [2]float32
}

type PhysicsState struct {
	previousPosition [2]float32
	nextPosition     [2]float32
	aabb             [4]float32
}

type World struct {
	physicsStates *sparse.Set[PhysicsState, entity.Entity]

	staticGrid *hash.Grid[entity.Entity] // Static world colliders
	bodyGrid   *hash.Grid[entity.Entity] // Dynamic and trigger body colliders

	contacts []Contact
}

func NewWorld() *World {
	return &World{
		physicsStates: sparse.NewSet[PhysicsState, entity.Entity](512),

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
	w.physicsStates.Insert(entity, PhysicsState{})
}

func (w *World) RemoveCollider(entity entity.Entity) {
	switch components.GetCollision(entity).Type {
	case components.CollisionTypeStatic:
		w.staticGrid.Remove(entity)
	default:
		w.bodyGrid.Remove(entity)
	}
	w.physicsStates.Remove(entity)
}

func (w *World) QueryStatic(region [4]float32) []entity.Entity {
	return w.staticGrid.Query(region[0], region[1], region[2], region[3])
}

func (w *World) QueryBody(region [4]float32) []entity.Entity {
	return w.bodyGrid.Query(region[0], region[1], region[2], region[3])
}

func (w *World) QueryStaticCells(region [4]float32) []uint64 {
	return w.staticGrid.QueryCells(region[0], region[1], region[2], region[3])
}

func (w *World) QueryBodyCells(region [4]float32) []uint64 {
	return w.bodyGrid.QueryCells(region[0], region[1], region[2], region[3])
}

func (w *World) Update(dt float64, region [4]float32) {
	components.EachPhysics(func(e entity.Entity, phys *components.Physics) {
		if ctype := components.GetCollision(e); ctype.Type == components.CollisionTypeIgnore {
			return
		}

		transform := components.GetTransform(e)
		bounds := components.GetBounds(e)

		x, y := transform.Position()
		aabb := calculateAABB([2]float32{x, y}, bounds)

		phys.Awake = aabb[2] >= region[0] && aabb[0] <= region[2] && aabb[3] >= region[1] && aabb[1] <= region[3]
		if !phys.Awake {
			return
		}

		state := w.physicsStates.Get(e)
		state.previousPosition = [2]float32{x, y}

		if phys.Gravity != 0 {
			velY := clamp(phys.Velocity[1]+Gravity*float32(dt), MaxVelocityRiseSpeed, MaxVelocityFallSpeed)

			_, grounded := w.isGrounded(e, components.GetCollision(e).Layer, aabb, velY*float32(dt))
			phys.OnGround = grounded

			if !phys.OnGround {
				phys.Velocity[1] = velY
			}
		}

		if abs(phys.Velocity[0]) < float32(MinimumVelocityThreshold) {
			phys.Velocity[0] = 0
		}
		if abs(phys.Velocity[1]) < float32(MinimumVelocityThreshold) {
			phys.Velocity[1] = 0
		}

		state.nextPosition[0] = x + phys.Velocity[0]*float32(dt)
		state.nextPosition[1] = y + phys.Velocity[1]*float32(dt)

		phys.Velocity[0] *= VelocityDamping

		state.aabb = calculateAABB(state.nextPosition, bounds)

		collision := components.GetCollision(e)

		if w.checkCollisions(e, collision.Shape, collision.Layer, state,
			w.staticGrid.Query(state.aabb[0], state.aabb[1], state.aabb[2], state.aabb[3])) {
			w.resolveStaticCollisions(state, phys, bounds)
		}
		if w.checkCollisions(e, collision.Shape, collision.Layer, state,
			w.bodyGrid.Query(state.aabb[0], state.aabb[1], state.aabb[2], state.aabb[3])) {
			w.resolveCollisions(state, phys)
		}

		state.aabb = calculateAABB(state.nextPosition, bounds)
		transform.SetPosition(state.nextPosition[0], state.nextPosition[1])

		if state.nextPosition == [2]float32{x, y} {
			return
		}

		w.bodyGrid.Remove(e)
		w.insert(e, w.bodyGrid)
	})
}

func (w *World) checkCollisions(e entity.Entity, shape components.CollisionShape, layer components.CollisionLayer, state *PhysicsState, others []entity.Entity) bool {
	w.contacts = w.contacts[:0]

	for _, other := range others {
		if other.Equals(e) {
			continue
		}

		collision := components.GetCollision(other)
		if collision.Type == components.CollisionTypeIgnore || !components.ShouldCollide(layer, collision.Layer) {
			continue
		}

		x, y := components.GetTransform(other).Position()
		aabb := calculateAABB([2]float32{x, y}, components.GetBounds(other))
		if state.aabb[2] <= aabb[0] || state.aabb[0] >= aabb[2] || state.aabb[3] <= aabb[1] || state.aabb[1] >= aabb[3] {
			continue
		}

		var contact Contact
		var overlaps bool

		switch shape {
		case components.CollisionShapeBox:
			if collision.Shape == components.CollisionShapeBox {
				contact, overlaps = checkAABBvsAABB(state.aabb, aabb)
			} else if collision.Shape == components.CollisionShapeTriangle {
				contact, overlaps = checkAABBvsTriangle(state.aabb, x, y, components.GetTriangleBody(other))
			}

		case components.CollisionShapeTriangle:
			triA := components.GetTriangleBody(e)
			x, y := components.GetTransform(e).Position()
			if collision.Shape == components.CollisionShapeBox {
				contact, overlaps = checkAABBvsTriangle(aabb, x, y, triA)
			} else if collision.Shape == components.CollisionShapeTriangle {
				contact, overlaps = checkTriangleVsTriangle(triA, components.GetTriangleBody(other))
			}
		}

		if overlaps {
			contact.shape = collision.Shape
			contact.isWall = collision.Role&components.CollisionRoleWall != 0
			w.contacts = append(w.contacts, contact)
		}
	}
	return len(w.contacts) > 0
}

func (w *World) isGrounded(id entity.Entity, layer components.CollisionLayer, aabb [4]float32, travelled float32) (float32, bool) {
	centerX := (aabb[0] + aabb[2]) / 2
	distance := max(travelled, GroundCheckDistance)

	for _, other := range w.staticGrid.Query(aabb[0], aabb[1], aabb[2], aabb[3]+distance) {
		if other.Equals(id) {
			continue
		}

		collision := components.GetCollision(other)
		if collision.Type == components.CollisionTypeIgnore || !components.ShouldCollide(layer, collision.Layer) {
			continue
		}

		if collision.Role&components.CollisionRoleFloor == 0 {
			continue
		}

		x, y := components.GetTransform(other).Position()
		otherAABB := calculateAABB([2]float32{x, y}, components.GetBounds(other))

		var surfaceY float32
		switch collision.Shape {
		case components.CollisionShapeBox:
			if aabb[2] <= otherAABB[0] || aabb[0] >= otherAABB[2] {
				continue
			}
			surfaceY = otherAABB[1]

		case components.CollisionShapeTriangle:
			if centerX < otherAABB[0] || centerX > otherAABB[2] {
				continue
			}
			tri := *components.GetTriangleBody(other)
			tri.X, tri.Y = x, y
			y, found := geom.FindTriangleSurfaceAt(centerX, &tri)
			if !found {
				continue
			}
			surfaceY = y

		default:
			continue
		}

		distanceToSurface := surfaceY - aabb[3]
		if distanceToSurface >= -GroundCheckTolerance && distanceToSurface <= distance {
			return surfaceY, true
		}
	}
	return 0, false
}

func (w *World) resolveStaticCollisions(state *PhysicsState, phys *components.Physics, bounds *components.Bounds) {
	var horizontal, vertical, slope *Contact

	for i := range w.contacts {
		contact := &w.contacts[i]

		switch contact.shape {
		case components.CollisionShapeTriangle:
			if slope == nil || contact.depth > slope.depth {
				slope = contact
			}
		default:
			if math.Abs(float64(contact.normal[1])) > math.Abs(float64(contact.normal[0])) {
				if vertical == nil || contact.depth > vertical.depth {
					vertical = contact
				}
			} else {
				if horizontal == nil || contact.depth > horizontal.depth {
					horizontal = contact
				}
			}
		}
	}

	if slope != nil {
		w.resolveStaticSlopeCollision(state, phys, slope)
	} else if vertical != nil {
		w.resolveStaticVerticalCollision(state, phys, vertical)
	}

	if horizontal != nil {
		if slope != nil {
			if horizontal.isWall {
				w.resolveStaticHorizontalCollision(state, phys, horizontal)
			}
		} else {
			w.resolveStaticHorizontalCollision(state, phys, horizontal)
		}
	}
}

func (w *World) resolveStaticSlopeCollision(state *PhysicsState, phys *components.Physics, contact *Contact) {
	state.nextPosition[0] += contact.normal[0] * contact.depth
	state.nextPosition[1] += contact.normal[1] * contact.depth

	dotProduct := phys.Velocity[0]*contact.normal[0] + phys.Velocity[1]*contact.normal[1]
	if dotProduct < 0 {
		phys.Velocity[1] -= contact.normal[1] * dotProduct
	}

	phys.OnGround = contact.normal[1] < 0
}

func (w *World) resolveStaticVerticalCollision(state *PhysicsState, phys *components.Physics, contact *Contact) {
	state.nextPosition[1] += contact.normal[1] * contact.depth
	if (contact.normal[1] < 0 && phys.Velocity[1] > 0) || (contact.normal[1] > 0 && phys.Velocity[1] < 0) {
		phys.Velocity[1] = 0
	}
}

func (w *World) resolveStaticHorizontalCollision(state *PhysicsState, phys *components.Physics, contact *Contact) {
	state.nextPosition[0] += contact.normal[0] * contact.depth
	if (contact.normal[0] < 0 && phys.Velocity[0] > 0) || (contact.normal[0] > 0 && phys.Velocity[0] < 0) {
		phys.Velocity[0] = 0
	}
}

func (w *World) resolveCollisions(state *PhysicsState, phys *components.Physics) {

}

func (w *World) insert(entity entity.Entity, grid *hash.Grid[entity.Entity]) {
	aabb := CalculateAABB(components.GetTransform(entity), components.GetBounds(entity))
	grid.Insert(entity, aabb[0], aabb[1], aabb[2], aabb[3], hash.NoGridPadding)
}
