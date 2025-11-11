package physics

import "github.com/adm87/deepdown/scripts/geom"

func checkAABBvsAABB(aabbA [4]float32, aabbB [4]float32) (Contact, bool) {
	var contact Contact

	overlapX := min(aabbA[2]-aabbB[0], aabbB[2]-aabbA[0])
	overlapY := min(aabbA[3]-aabbB[1], aabbB[3]-aabbA[1])

	if overlapX <= 0 || overlapY <= 0 {
		return contact, false
	}

	centerXA := (aabbA[0] + aabbA[2]) * 0.5
	centerYA := (aabbA[1] + aabbA[3]) * 0.5
	centerXB := (aabbB[0] + aabbB[2]) * 0.5
	centerYB := (aabbB[1] + aabbB[3]) * 0.5

	if overlapX < overlapY {
		if overlapX < 0.01 {
			return contact, false
		}
		if centerXA < centerXB {
			contact.normal = [2]float32{-1, 0}
		} else {
			contact.normal = [2]float32{1, 0}
		}
		contact.depth = overlapX
	} else {
		if overlapY < 0.01 {
			return contact, false
		}
		if centerYA < centerYB {
			contact.normal = [2]float32{0, -1}
		} else {
			contact.normal = [2]float32{0, 1}
		}
		contact.depth = overlapY
	}

	return contact, true
}

// func checkAABBvsTriangle(aabb [4]float32, tri *geom.Triangle) (Contact, bool) {
// 	var contact Contact

// 	// Get triangle vertices in world space
// 	v0x, v0y := tri.GetVertex(0)
// 	v1x, v1y := tri.GetVertex(1)
// 	v2x, v2y := tri.GetVertex(2)

// 	corners := [4][2]float32{
// 		{aabb[0], aabb[1]}, // top-left
// 		{aabb[2], aabb[1]}, // top-right
// 		{aabb[2], aabb[3]}, // bottom-right
// 		{aabb[0], aabb[3]}, // bottom-left
// 	}

// 	for _, corner := range corners {
// 		if tri.ContainsPoint(corner[0], corner[1]) {
// 			contact.normal = tri.SlopeNormal()
// 			surfaceY, found := geom.FindTriangleSurfaceAt(corner[0], tri)
// 			if found {
// 				contact.depth = abs(corner[1] - surfaceY)
// 				return contact, true
// 			}
// 		}
// 	}

// 	triVerts := [][2]float32{
// 		{v0x, v0y},
// 		{v1x, v1y},
// 		{v2x, v2y},
// 	}

// 	for _, vert := range triVerts {
// 		if vert[0] >= aabb[0] && vert[0] <= aabb[2] &&
// 			vert[1] >= aabb[1] && vert[1] <= aabb[3] {
// 			// Vertex inside AABB
// 			contact.normal = tri.SlopeNormal()
// 			// Calculate penetration depth
// 			dx := min(vert[0]-aabb[0], aabb[2]-vert[0])
// 			dy := min(vert[1]-aabb[1], aabb[3]-vert[1])
// 			contact.depth = min(dx, dy)
// 			return contact, true
// 		}
// 	}

// 	aabbCenterX := (aabb[0] + aabb[2]) * 0.5
// 	aabbBottomY := aabb[3]

// 	surfaceY, found := geom.FindTriangleSurfaceAt(aabbCenterX, tri)
// 	if found && aabbBottomY > surfaceY {
// 		contact.depth = aabbBottomY - surfaceY
// 		contact.normal = tri.SlopeNormal()
// 		return contact, true
// 	}

// 	return contact, false
// }

// func abs(v float32) float32 {
// 	if v < 0 {
// 		return -v
// 	}
// 	return v
// }

func checkAABBvsTriangle(aabb [4]float32, x, y float32, tri *geom.Triangle) (Contact, bool) {
	var contact Contact

	slope := tri.Slope()
	corner := tri.Corner()

	// Calculate AABB bottom-center
	aabbCenterX := (aabb[0] + aabb[2]) * 0.5
	aabbBottomY := aabb[3]

	// Check if bottom-center is within triangle's X range
	triMinX := x + min(slope[0][0], slope[1][0], corner[0])
	triMaxX := x + max(slope[0][0], slope[1][0], corner[0])

	if aabbCenterX < triMinX || aabbCenterX > triMaxX {
		return contact, false
	}

	tri.X, tri.Y = x, y
	defer func() {
		tri.X, tri.Y = 0, 0
	}()

	// Find the Y position on the slope at the AABB's center X
	surfaceY, found := geom.FindTriangleSurfaceAt(aabbCenterX, tri)
	if !found {
		return contact, false
	}

	// Check if AABB penetrates the slope
	if aabbBottomY > surfaceY {
		contact.depth = aabbBottomY - surfaceY
		contact.normal = tri.SlopeNormal()
		return contact, true
	}

	return Contact{}, false
}

func checkTriangleVsTriangle(triA *geom.Triangle, triB *geom.Triangle) (Contact, bool) {
	// TASK: Implement if needed later
	return Contact{}, false
}

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
