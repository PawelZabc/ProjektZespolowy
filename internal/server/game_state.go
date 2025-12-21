package server

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/entities"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/levels"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// GameState holds all server-side game state
type GameState struct {
	enemy   *entities.Enemy
	rooms   []levels.Room
	objects []colliders.Collider
}

func NewGameState() *GameState {
	enemy := &entities.Enemy{
		Collider: colliders.NewCylinderCollider(
			rl.NewVector3(20, 0, 15),
			1,
			1,
		),
		Speed: 0.05,
	}

	rooms := levels.ServerLoadRooms()
	objects := rooms[0].Colliders

	return &GameState{
		enemy:   enemy,
		rooms:   rooms,
		objects: objects,
	}
}

func (gs *GameState) GetEnemy() *entities.Enemy {
	return gs.enemy
}

func (gs *GameState) GetObjects() []colliders.Collider {
	return gs.objects
}
