package level

import (
	"errors"
	"log/slog"

	"github.com/adm87/deepdown/scripts/assets"
	"github.com/adm87/deepdown/scripts/collision"
	"github.com/adm87/tiled"
	"github.com/adm87/tiled/tilemap"
	"github.com/adm87/utilities/hash"
)

func BuildStaticCollision(logger *slog.Logger, world *collision.World, collisionGroup *tiled.ObjectGroup) error {
	if collisionGroup == nil || len(collisionGroup.Objects) == 0 {
		logger.Warn("No collision objects found")
		return nil
	}

	for i := range collisionGroup.Objects {
		var collider collision.Collider

		object := &collisionGroup.Objects[i]
		polygon := &object.Polygon

		switch {
		case len(collisionGroup.Objects[i].Polygon.Points) > 0:
			if len(collisionGroup.Objects[i].Polygon.Points) == 6 {
				collider = collision.NewTriangleCollider(object.X, object.Y, polygon.Points)
				collider.(*collision.TriangleCollider).SetType(collision.Static)
			} else {
				collider = collision.NewPolygonCollider(object.X, object.Y, polygon.Points)
				collider.(*collision.PolygonCollider).SetType(collision.Static)
			}
		default:
			collider = collision.NewBoxCollider(object.X, object.Y, object.Width, object.Height)
			collider.(*collision.BoxCollider).SetType(collision.Static)
		}

		world.AddCollider(collider, hash.GridCellPadding)
	}

	return nil
}

func BuildPlayer(logger *slog.Logger, world *collision.World, spawnGroup *tiled.ObjectGroup, tmx *tiled.Tmx) (*Player, error) {
	if spawnGroup == nil || len(spawnGroup.Objects) == 0 {
		logger.Warn("No player spawn object found")
		return nil, nil
	}

	for i := range spawnGroup.Objects {
		object := &spawnGroup.Objects[i]

		data, ok := tilemap.GetTileData(object.GID, tmx, 0, 0)
		if !ok {
			return nil, errors.New("Failed to get player tiled data")
		}

		player := &Player{
			BoxCollider: *collision.NewBoxCollider(object.X, object.Y, object.Width, object.Height),
			data:        data,
		}
		player.SetType(collision.Dynamic)

		tileID, _ := tiled.DecodeGID(object.GID)
		tileset, _, _ := tiled.TilesetByGID(tmx, tileID)

		tsx := assets.MustGet[*tiled.Tsx](assets.AssetHandle(tileset.Source))
		x, y := tiled.ObjectAlignmentAnchor(tsx.ObjectAlignment)

		player.Offset[0] = -x * object.Width
		player.Offset[1] = -y * object.Height

		return player, nil
	}

	return nil, nil
}
