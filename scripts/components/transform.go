package components

import (
	"github.com/adm87/deepdown/scripts/ecs/entity"
	"github.com/adm87/utilities/collection"
	"github.com/hajimehoshi/ebiten/v2"
)

// =========== Transform Component Storage ===========

var transformStorage = collection.NewSparseSet[Transform, entity.Entity](512)

func GetTransform(e entity.Entity) *Transform {
	return transformStorage.Get(e)
}

func HasTransform(e entity.Entity) bool {
	return transformStorage.Has(e)
}

func AddTransform(e entity.Entity, t Transform) {
	transformStorage.Insert(e, t)
}

func RemoveTransform(e entity.Entity) {
	transformStorage.Remove(e)
}

func GetOrAddTransform(e entity.Entity) *Transform {
	if !transformStorage.Has(e) {
		transformStorage.Insert(e, Transform{})
	}
	return transformStorage.UnsafeGet(e)
}

// =========== Transform Component ===========

type Transform struct {
	xy [2]float32
}

func (t *Transform) Position() (x, y float32) {
	return t.xy[0], t.xy[1]
}

func (t *Transform) SetPosition(x, y float32) {
	t.xy[0], t.xy[1] = x, y
}

func (t *Transform) Matrix() ebiten.GeoM {
	var m ebiten.GeoM
	m.Translate(float64(t.xy[0]), float64(t.xy[1]))
	return m
}
