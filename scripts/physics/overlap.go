package physics

// func CheckOverlap(colliderA, colliderB Collider) (Collision, bool) {
// 	switch a := colliderA.(type) {
// 	case *BoxCollider:
// 		switch b := colliderB.(type) {
// 		case *BoxCollider:
// 			return BoxVsBox(a, b)
// 		case *TriangleCollider:
// 			return BoxVsTriangle(a, b)
// 		}
// 	case *TriangleCollider:
// 		switch b := colliderB.(type) {
// 		case *BoxCollider:
// 			return BoxVsTriangle(b, a)
// 		case *TriangleCollider:
// 			return TriangleVsTriangle(a, b)
// 		}
// 	}
// 	return Collision{}, false
// }

// func BoxVsBox(b1, b2 *BoxCollider) (Collision, bool) {
// 	var contact Collision

// 	minXA, minYA, maxXA, maxYA := b1.AABB()
// 	minXB, minYB, maxXB, maxYB := b2.AABB()

// 	// Check for separation
// 	if minXA >= maxXB || maxXA <= minXB || minYA >= maxYB || maxYA <= minYB {
// 		return contact, false
// 	}

// 	// Calculate overlap depths
// 	overlapX := min(maxXA-minXB, maxXB-minXA)
// 	overlapY := min(maxYA-minYB, maxYB-minYA)

// 	// Compute centers
// 	centerXA := (minXA + maxXA) * 0.5
// 	centerYA := (minYA + maxYA) * 0.5
// 	centerXB := (minXB + maxXB) * 0.5
// 	centerYB := (minYB + maxYB) * 0.5

// 	// Resolve along axis of least penetration
// 	if overlapX < overlapY {
// 		if overlapX < 0.01 {
// 			return contact, false
// 		}
// 		// Normal points away from surface (standard convention)
// 		if centerXA < centerXB {
// 			contact.Normal = [2]float32{-1, 0}
// 		} else {
// 			contact.Normal = [2]float32{1, 0}
// 		}
// 		contact.Depth = overlapX
// 	} else {
// 		if overlapY < 0.01 {
// 			return contact, false
// 		}
// 		if centerYA < centerYB {
// 			contact.Normal = [2]float32{-1, 0}
// 		} else {
// 			contact.Normal = [2]float32{1, 0}
// 		}
// 		contact.Normal = [2]float32{0, contact.Normal[0]} // flip to Y axis
// 		contact.Depth = overlapY
// 	}

// 	contact.other = b2

// 	return contact, true
// }

// func BoxVsTriangle(box *BoxCollider, tri *TriangleCollider) (Collision, bool) {
// 	var contact Collision

// 	minXA, minYA, maxXA, maxYA := box.AABB()
// 	minXB, minYB, maxXB, maxYB := tri.AABB()

// 	// Quick AABB check
// 	if minXA >= maxXB || maxXA <= minXB || minYA >= maxYB || maxYA <= minYB {
// 		return contact, false
// 	}

// 	// Get triangle properties
// 	slope := tri.Slope()
// 	corner := tri.Corner()

// 	// Calculate box bottom-center
// 	boxCenterX := (minXA + maxXA) * 0.5
// 	boxBottomY := maxYA

// 	// Check if bottom-center is within triangle's X range
// 	triMinX := tri.X + min(slope[0][0], slope[1][0], corner[0])
// 	triMaxX := tri.X + max(slope[0][0], slope[1][0], corner[0])

// 	if boxCenterX < triMinX || boxCenterX > triMaxX {
// 		return contact, false
// 	}

// 	// Find the Y position on the slope at the box's center X
// 	surfaceY, found := geom.FindTriangleSurfaceAt(boxCenterX, &tri.Triangle)
// 	if !found {
// 		return contact, false
// 	}

// 	// Check if box penetrates the slope
// 	if boxBottomY > surfaceY {
// 		// Penetrating
// 		contact.Depth = boxBottomY - surfaceY
// 		contact.Normal = tri.SlopeNormal()
// 		contact.other = tri
// 		return contact, true
// 	}

// 	return contact, false
// }

// func TriangleVsTriangle(t1, t2 *TriangleCollider) (Collision, bool) {
// 	var contact Collision

// 	return contact, false
// }
