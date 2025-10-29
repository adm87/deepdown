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
		if depthX < 0.01 {
			return contact, false
		}

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

func BoxVsTriangle(a *BoxCollider, b *TriangleCollider) (Contact, bool) {
	var contact Contact
	minXA, minYA, maxXA, maxYA := a.Bounds()
	minXB, minYB, maxXB, maxYB := b.Bounds()

	if minXA >= maxXB || maxXA <= minXB || minYA >= maxYB || maxYA <= minYB {
		return contact, false
	}

	centerX := (minXA + maxXA) / 2

	if surfaceY, found := findTriangleSurfaceAt(centerX, &b.Triangle); found {
		if maxYA >= surfaceY && minYA < surfaceY {
			depth := maxYA - surfaceY
			if depth < 0.01 {
				return contact, false
			}
			contact = Contact{
				Point:  [2]float32{centerX, surfaceY},
				Normal: [2]float32{b.Triangle.Nx, b.Triangle.Ny},
				Depth:  depth,
			}
			return contact, true
		}
		if minYA <= surfaceY && maxYA > surfaceY {
			depth := surfaceY - minYA
			if depth < 0.01 {
				return contact, false
			}
			contact = Contact{
				Point:  [2]float32{centerX, surfaceY},
				Normal: [2]float32{b.Triangle.Nx, b.Triangle.Ny},
				Depth:  depth,
			}
			return contact, true
		}
	}

	return contact, false
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
