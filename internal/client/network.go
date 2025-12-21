package client

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/PawelZabc/ProjektZespolowy/internal/config"
	"github.com/PawelZabc/ProjektZespolowy/internal/protocol"
)

type Network struct {
	conn      *net.UDPConn
	gameState *GameState
	buffer    []byte
}

func NewNetwork(serverIP string, port int, gameState *GameState) (*Network, error) {
	serverAddr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(serverIP),
	}

	localAddr := net.UDPAddr{
		Port: 0, // random port - maybe change in the future?
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.DialUDP("udp", &localAddr, &serverAddr)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to server: %w", err)
	}

	return &Network{
		conn:      conn,
		gameState: gameState,
		buffer:    make([]byte, config.NetworkBufferSize),
	}, nil
}

func (n *Network) StartReceiving(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// timeout to 1s - maybe move to config?
			n.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			bytesRead, _, err := n.conn.ReadFromUDP(n.buffer)
			if err != nil {
				continue
			}

			serverData := protocol.DeserializeServerData(n.buffer[:bytesRead])
			n.gameState.UpdateFromServer(serverData)
		}
	}
}

func (n *Network) SendInput(input protocol.ClientData) error {
	data := protocol.SerializeClientData(input)
	_, err := n.conn.Write(data)
	return err
}

func (n *Network) SendHello() error {
	_, err := n.conn.Write([]byte("hejka"))
	return err
}

func (n *Network) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
}
