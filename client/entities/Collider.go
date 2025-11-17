package entities

import (
	types "github.com/PawelZabc/ProjektZespolowy/client/_types"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Collider interface {
	CollidesWith(Collider) bool
	GetPosition() rl.Vector3
	SetPosition(rl.Vector3)
	AddPosition(rl.Vector3)
	PushbackFrom(Collider) types.Direction
}
