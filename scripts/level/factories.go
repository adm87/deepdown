package level

import (
	"github.com/adm87/deepdown/scripts/components"
	"github.com/adm87/deepdown/scripts/ecs"
	"github.com/adm87/deepdown/scripts/ecs/entity"
	"github.com/adm87/deepdown/scripts/geom"
)

func NewPhysicalEntity(world *ecs.World, x, y, width, height float32) entity.Entity {
	entity := world.NewEntity()

	bounds := components.GetOrAddBounds(entity)
	bounds.SetSize(width, height)

	transform := components.GetOrAddTransform(entity)
	transform.SetPosition(x, y)

	collision := components.GetOrAddCollision(entity)
	collision.Layer = components.DefaultCollisionLayer
	collision.Shape = components.CollisionShapeBox

	return entity
}

// ========== Player Entity ===========

func NewPlayer(world *ecs.World, x, y, width, height float32) entity.Entity {
	player := NewPhysicalEntity(world, x, y, width, height)

	collision := components.GetCollision(player)
	collision.Type = components.CollisionTypeDynamic

	body := components.GetOrAddRectangleBody(player)
	body.Rectangle = geom.NewRectangle(0, 0, width, height)

	return player
}

// ========== Wall Entity ===========

func NewWall(world *ecs.World, x, y, width, height float32) entity.Entity {
	wall := NewPhysicalEntity(world, x, y, width, height)

	collision := components.GetOrAddCollision(wall)
	collision.Type = components.CollisionTypeStatic
	collision.Role = components.CollisionRoleWall

	body := components.GetOrAddRectangleBody(wall)
	body.Rectangle = geom.NewRectangle(0, 0, width, height)

	return wall
}

// ========== Flat Floor Entity ===========

func NewFlatFloor(world *ecs.World, x, y, width, height float32) entity.Entity {
	floor := NewPhysicalEntity(world, x, y, width, height)

	collision := components.GetOrAddCollision(floor)
	collision.Type = components.CollisionTypeStatic

	body := components.GetOrAddRectangleBody(floor)
	body.Rectangle = geom.NewRectangle(0, 0, width, height)
	collision.Role = components.CollisionRoleFloor

	return floor
}

// ========== Sloped Floor Entity ===========

func NewSlopedFloor(world *ecs.World, x, y float32, points [6]float32) entity.Entity {
	minX, minY, maxX, maxY := geom.ComputeAABB(points)

	slope := NewPhysicalEntity(world, x, y, maxX-minX, maxY-minY)

	bounds := components.GetOrAddBounds(slope)
	bounds.SetOffset(minX, minY)

	collision := components.GetOrAddCollision(slope)
	collision.Shape = components.CollisionShapeTriangle
	collision.Type = components.CollisionTypeStatic
	collision.Role = components.CollisionRoleFloor

	body := components.GetOrAddTriangleBody(slope)
	body.Triangle = geom.NewTriangle(0, 0, points)

	return slope
}

// ========== Platform Entity ===========

func NewPlatform(world *ecs.World, x, y, width, height float32) entity.Entity {
	platform := NewPhysicalEntity(world, x, y, width, height)

	collision := components.GetOrAddCollision(platform)
	collision.Type = components.CollisionTypeStatic
	collision.Role = components.CollisionRolePlatform

	return platform
}
