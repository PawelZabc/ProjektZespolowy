package colliders

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PlaneCollider struct {
	Position  rl.Vector3
	Direction physics.Direction
	Width     float32
	Height    float32
}

func (p *PlaneCollider) AddPosition(pos rl.Vector3) {
	p.Position = rl.Vector3Add(p.Position, pos)
}

func (p *PlaneCollider) CollidesWith(Collider) bool {
	return false
}

func (p *PlaneCollider) GetPosition() rl.Vector3 {
	return p.Position
}

func (p *PlaneCollider) PushbackFrom(Collider) physics.Direction {
	return physics.DirNone
}

func (p *PlaneCollider) SetPosition(pos rl.Vector3) {
	p.Position = pos
}

func NewPlaneCollider(position rl.Vector3, Width float32, Height float32, Direction physics.Direction) *PlaneCollider {
	return &PlaneCollider{
		Position:  position,
		Width:     Width,
		Height:    Height,
		Direction: Direction,
	}
}

var _ Collider = (*PlaneCollider)(nil)
