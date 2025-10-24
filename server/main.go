package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
)

type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

type Data struct {
	X     float32 `json:"x"`
	Y     float32 `json:"y"`
	Z     float32 `json:"z"`
	Frame int32   `json:"frame"`
}

func main() {
	addr := net.UDPAddr{
		Port: 9000,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Server listening on port 9000...")

	clients := make(map[string]*net.UDPAddr)

	pos := Position{X: 0, Y: 0, Z: 0}
	speed := float32(0.1)
	frame := int32(0)

	go func() {
		buffer := make([]byte, 1024)
		for {
			n, clientAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				continue
			}

			if _, ok := clients[clientAddr.String()]; !ok {
				clients[clientAddr.String()] = clientAddr
				fmt.Println("New client:", clientAddr)
			}

			input := strings.ToUpper(string(buffer[:n]))
			fmt.Printf("Received from %v: %s\n", clientAddr, input)
			for _, char := range input {
				switch char {
				case 'W':
					pos.Z -= speed
				case 'S':
					pos.Z += speed
				case 'A':
					pos.X -= speed
				case 'D':
					pos.X += speed
				}
			}

		}
	}()

	ticker := time.NewTicker(33 * time.Millisecond)
	for range ticker.C {
		frame += 1
		data, _ := json.Marshal(Data{X: pos.X, Y: pos.Y, Z: pos.Z, Frame: frame})
		for _, c := range clients {
			conn.WriteToUDP(data, c)
		}
	}
}
