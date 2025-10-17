package level

import (
	"fmt"
	"log/slog"

	"github.com/adm87/deepdown/scripts/collision"
	"github.com/adm87/tiled"
	"github.com/adm87/utilities/hashgrid"
)

func BuildCollision(logger *slog.Logger, world *collision.World, collisionGroup *tiled.ObjectGroup) error {
	if collisionGroup == nil || len(collisionGroup.Objects) == 0 {
		logger.Warn("No collision objects found")
		return nil
	}

	colliders := make([]collision.Collider, 0, len(collisionGroup.Objects))
	for i := range collisionGroup.Objects {
		var collider collision.Collider

		object := &collisionGroup.Objects[i]
		polygon := &object.Polygon

		switch {
		case len(collisionGroup.Objects[i].Polygon.Points) > 0:
			collider = collision.NewPolygonCollider(object.X, object.Y, polygon.Points)
		default:
			collider = collision.NewBoxCollider(object.X, object.Y, object.Width, object.Height)
		}

		colliders = append(colliders, collider)
		world.AddCollider(collider)
	}

	keys := make(map[collision.Collider][]hashgrid.GridKey)
	for _, collider := range colliders {
		keys[collider] = world.Grid.GetKeys(collider)
	}

	for key, value := range keys {
		minX, minY := key.Min()
		maxX, maxY := key.Max()
		nMinX, nMinY := int32(minX/64.0), int32(minY/64.0)
		nMaxX, nMaxY := int32(maxX/64.0), int32(maxY/64.0)

		println(fmt.Sprintf("Collider AABB: min(%.2f, %.2f), max(%.2f, %.2f)", minX, minY, maxX, maxY))
		println(fmt.Sprintf("Normalized: min(%d, %d), max(%d, %d)", nMinX, nMinY, nMaxX, nMaxY))
		println(fmt.Sprintf("Grid keys: %v", value))
		for _, k := range value {
			x, y := hashgrid.DecodeGridKey(k)
			println(fmt.Sprintf("  Key %d = 0x%X → Grid (%d,%d)", k, k, x, y))
		}
	}

	return nil
}
