package geom

import (
	"testing"
)

// Test case types for cleaner test organization
type polygonIntersectionTest struct {
	name     string
	poly1    *Polygon
	poly2    *Polygon
	expected bool
}

type polygonRectIntersectionTest struct {
	name     string
	poly     *Polygon
	rect     Rectangle
	expected bool
}

type polygonContainsPointTest struct {
	name     string
	poly     *Polygon
	x, y     float32
	expected bool
}

type polygonEdgeCaseTest struct {
	name        string
	poly        *Polygon
	description string
}

// Global test shapes with properly ordered vertices
var (
	testSquare          = NewPolygon(0, 0, []float32{0, 0, 10, 0, 10, 10, 0, 10})
	testTriangle        = NewPolygon(0, 0, []float32{5, 0, 10, 10, 0, 10})
	testDiamond         = NewPolygon(0, 0, []float32{5, 0, 10, 5, 5, 10, 0, 5})
	testLShape          = NewPolygon(0, 0, []float32{0, 0, 5, 0, 5, 5, 10, 5, 10, 10, 0, 10})
	testRectangle       = NewPolygon(5, 5, []float32{0, 0, 8, 0, 8, 6, 0, 6})
	overlappingSquare   = NewPolygon(5, 5, []float32{0, 0, 10, 0, 10, 10, 0, 10})
	separateSquare      = NewPolygon(20, 20, []float32{0, 0, 10, 0, 10, 10, 0, 10})
	adjacentSquare      = NewPolygon(10, 0, []float32{0, 0, 10, 0, 10, 10, 0, 10})
	overlappingTriangle = NewPolygon(5, 0, []float32{5, 0, 10, 10, 0, 10})
	separateTriangle    = NewPolygon(20, 20, []float32{0, 0, 6, 0, 6, 6})
	overlappingDiamond  = NewPolygon(5, 5, []float32{4, 0, 8, 4, 4, 8, 0, 4})
	separateDiamond     = NewPolygon(20, 20, []float32{3, 0, 6, 3, 3, 6, 0, 3})
	overlappingLShape   = NewPolygon(5, 5, []float32{0, 0, 5, 0, 5, 5, 10, 5, 10, 10, 0, 10})
	separateLShape      = NewPolygon(20, 20, []float32{0, 0, 4, 0, 4, 4, 8, 4, 8, 8, 0, 8})

	// Edge case polygons for degenerate testing
	emptyPolygon      = NewPolygon(0, 0, []float32{})
	singlePoint       = NewPolygon(0, 0, []float32{5, 5})
	twoPoints         = NewPolygon(0, 0, []float32{0, 0, 10, 10})
	collinearTriangle = NewPolygon(0, 0, []float32{0, 0, 5, 0, 10, 0})
	selfIntersecting  = NewPolygon(0, 0, []float32{0, 0, 10, 10, 10, 0, 0, 10})
	tinyPolygon       = NewPolygon(0, 0, []float32{0, 0, 1e-6, 0, 1e-6, 1e-6})
	precisionPolygon  = NewPolygon(0, 0, []float32{0.1, 0.1, 0.9, 0.1, 0.9, 0.9, 0.1, 0.9})

	// Containment testing polygons
	largeSquare      = NewPolygon(0, 0, []float32{0, 0, 20, 0, 20, 20, 0, 20})             // 20x20 square
	smallSquare      = NewPolygon(5, 5, []float32{0, 0, 5, 0, 5, 5, 0, 5})                 // 5x5 square inside large
	tinyTriangle     = NewPolygon(8, 8, []float32{0, 0, 2, 0, 1, 2})                       // Tiny triangle inside large
	identicalSquare  = NewPolygon(0, 0, []float32{0, 0, 20, 0, 20, 20, 0, 20})             // Identical to large square
	edgeSquare       = NewPolygon(15, 15, []float32{0, 0, 5, 0, 5, 5, 0, 5})               // Square at edge
	partialSquare    = NewPolygon(18, 18, []float32{0, 0, 5, 0, 5, 5, 0, 5})               // Partially outside
	outsideSquare    = NewPolygon(25, 25, []float32{0, 0, 5, 0, 5, 5, 0, 5})               // Completely outside
	smallSquarePoly  = NewPolygon(10, 10, []float32{0, 0, 5, 0, 5, 5, 0, 5})               // 5x5 square at (10,10)
	tinyTrianglePoly = NewPolygon(15, 15, []float32{0, 0, 3, 0, 1.5, 3})                   // Tiny triangle at (15,15)
	edgePoly         = NewPolygon(25, 25, []float32{0, 0, 5, 0, 5, 5, 0, 5})               // At rectangle edge
	partialPoly      = NewPolygon(28, 28, []float32{0, 0, 5, 0, 5, 5, 0, 5})               // Partially outside
	outsidePoly      = NewPolygon(35, 35, []float32{0, 0, 5, 0, 5, 5, 0, 5})               // Completely outside
	largePoly        = NewPolygon(0, 0, []float32{0, 0, 40, 0, 40, 40, 0, 40})             // 40x40 square
	hugePolygon      = NewPolygon(0, 0, []float32{0, 0, 1000, 0, 1000, 1000, 0, 1000})     // 1000x1000
	microPolygon     = NewPolygon(500, 500, []float32{0, 0, 0.1, 0, 0.1, 0.1, 0, 0.1})     // 0.1x0.1 at center
	nanoPolygon      = NewPolygon(500, 500, []float32{0, 0, 1e-3, 0, 1e-3, 1e-3, 0, 1e-3}) // Microscopic

	// Containment testing rectangles
	largeRect = NewRectangle(0, 0, 30, 30) // Large rectangle for containment tests
)

