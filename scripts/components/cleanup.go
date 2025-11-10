package components

import "github.com/adm87/deepdown/scripts/ecs/entity"

func DestroyEntity(e entity.Entity) {
	RemoveBounds(e)
	RemoveCollision(e)
	RemoveRectangleBody(e)
	RemoveTransform(e)
	RemoveTriangleBody(e)
}
