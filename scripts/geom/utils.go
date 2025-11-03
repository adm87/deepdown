package geom

import "math"

func ComputeNormal(p1, p2 [2]float32) [2]float32 {
	edgeX := float64(p2[0] - p1[0])
	edgeY := float64(p2[1] - p1[1])

	nx := -edgeY
	ny := edgeX

	return [2]float32{float32(nx), float32(ny)}
}

func Distance(p1, p2 [2]float32) float32 {
	dx := p2[0] - p1[0]
	dy := p2[1] - p1[1]
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}
