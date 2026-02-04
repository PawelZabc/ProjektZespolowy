package entities

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
)

// STARA RZECZ - powinno być entity
// wszędzie gdzie jest object powinno być entity
// This is valid object struct for client
// server just pack colliders to array
// only walls have multiple colliders for now
type Object struct {
	Colliders []colliders.Collider
	DrawPoint rl.Vector3 // brak, jest pozycja w entity
	Model     rl.Model   // w rednerable
	Color     rl.Color   // w rednerable
	Shader    rl.Shader  // w rednerable
}

func (o Object) Draw() {
	rl.DrawModel(o.Model, o.DrawPoint, 1, o.Color)
}
