package components

import (
	"github.com/adm87/deepdown/scripts/ecs/entity"
	"github.com/adm87/utilities/sparse"
)

// =========== Collision Component Storage ===========

var collisionMatrix [32][32]bool

func init() {
	for i := range 32 {
		collisionMatrix[DefaultCollisionLayer][i] = true
		collisionMatrix[i][DefaultCollisionLayer] = true
	}
}

var collisionStorage = sparse.NewSet[Collision, entity.Entity](512)

func GetCollision(e entity.Entity) *Collision {
	return collisionStorage.Get(e)
}

func HasCollision(e entity.Entity) bool {
	return collisionStorage.Has(e)
}

func AddCollision(e entity.Entity, c Collision) {
	collisionStorage.Insert(e, c)
}

func RemoveCollision(e entity.Entity) {
	collisionStorage.Remove(e)
}

func GetOrAddCollision(e entity.Entity) *Collision {
	if !collisionStorage.Has(e) {
		collisionStorage.Insert(e, Collision{})
	}
	return collisionStorage.UnsafeGet(e)
}

func EachCollision(f func(e entity.Entity, c *Collision)) {
	collisionStorage.Each(f)
}

// =========== Collision Component ==========

type Collision struct {
	Layer CollisionLayer
	Mode  CollisionMode
	Role  CollisionRole
	Shape CollisionShape
	Type  CollisionType
}

// =========== Collision Layer ==========

type CollisionLayer uint8

const (
	MaxCollisionLayers int = 32

	NoCollisionLayer      CollisionLayer = 0
	DefaultCollisionLayer CollisionLayer = iota
)

var nameByLayer = map[CollisionLayer]string{
	DefaultCollisionLayer: "Default",
}

func NewLayer(name string) CollisionLayer {
	if len(nameByLayer) >= MaxCollisionLayers {
		panic("maximum number of collision layers exceeded")
	}

	layer := CollisionLayer(len(nameByLayer))
	nameByLayer[layer] = name

	return layer
}

func (l CollisionLayer) String() string {
	if name, ok := nameByLayer[l]; ok {
		return name
	}
	return "unknown"
}

func (l CollisionLayer) IsValid() bool {
	_, ok := nameByLayer[l]
	return ok
}

func NameByLayer(layer CollisionLayer) (string, bool) {
	name, ok := nameByLayer[layer]
	return name, ok
}

// EnableCollision enables collision detection between the two specified layers.
func EnableCollision(layerA, layerB CollisionLayer) {
	collisionMatrix[layerA][layerB] = true
	collisionMatrix[layerB][layerA] = true
}

// DisableCollision disables collision detection between the two specified layers.
func DisableCollision(layerA, layerB CollisionLayer) {
	collisionMatrix[layerA][layerB] = false
	collisionMatrix[layerB][layerA] = false
}

// ShouldCollide returns true if collision detection is enabled between the two specified layers.
func ShouldCollide(layerA, layerB CollisionLayer) bool {
	return collisionMatrix[layerA][layerB]
}

// =========== Collision Type ==========

type CollisionType uint8

const (
	CollisionTypeIgnore CollisionType = 0
	CollisionTypeStatic CollisionType = iota
	CollisionTypeDynamic
)

// =========== Collision Role ==========

type CollisionRole uint8

const (
	CollisionRoleNone     CollisionRole = 0
	CollisionRoleWall     CollisionRole = 1 << 0
	CollisionRoleFloor    CollisionRole = 1 << 1
	CollisionRolePlatform CollisionRole = CollisionRoleFloor | CollisionRoleWall
)

// =========== Collision Mode ==========

type CollisionMode uint8

const (
	CollisionModeDiscrete CollisionMode = iota
	CollisionModeContinuous
)

// =========== Collision Shape ==========

type CollisionShape uint8

const (
	CollisionShapeBox CollisionShape = iota
	CollisionShapeTriangle
)
