package level

import (
	"github.com/adm87/deepdown/scripts/collision"
)

func (l *Level) OnCollisionEnter(colliderA, colliderB collision.Collider) {
	l.OnCollision(colliderA, colliderB) // Forward for now, handle specifics later
}

func (l *Level) OnCollision(colliderA, colliderB collision.Collider) {
	infoA := colliderA.Info()
	infoB := colliderB.Info()

	switch {
	// Static immovable world collision
	case infoA.Type == collision.Dynamic && infoB.Type == collision.Static:
		l.staticCollision(colliderA, colliderB)
	case infoA.Type == collision.Static && infoB.Type == collision.Dynamic:
		l.staticCollision(colliderB, colliderA)
	}
}

// Note: colliderA is dynamic, colliderB is static
// Note: colliderA will always be a BoxCollider for the time being
func (l *Level) staticCollision(colliderA, colliderB collision.Collider) {
	box := colliderA.(*collision.BoxCollider)

	switch other := colliderB.(type) {
	case *collision.BoxCollider:
		if contact, overlaps := collision.BoxVsBox(box, other); overlaps {
			minXA, _, maxXA, _ := box.Bounds()
			minXB, _, maxXB, _ := other.Bounds()

			center := (minXA + maxXA) / 2

			if contact.Normal[0] != 0 && center > minXB && center < maxXB {
				depth := min(center-minXB, maxXB-center)
				box.X -= contact.Normal[0] * depth
				box.Velocity[0] = 0
			}

			if contact.Normal[1] != 0 {
				if contact.Normal[1] < 0 {
					box.Y -= contact.Normal[1] * contact.Depth
				} else if center > minXB && center < maxXB && box.Velocity[1] > 0 {
					box.Y -= contact.Normal[1] * contact.Depth
				}
				box.Velocity[1] = 0
				if contact.Normal[1] > 0 {
					l.onGround = true
				}
			}
		}

	case *collision.TriangleCollider:
		if contact, overlaps := collision.BoxVsTriangle(box, other); overlaps {
			if contact.Normal[1] != 0 {
				box.Y -= contact.Normal[1] * contact.Depth
				if contact.Normal[1] > 0 {
					box.Velocity[1] = 0
					l.onGround = true
				}
			}
		}

	case *collision.PolygonCollider:
		if contact, overlaps := collision.BoxVsPolygon(box, other); overlaps {
			_ = contact // TODO handle polygon collision response
		}
	}
}
