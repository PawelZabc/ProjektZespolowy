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

	rl.InitWindow(800, 600, "Client: Send Inputs, Receive Position")
	rl.SetTargetFPS(60)
	defer rl.CloseWindow()
	camera := rl.Camera{
		Position:   rl.NewVector3(0, 4.0, 4.0),
		Target:     rl.NewVector3(0.0, 1.0, 0.0),
		Up:         rl.NewVector3(0.0, 1.0, 0.0),
		Fovy:       config.CameraFov,
		Projection: rl.CameraPerspective,
	}

	player := entities.CreateCylinderObject(rl.NewVector3(0, 0, 0), 0.5, 1)

	go func() {
		buffer := make([]byte, 1024)
		for {
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, _, err := conn.ReadFromUDP(buffer)
			if err != nil {
				continue
			}
			var data udp_data.ServerData = udp_data.DeserializeServerData(buffer[:n])
			player.Collider.SetPosition(data.Position)
		}
	}()

	objects := []*entities.Object{}

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

	conn.Write([]byte("hello"))
	rl.HideCursor()
	centerx := rl.GetScreenWidth() / 2
	centery := rl.GetScreenHeight() / 2
	cameraRotationx := float32(-math.Pi / 2)
	cameraRotationy := float32(-math.Pi / 2)
	rl.SetMousePosition(centerx, centery)

	udpSend := udp_data.ClientData{}
	for !rl.WindowShouldClose() {
		deltaMouse := rl.GetMousePosition()

		cameraRotationx += (deltaMouse.X - float32(centerx)) / 100 * config.CameraSensivity
		cameraRotationy -= (deltaMouse.Y - float32(centery)) / 100 * config.CameraSensivity
		if cameraRotationy > config.CameraLockMax {
			cameraRotationy = config.CameraLockMax
		} else if cameraRotationy < config.CameraLockMin {
			cameraRotationy = config.CameraLockMin
		}
		udpSend = udp_data.ClientData{
			RotationX: cameraRotationx,
			RotationY: cameraRotationy,
			Inputs:    make([]types.PlayerAction, 0, 5),
		}
		if rl.IsKeyDown(rl.KeyW) {
			udpSend.Inputs = append(udpSend.Inputs, types.MoveForward)
		}
		if rl.IsKeyDown(rl.KeyS) {
			udpSend.Inputs = append(udpSend.Inputs, types.MoveBackward)
		}
		if rl.IsKeyDown(rl.KeyA) {
			udpSend.Inputs = append(udpSend.Inputs, types.MoveLeft)
		}
		if rl.IsKeyDown(rl.KeyD) {
			udpSend.Inputs = append(udpSend.Inputs, types.MoveRight)
		}

		if rl.IsKeyDown(rl.KeySpace) {
			udpSend.Inputs = append(udpSend.Inputs, types.Jump)
		}

		target := rl.Vector3{X: float32(math.Sin(cameraRotationy) * math.Cos(cameraRotationx)),
			Z: float32(math.Sin(cameraRotationy) * math.Sin(cameraRotationx)),
			Y: float32(math.Cos(cameraRotationy))}
		target = rl.Vector3Normalize(target)

		camera.Position = rl.Vector3Add(player.Collider.GetPosition(), rl.NewVector3(0, 0.5, 0))
		playerRay := s_entities.Ray{Origin: camera.Position, Direction: target}
		target = rl.Vector3Add(target, camera.Position)
		camera.Target = target
		rl.SetMousePosition(centerx, centery)

		data := udp_data.SerializeClientData(udpSend)
		_, err := conn.Write(data)
		if err != nil {
			fmt.Println("Write error:", err)
		}

		var pointPosition *rl.Vector3 = nil
		var minLength = float32(0)
		for _, obj := range objects {
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
		}

		for _, obj := range objects {
			if obj != nil {
				if plane, ok := obj.Collider.(*s_entities.PlaneCollider); ok {
					switch plane.Direction {
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
				} else {
					rl.DrawModel(obj.Model, obj.Collider.GetPosition(), 1.0, rl.White)
				}

			}

		}

		rl.EndMode3D()
		rl.DrawText("Collision demo", 10, 10, 20, rl.Black)
		rl.EndDrawing()
	}
}
