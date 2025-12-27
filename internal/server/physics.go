package server

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/config"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/entities"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/state"
)

// Physics handles all physics simulation
type Physics struct {
	gameState     *GameState
	clientManager *ClientManager
}

func NewPhysics(gameState *GameState, clientManager *ClientManager) *Physics {
	return &Physics{
		gameState:     gameState,
		clientManager: clientManager,
	}
}

func (p *Physics) Update() {
	players := p.clientManager.GetActivePlayers()
	objects := p.gameState.GetObjects()
	enemy := p.gameState.GetEnemy()

	p.updatePlayers(players, objects, enemy)
	p.updateEnemy(players, objects, enemy)
}

func (p *Physics) updatePlayers(players []*entities.Player, objects []colliders.Collider, enemy *entities.Enemy) {
	for _, player := range players {
		player.Velocity.Y -= config.Gravity

		player.Move()
		player.IsOnFloor = false

		for _, obj := range objects {
			player.PushbackFrom(obj)
		}

		player.PushbackFrom(enemy.Collider)
	}
}

func (p *Physics) updateEnemy(players []*entities.Player, objects []colliders.Collider, enemy *entities.Enemy) {
	enemy.Update(players, &objects)

	for _, obj := range objects {
		if obj != nil {
			enemy.Collider.PushbackFrom(obj)
		}
	}

	for _, player := range players {
		if enemy.Collider.PushbackFrom(player.Collider) != physics.DirNone {
			enemy.SetState(state.Attacking)
		}
	}
}
