package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/PawelZabc/ProjektZespolowy/client/assets"
	"github.com/PawelZabc/ProjektZespolowy/client/config"
	entities "github.com/PawelZabc/ProjektZespolowy/client/entities"
	"github.com/PawelZabc/ProjektZespolowy/client/game"
	types "github.com/PawelZabc/ProjektZespolowy/shared/_types"
	s_entities "github.com/PawelZabc/ProjektZespolowy/shared/entities"
	udp_data "github.com/PawelZabc/ProjektZespolowy/shared/udp_data"
	math "github.com/chewxy/math32"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawTextOutlined(
	text string,
	x, y int32,
	fontSize int32,
	textColor, outlineColor rl.Color,
	thickness int32,
) {
	// Draw outline
	for ox := -thickness; ox <= thickness; ox++ {
		for oy := -thickness; oy <= thickness; oy++ {
			if ox == 0 && oy == 0 {
				continue
			}
			rl.DrawText(text, x+ox, y+oy, fontSize, outlineColor)
		}
	}

	// Draw main text
	rl.DrawText(text, x, y, fontSize, textColor)
}

func DrawHPBar(
	x, y int32,
	width, height int32,
	current, max int,
	backColor, borderColor rl.Color,
) {
	if max <= 0 {
		return
	}

	if current < 0 {
		current = 0
	}
	if current > max {
		current = max
	}

	pom := float32(current) / float32(max)

	var r, g uint8

	if pom >= 0.5 {
		pom = (1 - pom) * 2
		r = uint8(255 * pom)
		g = 255
	} else {
		r = 255
		pom = pom * 2
		g = uint8(255 * pom)
	}

	fillColor := rl.Color{R: r, G: g, B: 0, A: 255}
	percent := float32(current) / float32(max)
	fillWidth := int32(float32(width) * percent)

	// Border
	rl.DrawRectangle(x-2, y-2, width+4, height+4, borderColor)

	// Background
	rl.DrawRectangle(x, y, width, height, backColor)

	// Fill
	rl.DrawRectangle(x, y, fillWidth, height, fillColor)
}

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

	playerCollider := s_entities.NewCylinderCollider(rl.NewVector3(0, 0, 0), 0.5, 1)
	player := game.Object{Colliders: []types.Collider{playerCollider},
		Model: game.NewModelFromCollider(playerCollider),
	} //create player

	players := make(map[uint16]*entities.Actor)
	playerHp := 0
	createPlayer := func(Id uint16, Position rl.Vector3, Rotation float32) {
		cylinder := s_entities.NewCylinderCollider(Position, 0.5, 1) //if it doesnt create it
		players[Id] = entities.NewActor(cylinder, rl.Vector3{}, (Rotation*rl.Rad2deg)+90, assets.ModelPlayer)

	}
	cylinder := s_entities.NewCylinderCollider(rl.NewVector3(0, 0, 0), 0.5, 1)     //if it doesnt create it
	testPlayer := entities.NewActor(cylinder, rl.Vector3{}, 0, assets.ModelPlayer) //load player model
	fmt.Println(testPlayer)

	enemy := entities.NewActor(s_entities.NewCylinderCollider(rl.NewVector3(15, 0, 15), 1, 2), rl.Vector3{}, -45, assets.ModelGhost)

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
					player2Object.Object.Colliders[0].SetPosition(player2.Position) //if exists update position
					player2Object.Rotation = (player2.Rotation * rl.Rad2deg)
				} else {
					createPlayer(player2.Id, player2.Position, player2.Rotation)
				}
				updatedPlayers[player2.Id] = true //check player as updated

			}
			for id, _ := range players { //if a player wasnt updated,its no longer at the server so delete it
				if !updatedPlayers[id] {
					delete(players, id)
				}
			}
			player.Colliders[0].SetPosition(data.Position)
			playerHp = int(data.Hp)
			enemy.SetPosition(data.Enemy.Position)
			enemy.Rotation = -data.Enemy.Rotation
			enemy.SetAnimation(data.Enemy.Animation)
		}
	}()

	// objects := []*entities.Object{}
	//create objects

	shader := rl.LoadShader("assets/shaders/lighting_v2.vs", "assets/shaders/lighting_v2.fs")

	// Set view position uniform
	*shader.Locs = rl.GetShaderLocation(shader, "viewPos")

	ambientLoc := rl.GetShaderLocation(shader, "ambient")
	ambient := []float32{0.1, 0.1, 0.1, 1.0}
	rl.SetShaderValue(shader, ambientLoc, ambient, rl.ShaderUniformVec4)

	pointObject := game.Object{Model: game.NewModelFromCollider(s_entities.NewCubeCollider(rl.Vector3{}, 0.1, 0.1, 0.1)),
		Color: rl.Black,
	}
	side1Object := game.Object{Model: game.NewModelFromCollider(s_entities.NewCylinderCollider(rl.Vector3{}, 0.1, 0.2)),
		Color: rl.Black,
	}
	side2Object := game.Object{Model: game.NewModelFromCollider(s_entities.NewCylinderCollider(rl.Vector3{}, 0.1, 0.2)),
		Color: rl.Black,
	}
	//end of create objects

	conn.Write([]byte("hello")) //send hello to server to register address
	rl.HideCursor()
	centerx := rl.GetScreenWidth() / 2
	centery := rl.GetScreenHeight() / 2 //calculate center of the screen
	cameraRotationx := float32(-math.Pi / 2)
	cameraRotationy := float32(-math.Pi / 2) //setup camera rotation to look fowrward
	rl.SetMousePosition(centerx, centery)    //reset mouse to the middle of the screen

	rooms := game.LoadRooms(shader)
	fmt.Println(rooms)
	currentRoom := 0

	// white light
	light1 := entities.NewLight(
		entities.LightTypePoint,
		rl.NewVector3(0, 2.9, 0), // light position
		rl.NewVector3(0, 0, 0),   // target (unused for point light)
		rl.White,                 // light color
		shader,
	)
	light1.UpdateValues()

	// white light
	light2 := entities.NewLight(
		entities.LightTypePoint,
		rl.NewVector3(5, 2.9, -5), // light position
		rl.NewVector3(0, 0, 0),    // target (unused for point light)
		rl.White,                  // light color
		shader,
	)
	light2.UpdateValues()

	// white light
	light3 := entities.NewLight(
		entities.LightTypePoint,
		rl.NewVector3(5, 2.9, 5), // light position
		rl.NewVector3(0, 0, 0),   // target (unused for point light)
		rl.White,                 // light color
		shader,
	)
	light3.Enabled = 1 // ON is default
	light3.UpdateValues()

	lockMouse := false
	justClicked := false

	playerImg := rl.LoadTexture("assets/images/player.png")

	udpSend := udp_data.ClientData{}
	for !rl.WindowShouldClose() {
		if lockMouse {
			deltaMouse := rl.GetMousePosition() //check how much mouse has moved

			cameraRotationx += (deltaMouse.X - float32(centerx)) / 100 * config.CameraSensivity
			cameraRotationy -= (deltaMouse.Y - float32(centery)) / 100 * config.CameraSensivity //change camera rotation based on mouse movement
			if cameraRotationy > config.CameraLockMax {
				cameraRotationy = config.CameraLockMax
			} else if cameraRotationy < config.CameraLockMin {
				cameraRotationy = config.CameraLockMin
			}
			rl.SetMousePosition(centerx, centery)
		}
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
		if rl.IsKeyDown(rl.KeyR /*add to opts*/) && !justClicked {
			lockMouse = !lockMouse
			justClicked = true
			if lockMouse {
				rl.SetMousePosition(centerx, centery)
				rl.HideCursor()

			} else {
				rl.ShowCursor()

			}
		}
		if rl.IsKeyReleased(rl.KeyR) {
			justClicked = false
		}

		if rl.IsKeyDown(rl.KeySpace /*add to opts*/) {
			udpSend.Inputs = append(udpSend.Inputs, types.Jump)
		}

		target := rl.Vector3{X: float32(math.Sin(cameraRotationy) * math.Cos(cameraRotationx)),
			Z: float32(math.Sin(cameraRotationy) * math.Sin(cameraRotationx)),
			Y: float32(math.Cos(cameraRotationy))}
		target = rl.Vector3Normalize(target) //create a normal vector based on rotation

		camera.Position = rl.Vector3Add(player.Colliders[0].GetPosition(), rl.NewVector3(0, 0.5 /*ad to opts*/, 0)) //set camera to player position with height offset
		playerRay := s_entities.Ray{Origin: camera.Position, Direction: target}                                     //change player ray to have the same looking direction as the camera
		target = rl.Vector3Add(target, camera.Position)
		camera.Target = target //set camera target

		data := udp_data.SerializeClientData(udpSend) // send input and player data to the server
		_, err := conn.Write(data)
		if err != nil {
			fmt.Println("Write error:", err)
		}

		var pointPosition *rl.Vector3 = nil
		var minLength = float32(0)
		for _, object := range rooms[currentRoom].Objects { //check for nearest intersection point with the player ray
			if object != nil {
				for _, collider := range object.Colliders {
					point, length := playerRay.GetCollisionPoint(collider)
					if point != nil {
						if minLength == 0 || length < minLength {
							minLength = length
							pointPosition = point
						}
					}

				}

			}
		}
		cameraPos := []float32{
			camera.Position.X,
			camera.Position.Y,
			camera.Position.Z,
		}
		rl.SetShaderValue(shader, *shader.Locs, cameraPos, rl.ShaderUniformVec3)

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		// light spheres so it wouldnt shine out of nowhere
		rl.DrawSphereEx(light1.Position, 0.2, 8, 8, rl.White)
		rl.DrawSphereEx(light2.Position, 0.2, 8, 8, rl.White)
		rl.DrawSphereEx(light3.Position, 0.2, 8, 8, rl.White)
		if cylinder, ok := rooms[currentRoom].Objects[1].Colliders[0].(*s_entities.CylinderCollider); ok {
			pos1, pos2 := cylinder.GetSides(s_entities.GetVector2XZ(player.Colliders[0].GetPosition()))
			drawPoint1 := rl.Vector3Add(s_entities.GetVector3FromXZ(pos1), cylinder.Position)
			drawPoint1.Y += 0.5
			drawPoint2 := rl.Vector3Add(s_entities.GetVector3FromXZ(pos2), cylinder.Position)
			drawPoint2.Y += 0.5
			side1Object.DrawPoint = drawPoint1
			side1Object.Draw()
			side2Object.DrawPoint = drawPoint2
			side2Object.Draw()
		}

		if pointPosition != nil {
			pointObject.DrawPoint = rl.Vector3Add(*pointPosition, rl.NewVector3(-0.05, -0.05, -0.05))
			pointObject.Draw()
		} //draw the intersection point of player ray

		entities.DrawActorsMap(players) //draw players

		game.DrawRoom(&rooms[currentRoom]) //draw the room the player is currently in
		enemy.Draw()
		rl.EndMode3D()

		// draw UI

		rl.DrawTextureEx(
			playerImg,
			rl.Vector2{X: 10, Y: 490},
			0.0,  // rotation
			0.04, // scale
			rl.White,
		)

		DrawTextOutlined("G demo", 10, 10, 20, rl.Black, rl.LightGray, 2)

		DrawTextOutlined("Player hp:"+strconv.Itoa(playerHp), 160, 540, 20, rl.Black, rl.White, 2)

		DrawHPBar(160, 570, 150, 12, playerHp, 100, rl.Black, rl.White)

		rl.EndDrawing()
	}
}
