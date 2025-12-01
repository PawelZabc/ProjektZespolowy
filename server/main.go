package main

import (
	"fmt"
	"net"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"

	types "github.com/PawelZabc/ProjektZespolowy/shared/_types"
	s_entities "github.com/PawelZabc/ProjektZespolowy/shared/entities"
	leveldata "github.com/PawelZabc/ProjektZespolowy/shared/level_data"
	udp_data "github.com/PawelZabc/ProjektZespolowy/shared/udp_data"
)

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

	clients := make(map[string]*s_entities.Player)
	// player := s_entities.Player{
	// 	Velocity: rl.Vector3{},
	// 	Collider: s_entities.NewCylinderCollider(rl.NewVector3(0, 0, 0), 0.5, 1),
	// 	Speed:    0.1,
	// }

	objects := make([]types.Collider, 0, 100)
	floor := s_entities.NewPlaneCollider(rl.NewVector3(-25, 0, -25), 50, 50, types.DirY)
	objects = append(objects, floor)
	objects = append(objects, s_entities.CreateRoomWallsFromChanges(rl.NewVector3(-10, 0, -10), leveldata.Changes, 3)...)
	object := s_entities.NewCylinderCollider(rl.NewVector3(1, 1, 0), 0.5, 1)
	objects = append(objects, object)
	object2 := s_entities.NewCubeCollider(rl.NewVector3(-3, 0, 6), 6, 1, 2)
	objects = append(objects, object2)
	ceiling := s_entities.NewPlaneCollider(rl.NewVector3(-25, 3, -25), 50, 50, types.DirYminus)
	objects = append(objects, ceiling)
	go func() {
		buffer := make([]byte, 1024)
		for {

			n, clientAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				fmt.Print("error")
				continue
			}

			if _, ok := clients[clientAddr.String()]; !ok {
				clients[clientAddr.String()] = &s_entities.Player{
					Velocity: rl.Vector3{},
					Collider: s_entities.NewCylinderCollider(rl.NewVector3(0, 0, 0), 0.5, 1),
					Speed:    0.1,
					Address:  clientAddr,
				}
				fmt.Println("New client:", clientAddr)
			} else {
				player := clients[clientAddr.String()]
				var data udp_data.ClientData = udp_data.DeserializeClientData(buffer[:n])
				player.RotationX = data.RotationX
				player.RotationY = data.RotationY
				for _, input := range data.Inputs {
					switch input {
					case types.MoveForward:
						player.Movement.Y = 1
					case types.MoveBackward:
						player.Movement.Y = -1
					case types.MoveLeft:
						player.Movement.X = 1
					case types.MoveRight:
						player.Movement.X = -1
					case types.Jump:
						if player.IsOnFloor {
							player.Velocity.Y = 0.1
						}
					}
				}
			}
		}
	}()

	gravity := float32(0.005)
	physicsUpdate := func() {
		for _, player := range clients {
			player.Velocity.Y -= gravity
			player.Move()
			player.IsOnFloor = false
			for _, obj := range objects {
				if obj != nil {
					direction := player.Collider.PushbackFrom(obj)
					if direction == types.DirYminus {
						player.IsOnFloor = true
						player.Velocity.Y = 0
					} else if direction == types.DirY {
						player.Velocity.Y = 0
					}
				}
			}
			// fmt.Println(player.GetPosition())
		}

	}

	go func() {
		ticker := time.NewTicker(time.Second / 60)
		for range ticker.C {
			physicsUpdate()
		}

	}()

	ticker := time.NewTicker(time.Second / 30)
	for range ticker.C {

		for _, player := range clients {
			udpSend := udp_data.ServerData{}
			udpSend.Position = player.GetPosition()
			fmt.Println(player.Address, player.GetPosition())
			conn.WriteToUDP(udp_data.SerializeServerData(udpSend), player.Address)
		}
	}
}
