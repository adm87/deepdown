package collision

import (
	"github.com/adm87/deepdown/scripts/deepdown"
	"github.com/adm87/utilities/hash"
)

const (
	GridCellSize float32 = 8.0
)

var collisionMatrix [256][256]bool

func init() {
	for i := range 256 {
		collisionMatrix[DefaultLayer][i] = true
		collisionMatrix[i][DefaultLayer] = true
	}
}

func EnableCollision(layerA, layerB Layer) {
	collisionMatrix[layerA][layerB] = true
	collisionMatrix[layerB][layerA] = true
}

func DisableCollision(layerA, layerB Layer) {
	collisionMatrix[layerA][layerB] = false
	collisionMatrix[layerB][layerA] = false
}

func ShouldCollide(layerA, layerB Layer) bool {
	return collisionMatrix[layerA][layerB]
}

type CollisionPair uint64

func EncodePair(id1, id2 uint32) CollisionPair {
	if id1 > id2 {
		id1, id2 = id2, id1
	}
	return CollisionPair(uint64(id1)<<32 | uint64(id2))
}

func DecodePair(pair CollisionPair) (id1, id2 uint32) {
	id1 = uint32(pair >> 32)
	id2 = uint32(pair & 0xFFFFFFFF)
	return id1, id2
}

type World struct {
	ctx deepdown.Context

	grid *hash.Grid[Collider]

	staticColliders  []Collider
	dynamicColliders []Collider

	activePairs hash.Set[CollisionPair]
}

func NewWorld(ctx deepdown.Context) *World {
	return &World{
		grid:        hash.NewGrid[Collider](GridCellSize, GridCellSize),
		activePairs: hash.NewSet[CollisionPair](),
		ctx:         ctx,
	}
}

func (w *World) AddCollider(c Collider, padding hash.GridItemPadding) {
	w.insert(c, padding)

	switch c.Info().Type {
	case Static:
		w.staticColliders = append(w.staticColliders, c)
	case Dynamic:
		w.dynamicColliders = append(w.dynamicColliders, c)
	}
}

func (w *World) RemoveCollider(c Collider) {
	w.grid.Remove(c)

	switch c.Info().Type {
	case Static:
		w.staticColliders = w.removeFrom(w.staticColliders, c)
	case Dynamic:
		w.dynamicColliders = w.removeFrom(w.dynamicColliders, c)
	}
}

func (w *World) GetCells() []uint64 {
	return w.grid.Cells()
}

func (w *World) GetCellSize() (cellWidth, cellHeight float32) {
	return w.grid.CellSize()
}

func (w *World) QueryCells(minX, minY, maxX, maxY float32) []uint64 {
	return w.grid.QueryCells(minX, minY, maxX, maxY)
}

func (w *World) Query(minX, minY, maxX, maxY float32) []Collider {
	return w.grid.Query(minX, minY, maxX, maxY)
}

func (w *World) UpdateCollider(c Collider, padding hash.GridItemPadding) {
	w.grid.Remove(c)
	w.insert(c, padding)
}

func (w *World) CheckCollisions() {
	newPairs := hash.NewSet[CollisionPair]()

	for i := range w.dynamicColliders {
		dynamic := w.dynamicColliders[i]
		dynamicID := dynamic.Info().ID()
		dynamicLayer := dynamic.Info().Layer

		if dynamic.Info().Type == Ignore {
			continue
		}

		minX, minY, maxX, maxY := dynamic.Bounds()
		others := w.grid.Query(minX, minY, maxX, maxY)

		for j := range others {
			other := others[j]
			otherID := other.Info().ID()
			otherLayer := other.Info().Layer

			if !ShouldCollide(dynamicLayer, otherLayer) {
				continue
			}
			if dynamicID == otherID || other.Info().Type == Ignore {
				continue
			}
		}
	}

	w.activePairs = newPairs
}

func (w *World) insert(c Collider, padding hash.GridItemPadding) {
	minX, minY, maxX, maxY := c.Bounds()
	if polygon, ok := c.(*PolygonCollider); ok {
		w.grid.InsertFunc(c, minX, minY, maxX, maxY, padding, w.insertPolygon(polygon))
	} else {
		w.grid.Insert(c, minX, minY, maxX, maxY, padding)
	}
}

func (w *World) insertPolygon(polygon *PolygonCollider) hash.GridInsertionFunc[Collider] {
	return func(cellMinX, cellMinY, cellMaxX, cellMaxY float32) bool {
		if polygon.IntersectsAABB(cellMinX, cellMinY, cellMaxX, cellMaxY) {
			return true
		}
		return false
	}
}

func (w *World) removeFrom(collection []Collider, c Collider) []Collider {
	for i, col := range collection {
		if col == c {
			return append(collection[:i], collection[i+1:]...)
		}
	}
	return collection
}
