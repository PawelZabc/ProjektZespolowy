package server

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"

	"github.com/PawelZabc/ProjektZespolowy/internal/config"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/entities"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	"github.com/PawelZabc/ProjektZespolowy/internal/protocol"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Network handles all server network operations
// receiving input data and distributing them to players to handle
// adding clients to the server
// removing disconnected clients
// broadcasting updated game state to clients
type Network struct {
	conn        *net.UDPConn
	gameState   *GameState
	buffer      []byte
	updateCount int64
}

func NewNetwork(port int, gameState *GameState) (*Network, error) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	// TODO: move to logging
	fmt.Printf("Server listening on port %d\n", port)

	return &Network{
		conn:        conn,
		gameState:   gameState,
		buffer:      make([]byte, config.NetworkBufferSize),
		updateCount: 0,
	}, nil
}

func (n *Network) StartReceiving(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			bytesRead, clientAddr, err := n.conn.ReadFromUDP(n.buffer)
			if err != nil {
				fmt.Printf("Read error: %v\n", err)
				continue
			}

			currentUpdate := atomic.LoadInt64(&n.updateCount)
			n.handleClientMessage(n.buffer[:bytesRead], clientAddr, currentUpdate)
		}
	}
}

func (n *Network) SetUpdateCount(count int64) {
	atomic.StoreInt64(&n.updateCount, count)
}

// Sends updated game state to clients
func (n *Network) BroadcastGameState() {
	clients := n.gameState.clients
	enemy := n.gameState.enemy

	for _, player := range clients {
		playerData := make([]protocol.PlayerData, 0, len(clients)-1)
		
		// Sending your player data to OTHER players, not to yourself
		for _, otherPlayer := range clients {
			if otherPlayer.Address.String() != player.Address.String() {
				playerData = append(playerData, protocol.PlayerData{
					Position: otherPlayer.Collider.GetPosition(),
					Rotation: otherPlayer.RotationX,
					Id:       otherPlayer.Id,
				})
			}
		}

		serverData := protocol.ServerData{
			Position: player.GetPosition(),
			Players:  playerData,
			Enemy:    protocol.EnemyData{Position: enemy.Collider.GetPosition(), Rotation: enemy.RotationX, AnimationFrame: uint8(enemy.State)},
			PlayerHp: player.Hp,
		}

		data := protocol.SerializeServerData(serverData)
		_, err := n.conn.WriteToUDP(data, player.Address)
		if err != nil {
			fmt.Printf("Failed to send to %s: %v\n", player.Address, err)
		}
	}
}

// Removes players that were not connected in some time
func (n *Network) RemoveDisconnectedClients(currentUpdate int64) {
	for addr, player := range n.gameState.clients {
		if currentUpdate-player.LastMessage > config.ClientTimeoutTicks {
			fmt.Printf("Client disconnected: %s (ID: %d)\n", addr, player.Id)
			delete(n.gameState.clients, addr)
		}
	}
}

// Handles input data received from client. If the input addr is new - add a client.
// If the input addr is known it updates player data
// LIVES IN GOROUTINE
func (n *Network) handleClientMessage(data []byte, addr *net.UDPAddr, updateCount int64) {
	addrStr := addr.String()

	player, exists := n.gameState.clients[addrStr]
	if !exists {
		n.addClient(addr, updateCount)
		return
	}

	// Prevent duplicate processing in same update
	if player.LastMessage == updateCount {
		return
	}

	player.LastMessage = updateCount
	player.ProcessInput(protocol.DeserializeClientData(data))
}

// Adding client to the game
// LIVES IN GOROUTINE
// TODO: Maybe refactor after "player renovation"
func (n *Network) addClient(addr *net.UDPAddr, updateCount int64) {
	player := &entities.Player{
		Velocity: rl.Vector3{},
		Collider: colliders.NewCylinderCollider(
			rl.NewVector3(0, 0, 0),
			config.PlayerRadius,
			config.PlayerHeight,
		),
		Speed:       config.PlayerSpeed,
		Address:     addr,
		Id:          n.gameState.nextPlayerId,
		LastMessage: updateCount,
		Hp:          100,
	}

	n.gameState.clients[addr.String()] = player
	n.gameState.nextPlayerId++

	fmt.Printf("New client connected: %s (ID: %d)\n", addr.String(), player.Id)
}

func (n *Network) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
}
