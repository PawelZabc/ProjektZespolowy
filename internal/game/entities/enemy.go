package entities

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Enemy struct {
	RotationX float32
	Collider  colliders.Collider
	Target    *rl.Vector2
	Speed     float32
}

func (e *Enemy) Move() {
	if e.Target != nil {
		position2 := physics.GetVector2XZ(e.Collider.GetPosition())
		difference := rl.Vector2Subtract(*e.Target, position2)
		direction := rl.Vector2Normalize(difference)
		e.RotationX = GetRotationX(direction)
		movement := rl.Vector2Scale(direction, e.Speed)
		velocity := physics.GetVector3FromXZ(movement)
		e.Collider.AddPosition(velocity)
		if rl.Vector2Length(rl.Vector2Subtract(*e.Target, GetVector2XZ(e.Collider.GetPosition()))) < 0.5 {
			e.Target = nil
		}

	} else {
		e.RotationX += 1
		if e.RotationX > 360 {
			e.RotationX -= 360
		}
	}

}

func (e *Enemy) UpdateTarget(players []*Player, colliders *[]types.Collider) {
	minLength := float32(0)
	minId := -1
	for i, player := range players {

		playerHeight := float32(0)
		enemyHeight := float32(0)
		if cylinder, ok := player.Collider.(*CylinderCollider); ok {
			playerHeight = cylinder.Height
		}
		if cylinder2, ok := e.Collider.(*CylinderCollider); ok {
			enemyHeight = cylinder2.Height
		}
		ray := Ray{Origin: rl.Vector3Add(e.Collider.GetPosition(), rl.NewVector3(0, enemyHeight, 0)),
			Direction: rl.Vector3Normalize(rl.Vector3Subtract(rl.Vector3Add(player.GetPosition(), rl.NewVector3(0, playerHeight, 0)), rl.Vector3Add(e.Collider.GetPosition(), rl.NewVector3(0, enemyHeight, 0))))}
		point, length := ray.GetCollisionPoint(player.Collider)

		difference := ray.GetRotationX() - e.RotationX
		if point != nil && length < 15 && (math.Abs(difference) < 45 || math.Abs(difference) > 315) {
			minId2 := -1
			minLength2 := float32(0)
			for i, collider := range *colliders {
				point2, length2 := ray.GetCollisionPoint(collider)
				if point2 != nil && length2 < length && (minId2 == -1 || minLength2 > length2) {
					minId2 = i
					minLength2 = length2

				}
			}
			if minId2 == -1 {
				if minLength > length || minId == -1 {
					minId = i
					minLength = length
				}

				// } else {
				// 	if cube, ok := (*colliders)[minId2].(*CubeCollider); ok {
				// 		playerPos := GetVector2XZ(player.GetPosition())
				// 		// objectPos := cube.GetPosition()
				// 		closerWallXToPlayer := GetCloserWall(playerPos.X, cube.Position.X, cube.SizeX)
				// 		if closerWallXToPlayer > 0 {

				// 		}
				// 	}

			}

		}
		if minId != -1 {
			vector := GetVector2XZ(players[minId].Collider.GetPosition())
			e.Target = &vector

		}

	}

}

func (e *Enemy) GetDistanceToCollider(collider colliders.Collider) float32 {
	target := physics.GetVector2XZ(collider.GetPosition())
	position := physics.GetVector2XZ(e.Collider.GetPosition())
	return rl.Vector2Distance(target, position)
}
