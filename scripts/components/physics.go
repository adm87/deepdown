package components

import (
	"github.com/adm87/deepdown/scripts/ecs/entity"
	"github.com/adm87/utilities/sparse"
)

var physicsStorage = sparse.NewSet[Physics, entity.Entity](512)

func GetPhysics(e entity.Entity) *Physics {
	return physicsStorage.Get(e)
}

func HasPhysics(e entity.Entity) bool {
	return physicsStorage.Has(e)
}

func AddPhysics(e entity.Entity, p Physics) {
	physicsStorage.Insert(e, p)
}

func RemovePhysics(e entity.Entity) {
	physicsStorage.Remove(e)
}

func GetOrAddPhysics(e entity.Entity) *Physics {
	if !physicsStorage.Has(e) {
		physicsStorage.Insert(e, Physics{})
	}
	return physicsStorage.UnsafeGet(e)
}

func EachPhysics(f func(e entity.Entity, p *Physics)) {
	physicsStorage.Each(f)
}

// =========== Physics Component ===========

type Physics struct {
	Awake    bool
	OnGround bool
	Gravity  float32
	Velocity [2]float32
}
