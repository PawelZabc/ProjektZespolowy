package server

import (
	"maps"
	"slices"

	"github.com/PawelZabc/ProjektZespolowy/internal/game/entities/server"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/levels"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Usefull type aliases
type Clients map[string]*server.Player // map IP:PlayerPointer
type Players []*server.Player          // slice of PlayerPointers

// GameState holds all server-side game state
type GameState struct {
	enemy *server.Enemy
	// rooms   []levels.Room
	objects       []*server.CollisionObject
	effectObjects []*server.EffectObject
	items         []*server.Item

	// Players as clients
	clients      Clients
	progress     int16
	nextPlayerId uint16
}

// Returns new GameState, used for initialisation of game
func NewGameState() *GameState {
	enemy := &server.Enemy{
		Collider: colliders.NewCylinderCollider(
			rl.NewVector3(20, 0, 15),
			1,
			1,
		),
		Speed: 0.05,
	}
	clients := make(Clients, 2)

	rooms := levels.ServerLoadRooms()
	objects := make([]*server.CollisionObject, 0, len(rooms[0].Colliders))
	for _, collider := range rooms[0].Colliders {
		object := server.CollisionObject{Collider: &collider}
		objects = append(objects, &object)
	}
	gameState := GameState{
		enemy: enemy,
		// rooms:        rooms,
		objects: objects,
		// items:         items,
		// effectObjects: effectObjects,
		clients:      clients,
		nextPlayerId: 0,
	}

	items := make([]*server.Item, 0, 10)
	effectObjects := make([]*server.EffectObject, 0, 10)
	// enemies := make(*server.Enemy,0,4)
	// enemies = append(enemies,enemy)

	// effect := func(playerIp string) {
	// 	println(playerIp, " touched me")
	// }

	//for every item
	item := &server.Item{Type: server.ItemRepair, Id: 0}

	effect := ReturnItemPickupEffect(item, &clients)

	effectObject := &server.EffectObject{
		Effect:   effect,
		Collider: colliders.NewCylinderCollider(rl.NewVector3(-9, 0, -9), 0.5, 1),
		Active:   true,
	}
	effectObjects = append(effectObjects, effectObject)
	item.EffectObject = effectObject
	items = append(items, item)
	//end for every item

	effect2 := ReturnItemPutdownEffect(&clients)
	effectObject2 := &server.EffectObject{
		Effect:   effect2,
		Collider: colliders.NewCubeCollider(rl.NewVector3(-3, 0, 14), 6, 3, 6),
		Active:   true,
	}
	effectObjects = append(effectObjects, effectObject2)

	gameState.effectObjects = effectObjects
	gameState.items = items
	return &gameState
}

func ReturnItemPickupEffect(item *server.Item, clients *Clients) *func(playerIp string) {
	fun := func(playerIp string) {
		if player, exists := (*clients)[playerIp]; exists && player.Item == nil {
			println("player", playerIp, "picked up item type", item.Type)
			player.Item = item
			item.EffectObject.Active = false
		}

	}
	return &fun

}

func ReturnItemPutdownEffect(clients *Clients) *func(playerIp string) {
	fun := func(playerIp string) {
		if player, exists := (*clients)[playerIp]; exists && player.Item != nil {
			println("player", playerIp, "put down item type", player.Item.Type)
			player.Item.EffectObject.Active = true
			player.Item = nil
		}

	}
	return &fun

}

// Converts Clients to players slice (for updating in certain places - enemy f.e)
func (gs *GameState) GetClientsAsPlayerSlice() Players {
	return slices.Collect(maps.Values(gs.clients))
}
