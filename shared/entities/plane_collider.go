package entities

import (
	types "github.com/PawelZabc/ProjektZespolowy/shared/_types"
	math "github.com/chewxy/math32"

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

func NewPlaneCollider(position rl.Vector3, Width float32, Height float32, Direction types.Direction) *PlaneCollider {
	return &PlaneCollider{
		Position:  position,
		Width:     Width,
		Height:    Height,
		Direction: Direction,
	}
}

var _ types.Collider = (*PlaneCollider)(nil)

func CreateRoomWallsFromChanges(StartPoint rl.Vector3, Changes []types.Change, Height float32) []types.Collider {

	count := 0
	for _, change := range Changes {
		if !change.Skip {
			count++
		}
	}
	walls := make([]types.Collider, len(Changes))
	skipped := 0
	for i, change := range Changes {
		if change.Axis == types.DirX {
			change.Axis = types.DirZ
		} else {
			change.Axis = types.DirX
		}

		var object PlaneCollider
		if change.Value < 0 {
			if change.Axis == types.DirX {
				StartPoint = rl.Vector3Add(StartPoint, rl.NewVector3(0, 0, change.Value))
			} else {
				StartPoint = rl.Vector3Add(StartPoint, rl.NewVector3(change.Value, 0, 0))
			}
		}
		if !change.Skip {
			object = *NewPlaneCollider(StartPoint, math.Abs(change.Value), Height, change.Axis)
			walls[i-skipped] = &object
		} else {
			skipped += 1
		}
		if change.Value > 0 {
			if change.Axis == types.DirX {
				StartPoint = rl.Vector3Add(StartPoint, rl.NewVector3(0, 0, change.Value))
			} else {
				StartPoint = rl.Vector3Add(StartPoint, rl.NewVector3(change.Value, 0, 0))
			}
		}

	}

	return walls
}
