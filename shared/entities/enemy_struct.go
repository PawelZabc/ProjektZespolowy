package entities

import (
	types "github.com/PawelZabc/ProjektZespolowy/shared/_types"
	math "github.com/chewxy/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Enemy struct {
	RotationX      float32
	Collider       types.Collider
	Target         *rl.Vector2
	Speed          float32
	State          types.EnemyState
	AttackTimer    uint8
	AttackCooldown uint8
}

func (e *Enemy) Attack(players []*Player, colliders *[]types.Collider) {
	for _, player := range players {
		if e.GetDistanceToCollider(player.Collider) < 15 {
			ray := e.GetRayFromTopToColliderTop(player.Collider)
			point, length := ray.GetCollisionPoint(player.Collider)
			if point != nil && length < 10 {
				direct := true
				for _, collider := range *colliders {
					point2, length2 := ray.GetCollisionPoint(collider)
					if point2 != nil && length2 < length {
						direct = false
						break
					}
				}
				if direct {
					player.Hit(5)
				}
			}

		}
	}
	e.AttackCooldown = 20
}

func (e *Enemy) Update(players []*Player, colliders *[]types.Collider) {
	if e.AttackCooldown > 0 {
		e.AttackCooldown -= 1
	}
	switch e.State {
	case types.Walking:
		e.UpdateTarget(players, colliders)
		e.Move()
	case types.Attacking:
		e.AttackTimer -= 1
		if e.AttackTimer <= 0 {
			e.Attack(players, colliders)
			e.State = types.Walking
		}
	}
}

func (e *Enemy) SetState(state types.EnemyState) {
	if state != e.State {
		if state == types.Attacking {
			if e.AttackCooldown == 0 {
				e.State = types.Attacking
				e.AttackTimer = 50
			}
		} else {
			e.State = state
		}
	}
}

func (e *Enemy) Move() {
	if e.Target != nil {
		position2 := GetVector2XZ(e.Collider.GetPosition())
		difference := rl.Vector2Subtract(*e.Target, position2)
		direction := rl.Vector2Normalize(difference)
		e.RotationX = GetRotationX(direction)
		movement := rl.Vector2Scale(direction, e.Speed)
		velocity := GetVector3FromXZ(movement)
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

		ray := e.GetRayFromTopToColliderTop(player.Collider)
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

func (e *Enemy) GetDistanceToCollider(collider types.Collider) float32 {
	target := GetVector2XZ(collider.GetPosition())
	position := GetVector2XZ(e.Collider.GetPosition())
	return rl.Vector2Distance(target, position)
}

func (e *Enemy) GetRayFromTopToColliderTop(collider types.Collider) Ray { //think of a better function name
	colliderHeight := collider.GetSizeOnAxis(types.DirY)
	enemyHeight := e.Collider.GetSizeOnAxis(types.DirY)
	return Ray{Origin: rl.Vector3Add(e.Collider.GetPosition(), rl.NewVector3(0, enemyHeight, 0)),
		Direction: rl.Vector3Normalize(rl.Vector3Subtract(rl.Vector3Add(collider.GetPosition(), rl.NewVector3(0, colliderHeight, 0)), rl.Vector3Add(e.Collider.GetPosition(), rl.NewVector3(0, enemyHeight, 0))))}
}
