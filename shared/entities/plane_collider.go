package entities

import (
	types "github.com/PawelZabc/ProjektZespolowy/shared/_types"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PlaneCollider struct {
	Position  rl.Vector3
	Direction types.Direction
	Width     float32
	Height    float32
}

func (p *PlaneCollider) AddPosition(pos rl.Vector3) {
	p.Position = rl.Vector3Add(p.Position, pos)
}

func (p *PlaneCollider) CollidesWith(types.Collider) bool {
	return false
}

func (p *PlaneCollider) GetPosition() rl.Vector3 {
	return p.Position
}

func (p *PlaneCollider) PushbackFrom(types.Collider) types.Direction {
	return types.DirNone
}

func (p *PlaneCollider) SetPosition(pos rl.Vector3) {
	p.Position = pos
}

var _ types.Collider = (*PlaneCollider)(nil)
