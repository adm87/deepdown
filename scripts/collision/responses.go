package collision

// Contact represents a collision contact between two colliders.
type Contact struct {
	NormalX, NormalY float32
	Depth            float32
	ColliderA        Collider
	ColliderB        Collider
}

// ResponseFunc is a function that handles a collision contact.
type ResponseFunc func(contact Contact)

// Interactions maps collision layers to their corresponding response functions.
type Interactions map[Layer]ResponseFunc

// Profiles maps collision layers to their interaction mappings.
// For example, a profile can define how a "Player" layer interacts with "Enemy" and "Wall" layers.
//
// Builds collision response table example:
//
//	         	Player  Enemy   Wall    Pickup
//			Player     -     Attack  Block   Collect
//			Enemy    Attack    -     Block     -
//			Wall     Block   Block     -       -
//			Pickup   Collect   -       -       -
//
// Example usage:
//
//	var PlayerProfile = collision.Profiles{
//	    LayerPlayer: {
//	        LayerEnemy: handlePlayerEnemyCollision,
//	        LayerWall:  handlePlayerWallCollision,
//	    },
//	}
type Profiles map[Layer]Interactions
