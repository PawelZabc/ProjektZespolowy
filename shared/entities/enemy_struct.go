package entities

import (
	types "github.com/PawelZabc/ProjektZespolowy/shared/_types"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Enemy struct {
	RotationX float32
	Collider  types.Collider
	Target    *rl.Vector2
	Speed     float32
}

func (e *Enemy) Move() {
	if e.Target != nil {
		position2 := GetVector2XZ(e.Collider.GetPosition())
		difference := rl.Vector2Subtract(*e.Target, position2)
		direction := rl.Vector2Normalize(difference)
		e.RotationX = rl.Vector2Angle(direction, rl.NewVector2(-1, 0)) * rl.Rad2deg
		movement := rl.Vector2Scale(direction, e.Speed)
		velocity := GetVector3FromXZ(movement)
		e.Collider.AddPosition(velocity)

	}

}

func (e *Enemy) UpdateTarget(players []*Player) {
	minLength := float32(0)
	minId := -1
	for i, player := range players {
		length := e.GetDistanceToCollider(player.Collider)
		if minLength > length || minId == -1 {
			minId = i
			minLength = length
		}
	}
	if minId != -1 {
		vector := GetVector2XZ(players[minId].Collider.GetPosition())
		e.Target = &vector

	} else {
		e.Target = nil
	}

}

func (e *Enemy) GetDistanceToCollider(collider types.Collider) float32 {
	target := GetVector2XZ(collider.GetPosition())
	position := GetVector2XZ(e.Collider.GetPosition())
	return rl.Vector2Distance(target, position)
}
