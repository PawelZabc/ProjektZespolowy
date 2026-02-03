package server

import (
	"maps"
	"slices"

	"github.com/PawelZabc/ProjektZespolowy/internal/game/entities"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/entities/server"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/levels"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Usefull type aliases
type Clients map[string]*entities.Player // map IP:PlayerPointer
type Players []*entities.Player          // slice of PlayerPointers

// GameState holds all server-side game state
type GameState struct {
	enemy *entities.Enemy
	// rooms   []levels.Room
	objects       []*server.CollisionObject
	effectObjects []*server.EffectObject

	// Players as clients
	clients      Clients
	nextPlayerId uint16
}

// Returns new GameState, used for initialisation of game
func NewGameState() *GameState {
	enemy := &entities.Enemy{
		Collider: colliders.NewCylinderCollider(
			rl.NewVector3(20, 0, 15),
			1,
			1,
		),
		Speed: 0.05,
	}

	// enemies := make(*entities.Enemy,0,4)
	// enemies = append(enemies,enemy)

	effect := func(playerId uint16) {
		println(playerId, " touched me")
	}

	effectObject := server.EffectObject{
		Effect:   &effect,
		Collider: colliders.NewCubeCollider(rl.NewVector3(0, 0, 0), 2, 2, 2),
	}

	effectObjects := make([]*server.EffectObject, 0, 10)
	effectObjects = append(effectObjects, &effectObject)

	rooms := levels.ServerLoadRooms()
	objects := make([]*server.CollisionObject, 0, len(rooms[0].Colliders))
	for _, collider := range rooms[0].Colliders {
		object := server.CollisionObject{Collider: &collider}
		objects = append(objects, &object)
	}

	return &GameState{
		enemy: enemy,
		// rooms:        rooms,
		objects:       objects,
		effectObjects: effectObjects,
		clients:       make(Clients, 2),
		nextPlayerId:  0,
	}
}

// Converts Clients to players slice (for updating in certain places - enemy f.e)
func (gs *GameState) GetClientsAsPlayerSlice() Players {
	return slices.Collect(maps.Values(gs.clients))
}
