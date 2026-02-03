package client

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/animation"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/state"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func NewEnemy(model rl.Model, shader rl.Shader) *Actor {
	entity := CEntity{
		Renderable: &BasicRenderable{
			// Initialize with enemy model, color, shader, offset, rotation
			Model:    model,
			Color:    rl.White,
			Shader:   shader,
			Offset:   rl.NewVector3(0, 0, 0),
			Rotation: 0,
		},
		Position: rl.NewVector3(15, 0, 15),
		Collider: nil,
	}

	animationHandler := animation.AnimationHandler{}

	return &Actor{
		Entity:           entity,
		AnimationHandler: animationHandler,
		State:            state.Walking, // Initial state
		animationBase:    animation.NewEnemyAnimationMap(),
	}
}
