package components

import (
	"github.com/adm87/deepdown/scripts/ecs/entity"
	"github.com/adm87/deepdown/scripts/geom"
	"github.com/adm87/utilities/sparse"
)

// =========== Body Component Storage ===========

var (
	rectangleBodyStorage = sparse.NewSet[geom.Rectangle, entity.Entity](512)
	triangleBodyStorage  = sparse.NewSet[geom.Triangle, entity.Entity](512)
)

// =========== Rectangle Body Component ===========

func GetRectangleBody(e entity.Entity) *geom.Rectangle {
	return rectangleBodyStorage.Get(e)
}

func HasRectangleBody(e entity.Entity) bool {
	return rectangleBodyStorage.Has(e)
}

func AddRectangleBody(e entity.Entity, b geom.Rectangle) {
	rectangleBodyStorage.Insert(e, b)
}

func RemoveRectangleBody(e entity.Entity) {
	rectangleBodyStorage.Remove(e)
}

func GetOrAddRectangleBody(e entity.Entity) *geom.Rectangle {
	if !rectangleBodyStorage.Has(e) {
		rectangleBodyStorage.Insert(e, geom.Rectangle{})
	}
	return rectangleBodyStorage.UnsafeGet(e)
}

func EachRectangleBody(f func(e entity.Entity, b *geom.Rectangle)) {
	rectangleBodyStorage.Each(f)
}

// =========== Triangle Body Component ===========

func GetTriangleBody(e entity.Entity) *geom.Triangle {
	return triangleBodyStorage.Get(e)
}

func HasTriangleBody(e entity.Entity) bool {
	return triangleBodyStorage.Has(e)
}

func AddTriangleBody(e entity.Entity, b geom.Triangle) {
	triangleBodyStorage.Insert(e, b)
}

func RemoveTriangleBody(e entity.Entity) {
	triangleBodyStorage.Remove(e)
}

func GetOrAddTriangleBody(e entity.Entity) *geom.Triangle {
	if !triangleBodyStorage.Has(e) {
		triangleBodyStorage.Insert(e, geom.Triangle{})
	}
	return triangleBodyStorage.UnsafeGet(e)
}

func EachTriangleBody(f func(e entity.Entity, b *geom.Triangle)) {
	triangleBodyStorage.Each(f)
}
