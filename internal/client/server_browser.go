package client

import (
	"net"
	"time"

	"github.com/PawelZabc/ProjektZespolowy/internal/protocol"
)

type PlayerInfo struct {
	ID   uint16
	Name string
}

// Connects to server and browse players
type ServerBrowser struct {
	serverIP string
	port     int
	Players  []PlayerInfo
	Online   bool
	lastPing time.Time
}

func NewServerBrowser(serverIP string, port int) *ServerBrowser {
	return &ServerBrowser{
		serverIP: serverIP,
		port:     port,
		Players:  make([]PlayerInfo, 0),
		Online:   false,
	}
}

func (sb *ServerBrowser) Refresh() {
	if time.Since(sb.lastPing) < 2*time.Second {
		return
	}
	sb.lastPing = time.Now()

	serverAddr := net.UDPAddr{
		Port: sb.port,
		IP:   net.ParseIP(sb.serverIP),
	}

	conn, err := net.DialUDP("udp", nil, &serverAddr)
	if err != nil {
		sb.Online = false
		sb.Players = nil
		return
	}
	defer conn.Close()

	conn.SetWriteDeadline(time.Now().Add(500 * time.Millisecond))
	_, err = conn.Write(protocol.SerializePing())
	if err != nil {
		sb.Online = false
		sb.Players = nil
		return
	}

	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	buffer := make([]byte, 256)
	n, err := conn.Read(buffer)
	if err != nil {
		sb.Online = false
		sb.Players = nil
		return
	}

	pong := protocol.DeserializePong(buffer[:n])
	if pong.Type != protocol.MsgPong {
		sb.Online = false
		sb.Players = nil
		return
	}

	sb.Online = true
	sb.Players = make([]PlayerInfo, len(pong.PlayerIDs))
	for i, id := range pong.PlayerIDs {
		sb.Players[i] = PlayerInfo{ID: id, Name: "Player"}
	}
}
