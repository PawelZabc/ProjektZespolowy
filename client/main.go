package main

import (
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/PawelZabc/ProjektZespolowy/client/assets"
	"github.com/PawelZabc/ProjektZespolowy/client/config"
	entities "github.com/PawelZabc/ProjektZespolowy/client/entities"
	types "github.com/PawelZabc/ProjektZespolowy/shared/_types"
	s_entities "github.com/PawelZabc/ProjektZespolowy/shared/entities"
	leveldata "github.com/PawelZabc/ProjektZespolowy/shared/level_data"
	udp_data "github.com/PawelZabc/ProjektZespolowy/shared/udp_data"
	math "github.com/chewxy/math32"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func init() {
	assets.Init()
}

//go:generate go run ./utils/assetgen/main.go

func main() {
	//create connection
	serverIP := flag.String("ip", "127.0.0.1", "Server IP address")
	flag.Parse()

	println(*serverIP)
	serverAddr := net.UDPAddr{
		Port: 9000,
		IP:   net.ParseIP(*serverIP),
	}

	localAddr := net.UDPAddr{
		Port: 0,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.DialUDP("udp", &localAddr, &serverAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	//end of create connection

	rl.InitWindow(800, 600, "Client: Send Inputs, Receive Position") //setup window
	rl.SetTargetFPS(60 /*add to opts*/)
	defer rl.CloseWindow()
	camera := rl.Camera{ //setup camera
		Position:   rl.NewVector3(0, 4.0, 4.0),
		Target:     rl.NewVector3(0.0, 1.0, 0.0),
		Up:         rl.NewVector3(0.0, 1.0, 0.0),
		Fovy:       config.CameraFov,
		Projection: rl.CameraPerspective,
	}

	player := entities.CreateCylinderObject(rl.NewVector3(0, 0, 0), 0.5, 1) //create player
	players := make(map[uint16]*entities.Object)

	go func() { //go routine for receving messages
		buffer := make([]byte, 1024)
		for {
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, _, err := conn.ReadFromUDP(buffer) //check if  theres a new message
			if err != nil {
				continue
			}
			var data udp_data.ServerData = udp_data.DeserializeServerData(buffer[:n]) //deserialize data
			updatedPlayers := make(map[uint16]bool)                                   //create a map to check which players were sent
			for _, player2 := range data.Players {                                    //update players slice with received players
				if player2Object, ok := players[player2.Id]; ok { //check if player with that id existed before
					player2Object.Collider.SetPosition(player2.Position) //if exists update position
				} else {
					cylinder := entities.CreateCylinderObject(player2.Position, 0.5, 1) //if it doesnt create it
					players[player2.Id] = &cylinder
				}
				updatedPlayers[player2.Id] = true //check player as updated

			}
			for id, _ := range players { //if a player wasnt updated,its no longer at the server so delete it
				if !updatedPlayers[id] {
					delete(players, id)
				}
			}
			player.Collider.SetPosition(data.Position)
		}
	}()

	objects := []*entities.Object{}
	//create objects
	object := entities.CreateCylinderObject(rl.NewVector3(1, 1, 0), 0.5, 1)
	objects = append(objects, &object)
	object2 := entities.CreateCubeObject(rl.NewVector3(-3, 0, 6), 6, 1, 2)
	objects = append(objects, &object2)
	floor := entities.CreatePlaneObject(rl.NewVector3(-25, 0, -25), 50, 50, types.DirY)
	objects = append(objects, &floor)
	ceiling := entities.CreatePlaneObject(rl.NewVector3(-25, 3, -25), 50, 50, types.DirYminus)
	objects = append(objects, &ceiling)

	objects = append(objects, entities.CreateRoomWallsFromChanges(rl.NewVector3(-10, 0, -10), leveldata.Changes, 3)...)
	pointObject := entities.CreateCubeObject(rl.Vector3{}, 0.1, 0.1, 0.1)
	//end of create objects

	conn.Write([]byte("hello")) //send hello to server to register address
	rl.HideCursor()
	centerx := rl.GetScreenWidth() / 2
	centery := rl.GetScreenHeight() / 2 //calculate center of the screen
	cameraRotationx := float32(-math.Pi / 2)
	cameraRotationy := float32(-math.Pi / 2) //setup camera rotation to look fowrward
	rl.SetMousePosition(centerx, centery)    //reset mouse to the middle of the screen

	udpSend := udp_data.ClientData{}
	for !rl.WindowShouldClose() {
		deltaMouse := rl.GetMousePosition() //check how much mouse has moved

		cameraRotationx += (deltaMouse.X - float32(centerx)) / 100 * config.CameraSensivity
		cameraRotationy -= (deltaMouse.Y - float32(centery)) / 100 * config.CameraSensivity //change camera rotation based on mouse movement
		if cameraRotationy > config.CameraLockMax {
			cameraRotationy = config.CameraLockMax
		} else if cameraRotationy < config.CameraLockMin {
			cameraRotationy = config.CameraLockMin
		}
		rl.SetMousePosition(centerx, centery)
		udpSend = udp_data.ClientData{ //create object to send
			RotationX: cameraRotationx,
			RotationY: cameraRotationy,
			Inputs:    make([]types.PlayerAction, 0, 5),
		}
		if rl.IsKeyDown(rl.KeyW /*add to opts*/) { //add player actions to udpSend based on inputs
			udpSend.Inputs = append(udpSend.Inputs, types.MoveForward)
		}
		if rl.IsKeyDown(rl.KeyS /*add to opts*/) {
			udpSend.Inputs = append(udpSend.Inputs, types.MoveBackward)
		}
		if rl.IsKeyDown(rl.KeyA /*add to opts*/) {
			udpSend.Inputs = append(udpSend.Inputs, types.MoveLeft)
		}
		if rl.IsKeyDown(rl.KeyD /*add to opts*/) {
			udpSend.Inputs = append(udpSend.Inputs, types.MoveRight)
		}

		if rl.IsKeyDown(rl.KeySpace /*add to opts*/) {
			udpSend.Inputs = append(udpSend.Inputs, types.Jump)
		}

		target := rl.Vector3{X: float32(math.Sin(cameraRotationy) * math.Cos(cameraRotationx)),
			Z: float32(math.Sin(cameraRotationy) * math.Sin(cameraRotationx)),
			Y: float32(math.Cos(cameraRotationy))}
		target = rl.Vector3Normalize(target) //create a normal vector based on rotation

		camera.Position = rl.Vector3Add(player.Collider.GetPosition(), rl.NewVector3(0, 0.5 /*ad to opts*/, 0)) //set camera to player position with height offset
		playerRay := s_entities.Ray{Origin: camera.Position, Direction: target}                                 //change player ray to have the same looking direction as the camera
		target = rl.Vector3Add(target, camera.Position)
		camera.Target = target //set camera target

		data := udp_data.SerializeClientData(udpSend) // send input and player data to the server
		_, err := conn.Write(data)
		if err != nil {
			fmt.Println("Write error:", err)
		}

		var pointPosition *rl.Vector3 = nil
		var minLength = float32(0)
		for _, obj := range objects { //check for nearest intersection point with the player ray
			if obj != nil {
				point, length := playerRay.GetCollisionPoint(obj.Collider)
				if point != nil {
					if minLength == 0 || length < minLength {
						minLength = length
						pointPosition = point
					}
				}
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		if pointPosition != nil {
			rl.DrawModel(pointObject.Model, rl.Vector3Add(*pointPosition, rl.NewVector3(-0.05, -0.05, -0.05)), 1.0, rl.Black)
		} //draw the intersection point of player ray

		for _, obj := range objects { //draw every object
			if obj != nil {
				if plane, ok := obj.Collider.(*s_entities.PlaneCollider); ok {
					switch plane.Direction { //check which color to draw the plane as
					case types.DirX:
						{
							rl.DrawModel(obj.Model, obj.Collider.GetPosition(), 1.0, rl.Red)
						}
					case types.DirY:
						{
							rl.DrawModel(obj.Model, obj.Collider.GetPosition(), 1.0, rl.Orange)
						}
					case types.DirYminus:
						{
							rl.DrawModel(obj.Model, obj.Collider.GetPosition(), 1.0, rl.Green)
						}
					case types.DirZ:
						{
							rl.DrawModel(obj.Model, obj.Collider.GetPosition(), 1.0, rl.Yellow)
						}
					}
				} else { //if not plane color white
					rl.DrawModel(obj.Model, obj.Collider.GetPosition(), 1.0, rl.White)
				}

			}

		}

		for _, obj := range players { //draw players
			if obj != nil {
				rl.DrawModel(obj.Model, obj.Collider.GetPosition(), 1.0, rl.White)
			}
		}

		rl.EndMode3D()
		rl.DrawText("Collision demo", 10, 10, 20, rl.Black)
		rl.EndDrawing()
	}
}
