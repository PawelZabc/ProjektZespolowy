package udp_data

import (
	"bytes"
	"encoding/binary"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ServerData struct {
	Players  []PlayerData
	Position rl.Vector3
}

type PlayerData struct {
	Position rl.Vector3
	Id       uint16
}

func SerializeServerData(s ServerData) []byte {
	// 1 byte for number of players + 12 bytes for position + 14 bytes for position and id * players
	size := 1 + 12 + len(s.Players)*14
	buf := make([]byte, 0, size)
	b := bytes.NewBuffer(buf)

	binary.Write(b, binary.LittleEndian, uint8(len(s.Players)))
	binary.Write(b, binary.LittleEndian, s.Position.X)
	binary.Write(b, binary.LittleEndian, s.Position.Y)
	binary.Write(b, binary.LittleEndian, s.Position.Z)
	for _, p := range s.Players {
		binary.Write(b, binary.LittleEndian, p.Position.X)
		binary.Write(b, binary.LittleEndian, p.Position.Y)
		binary.Write(b, binary.LittleEndian, p.Position.Z)
		binary.Write(b, binary.LittleEndian, p.Id)
	}
	bytes := b.Bytes()
	fmt.Println(bytes)
	return bytes
}

func DeserializeServerData(data []byte) ServerData {
	b := bytes.NewReader(data)
	var s ServerData

	var count uint8
	binary.Read(b, binary.LittleEndian, &count)

	binary.Read(b, binary.LittleEndian, &s.Position.X)
	binary.Read(b, binary.LittleEndian, &s.Position.Y)
	binary.Read(b, binary.LittleEndian, &s.Position.Z)

	s.Players = make([]PlayerData, count)
	for i := 0; i < int(count); i++ {
		binary.Read(b, binary.LittleEndian, &s.Players[i].Position.X)
		binary.Read(b, binary.LittleEndian, &s.Players[i].Position.Y)
		binary.Read(b, binary.LittleEndian, &s.Players[i].Position.Z)
		binary.Read(b, binary.LittleEndian, &s.Players[i].Id)
	}

	return s
}