// TestShapeCreation verifies basic polygon properties
func TestShapeCreation(t *testing.T) {
	// Verify basic properties
	if testSquare.VertexCount() != 4 {
		t.Errorf("Square should have 4 vertices, got %d", testSquare.VertexCount())
	}
	if testRectangle.VertexCount() != 4 {
		t.Errorf("Rectangle should have 4 vertices, got %d", testRectangle.VertexCount())
	}
	if testTriangle.VertexCount() != 3 {
		t.Errorf("Triangle should have 3 vertices, got %d", testTriangle.VertexCount())
	}
	if testLShape.VertexCount() != 6 {
		t.Errorf("L-shape should have 6 vertices, got %d", testLShape.VertexCount())
	}
	if testDiamond.VertexCount() != 4 {
		t.Errorf("Diamond should have 4 vertices, got %d", testDiamond.VertexCount())
	}

	// Verify shapes are not empty
	if testSquare.IsEmpty() {
		t.Error("Square should not be empty")
	}
	if testRectangle.IsEmpty() {
		t.Error("Rectangle should not be empty")
	}
	if testTriangle.IsEmpty() {
		t.Error("Triangle should not be empty")
	}
}

// TestBasicPolygonIntersections tests basic polygon-to-polygon intersections
func TestBasicPolygonIntersections(t *testing.T) {
	tests := []polygonIntersectionTest{
		{"Overlapping squares", &testSquare, &overlappingSquare, true},
		{"Non-overlapping squares", &testSquare, &separateSquare, false},
		{"Adjacent squares (touching edge)", &testSquare, &adjacentSquare, true},
		{"Triangle overlapping square", &testSquare, &overlappingTriangle, true},
		{"Triangle outside square", &testSquare, &separateTriangle, false},
		{"Diamond overlapping square", &testSquare, &overlappingDiamond, true},
		{"Diamond outside square", &testSquare, &separateDiamond, false},
		{"L-shape overlapping square", &testSquare, &overlappingLShape, true},
		{"L-shape outside square", &testSquare, &separateLShape, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.poly1.IntersectsPolygon(tt.poly2)
			if result != tt.expected {
				t.Errorf("IntersectsPolygon() = %v, expected %v", result, tt.expected)
			}

			// Test symmetry (intersection should be commutative)
			result2 := tt.poly2.IntersectsPolygon(tt.poly1)
			if result2 != tt.expected {
				t.Errorf("IntersectsPolygon() symmetry failed: %v != %v", result, result2)
			}
		})
	}
}

