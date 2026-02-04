package protocol

import (
	"bytes"
	"encoding/binary"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ServerData struct {
	Players  []PlayerData
	Items    []ItemData
	Position rl.Vector3
	Enemy    EnemyData
	PlayerHp uint8
	ItemHeld uint8
}

type PlayerData struct {
	Position rl.Vector3
	Rotation float32
	Id       uint16
}

type EnemyData struct {
	Position       rl.Vector3
	Rotation       float32
	AnimationFrame uint8
}

type ItemData struct {
	Position rl.Vector2
	Type     uint8
	Id       uint8
}

func SerializeServerData(s ServerData) []byte {
	// 4 bytes for number of players + 12 bytes for position * players + 12 bytes + 4 bytes for enemy
	size := 4 + 12 + len(s.Players)*12 + 16
	buf := make([]byte, 0, size)
	b := bytes.NewBuffer(buf)

	binary.Write(b, binary.LittleEndian, uint8(len(s.Players)))
	binary.Write(b, binary.LittleEndian, s.Position.X)
	binary.Write(b, binary.LittleEndian, s.Position.Y)
	binary.Write(b, binary.LittleEndian, s.Position.Z)
	binary.Write(b, binary.LittleEndian, s.PlayerHp)
	for _, p := range s.Players {
		binary.Write(b, binary.LittleEndian, p.Position.X)
		binary.Write(b, binary.LittleEndian, p.Position.Y)
		binary.Write(b, binary.LittleEndian, p.Position.Z)
		binary.Write(b, binary.LittleEndian, p.Rotation)
		binary.Write(b, binary.LittleEndian, p.Id)
	}
	binary.Write(b, binary.LittleEndian, s.Enemy.Position.X)
	binary.Write(b, binary.LittleEndian, s.Enemy.Position.Y)
	binary.Write(b, binary.LittleEndian, s.Enemy.Position.Z)
	binary.Write(b, binary.LittleEndian, s.Enemy.Rotation)
	binary.Write(b, binary.LittleEndian, s.Enemy.AnimationFrame)
	binary.Write(b, binary.LittleEndian, uint8(len(s.Items)))
	for _, i := range s.Items {
		binary.Write(b, binary.LittleEndian, i.Position.X)
		binary.Write(b, binary.LittleEndian, i.Position.Y)
		binary.Write(b, binary.LittleEndian, i.Type)
		binary.Write(b, binary.LittleEndian, i.Id)
	}
	binary.Write(b, binary.LittleEndian, s.ItemHeld)

	return b.Bytes()
}

func DeserializeServerData(data []byte) ServerData {
	b := bytes.NewReader(data)
	var s ServerData

	var count uint8
	binary.Read(b, binary.LittleEndian, &count)

	binary.Read(b, binary.LittleEndian, &s.Position.X)
	binary.Read(b, binary.LittleEndian, &s.Position.Y)
	binary.Read(b, binary.LittleEndian, &s.Position.Z)
	binary.Read(b, binary.LittleEndian, &s.PlayerHp)

	s.Players = make([]PlayerData, count)
	for i := 0; i < int(count); i++ {
		binary.Read(b, binary.LittleEndian, &s.Players[i].Position.X)
		binary.Read(b, binary.LittleEndian, &s.Players[i].Position.Y)
		binary.Read(b, binary.LittleEndian, &s.Players[i].Position.Z)
		binary.Read(b, binary.LittleEndian, &s.Players[i].Rotation)
		binary.Read(b, binary.LittleEndian, &s.Players[i].Id)
	}
	binary.Read(b, binary.LittleEndian, &s.Enemy.Position.X)
	binary.Read(b, binary.LittleEndian, &s.Enemy.Position.Y)
	binary.Read(b, binary.LittleEndian, &s.Enemy.Position.Z)
	binary.Read(b, binary.LittleEndian, &s.Enemy.Rotation)
	binary.Read(b, binary.LittleEndian, &s.Enemy.AnimationFrame)
	var itemCount uint8
	binary.Read(b, binary.LittleEndian, &itemCount)
	s.Items = make([]ItemData, itemCount)
	for i := 0; i < int(itemCount); i++ {
		binary.Read(b, binary.LittleEndian, &s.Items[i].Position.X)
		binary.Read(b, binary.LittleEndian, &s.Items[i].Position.Y)
		binary.Read(b, binary.LittleEndian, &s.Items[i].Type)
		binary.Read(b, binary.LittleEndian, &s.Items[i].Id)
	}
	binary.Read(b, binary.LittleEndian, &s.ItemHeld)
	return s
}
