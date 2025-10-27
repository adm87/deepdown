package collision

import "github.com/adm87/deepdown/scripts/geom"

type Contact struct {
	Point  [2]float32
	Normal [2]float32
	Depth  float32
}

func BoxVsBox(a, b *BoxCollider) (Contact, bool) {
	var contact Contact

	minXA, minYA, maxXA, maxYA := a.Bounds()
	minXB, minYB, maxXB, maxYB := b.Bounds()

	// Check for overlap
	if minXA >= maxXB || maxXA <= minXB || minYA >= maxYB || maxYA <= minYB {
		return contact, false
	}

	// Calculate overlap depths
	depthX := min(maxXA-minXB, maxXB-minXA)
	depthY := min(maxYA-minYB, maxYB-minYA)

	// Choose separation axis based on minimum depth
	if depthX < depthY {
		normalX := float32(1)
		if (minXA + maxXA) > (minXB + maxXB) {
			normalX = -1
		}
		contact = Contact{
			Point:  [2]float32{(minXA + maxXA) / 2, (minYA + maxYA) / 2},
			Normal: [2]float32{normalX, 0},
			Depth:  depthX,
		}
	} else {
		normalY := float32(1)
		if (minYA + maxYA) > (minYB + maxYB) {
			normalY = -1
		}
		contact = Contact{
			Point:  [2]float32{(minXA + maxXA) / 2, (minYA + maxYA) / 2},
			Normal: [2]float32{0, normalY},
			Depth:  depthY,
		}
	}

	return contact, true
}

func BoxVsTriangle(a *BoxCollider, b *TriangleCollider, velocityY float32) (Contact, bool) {
	var contact Contact
	minXA, minYA, maxXA, maxYA := a.Bounds()
	minXB, minYB, maxXB, maxYB := b.Bounds()

	if minXA >= maxXB || maxXA <= minXB || minYA >= maxYB || maxYA <= minYB {
		return contact, false
	}

	triangle := b.Triangle
	testX := (minXA + maxXA) * 0.5
	var testY float32
	var normal [2]float32

	if velocityY >= 0 {
		// Falling or stationary: test bottom
		testY = maxYA
		normal = [2]float32{triangle.Nx, triangle.Ny}
	} else {
		// Jumping: test top
		testY = minYA
		normal = [2]float32{-triangle.Nx, -triangle.Ny} // Flip normal for ceiling
	}

	if testX < minXB || testX > maxXB {
		return contact, false
	}

	surfaceY, found := findTriangleSurfaceAt(testX, &triangle)
	if !found {
		return contact, false
	}
	const tolerance = 0.1
	if testY <= surfaceY+tolerance {
		return contact, false
	}

	if velocityY >= 0 {
		if testY <= surfaceY+tolerance {
			return contact, false
		}
		contact.Depth = testY - surfaceY
	} else {
		if testY >= surfaceY-tolerance {
			return contact, false
		}
		contact.Depth = surfaceY - testY
	}

	contact.Point = [2]float32{testX, surfaceY}
	contact.Normal = normal

	return contact, true
}

func abs(a float32) float32 {
	if a < 0 {
		return -a
	}
	return a
}

func findTriangleSurfaceAt(x float32, triangle *geom.Triangle) (surfaceY float32, found bool) {
	points := triangle.Points()

	// Test each edge
	for i := range 3 {
		x1, y1 := triangle.X+points[i*2], triangle.Y+points[i*2+1]
		x2, y2 := triangle.X+points[((i+1)%3)*2], triangle.Y+points[((i+1)%3)*2+1]

		// Skip vertical or horizontal edges
		if abs(x2-x1) < 0.1 || abs(y2-y1) < 0.1 {
			continue
		}

		// Check if edge spans x coordinate
		if x >= min(x1, x2) && x <= max(x1, x2) {
			t := (x - x1) / (x2 - x1)
			return y1 + t*(y2-y1), true
		}
	}

	return 0, false
}

func BoxVsPolygon(a *BoxCollider, b *PolygonCollider) (Contact, bool) {
	var contact Contact

	minXA, minYA, maxXA, maxYA := a.Bounds()
	minXB, minYB, maxXB, maxYB := b.Bounds()

	if minXA >= maxXB || maxXA <= minXB || minYA >= maxYB || maxYA <= minYB {
		return contact, false
	}

	return contact, contact.Depth > 0
}
