package collision

// =========== Layer ==========

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
	DefaultLayer Layer = iota // Default collision layer
)

var nameByLayer = map[Layer]string{
	DefaultLayer: "Default",
}

// NewLayer creates a new collision layer with the given name.
// It panics if the maximum number of layers (255) is exceeded.
func NewLayer(name string) Layer {
	if len(nameByLayer) >= 255 {
		panic("maximum number of layers reached")
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

// =========== Behaviour ==========

// Behaviour represents the collision behavior type for a collider.
// It defines how the collider interacts with other colliders in the physics simulation.
//
// Types:
//   - StaticCollision: Immovable objects (e.g., walls, floors).
//   - DynamicCollision: Movable objects affected by physics (e.g., players, enemies).
//   - SensorCollision: Trigger volumes that detect overlaps but do not cause physical responses (e.g., pickups, zones).
type Behaviour uint8

const (
	StaticBehaviour  Behaviour = iota // Immovable object (e.g., walls, floors)
	DynamicBehaviour                  // Movable object affected by physics (e.g., players, enemies)
	SensorBehaviour                   // Trigger volume, detects overlaps but no physical response (e.g., pickups, zones)
)

// String returns the string representation of the Behaviour.
func (ct Behaviour) String() string {
	switch ct {
	case StaticBehaviour:
		return "Static"
	case DynamicBehaviour:
		return "Dynamic"
	case SensorBehaviour:
		return "Sensor"
	default:
		return "Unknown"
	}
}

// IsValid reports whether the Behaviour value is a valid collision type.
func (ct Behaviour) IsValid() bool {
	return ct <= SensorBehaviour
}

// =========== Detection ==========

// Detection represents the collision detection mode for a collider.
// It defines how collisions are detected and handled during the physics simulation.
//
// Types:
//   - DiscreteDetection: Standard collision detection performed at discrete time steps.
//   - ContinuousDetection: Advanced collision detection that prevents fast-moving objects from tunneling through other colliders.
type Detection uint8

const (
	DiscreteDetection   Detection = iota // Standard discrete collision detection
	ContinuousDetection                  // Continuous collision detection to prevent tunneling
)

func (dt Detection) String() string {
	switch dt {
	case DiscreteDetection:
		return "Discrete"
	case ContinuousDetection:
		return "Continuous"
	default:
		return "Unknown"
	}
}

func (dt Detection) IsValid() bool {
	return dt <= ContinuousDetection
}
