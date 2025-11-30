package entities

import (
	"net"

	types "github.com/PawelZabc/ProjektZespolowy/shared/_types"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	Position  rl.Vector3
	IsOnFloor bool
	Velocity  rl.Vector3
	Collider  types.Collider
	Address   net.UDPAddr
}
