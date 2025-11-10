package components

import (
	"github.com/adm87/deepdown/scripts/ecs/entity"
	"github.com/adm87/deepdown/scripts/geom"
	"github.com/adm87/utilities/collection"
)

// =========== Body Component Storage ===========

var (
	rectangleBodyStorage = collection.NewSparseSet[RectangleBody, entity.Entity](512)
	triangleBodyStorage  = collection.NewSparseSet[TriangleBody, entity.Entity](512)
)

func GetRectangleBody(e entity.Entity) *RectangleBody {
	return rectangleBodyStorage.Get(e)
}

func HasRectangleBody(e entity.Entity) bool {
	return rectangleBodyStorage.Has(e)
}

func AddRectangleBody(e entity.Entity, b RectangleBody) {
	rectangleBodyStorage.Insert(e, b)
}

func RemoveRectangleBody(e entity.Entity) {
	rectangleBodyStorage.Remove(e)
}

func GetOrAddRectangleBody(e entity.Entity) *RectangleBody {
	if !rectangleBodyStorage.Has(e) {
		rectangleBodyStorage.Insert(e, RectangleBody{})
	}
	return rectangleBodyStorage.UnsafeGet(e)
}

func GetTriangleBody(e entity.Entity) *TriangleBody {
	return triangleBodyStorage.Get(e)
}

func HasTriangleBody(e entity.Entity) bool {
	return triangleBodyStorage.Has(e)
}

func AddTriangleBody(e entity.Entity, b TriangleBody) {
	triangleBodyStorage.Insert(e, b)
}

func RemoveTriangleBody(e entity.Entity) {
	triangleBodyStorage.Remove(e)
}

func GetOrAddTriangleBody(e entity.Entity) *TriangleBody {
	if !triangleBodyStorage.Has(e) {
		triangleBodyStorage.Insert(e, TriangleBody{})
	}
	return triangleBodyStorage.UnsafeGet(e)
}

// =========== Body Component ===========

type RectangleBody struct {
	Rectangle geom.Rectangle
}

type TriangleBody struct {
	Triangle geom.Triangle
}
