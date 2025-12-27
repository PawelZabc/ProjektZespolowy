package server

import (
	"fmt"
	"net"

	"github.com/PawelZabc/ProjektZespolowy/internal/config"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/entities"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/input"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	"github.com/PawelZabc/ProjektZespolowy/internal/protocol"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// ClientManager manages all connected clients
type ClientManager struct {
	clients      map[string]*entities.Player
	gameState    *GameState
	serverConfig config.ServerConfig
	nextPlayerId uint16
}

func NewClientManager(gameState *GameState, cfg config.ServerConfig) *ClientManager {
	return &ClientManager{
		clients:      make(map[string]*entities.Player),
		gameState:    gameState,
		serverConfig: cfg,
		nextPlayerId: 0,
	}
}

func (cm *ClientManager) HandleMessage(data []byte, addr *net.UDPAddr, updateCount int64) {
	addrStr := addr.String()

	player, exists := cm.clients[addrStr]
	if !exists {
		cm.addClient(addr, updateCount)
		return
	}

	// Prevent duplicate processing in same update
	if player.LastMessage == updateCount {
		return
	}
	player.LastMessage = updateCount

	cm.processClientInput(player, data)
}

func (cm *ClientManager) addClient(addr *net.UDPAddr, updateCount int64) {
	player := &entities.Player{
		Velocity: rl.Vector3{},
		Collider: colliders.NewCylinderCollider(
			rl.NewVector3(0, 0, 0),
			config.PlayerRadius,
			config.PlayerHeight,
		),
		Speed:       config.PlayerSpeed,
		Address:     addr,
		Id:          cm.nextPlayerId,
		LastMessage: updateCount,
		Hp:          100,
	}

	cm.clients[addr.String()] = player
	cm.nextPlayerId++

	fmt.Printf("New client connected: %s (ID: %d)\n", addr.String(), player.Id)
}

func (cm *ClientManager) processClientInput(player *entities.Player, data []byte) {
	clientData := protocol.DeserializeClientData(data)

	player.RotationX = clientData.RotationX
	player.RotationY = clientData.RotationY

	player.Movement = rl.Vector2{}

	for _, i := range clientData.Inputs {
		switch i {
		case input.MoveForward:
			player.Movement.Y = 1
		case input.MoveBackward:
			player.Movement.Y = -1
		case input.MoveLeft:
			player.Movement.X = 1
		case input.MoveRight:
			player.Movement.X = -1
		case input.Jump:
			if player.IsOnFloor {
				player.Velocity.Y = config.JumpStrength
			}
		}
	}
}

func (cm *ClientManager) RemoveDisconnected(currentUpdate int64) {
	for addr, player := range cm.clients {
		if currentUpdate-player.LastMessage > config.ClientTimeoutTicks {
			fmt.Printf("Client disconnected: %s (ID: %d)\n", addr, player.Id)
			delete(cm.clients, addr)
		}
	}
}

func (cm *ClientManager) GetAllClients() map[string]*entities.Player {
	return cm.clients
}

func (cm *ClientManager) GetActivePlayers() []*entities.Player {
	players := make([]*entities.Player, 0, len(cm.clients))
	for _, player := range cm.clients {
		players = append(players, player)
	}
	return players
}
