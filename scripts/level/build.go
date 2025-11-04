package level

import (
	"strconv"

	"github.com/adm87/deepdown/scripts/physics"
	"github.com/adm87/tiled"
	"github.com/adm87/tiled/tilemap"
)

func (l *Level) BuildStaticCollision(collisionGroup *tiled.ObjectGroup) error {
	if collisionGroup == nil || len(collisionGroup.Objects) == 0 {
		l.ctx.Logger().Warn("No collision objects found")
		return nil
	}

	for i := range collisionGroup.Objects {
		obj := &collisionGroup.Objects[i]

		var collider physics.Collider

		if len(obj.Polygon.Points) > 0 {
			collider = physics.GetTriangleCollider(obj.X, obj.Y, [6]float32(obj.Polygon.Points))
		} else {
			collider = physics.GetBoxCollider(obj.X, obj.Y, obj.Width, obj.Height)
		}

		var role physics.Role
		if prop := tiled.PropertyByType(obj.Properties, "CollisionRole"); prop != nil {
			bit, err := strconv.Atoi(prop.Value)
			if err != nil {
				return err
			}
			role = physics.Role(bit >> 1)
		}

		collider.Info().Role = role
		collider.Info().State = physics.ColliderStateStatic

		l.world.AddCollider(collider)
	}

	return nil
}

func (l *Level) BuildPlayer(spawnGroup *tiled.ObjectGroup, tmx *tiled.Tmx) error {
	if spawnGroup == nil || len(spawnGroup.Objects) == 0 {
		l.ctx.Logger().Warn("No player spawn object found")
		return nil
	}

	for i := range spawnGroup.Objects {
		obj := &spawnGroup.Objects[i]

		l.player = &Player{}
		l.player.X = obj.X
		l.player.Y = obj.Y
		l.player.Width = obj.Width * 0.5
		l.player.Height = obj.Height

		l.player.BoxCollider = *physics.GetBoxCollider(obj.X, obj.Y, l.player.Width, l.player.Height)
		l.player.BoxCollider.Info().State = physics.ColliderStateDynamic
		l.player.Offset[0] = (obj.Width - l.player.Width) * 0.5

		data, ok := tilemap.GetTileData(obj.GID, tmx, obj.X, obj.Y)
		if !ok {
			l.ctx.Logger().Warn("No tile data found for player spawn")
		}
		l.player.Data = data

		l.world.AddCollider(&l.player.BoxCollider)

		l.ctx.Logger().Info("Player spawn created at ", obj.X, ", ", obj.Y)
	}

	return nil
}
