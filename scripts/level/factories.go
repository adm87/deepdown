package level

import (
	"github.com/adm87/deepdown/scripts/components"
	"github.com/adm87/deepdown/scripts/ecs"
	"github.com/adm87/deepdown/scripts/ecs/entity"
	"github.com/adm87/deepdown/scripts/geom"
	"github.com/adm87/deepdown/scripts/physics"
)

// ========== Player Entity ===========

func NewPlayer(world *ecs.World, x, y, width, height float32) entity.Entity {
	player := world.NewEntity()

	transform := components.GetOrAddTransform(player)
	transform.SetPosition(x, y)

	bounds := components.GetOrAddBounds(player)
	bounds.SetOffset(0, 0)
	bounds.SetWidth(width)
	bounds.SetHeight(height)

	collision := components.GetOrAddCollision(player)
	collision.Layer = components.DefaultCollisionLayer
	collision.Type = components.CollisionTypeDynamic

	phys := components.GetOrAddPhysics(player)
	phys.Gravity = physics.Gravity

	rect := components.GetOrAddRectangleBody(player)
	rect.Width, rect.Height = width, height

	return player
}

// ========== Wall Entity ===========

func NewWall(world *ecs.World, x, y, width, height float32) entity.Entity {
	wall := world.NewEntity()

	transform := components.GetOrAddTransform(wall)
	transform.SetPosition(x, y)

	bounds := components.GetOrAddBounds(wall)
	bounds.SetOffset(0, 0)
	bounds.SetWidth(width)
	bounds.SetHeight(height)

	collision := components.GetOrAddCollision(wall)
	collision.Layer = components.DefaultCollisionLayer
	collision.Type = components.CollisionTypeStatic
	collision.Role = components.CollisionRoleWall

	rect := components.GetOrAddRectangleBody(wall)
	rect.Width, rect.Height = width, height

	return wall
}

// ========== Flat Floor Entity ===========

func NewFlatFloor(world *ecs.World, x, y, width, height float32) entity.Entity {
	floor := world.NewEntity()

	transform := components.GetOrAddTransform(floor)
	transform.SetPosition(x, y)

	bounds := components.GetOrAddBounds(floor)
	bounds.SetOffset(0, 0)
	bounds.SetWidth(width)
	bounds.SetHeight(height)

	collision := components.GetOrAddCollision(floor)
	collision.Layer = components.DefaultCollisionLayer
	collision.Type = components.CollisionTypeStatic
	collision.Role = components.CollisionRoleFloor

	rect := components.GetOrAddRectangleBody(floor)
	rect.Width, rect.Height = width, height

	return floor
}

// ========== Sloped Floor Entity ===========

func NewSlopedFloor(world *ecs.World, x, y float32, points [6]float32) entity.Entity {
	minX, minY, maxX, maxY := geom.ComputeAABB(points)

	slope := world.NewEntity()

	transform := components.GetOrAddTransform(slope)
	transform.SetPosition(x, y)

	bounds := components.GetOrAddBounds(slope)
	bounds.SetOffset(minX, minY)
	bounds.SetWidth(maxX - minX)
	bounds.SetHeight(maxY - minY)

	collision := components.GetOrAddCollision(slope)
	collision.Layer = components.DefaultCollisionLayer
	collision.Shape = components.CollisionShapeTriangle
	collision.Type = components.CollisionTypeStatic
	collision.Role = components.CollisionRoleFloor

	triangle := components.GetOrAddTriangleBody(slope)
	triangle.SetPoints(points)

	return slope
}

// ========== Platform Entity ===========

func NewPlatform(world *ecs.World, x, y, width, height float32) entity.Entity {
	platform := world.NewEntity()

	transform := components.GetOrAddTransform(platform)
	transform.SetPosition(x, y)

	bounds := components.GetOrAddBounds(platform)
	bounds.SetOffset(0, 0)
	bounds.SetWidth(width)
	bounds.SetHeight(height)

	collision := components.GetOrAddCollision(platform)
	collision.Layer = components.DefaultCollisionLayer
	collision.Type = components.CollisionTypeStatic
	collision.Role = components.CollisionRolePlatform

	rect := components.GetOrAddRectangleBody(platform)
	rect.Width, rect.Height = width, height

	return platform
}
