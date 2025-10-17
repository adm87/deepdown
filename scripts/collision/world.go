package collision

import (
	"github.com/adm87/deepdown/scripts/deepdown"
	"github.com/adm87/utilities/hashgrid"
)

const (
	GridCellSize float32 = 4.0
)

type World struct {
	ctx      deepdown.Context
	profiles Profiles
	Grid     *hashgrid.Grid[Collider]
}

func NewWorld(ctx deepdown.Context) *World {
	return &World{
		ctx:      ctx,
		Grid:     hashgrid.NewWithPadding[Collider](GridCellSize, hashgrid.AllPadding),
		profiles: make(Profiles),
	}
}

func (w *World) AddProfile(layer Layer, interactions Interactions) {
	w.profiles[layer] = interactions
}

func (w *World) AddCollider(c Collider) {
	w.Grid.Insert(c)
}

func (w *World) RemoveCollider(c Collider) {
	w.Grid.Remove(c)
}
