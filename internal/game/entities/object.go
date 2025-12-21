package entities

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
)

// This is valid object struct for client
// server just pack colliders to array
// only walls have multiple colliders for now
type Object struct {
	Colliders []colliders.Collider
	DrawPoint rl.Vector3
	Model     rl.Model
	Color     rl.Color
	Shader    rl.Shader
}

func (o Object) Draw() {
	rl.DrawModel(o.Model, o.DrawPoint, 1, o.Color)
}
