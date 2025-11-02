package physics

var collisionMatrix [32][32]bool

func init() {
	for i := range 32 {
		collisionMatrix[CollisionLayerDefault][i] = true
		collisionMatrix[i][CollisionLayerDefault] = true
	}
}

// EnableCollision enables collision detection between the two specified layers.
func EnableCollision(layerA, layerB Layer) {
	collisionMatrix[layerA][layerB] = true
	collisionMatrix[layerB][layerA] = true
}

// DisableCollision disables collision detection between the two specified layers.
func DisableCollision(layerA, layerB Layer) {
	collisionMatrix[layerA][layerB] = false
	collisionMatrix[layerB][layerA] = false
}

// ShouldCollide returns true if collision detection is enabled between the two specified layers.
func ShouldCollide(layerA, layerB Layer) bool {
	return collisionMatrix[layerA][layerB]
}

// Layer represents a collision layer.
// Layers are used to categorize colliders and manage their interactions.
//
// Layers are defined at runtime using NewLayer.
// The first layer created is assigned the value 0, the second layer 1, and so on.
// A maximum of 255 layers can be created.
//
// Example usage:
//
//	var (
//	    LayerPlayer = collision.NewLayer("Player")
//	    LayerEnemy  = collision.NewLayer("Enemy")
//	    LayerWall   = collision.NewLayer("Wall")
//	)
type Layer uint8

const (
	MaxCollisionLayers = 32

	CollisionLayerDefault Layer = iota // Default collision layer
)

var nameByLayer = map[Layer]string{
	CollisionLayerDefault: "Default",
}

// NewLayer creates a new collision layer with the given name.
// It panics if the maximum number of layers MaxCollisionLayers is exceeded.
func NewLayer(name string) Layer {
	if len(nameByLayer) >= MaxCollisionLayers {
		panic("maximum number of collision layers exceeded")
	}

	layer := Layer(len(nameByLayer))
	nameByLayer[layer] = name

	return layer
}

func (l Layer) String() string {
	if name, ok := nameByLayer[l]; ok {
		return name
	}
	return "unknown"
}

func (l Layer) IsValid() bool {
	_, ok := nameByLayer[l]
	return ok
}

func NameByLayer(layer Layer) (string, bool) {
	name, ok := nameByLayer[layer]
	return name, ok
}
