package entities

import (
	"net"

	types "github.com/PawelZabc/ProjektZespolowy/shared/_types"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	IsOnFloor   bool
	Velocity    rl.Vector3
	Collider    types.Collider
	Address     *net.UDPAddr
	RotationX   float32
	RotationY   float32
	Movement    rl.Vector2
	Speed       float32
	LastMessage int64
	Id          uint16
	Hp          uint8
}

func (p *Player) Hit(damage uint8) {
	if p.Hp <= damage {
		p.Hp = 0
	} else {
		p.Hp -= damage
	}
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

func (p *Player) PushbackFrom(collider types.Collider) {
	if collider != nil {
		direction := p.Collider.PushbackFrom(collider)
		if direction == types.DirYminus {
			p.IsOnFloor = true
			p.Velocity.Y = 0
		} else if direction == types.DirY {
			p.Velocity.Y = 0
		}
	}
}
