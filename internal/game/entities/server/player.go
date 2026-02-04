package server

import (
	"net"

	"github.com/PawelZabc/ProjektZespolowy/internal/config"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/input"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	"github.com/PawelZabc/ProjektZespolowy/internal/protocol"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	IsOnFloor   bool
	Velocity    rl.Vector3
	Collider    colliders.Collider
	Address     *net.UDPAddr
	RotationX   float32
	RotationY   float32
	Movement    rl.Vector2
	Speed       float32
	LastMessage int64
	Item        *Item
	Id          uint16
	Hp          uint8
}

func (p *Player) Move() {
	p.Movement = rl.Vector2Normalize(p.Movement)
	p.Movement = rl.Vector2Rotate(p.Movement, rl.Deg2rad*90+p.RotationX)
	p.Movement = rl.Vector2Scale(p.Movement, p.Speed)
	p.Collider.AddPosition(rl.Vector3Add(p.Velocity, physics.GetVector3FromXZ(p.Movement)))
	p.Movement = rl.Vector2{} //TODO: change movement resseting every tick
}

func (p *Player) GetPosition() rl.Vector3 {
	return p.Collider.GetPosition()
}

func (p *Player) AddPosition(vec rl.Vector3) {
	p.Collider.AddPosition(vec)
}

func (p *Player) PushbackFrom(collider colliders.Collider) {
	if collider != nil {
		direction := p.Collider.PushbackFrom(collider)
		switch direction {
		case physics.DirYminus:
			p.IsOnFloor = true
			p.Velocity.Y = 0
		case physics.DirY:
			p.Velocity.Y = 0
		}
	}
}

func (p *Player) DropItem() {
	if p.Item != nil {
		p.Item.EffectObject.Collider.SetPosition(p.GetPosition())
		p.Item.EffectObject.Active = true
	}
}

// Changes player position based on data received from client
// LIVES IN GOROUTINE
func (p *Player) ProcessInput(data protocol.ClientData) {

	p.RotationX = data.RotationX
	p.RotationY = data.RotationY

	p.Movement = rl.Vector2{}

	for _, i := range data.Inputs {
		switch i {
		case input.MoveForward:
			p.Movement.Y = 1
		case input.MoveBackward:
			p.Movement.Y = -1
		case input.MoveLeft:
			p.Movement.X = 1
		case input.MoveRight:
			p.Movement.X = -1
		case input.Jump:
			if p.IsOnFloor {
				p.Velocity.Y = config.JumpStrength
			}
		}
	}
}

func (p *Player) Hit(damage uint8) {
	if p.Hp <= damage {
		p.Hp = 0
		p.DropItem()
	} else {
		p.Hp -= damage
	}
}
