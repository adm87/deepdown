package collision

import (
	"github.com/adm87/deepdown/scripts/deepdown"
	"github.com/adm87/utilities/hash"
)

const (
	GridCellSize float32 = 8.0
)

var collisionMatrix [32][32]bool

func init() {
	for i := range 32 {
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

// CollisionHandler is a function signature for when two colliders enter/exit a potential phase
type CollisionHandler func(colliderA, colliderB Collider)

var DefaultCollisionHandler = func(colliderA, colliderB Collider) {}

type World struct {
	ctx deepdown.Context

	grid *hash.Grid[Collider] // Spatial indexing grid

	dynamicColliders []Collider

	activePairs map[CollisionPair][2]Collider
	pairs       map[CollisionPair][2]Collider

	OnEnter CollisionHandler
	OnStay  CollisionHandler
	OnExit  CollisionHandler
}

func NewWorld(ctx deepdown.Context) *World {
	return &World{
		OnEnter:     DefaultCollisionHandler,
		OnStay:      DefaultCollisionHandler,
		OnExit:      DefaultCollisionHandler,
		grid:        hash.NewGrid[Collider](GridCellSize, GridCellSize),
		activePairs: make(map[CollisionPair][2]Collider, 10),
		pairs:       make(map[CollisionPair][2]Collider, 10),
		ctx:         ctx,
	}
}

func (w *World) AddCollider(c Collider, padding hash.GridItemPadding) {
	w.insert(c, padding)

	if c.Info().Type == Dynamic {
		w.dynamicColliders = append(w.dynamicColliders, c)
	}
}

func (w *World) RemoveCollider(c Collider) {
	w.grid.Remove(c)

	if c.Info().Type == Dynamic {
		w.dynamicColliders = w.removeFrom(w.dynamicColliders, c)
	}
}

func (w *World) OnCollisionEnter(handler CollisionHandler) {
	w.OnEnter = handler
}

func (w *World) OnCollisionExit(handler CollisionHandler) {
	w.OnExit = handler
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

func (w *World) UpdateColliders(dt float64) {
	for i := range w.dynamicColliders {
		if w.dynamicColliders[i].applyVelocity(dt) {
			w.grid.Remove(w.dynamicColliders[i])
			w.insert(w.dynamicColliders[i], hash.NoGridPadding)
		}
	}
}

func (w *World) CheckCollisions() {
	w.activePairs, w.pairs = w.pairs, w.activePairs
	clear(w.pairs)

	w.collectPairs()
	if len(w.pairs) > 0 {
		for pairKey, colliders := range w.pairs {
			if _, wasActive := w.activePairs[pairKey]; !wasActive {
				w.OnEnter(colliders[0], colliders[1])
				continue
			}
			w.OnStay(colliders[0], colliders[1])
		}
	}

	for pairKey, colliders := range w.activePairs {
		if _, stillActive := w.pairs[pairKey]; !stillActive {
			w.OnExit(colliders[0], colliders[1])
		}
	}
}

// Broadphase collision detection: collect potential collision pairs
func (w *World) collectPairs() {
	for i := range w.dynamicColliders {
		infoA := w.dynamicColliders[i].Info()
		if infoA.Type == Ignore {
			continue
		}

		colliderAID := infoA.id
		colliderALayer := infoA.Layer

		minX, minY, maxX, maxY := w.dynamicColliders[i].Bounds()
		others := w.grid.Query(minX, minY, maxX, maxY)

		for j := range others {
			infoB := others[j].Info()

			colliderBID := infoB.id
			colliderBLayer := infoB.Layer

			if colliderBID == colliderAID || infoB.Type == Ignore || !ShouldCollide(colliderALayer, colliderBLayer) {
				continue // Prevent self-collision and ignored types
			}

			pairKey := EncodePair(colliderAID, colliderBID)
			if _, exists := w.pairs[pairKey]; exists {
				continue // Prevent duplicate pairs
			}

			w.pairs[pairKey] = [2]Collider{w.dynamicColliders[i], others[j]}
		}
	}
}

func (w *World) insert(c Collider, padding hash.GridItemPadding) {
	minX, minY, maxX, maxY := c.Bounds()
	switch c.Info().Shape {
	case Polygon:
		w.grid.InsertFunc(c, minX, minY, maxX, maxY, padding, w.insertPolygon(c.(*PolygonCollider)))
	case Triangle:
		w.grid.InsertFunc(c, minX, minY, maxX, maxY, padding, w.insertTriangle(c.(*TriangleCollider)))
	default:
		w.grid.Insert(c, minX, minY, maxX, maxY, padding)
	}
}

func (w *World) insertPolygon(polygon *PolygonCollider) hash.GridInsertionFunc[Collider] {
	return func(cellMinX, cellMinY, cellMaxX, cellMaxY float32) bool {
		return polygon.IntersectsAABB(cellMinX, cellMinY, cellMaxX, cellMaxY)
	}
}

func (w *World) insertTriangle(triangle *TriangleCollider) hash.GridInsertionFunc[Collider] {
	return func(cellMinX, cellMinY, cellMaxX, cellMaxY float32) bool {
		return triangle.IntersectsAABB(cellMinX, cellMinY, cellMaxX, cellMaxY)
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
