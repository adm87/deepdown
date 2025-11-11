package level

import (
	"log/slog"
	"strconv"

	"github.com/adm87/deepdown/scripts/components"
	"github.com/adm87/deepdown/scripts/ecs/entity"
	"github.com/adm87/tiled"
)

func (l *Level) BuildStaticCollision(collisionGroup *tiled.ObjectGroup) error {
	if collisionGroup == nil || len(collisionGroup.Objects) == 0 {
		l.ctx.Logger().Warn("No collision objects found")
		return nil
	}

	for i := range collisionGroup.Objects {
		obj := &collisionGroup.Objects[i]

		var role components.CollisionRole
		if prop := tiled.PropertyByType(obj.Properties, "CollisionRole"); prop != nil {
			bit, err := strconv.Atoi(prop.Value)
			if err != nil {
				return err
			}
			role = components.CollisionRole(bit >> 1)
		}

		var e entity.Entity

		switch role {
		case components.CollisionRoleWall:
			e = NewWall(l.ecs, obj.X, obj.Y, obj.Width, obj.Height)
		case components.CollisionRoleFloor:
			if len(obj.Polygon.Points) > 0 {
				e = NewSlopedFloor(l.ecs, obj.X, obj.Y, [6]float32(obj.Polygon.Points))
			} else {
				e = NewFlatFloor(l.ecs, obj.X, obj.Y, obj.Width, obj.Height)
			}
		case components.CollisionRolePlatform:
			e = NewPlatform(l.ecs, obj.X, obj.Y, obj.Width, obj.Height)
		default:
			l.ctx.Logger().Warn("Unknown collision role for object", slog.String("name", obj.Name))
			continue
		}

		l.physics.Add(e)
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

		l.player = NewPlayer(l.ecs, obj.X, obj.Y, obj.Width, obj.Height)

		// data, ok := tilemap.GetTileData(obj.GID, tmx, obj.X, obj.Y)
		// if !ok {
		// 	l.ctx.Logger().Warn("No tile data found for player spawn")
		// }

		l.physics.Add(l.player)

		l.ctx.Logger().Info("Player spawn created at ", obj.X, ", ", obj.Y)
	}

	return nil
}
