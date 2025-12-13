package colliders

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/shared"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Collider interface {
	CollidesWith(Collider) bool
	GetPosition() rl.Vector3
	SetPosition(rl.Vector3)
	AddPosition(rl.Vector3)
	PushbackFrom(Collider) shared.Direction
}
