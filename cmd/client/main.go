package main

import (
	"log"

	"github.com/PawelZabc/ProjektZespolowy/assets"
	"github.com/PawelZabc/ProjektZespolowy/internal/client"
	"github.com/PawelZabc/ProjektZespolowy/internal/config"
)

//go:generate go run ../../pkg/assets_name_gen/main.go

func init() {
	assets.Init()
}

func main() {
	cfg := config.DefaultClientConfig() // default is good for now, but it can be overwritten below

	app := client.NewApp(cfg)
	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

// for now i will leave this code below
// package main

// import (
// 	"flag"
// 	"fmt"
// 	"net"
// 	"time"

// 	"github.com/PawelZabc/ProjektZespolowy/assets"
// 	"github.com/PawelZabc/ProjektZespolowy/internal/config"
// 	"github.com/PawelZabc/ProjektZespolowy/internal/game/entities"
// 	"github.com/PawelZabc/ProjektZespolowy/internal/game/levels"
// 	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics"
// 	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
// 	"github.com/PawelZabc/ProjektZespolowy/internal/protocol"
// 	"github.com/PawelZabc/ProjektZespolowy/internal/shared"
// 	"github.com/chewxy/math32"

// 	rl "github.com/gen2brain/raylib-go/raylib"
// )

// func init() {
// 	assets.Init()
// }

// //go:generate go run ./utils/assetgen/main.go

// func main() {
// 	//create connection
// 	serverIP := flag.String("ip", "127.0.0.1", "Server IP address")
// 	flag.Parse()

// 	println(*serverIP)
// 	serverAddr := net.UDPAddr{
// 		Port: 9000,
// 		IP:   net.ParseIP(*serverIP),
// 	}

// 	localAddr := net.UDPAddr{
// 		Port: 0,
// 		IP:   net.ParseIP("0.0.0.0"),
// 	}

// 	conn, err := net.DialUDP("udp", &localAddr, &serverAddr)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer conn.Close()
// 	//end of create connection

// 	rl.InitWindow(800, 600, "Client: Send Inputs, Receive Position") //setup window
// 	rl.SetTargetFPS(60 /*add to opts*/)
// 	defer rl.CloseWindow()
// 	camera := rl.Camera{ //setup camera
// 		Position:   rl.NewVector3(0, 4.0, 4.0),
// 		Target:     rl.NewVector3(0.0, 1.0, 0.0),
// 		Up:         rl.NewVector3(0.0, 1.0, 0.0),
// 		Fovy:       config.CameraFov,
// 		Projection: rl.CameraPerspective,
// 	}

// 	playerCollider := colliders.NewCylinderCollider(rl.NewVector3(0, 0, 0), 0.5, 1)
// 	player := entities.Object{Colliders: []colliders.Collider{playerCollider}, // TODO: fix object
// 		Model: levels.NewModelFromCollider(playerCollider), // TODO: fix object
// 	} //create player

// 	players := make(map[uint16]*entities.Actor)
// 	createPlayer := func(Id uint16, Position rl.Vector3, Rotation float32) {
// 		cylinder := colliders.NewCylinderCollider(Position, 0.5, 1) //if it doesnt create it
// 		players[Id] = entities.NewActor(cylinder, rl.Vector3{}, (Rotation*rl.Rad2deg)+90, assets.ModelPlayer)

// 	}
// 	cylinder := colliders.NewCylinderCollider(rl.NewVector3(0, 0, 0), 0.5, 1)      //if it doesnt create it
// 	testPlayer := entities.NewActor(cylinder, rl.Vector3{}, 0, assets.ModelPlayer) //load player model
// 	fmt.Println(testPlayer)

// 	enemy := entities.NewActor(colliders.NewCylinderCollider(rl.NewVector3(15, 0, 15), 1, 2), rl.Vector3{}, -45, assets.ModelGhost)

// 	go func() { //go routine for receving messages
// 		buffer := make([]byte, 1024)
// 		for {
// 			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
// 			n, _, err := conn.ReadFromUDP(buffer) //check if  theres a new message
// 			if err != nil {
// 				continue
// 			}
// 			var data protocol.ServerData = protocol.DeserializeServerData(buffer[:n]) //deserialize data
// 			updatedPlayers := make(map[uint16]bool)                                   //create a map to check which players were sent
// 			for _, player2 := range data.Players {                                    //update players slice with received players
// 				if player2Object, ok := players[player2.Id]; ok { //check if player with that id existed before
// 					player2Object.Object.Colliders[0].SetPosition(player2.Position) //if exists update position
// 					player2Object.Rotation = (player2.Rotation * rl.Rad2deg)
// 				} else {
// 					createPlayer(player2.Id, player2.Position, player2.Rotation)
// 				}
// 				updatedPlayers[player2.Id] = true //check player as updated

// 			}
// 			for id, _ := range players { //if a player wasnt updated,its no longer at the server so delete it
// 				if !updatedPlayers[id] {
// 					delete(players, id)
// 				}
// 			}
// 			player.Colliders[0].SetPosition(data.Position)
// 			enemy.SetPosition(data.Enemy.Position)
// 			enemy.Rotation = -data.Enemy.Rotation
// 		}
// 	}()

// 	// objects := []*entities.Object{}
// 	//create objects

// 	// TODO: fix Object
// 	pointObject := entities.Object{Model: levels.NewModelFromCollider(colliders.NewCubeCollider(rl.Vector3{}, 0.1, 0.1, 0.1)),
// 		Color: rl.Black,
// 	}
// 	side1Object := entities.Object{Model: levels.NewModelFromCollider(colliders.NewCylinderCollider(rl.Vector3{}, 0.1, 0.2)),
// 		Color: rl.Black,
// 	}
// 	side2Object := entities.Object{Model: levels.NewModelFromCollider(colliders.NewCylinderCollider(rl.Vector3{}, 0.1, 0.2)),
// 		Color: rl.Black,
// 	}
// 	//end of create objects

// 	conn.Write([]byte("hello")) //send hello to server to register address
// 	rl.HideCursor()
// 	centerx := rl.GetScreenWidth() / 2
// 	centery := rl.GetScreenHeight() / 2 //calculate center of the screen
// 	cameraRotationx := float32(-math32.Pi / 2)
// 	cameraRotationy := float32(-math32.Pi / 2) //setup camera rotation to look fowrward
// 	rl.SetMousePosition(centerx, centery)      //reset mouse to the middle of the screen
// 	rooms := levels.LoadRooms()                // TODO: fix object

// 	currentRoom := 0
// 	lockMouse := false
// 	justClicked := false
// 	udpSend := protocol.ClientData{}
// 	for !rl.WindowShouldClose() {
// 		if lockMouse {
// 			deltaMouse := rl.GetMousePosition() //check how much mouse has moved

// 			cameraRotationx += (deltaMouse.X - float32(centerx)) / 100 * config.CameraSensivity
// 			cameraRotationy -= (deltaMouse.Y - float32(centery)) / 100 * config.CameraSensivity //change camera rotation based on mouse movement
// 			if cameraRotationy > config.CameraLockMax {
// 				cameraRotationy = config.CameraLockMax
// 			} else if cameraRotationy < config.CameraLockMin {
// 				cameraRotationy = config.CameraLockMin
// 			}
// 			rl.SetMousePosition(centerx, centery)
// 		}
// 		udpSend = protocol.ClientData{ //create object to send
// 			RotationX: cameraRotationx,
// 			RotationY: cameraRotationy,
// 			Inputs:    make([]shared.PlayerAction, 0, 5),
// 		}
// 		if rl.IsKeyDown(rl.KeyW /*add to opts*/) { //add player actions to udpSend based on inputs
// 			udpSend.Inputs = append(udpSend.Inputs, shared.MoveForward)
// 		}
// 		if rl.IsKeyDown(rl.KeyS /*add to opts*/) {
// 			udpSend.Inputs = append(udpSend.Inputs, shared.MoveBackward)
// 		}
// 		if rl.IsKeyDown(rl.KeyA /*add to opts*/) {
// 			udpSend.Inputs = append(udpSend.Inputs, shared.MoveLeft)
// 		}
// 		if rl.IsKeyDown(rl.KeyD /*add to opts*/) {
// 			udpSend.Inputs = append(udpSend.Inputs, shared.MoveRight)
// 		}
// 		if rl.IsKeyDown(rl.KeyR /*add to opts*/) && !justClicked {
// 			lockMouse = !lockMouse
// 			justClicked = true
// 			if lockMouse {
// 				rl.SetMousePosition(centerx, centery)
// 				rl.HideCursor()

// 			} else {
// 				rl.ShowCursor()

// 			}
// 		}
// 		if rl.IsKeyReleased(rl.KeyR) {
// 			justClicked = false
// 		}

// 		if rl.IsKeyDown(rl.KeySpace /*add to opts*/) {
// 			udpSend.Inputs = append(udpSend.Inputs, shared.Jump)
// 		}

// 		target := rl.Vector3{X: float32(math32.Sin(cameraRotationy) * math32.Cos(cameraRotationx)),
// 			Z: float32(math32.Sin(cameraRotationy) * math32.Sin(cameraRotationx)),
// 			Y: float32(math32.Cos(cameraRotationy))}
// 		target = rl.Vector3Normalize(target) //create a normal vector based on rotation

// 		camera.Position = rl.Vector3Add(player.Colliders[0].GetPosition(), rl.NewVector3(0, 0.5 /*ad to opts*/, 0)) //set camera to player position with height offset
// 		playerRay := colliders.Ray{Origin: camera.Position, Direction: target}                                      //change player ray to have the same looking direction as the camera
// 		target = rl.Vector3Add(target, camera.Position)
// 		camera.Target = target //set camera target

// 		data := protocol.SerializeClientData(udpSend) // send input and player data to the server
// 		_, err := conn.Write(data)
// 		if err != nil {
// 			fmt.Println("Write error:", err)
// 		}

// 		var pointPosition *rl.Vector3 = nil
// 		var minLength = float32(0)
// 		for _, object := range rooms[currentRoom].Objects { //check for nearest intersection point with the player ray
// 			if object != nil {
// 				for _, collider := range object.Colliders {
// 					point, length := playerRay.GetCollisionPoint(collider)
// 					if point != nil {
// 						if minLength == 0 || length < minLength {
// 							minLength = length
// 							pointPosition = point
// 						}
// 					}

// 				}

// 			}
// 		}

// 		rl.BeginDrawing()
// 		rl.ClearBackground(rl.RayWhite)

// 		rl.BeginMode3D(camera)

// 		if cylinder, ok := rooms[currentRoom].Objects[1].Colliders[0].(*colliders.CylinderCollider); ok {
// 			pos1, pos2 := cylinder.GetSides(physics.GetVector2XZ(player.Colliders[0].GetPosition()))
// 			drawPoint1 := rl.Vector3Add(physics.GetVector3FromXZ(pos1), cylinder.Position)
// 			drawPoint1.Y += 0.5
// 			drawPoint2 := rl.Vector3Add(physics.GetVector3FromXZ(pos2), cylinder.Position)
// 			drawPoint2.Y += 0.5
// 			side1Object.DrawPoint = drawPoint1
// 			side1Object.Draw()
// 			side2Object.DrawPoint = drawPoint2
// 			side2Object.Draw()
// 		}

// 		if pointPosition != nil {
// 			pointObject.DrawPoint = rl.Vector3Add(*pointPosition, rl.NewVector3(-0.05, -0.05, -0.05))
// 			pointObject.Draw()
// 		} //draw the intersection point of player ray

// 		entities.DrawActorsMap(players)      //draw players
// 		levels.DrawRoom(&rooms[currentRoom]) //draw the room the player is currently in
// 		enemy.Draw()
// 		rl.EndMode3D()
// 		rl.DrawText("Collision demo", 10, 10, 20, rl.Black)
// 		rl.EndDrawing()
// 	}
// }
