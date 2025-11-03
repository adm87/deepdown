package geom

import "math"

func ComputeNormal(p1, p2 [2]float32) [2]float32 {
	edgeX := float64(p2[0] - p1[0])
	edgeY := float64(p2[1] - p1[1])

	nx := -edgeY
	ny := edgeX

	return [2]float32{float32(nx), float32(ny)}
}

func ComputeSlopeNormal(p1, p2 [2]float32) [2]float32 {
	edgeX := float64(p2[0] - p1[0])
	edgeY := float64(p2[1] - p1[1])

	nx := -edgeY
	ny := edgeX

	length := math.Sqrt(nx*nx + ny*ny)
	if length == 0 {
		return [2]float32{0, 0}
	}

	return [2]float32{-float32(nx / length), -float32(ny / length)}
}

func FindSlope(points [6]float32) ([2][2]float32, [2]float32, SlopeType) {
	var slope [2][2]float32
	longEdgeIdx := 0
	longestDist := float32(0)

	for i := range 3 {
		j := (i + 1) % 3
		p1 := [2]float32{points[i*2], points[i*2+1]}
		p2 := [2]float32{points[j*2], points[j*2+1]}
		dist := Distance(p1, p2)
		if dist > longestDist {
			longestDist = dist
			longEdgeIdx = i
		}
	}

	i := longEdgeIdx
	j := (i + 1) % 3
	k := (i + 2) % 3

	p1 := [2]float32{points[i*2], points[i*2+1]}
	p2 := [2]float32{points[j*2], points[j*2+1]}
	corner := [2]float32{points[k*2], points[k*2+1]}

	slope[0], slope[1] = p1, p2

	var slopeType SlopeType
	switch {
	case p1[1] == p2[1]:
		slopeType = SlopeFlat
	case (p1[0] < p2[0] && p1[1] < p2[1]) || (p1[0] > p2[0] && p1[1] > p2[1]):
		slopeType = SlopeAcending
	default:
		slopeType = SlopeDecending
	}

	return slope, corner, slopeType
}

func Distance(p1, p2 [2]float32) float32 {
	dx := p2[0] - p1[0]
	dy := p2[1] - p1[1]
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}
