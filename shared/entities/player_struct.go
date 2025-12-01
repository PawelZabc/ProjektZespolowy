package entities

import (
	"net"

	types "github.com/PawelZabc/ProjektZespolowy/shared/_types"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	IsOnFloor bool
	Velocity  rl.Vector3
	Collider  types.Collider
	Address   *net.UDPAddr
	RotationX float32
	RotationY float32
	Movement  rl.Vector2
	Speed     float32
}

func (p *Player) Move() {
	p.Movement = rl.Vector2Normalize(p.Movement)
	p.Movement = rl.Vector2Rotate(p.Movement, rl.Deg2rad*90+p.RotationX)
	p.Movement = rl.Vector2Scale(p.Movement, p.Speed)
	p.Collider.AddPosition(rl.Vector3Add(p.Velocity, GetVector3FromXZ(p.Movement)))
	p.Movement = rl.Vector2{}
}

func (p *Player) GetPosition() rl.Vector3 {
	return p.Collider.GetPosition()
}

func (p *Player) AddPosition(vec rl.Vector3) {
	p.Collider.AddPosition(vec)
}
