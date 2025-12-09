package main

import (
	"fmt"
	"net"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/PawelZabc/ProjektZespolowy/server/game"
	types "github.com/PawelZabc/ProjektZespolowy/shared/_types"
	s_entities "github.com/PawelZabc/ProjektZespolowy/shared/entities"
	udp_data "github.com/PawelZabc/ProjektZespolowy/shared/udp_data"
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

	clients := make(map[string]*s_entities.Player) //create player map

	rooms := game.LoadRooms() //load rooms
	objects := rooms[0].Colliders

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
				clients[clientAddr.String()] = &s_entities.Player{ //add new client to player map
					Velocity: rl.Vector3{},
					Collider: s_entities.NewCylinderCollider(rl.NewVector3(0, 0, 0), 0.5 /*add to opts*/, 1 /*add to opts*/), //add to opts
					Speed:    0.1,                                                                                            //add to opts
					Address:  clientAddr,
				}
				fmt.Println("New client:", clientAddr)
			} else {
				player := clients[clientAddr.String()]     //get current player from address
				if player.LastMessage != numberOFUpdates { //check if there was already an update from the player
					player.LastMessage = numberOFUpdates
					var data udp_data.ClientData = udp_data.DeserializeClientData(buffer[:n]) // deserialize data
					player.RotationX = data.RotationX
					player.RotationY = data.RotationY
					for _, input := range data.Inputs { //decide what to do with the inputs
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
					if direction == types.DirYminus {
						player.IsOnFloor = true
						player.Velocity.Y = 0
					} else if direction == types.DirY {
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

	players := make([]udp_data.PlayerData, 0, 10 /*add to opts*/)
	ticker := time.NewTicker(time.Second / 30 /*add to opts*/)
	for range ticker.C { //send data 30 times a second

		for _, player := range clients {
			players = make([]udp_data.PlayerData, 0, 10) //create a list with every player except itself
			for _, player2 := range clients {
				if player2.Address != player.Address {
					players = append(players, udp_data.PlayerData{
						Position: player2.Collider.GetPosition(),
					})
				}
			}
			udpSend := udp_data.ServerData{}
			udpSend.Position = player.GetPosition()
			udpSend.Players = players
			conn.WriteToUDP(udp_data.SerializeServerData(udpSend), player.Address) //send data
		}
	}
}
