package server

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"

	"github.com/PawelZabc/ProjektZespolowy/internal/config"
	"github.com/PawelZabc/ProjektZespolowy/internal/protocol"
)

// Network handles all server network operations
type Network struct {
	conn          *net.UDPConn
	clientManager *ClientManager
	gameState     *GameState
	buffer        []byte
	updateCount   int64
}

func NewNetwork(port int, clientManager *ClientManager, gameState *GameState) (*Network, error) {
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
		conn:          conn,
		clientManager: clientManager,
		gameState:     gameState,
		buffer:        make([]byte, config.NetworkBufferSize),
		updateCount:   0,
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
			n.clientManager.HandleMessage(n.buffer[:bytesRead], clientAddr, currentUpdate)
		}
	}
}

func (n *Network) SetUpdateCount(count int64) {
	atomic.StoreInt64(&n.updateCount, count)
}

func (n *Network) BroadcastGameState() {
	clients := n.clientManager.GetAllClients()
	enemy := n.gameState.GetEnemy()

	for _, player := range clients {
		playerData := make([]protocol.PlayerData, 0, len(clients)-1)
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

func (n *Network) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
}
