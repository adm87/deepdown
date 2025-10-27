package level

import (
	"log/slog"
	"sync"

	"github.com/adm87/deepdown/scripts/assets"
	"github.com/adm87/deepdown/scripts/collision"
	"github.com/adm87/tiled"
	"github.com/adm87/utilities/hash"
)

func BuildLevel(logger *slog.Logger, world *collision.World, tmx *tiled.Tmx) (collision.Collider, error) {
	player, err := BuildPlayer(logger, world, tiled.ObjectGroupByName(tmx, "Player"), tmx)
	if err != nil {
		return nil, err
	}
	player.Layer = collision.NewLayer("Player")

	errch := make(chan error, 1)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := BuildStaticCollision(logger, world, tiled.ObjectGroupByName(tmx, "Static")); err != nil {
			errch <- err
		}

	}()
	wg.Wait()

	close(errch)

	if err, ok := <-errch; ok {
		return nil, err
	}

	return player, nil
}

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

func BuildPlayer(logger *slog.Logger, world *collision.World, spawnGroup *tiled.ObjectGroup, tmx *tiled.Tmx) (*collision.BoxCollider, error) {
	if spawnGroup == nil || len(spawnGroup.Objects) == 0 {
		logger.Warn("No player spawn object found")
		return nil, nil
	}

	for i := range spawnGroup.Objects {
		object := &spawnGroup.Objects[i]
		box := collision.NewBoxCollider(object.X, object.Y, object.Width, object.Height)
		box.SetType(collision.Dynamic)

		tileID, _ := tiled.DecodeGID(object.GID)
		tileset, _, _ := tiled.TilesetByGID(tmx, tileID)

		tsx := assets.MustGet[*tiled.Tsx](assets.AssetHandle(tileset.Source))
		x, y := tiled.ObjectAlignmentAnchor(tsx.ObjectAlignment)

		box.Offset[0] = -x * object.Width
		box.Offset[1] = -y * object.Height

		world.AddCollider(box, hash.NoGridPadding)
		return box, nil
	}

	return nil, nil
}
