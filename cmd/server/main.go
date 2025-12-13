package main

import (
	"fmt"
	"net"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/PawelZabc/ProjektZespolowy/internal/game/entities"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/levels"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	"github.com/PawelZabc/ProjektZespolowy/internal/protocol"
	shared "github.com/PawelZabc/ProjektZespolowy/internal/shared"
)

func main() {
	//create connection
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
	//end creating connection

	clients := make(map[string]*entities.Player) //create player map

	//create objects
	objects := make([]colliders.Collider, 0, 100)
	floor := colliders.NewPlaneCollider(rl.NewVector3(-25, 0, -25), 50, 50, shared.DirY)
	objects = append(objects, floor)
	objects = append(objects, levels.CreateRoomWallsFromChanges(rl.NewVector3(-10, 0, -10), levels.BasicLevel, 3)...)
	object := colliders.NewCylinderCollider(rl.NewVector3(1, 1, 0), 0.5, 1)
	objects = append(objects, object)
	object2 := colliders.NewCubeCollider(rl.NewVector3(-3, 0, 6), 6, 1, 2)
	objects = append(objects, object2)
	ceiling := colliders.NewPlaneCollider(rl.NewVector3(-25, 3, -25), 50, 50, shared.DirYminus)
	objects = append(objects, ceiling)
	//end of create objects
	numberOFUpdates := int64(0)
	go func() {
		buffer := make([]byte, 1024)
		for {

			n, clientAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				fmt.Print("error")
				continue
			} //check if there is a new message, if not continue

			if _, ok := clients[clientAddr.String()]; !ok { //check if the address is new
				clients[clientAddr.String()] = &entities.Player{ //add new client to player map
					Velocity: rl.Vector3{},
					Collider: colliders.NewCylinderCollider(rl.NewVector3(0, 0, 0), 0.5 /*add to opts*/, 1 /*add to opts*/), //add to opts
					Speed:    0.1,                                                                                           //add to opts
					Address:  clientAddr,
				}
				fmt.Println("New client:", clientAddr)
			} else {
				player := clients[clientAddr.String()]     //get current player from address
				if player.LastMessage != numberOFUpdates { //check if there was already an update from the player
					player.LastMessage = numberOFUpdates
					var data protocol.ClientData = protocol.DeserializeClientData(buffer[:n]) // deserialize data
					player.RotationX = data.RotationX
					player.RotationY = data.RotationY
					for _, input := range data.Inputs { //decide what to do with the inputs
						switch input {
						case shared.MoveForward:
							player.Movement.Y = 1
						case shared.MoveBackward:
							player.Movement.Y = -1
						case shared.MoveLeft:
							player.Movement.X = 1
						case shared.MoveRight:
							player.Movement.X = -1
						case shared.Jump:
							if player.IsOnFloor {
								player.Velocity.Y = 0.1
							}
						}
					}
				}

			}
		}
	}()

	gravity := float32(0.005) //set gravity ,add to opts
	physicsUpdate := func() {
		numberOFUpdates++
		for _, player := range clients {
			if numberOFUpdates-player.LastMessage > 200 /*add to opts*/ { //if last message was 200 updates ago disconnect player
				fmt.Println("Client disconnected: ", player.Address.String())
				delete(clients, player.Address.String())
				continue
			}
			player.Velocity.Y -= gravity //aply gravity
			player.Move()
			player.IsOnFloor = false
			for _, obj := range objects { //collide with every object
				if obj != nil {
					direction := player.Collider.PushbackFrom(obj)
					if direction == shared.DirYminus {
						player.IsOnFloor = true
						player.Velocity.Y = 0
					} else if direction == shared.DirY {
						player.Velocity.Y = 0
					}
				}
			}
		}

	}

	go func() { //update 60 times a second
		ticker := time.NewTicker(time.Second / 60 /*add to opts*/)
		for range ticker.C {
			physicsUpdate()
		}

	}()

	players := make([]protocol.PlayerData, 0, 10 /*add to opts*/)
	ticker := time.NewTicker(time.Second / 30 /*add to opts*/)
	for range ticker.C { //send data 30 times a second

		for _, player := range clients {
			players = make([]protocol.PlayerData, 0, 10) //create a list with every player except itself
			for _, player2 := range clients {
				if player2.Address != player.Address {
					players = append(players, protocol.PlayerData{
						Position: player2.Collider.GetPosition(),
					})
				}
			}
			udpSend := protocol.ServerData{}
			udpSend.Position = player.GetPosition()
			udpSend.Players = players
			conn.WriteToUDP(protocol.SerializeServerData(udpSend), player.Address) //send data
		}
	}
}
