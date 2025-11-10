package components

import (
	"github.com/adm87/deepdown/scripts/ecs/entity"
	"github.com/adm87/utilities/collection"
)

// =========== Bounds Component Storage ===========

var boundsStorage = collection.NewSparseSet[Bounds, entity.Entity](512)

func GetBounds(e entity.Entity) *Bounds {
	return boundsStorage.Get(e)
}

func HasBounds(e entity.Entity) bool {
	return boundsStorage.Has(e)
}

func AddBounds(e entity.Entity, b Bounds) {
	boundsStorage.Insert(e, b)
}

func RemoveBounds(e entity.Entity) {
	boundsStorage.Remove(e)
}

func GetOrAddBounds(e entity.Entity) *Bounds {
	if !boundsStorage.Has(e) {
		boundsStorage.Insert(e, Bounds{})
	}
	return boundsStorage.UnsafeGet(e)
}

// =========== Bounds Component ===========

type Bounds struct {
	size   [2]float32
	offset [2]float32
}

func (b *Bounds) Width() float32 {
	return b.size[0]
}

func (b *Bounds) SetWidth(width float32) {
	b.size[0] = width
}

func (b *Bounds) Height() float32 {
	return b.size[1]
}

func (b *Bounds) SetHeight(height float32) {
	b.size[1] = height
}

func (b *Bounds) Size() (width, height float32) {
	return b.size[0], b.size[1]
}

func (b *Bounds) SetSize(width, height float32) {
	b.size[0] = width
	b.size[1] = height
}

func (b *Bounds) Offset() (offsetX, offsetY float32) {
	return b.offset[0], b.offset[1]
}

func (b *Bounds) SetOffset(offsetX, offsetY float32) {
	b.offset[0] = offsetX
	b.offset[1] = offsetY
}

func (b *Bounds) Min() (minX, minY float32) {
	return b.offset[0], b.offset[1]
}

func (b *Bounds) Max() (maxX, maxY float32) {
	return b.offset[0] + b.size[0], b.offset[1] + b.size[1]
}