// TestBasicRectIntersections tests polygon-to-rectangle intersections
func TestBasicRectIntersections(t *testing.T) {
	tests := []polygonRectIntersectionTest{
		{"Square overlapping rectangle", &testSquare, NewRectangle(5, 5, 10, 10), true},
		{"Square outside rectangle", &testSquare, NewRectangle(20, 20, 10, 10), false},
		{"Square touching rectangle edge", &testSquare, NewRectangle(10, 0, 10, 10), true},
		{"Triangle overlapping rectangle", &testTriangle, NewRectangle(2, 2, 6, 6), true},
		{"Triangle outside rectangle", &testTriangle, NewRectangle(20, 20, 6, 6), false},
		{"Diamond overlapping rectangle", &testDiamond, NewRectangle(2, 2, 6, 6), true},
		{"Diamond outside rectangle", &testDiamond, NewRectangle(20, 20, 6, 6), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.poly.IntersectsRect(&tt.rect)
			if result != tt.expected {
				t.Errorf("IntersectsRect() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestContainsPoint tests point-in-polygon detection
func TestContainsPoint(t *testing.T) {
	tests := []polygonContainsPointTest{
		{"Square center", &testSquare, 5, 5, true},
		{"Square corner", &testSquare, 0, 0, true},
		{"Square edge", &testSquare, 5, 0, true},
		{"Square outside", &testSquare, 15, 15, false},
		{"Triangle center", &testTriangle, 5, 7, true},
		{"Triangle vertex", &testTriangle, 5, 0, true},
		{"Triangle outside", &testTriangle, 15, 15, false},
		{"Diamond center", &testDiamond, 5, 5, true},
		{"Diamond outside", &testDiamond, 0, 0, false},
		{"L-shape inside corner", &testLShape, 2, 2, true},
		{"L-shape outside corner", &testLShape, 7, 2, false},
		{"L-shape inner notch", &testLShape, 7, 7, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.poly.ContainsPoint(tt.x, tt.y)
			if result != tt.expected {
				t.Errorf("ContainsPoint(%v, %v) = %v, expected %v", tt.x, tt.y, result, tt.expected)
			}
		})
	}
}

// TestEdgeBoundaryConditions tests points exactly on polygon edges
func TestEdgeBoundaryConditions(t *testing.T) {
	tests := []polygonContainsPointTest{
		{"Square bottom edge middle", &testSquare, 5, 0, true},
		{"Square right edge middle", &testSquare, 10, 5, true},
		{"Square top edge middle", &testSquare, 5, 10, true},
		{"Square left edge middle", &testSquare, 0, 5, true},
		{"Triangle bottom edge middle", &testTriangle, 7.5, 10, true},
		{"Triangle right edge middle", &testTriangle, 7.5, 5, true},
		{"Triangle left edge middle", &testTriangle, 2.5, 5, true},
		{"Diamond top edge middle", &testDiamond, 7.5, 2.5, true},
		{"Diamond right edge middle", &testDiamond, 7.5, 7.5, true},
		{"Diamond bottom edge middle", &testDiamond, 2.5, 7.5, true},
		{"Diamond left edge middle", &testDiamond, 2.5, 2.5, true},
		{"L-shape bottom edge", &testLShape, 2.5, 0, true},
		{"L-shape inner vertical edge", &testLShape, 5, 2.5, true},
		{"L-shape inner horizontal edge", &testLShape, 7.5, 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.poly.ContainsPoint(tt.x, tt.y)
			if result != tt.expected {
				t.Errorf("ContainsPoint(%v, %v) = %v, expected %v", tt.x, tt.y, result, tt.expected)
			}
		})
	}
}

// TestPolygonEdgeCases tests degenerate and edge case polygons
func TestPolygonEdgeCases(t *testing.T) {
	tests := []polygonEdgeCaseTest{
		{"Empty polygon", &emptyPolygon, "No vertices"},
		{"Single point", &singlePoint, "Only one vertex"},
		{"Two points", &twoPoints, "Line segment"},
		{"Collinear triangle", &collinearTriangle, "All vertices on same line"},
		{"Self-intersecting", &selfIntersecting, "Bowtie shape"},
		{"Tiny polygon", &tinyPolygon, "Microscopic dimensions"},
		{"Precision polygon", &precisionPolygon, "Decimal coordinates"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isEmpty := tt.poly.IsEmpty()
			vertCount := tt.poly.VertexCount()
			result := tt.poly.ContainsPoint(0.5, 0.5)
			t.Logf("%s: isEmpty=%v, vertCount=%d, contains(0.5,0.5)=%v",
				tt.description, isEmpty, vertCount, result)
		})
	}
}

// TestDegeneratePolygonBehavior tests specific behaviors of degenerate polygons
func TestDegeneratePolygonBehavior(t *testing.T) {
	if !emptyPolygon.IsEmpty() {
		t.Error("Empty polygon should report as empty")
	}
	if emptyPolygon.VertexCount() != 0 {
		t.Errorf("Empty polygon vertex count should be 0, got %d", emptyPolygon.VertexCount())
	}
	if emptyPolygon.ContainsPoint(0, 0) {
		t.Error("Empty polygon should not contain any points")
	}
	if singlePoint.ContainsPoint(5, 5) {
		t.Error("Single point should not contain points (insufficient vertices)")
	}
	if twoPoints.ContainsPoint(5, 5) {
		t.Error("Line segment should not contain points (insufficient vertices)")
	}
}

// TestFloatingPointPrecision tests precision edge cases
func TestFloatingPointPrecision(t *testing.T) {
	tests := []polygonContainsPointTest{
		{"Tiny polygon center", &tinyPolygon, 5e-7, 5e-7, true},
		{"Tiny polygon outside", &tinyPolygon, 2e-6, 2e-6, false},
		{"Precision polygon center", &precisionPolygon, 0.5, 0.5, true},
		{"Precision polygon outside", &precisionPolygon, 0.05, 0.05, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.poly.ContainsPoint(tt.x, tt.y)
			if result != tt.expected {
				t.Errorf("ContainsPoint(%v, %v) = %v, expected %v", tt.x, tt.y, result, tt.expected)
			}
		})
	}
}

// ========== Rectangle-specific tests ==========

type rectangleContainsTest struct {
	name     string
	rect     Rectangle
	x, y     float32
	expected bool
}

type rectangleIntersectionTest struct {
	name     string
	rect1    Rectangle
	rect2    Rectangle
	expected bool
}

type rectangleCenterTest struct {
	name      string
	rect      Rectangle
	expectedX float32
	expectedY float32
}

// TestRectangleContains tests the Contains method
func TestRectangleContains(t *testing.T) {
	rect := NewRectangle(10, 20, 30, 40)

	tests := []rectangleContainsTest{
		{"Center point", rect, 25, 40, true},
		{"Top-left corner", rect, 10, 20, true},
		{"Top-right corner", rect, 40, 20, true},
		{"Bottom-left corner", rect, 10, 60, true},
		{"Bottom-right corner", rect, 40, 60, true},
		{"Left edge middle", rect, 10, 40, true},
		{"Right edge middle", rect, 40, 40, true},
		{"Top edge middle", rect, 25, 20, true},
		{"Bottom edge middle", rect, 25, 60, true},
		{"Inside near center", rect, 25, 35, true},
		{"Outside left", rect, 5, 40, false},
		{"Outside right", rect, 45, 40, false},
		{"Outside above", rect, 25, 15, false},
		{"Outside below", rect, 25, 65, false},
		{"Outside diagonal", rect, 5, 15, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rect.Contains(tt.x, tt.y)
			if result != tt.expected {
				t.Errorf("Contains(%v, %v) = %v, expected %v", tt.x, tt.y, result, tt.expected)
			}
		})
	}
}

