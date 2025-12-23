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

func NewPlaneCollider(position rl.Vector3, Width float32, Height float32, Direction physics.Direction) *PlaneCollider {
	return &PlaneCollider{
		Position:  position,
		Width:     Width,
		Height:    Height,
		Direction: Direction,
	}
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

func (p PlaneCollider) GetSizeOnAxis(axis physics.Direction) float32 { //test this
	if axis == p.Direction || axis == -p.Direction {
		return 0
	}
	switch axis {
	case physics.DirX:
		return p.Width
	case physics.DirY:
		return p.Height
	case physics.DirZ:
		if p.Direction == physics.DirX {
			return p.Width
		} else {
			return p.Height
		}
	}
	return 0
}

var _ Collider = (*PlaneCollider)(nil)
