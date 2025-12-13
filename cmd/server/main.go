package main

import (
	"fmt"
	"net"
	"time"

	"github.com/PawelZabc/ProjektZespolowy/internal/game/entities"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/levels"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	"github.com/PawelZabc/ProjektZespolowy/internal/protocol"
	"github.com/PawelZabc/ProjektZespolowy/internal/shared"
	rl "github.com/gen2brain/raylib-go/raylib"
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
	enemy := &entities.Enemy{                    //create the enemy
		Collider: colliders.NewCylinderCollider(rl.NewVector3(20, 0, 15), 1 /*add to opts*/, 2 /*add to opts*/), //add to opts
		Speed:    0.05,                                                                                          //add to opts
	}

	rooms := levels.ServerLoadRooms() //load rooms // TODO: Fix object
	objects := rooms[0].Colliders

	numberOFUpdates := int64(0)
	newPlayerId := uint16(0) //create counter for new player id
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
					Id:       newPlayerId,
				}
				newPlayerId++
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
		players := make([]*entities.Player, 0, len(clients))
		for _, player := range clients {
			if numberOFUpdates-player.LastMessage > 200 /*add to opts*/ { //if last message was 200 updates ago disconnect player
				fmt.Println("Client disconnected: ", player.Address.String())
				delete(clients, player.Address.String())
				continue
			}
			players = append(players, player)
			player.Velocity.Y -= gravity //apply gravity
			player.Move()
			player.IsOnFloor = false
			for _, obj := range objects { //collide with every object
				player.PushbackFrom(obj)
			}
			player.PushbackFrom(enemy.Collider)
		}
		enemy.UpdateTarget(players, &objects)
		enemy.Move()
		for _, obj := range objects { //collide with every object
			if obj != nil {
				enemy.Collider.PushbackFrom(obj)
			}
		}
		for _, player := range clients {
			enemy.Collider.PushbackFrom(player.Collider)
		}

	}

	updateFrequency := float64(60)      /*add to opts*/
	sendUpdatesFrequency := float64(30) /*add to opts*/
	lastSend := 0
	ratio := sendUpdatesFrequency / updateFrequency
	players := make([]protocol.PlayerData, 0, 10 /*add to opts*/) //max players

	ticker := time.NewTicker(time.Second / time.Duration(updateFrequency))
	for range ticker.C { //update 60 times a second
		physicsUpdate()

		if (ratio*float64(numberOFUpdates))-float64(lastSend) >= 1 { //send every 60 seconds
			lastSend++
			if numberOFUpdates%2 == 0 {
				for _, player := range clients {
					players = make([]protocol.PlayerData, 0, 10) //create a list with every player except itself
					for _, player2 := range clients {
						if player2.Address != player.Address {
							players = append(players, protocol.PlayerData{
								Position: player2.Collider.GetPosition(),
								Rotation: player2.RotationX,
								Id:       player2.Id,
							})
						}
					}
					udpSend := protocol.ServerData{}
					udpSend.Position = player.GetPosition()
					udpSend.Players = players
					udpSend.Enemy = protocol.EnemyData{Position: enemy.Collider.GetPosition(), Rotation: enemy.RotationX}
					conn.WriteToUDP(protocol.SerializeServerData(udpSend), player.Address) //send data
				}
			}

		}

	}

}
