package collision

import "github.com/adm87/utilities/hashgrid"

const (
	GridCellSize float32 = 64.0
)

type World struct {
	grid *hashgrid.Grid[Collider]

	profiles Profiles
}

func NewWorld() *World {
	return &World{
		grid:     hashgrid.New[Collider](GridCellSize),
		profiles: make(Profiles),
	}
}

func (w *World) AddProfile(layer Layer, interactions Interactions) {
	w.profiles[layer] = interactions
}

func (w *World) AddCollider(c Collider) {

}

func (w *World) RemoveCollider(c Collider) {

}
