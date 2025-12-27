package entities

import (
	"log"

	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/state"
	"github.com/chewxy/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Enemy struct {
	RotationX      float32
	Collider       colliders.Collider
	Target         *rl.Vector2
	Speed          float32
	State          state.EnemyState
	AttackTimer    uint8
	AttackCooldown uint8
}

func (e *Enemy) Attack(players []*Player, colliders *[]colliders.Collider) {
	for _, player := range players {
		if e.GetDistanceToCollider(player.Collider) < 15 {
			ray := e.GetRayFromTopToColliderTop(player.Collider)
			point, length := ray.GetCollisionPoint(player.Collider)
			if point != nil && length < 10 {
				direct := true // ASK (to Pabox): Is this flag needed?
				for _, collider := range *colliders {
					point2, length2 := ray.GetCollisionPoint(collider)
					if point2 != nil && length2 < length {
						direct = false
						break
					}
				}
				if direct {
					player.Hit(5)
					log.Println("Ała! Nie w szczepionkę!")
				}
			}

		}
	}
	e.AttackCooldown = 20
}

func (e *Enemy) Update(players []*Player, colliders *[]colliders.Collider) {
	if e.AttackCooldown > 0 {
		e.AttackCooldown -= 1
	}
	switch e.State {
	case state.Walking:
		e.UpdateTarget(players, colliders)
		e.Move()
	case state.Attacking:
		e.AttackTimer -= 1
		if e.AttackTimer <= 0 {
			e.Attack(players, colliders)
			e.State = state.Walking
		}
	}
}

func (e *Enemy) SetState(s state.EnemyState) {
	if s != e.State {
		if s == state.Attacking {
			if e.AttackCooldown == 0 {
				e.State = state.Attacking
				e.AttackTimer = 50
			}
		} else {
			e.State = s
		}
	}
}

func (e *Enemy) Move() {
	if e.Target != nil {
		position2 := physics.GetVector2XZ(e.Collider.GetPosition())
		difference := rl.Vector2Subtract(*e.Target, position2)
		direction := rl.Vector2Normalize(difference)
		e.RotationX = physics.GetRotationX(direction)
		movement := rl.Vector2Scale(direction, e.Speed)
		velocity := physics.GetVector3FromXZ(movement)
		e.Collider.AddPosition(velocity)
		if rl.Vector2Length(rl.Vector2Subtract(*e.Target, physics.GetVector2XZ(e.Collider.GetPosition()))) < 1 {
			e.Target = nil
		}

	} else {
		e.RotationX += 1
		if e.RotationX > 360 {
			e.RotationX -= 360
		}
	}

}

func (e *Enemy) UpdateTarget(players []*Player, cols *[]colliders.Collider) {
	minLength := float32(0)
	minId := -1
	for i, player := range players {

		ray := e.GetRayFromTopToColliderTop(player.Collider)
		point, length := ray.GetCollisionPoint(player.Collider)

		difference := ray.GetRotationX() - e.RotationX
		if point != nil && length < 15 && (math32.Abs(difference) < 45 || math32.Abs(difference) > 315) {
			minId2 := -1
			minLength2 := float32(0)
			for i, collider := range *cols {
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
			vector := physics.GetVector2XZ(players[minId].Collider.GetPosition())
			e.Target = &vector

		}

	}

}

func (e *Enemy) GetDistanceToCollider(collider colliders.Collider) float32 {
	target := physics.GetVector2XZ(collider.GetPosition())
	position := physics.GetVector2XZ(e.Collider.GetPosition())
	return rl.Vector2Distance(target, position)
}

// TODO: think of a better function name
func (e *Enemy) GetRayFromTopToColliderTop(collider colliders.Collider) colliders.Ray {
	colliderHeight := collider.GetSizeOnAxis(physics.DirY)
	enemyHeight := e.Collider.GetSizeOnAxis(physics.DirY)
	return colliders.Ray{Origin: rl.Vector3Add(e.Collider.GetPosition(), rl.NewVector3(0, enemyHeight, 0)),
		Direction: rl.Vector3Normalize(rl.Vector3Subtract(rl.Vector3Add(collider.GetPosition(), rl.NewVector3(0, colliderHeight, 0)), rl.Vector3Add(e.Collider.GetPosition(), rl.NewVector3(0, enemyHeight, 0))))}
}
