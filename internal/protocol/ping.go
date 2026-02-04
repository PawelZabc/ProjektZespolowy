package protocol

import (
	"bytes"
	"encoding/binary"
)

const (
	MsgPing byte = 0x01
	MsgPong byte = 0x02
)

// Ping pong
type PingData struct {
	Type byte
}

type PongData struct {
	Type        byte
	PlayerCount uint8
	PlayerIDs   []uint16
}

func SerializePing() []byte {
	return []byte{MsgPing}
}

func IsPingMessage(data []byte) bool {
	return len(data) == 1 && data[0] == MsgPing
}

func SerializePong(playerIDs []uint16) []byte {
	buf := make([]byte, 0, 2+len(playerIDs)*2)
	b := bytes.NewBuffer(buf)

	binary.Write(b, binary.LittleEndian, MsgPong)
	binary.Write(b, binary.LittleEndian, uint8(len(playerIDs)))
	for _, id := range playerIDs {
		binary.Write(b, binary.LittleEndian, id)
	}

	return b.Bytes()
}

func DeserializePong(data []byte) PongData {
	if len(data) < 2 || data[0] != MsgPong {
		return PongData{}
	}

	b := bytes.NewReader(data)
	var pong PongData

	binary.Read(b, binary.LittleEndian, &pong.Type)
	binary.Read(b, binary.LittleEndian, &pong.PlayerCount)

	pong.PlayerIDs = make([]uint16, pong.PlayerCount)
	for i := uint8(0); i < pong.PlayerCount; i++ {
		binary.Read(b, binary.LittleEndian, &pong.PlayerIDs[i])
	}

	return pong
}
