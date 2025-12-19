package server

import (
	"fmt"
	"net"

	"github.com/PawelZabc/ProjektZespolowy/internal/config"
)

type Network struct {
	conn   *net.UDPConn
	buffer []byte
}

func NewNetwork(port int) (*Network, error) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return &Network{}, err
	}

	// TODO: move to logging
	fmt.Printf("Server listening on port %d\n", port)

	return &Network{
		conn:   conn,
		buffer: make([]byte, config.NetworkBufferSize),
	}, nil
}

func (n *Network) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
}