// TestRectangleIntersects tests the Intersects method
func TestRectangleIntersects(t *testing.T) {
	baseRect := NewRectangle(10, 10, 20, 20)

	tests := []rectangleIntersectionTest{
		{"Identical rectangles", baseRect, NewRectangle(10, 10, 20, 20), true},
		{"Overlapping rectangles", baseRect, NewRectangle(20, 20, 20, 20), true},
		{"Touching edge (right)", baseRect, NewRectangle(30, 10, 10, 20), false},  // Strict inequality in Intersects
		{"Touching edge (bottom)", baseRect, NewRectangle(10, 30, 20, 10), false}, // Strict inequality in Intersects
		{"Touching corner", baseRect, NewRectangle(30, 30, 10, 10), false},        // Strict inequality in Intersects
		{"Inside rectangle", baseRect, NewRectangle(15, 15, 10, 10), true},
		{"Containing rectangle", baseRect, NewRectangle(5, 5, 30, 30), true},
		{"Non-overlapping (right)", baseRect, NewRectangle(35, 10, 10, 20), false},
		{"Non-overlapping (left)", baseRect, NewRectangle(0, 10, 5, 20), false},
		{"Non-overlapping (above)", baseRect, NewRectangle(10, 0, 20, 5), false},
		{"Non-overlapping (below)", baseRect, NewRectangle(10, 35, 20, 10), false},
		{"Non-overlapping diagonal", baseRect, NewRectangle(35, 35, 10, 10), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rect1.Intersects(&tt.rect2)
			if result != tt.expected {
				t.Errorf("Intersects() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestRectangleCenter tests the Center method
func TestRectangleCenter(t *testing.T) {
	tests := []rectangleCenterTest{
		{"Standard rectangle", NewRectangle(10, 20, 30, 40), 25, 40},
		{"Origin rectangle", NewRectangle(0, 0, 10, 10), 5, 5},
		{"Unit rectangle", NewRectangle(5, 5, 1, 1), 5.5, 5.5},
		{"Wide rectangle", NewRectangle(0, 0, 100, 10), 50, 5},
		{"Tall rectangle", NewRectangle(0, 0, 10, 100), 5, 50},
		{"Negative coordinates", NewRectangle(-10, -20, 20, 40), 0, 0},
		{"Fractional dimensions", NewRectangle(1.5, 2.5, 3.0, 4.0), 3.0, 4.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y := tt.rect.Center()
			if x != tt.expectedX || y != tt.expectedY {
				t.Errorf("Center() = (%v, %v), expected (%v, %v)", x, y, tt.expectedX, tt.expectedY)
			}
		})
	}
}

// TestRectangleEdgeCases tests edge cases and degenerate rectangles
func TestRectangleEdgeCases(t *testing.T) {
	t.Run("Zero width rectangle", func(t *testing.T) {
		rect := NewRectangle(10, 10, 0, 20)

		// Should still contain points on the line
		if !rect.Contains(10, 15) {
			t.Error("Zero width rectangle should contain points on its line")
		}
		if rect.Contains(11, 15) {
			t.Error("Zero width rectangle should not contain points outside its line")
		}

		// Center should still work
		x, y := rect.Center()
		if x != 10 || y != 20 {
			t.Errorf("Zero width rectangle center = (%v, %v), expected (10, 20)", x, y)
		}
	})

	t.Run("Zero height rectangle", func(t *testing.T) {
		rect := NewRectangle(10, 10, 20, 0)

		// Should still contain points on the line
		if !rect.Contains(15, 10) {
			t.Error("Zero height rectangle should contain points on its line")
		}
		if rect.Contains(15, 11) {
			t.Error("Zero height rectangle should not contain points outside its line")
		}

		// Center should still work
		x, y := rect.Center()
		if x != 20 || y != 10 {
			t.Errorf("Zero height rectangle center = (%v, %v), expected (20, 10)", x, y)
		}
	})

	t.Run("Zero area rectangle (point)", func(t *testing.T) {
		rect := NewRectangle(5, 5, 0, 0)

		// Should contain only its own point
		if !rect.Contains(5, 5) {
			t.Error("Zero area rectangle should contain its own point")
		}
		if rect.Contains(5.1, 5) {
			t.Error("Zero area rectangle should not contain nearby points")
		}

		// Center should be the point itself
		x, y := rect.Center()
		if x != 5 || y != 5 {
			t.Errorf("Zero area rectangle center = (%v, %v), expected (5, 5)", x, y)
		}
	})

	t.Run("Negative dimensions", func(t *testing.T) {
		rect := NewRectangle(10, 10, -5, -5)

		// For Contains: x >= 10 && x <= 10 + (-5) && y >= 10 && y <= 10 + (-5)
		// Which is: x >= 10 && x <= 5 && y >= 10 && y <= 5
		// This is impossible (x cannot be >= 10 AND <= 5), so nothing should be contained
		if rect.Contains(10, 10) {
			t.Error("Rectangle with negative dimensions should not contain any points")
		}
		if rect.Contains(7, 7) {
			t.Error("Rectangle with negative dimensions should not contain any points")
		}
		if rect.Contains(5, 5) {
			t.Error("Rectangle with negative dimensions should not contain any points")
		}

		// Center calculation should still work correctly
		x, y := rect.Center()
		if x != 7.5 || y != 7.5 {
			t.Errorf("Rectangle with negative dimensions center = (%v, %v), expected (7.5, 7.5)", x, y)
		}
	})
}

// TestRectangleAABBInterface tests the AABB interface methods
func TestRectangleAABBInterface(t *testing.T) {
	rect := NewRectangle(10, 20, 30, 40)

	t.Run("Min coordinates", func(t *testing.T) {
		minX, minY := rect.Min()
		if minX != 10 || minY != 20 {
			t.Errorf("Min() = (%v, %v), expected (10, 20)", minX, minY)
		}
	})

	t.Run("Max coordinates", func(t *testing.T) {
		maxX, maxY := rect.Max()
		if maxX != 40 || maxY != 60 {
			t.Errorf("Max() = (%v, %v), expected (40, 60)", maxX, maxY)
		}
	})
}

// ========== Full Containment Intersection Tests ==========

// TestPolygonContainmentIntersections tests cases where one polygon is fully contained within another
func TestPolygonContainmentIntersections(t *testing.T) {
	tests := []polygonIntersectionTest{
		// Small polygon fully inside large polygon
		{"Small square inside large square", &largeSquare, &smallSquare, true},
		{"Large square contains small square", &smallSquare, &largeSquare, true}, // Symmetry test

		// Tiny polygon fully inside large polygon
		{"Tiny triangle inside large square", &largeSquare, &tinyTriangle, true},
		{"Large square contains tiny triangle", &tinyTriangle, &largeSquare, true},

		// Identical polygons (perfect containment)
		{"Identical squares", &largeSquare, &identicalSquare, true},

		// Edge case: Small polygon at edge of large polygon (still contained)
		{"Small square at edge", &largeSquare, &edgeSquare, true},

		// Small polygon partially outside (not fully contained, but intersecting)
		{"Small square partially outside", &largeSquare, &partialSquare, true},

		// Small polygon completely outside
		{"Small square completely outside", &largeSquare, &outsideSquare, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.poly1.IntersectsPolygon(tt.poly2)
			if result != tt.expected {
				t.Errorf("IntersectsPolygon() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestPolygonRectContainmentIntersections tests cases where polygons and rectangles fully contain each other
func TestPolygonRectContainmentIntersections(t *testing.T) {
	tests := []polygonRectIntersectionTest{
		// Small polygons fully inside large rectangle
		{"Small square polygon inside large rectangle", &smallSquarePoly, largeRect, true},
		{"Tiny triangle polygon inside large rectangle", &tinyTrianglePoly, largeRect, true},

		// Polygon at edge of rectangle (still contained)
		{"Polygon at rectangle edge", &edgePoly, largeRect, true},

		// Polygon partially outside rectangle
		{"Polygon partially outside rectangle", &partialPoly, largeRect, true},

		// Polygon completely outside rectangle
		{"Polygon completely outside rectangle", &outsidePoly, largeRect, false},

		// Small rectangle fully inside large polygon
		{"Large polygon contains small rectangle", &largePoly, NewRectangle(10, 10, 8, 8), true},
		{"Large polygon contains tiny rectangle", &largePoly, NewRectangle(20, 20, 2, 2), true},

		// Rectangle at edge of polygon (still contained)
		{"Rectangle at polygon edge", &largePoly, NewRectangle(35, 35, 5, 5), true},

		// Rectangle partially outside polygon
		{"Rectangle partially outside polygon", &largePoly, NewRectangle(38, 38, 5, 5), true},

		// Rectangle completely outside polygon
		{"Rectangle completely outside polygon", &largePoly, NewRectangle(45, 45, 5, 5), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.poly.IntersectsRect(&tt.rect)
			if result != tt.expected {
				t.Errorf("IntersectsRect() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestMicroContainmentCases tests very small shapes contained within larger ones (precision edge cases)
func TestMicroContainmentCases(t *testing.T) {
	tests := []polygonIntersectionTest{
		{"Huge polygon contains micro polygon", &hugePolygon, &microPolygon, true},
		{"Micro polygon in huge polygon", &microPolygon, &hugePolygon, true},
		{"Huge polygon contains nano polygon", &hugePolygon, &nanoPolygon, true},
		{"Nano polygon in huge polygon", &nanoPolygon, &hugePolygon, true},

		// Edge case: Micro polygons with each other
		{"Micro polygons intersect", &microPolygon, &nanoPolygon, true}, // Both at same location
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.poly1.IntersectsPolygon(tt.poly2)
			if result != tt.expected {
				t.Errorf("IntersectsPolygon() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
