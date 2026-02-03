package client

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/animation"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/state"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func NewEnemy() *Actor {
	entity := CEntity{
		Renderable: &BasicRenderable{
			// Initialize with enemy model, color, shader, offset, rotation
		},
		Position: rl.Vector3{X: 0, Y: 0, Z: 0},
		Collider: nil, // Initialize with appropriate collider
	}

	animationHandler := animation.AnimationHandler{}
	// Initialize animationHandler with enemy animations

	return &Actor{
		Entity:           entity,
		AnimationHandler: animationHandler,
		State:            state.Walking, // Initial state
	}
}
