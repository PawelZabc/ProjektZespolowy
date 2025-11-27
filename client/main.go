package main

import (
	types "github.com/PawelZabc/ProjektZespolowy/client/_types"
	"github.com/PawelZabc/ProjektZespolowy/client/assets"
	"github.com/PawelZabc/ProjektZespolowy/client/config"
	"github.com/PawelZabc/ProjektZespolowy/client/entities"
	math "github.com/chewxy/math32"

	rl "github.com/gen2brain/raylib-go/raylib"
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

func init() {
	assets.Init()
}

//go:generate go run ./utils/assetgen/main.go

func main() {
	// serverIP := flag.String("ip", "127.0.0.1", "Server IP address")
	// flag.Parse()

	// println(*serverIP)
	// serverAddr := net.UDPAddr{
	// 	Port: 9000,
	// 	IP:   net.ParseIP(*serverIP),
	// }

	// localAddr := net.UDPAddr{
	// 	Port: 0,
	// 	IP:   net.ParseIP("0.0.0.0"),
	// }

	// conn, err := net.DialUDP("udp", &localAddr, &serverAddr)
	// if err != nil {
	// 	panic(err)
	// }
	// defer conn.Close()

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

	// go func() {
	// 	buffer := make([]byte, 1024)
	// 	for {
	// 		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	// 		n, _, err := conn.ReadFromUDP(buffer)
	// 		if err != nil {
	// 			continue
	// 		}
	// 		// var pos Position
	// 		var data Data
	// 		err = json.Unmarshal(buffer[:n], &data)
	// 		// if len(animations) > 0 && animations[0].FrameCount > 0 && model.BoneCount > 0 {
	// 		// 	rl.UpdateModelAnimation(model, animations[0], data.Frame%animations[0].FrameCount)
	// 		// }
	// 		// println(data.Frame % animations[0].FrameCount)
	// 		if err == nil {
	// 			position = rl.NewVector3(data.X, data.Y, data.Z)
	// 		}
	// 	}
	// }()

	objects := []*entities.Object{}

	player := entities.CreateCylinderObject(rl.NewVector3(0, 0, 0), 0.5, 1)

	object := entities.CreateCylinderObject(rl.NewVector3(1, 1, 0), 0.5, 1)
	objects = append(objects, &object)
	object2 := entities.CreateCubeObject(rl.NewVector3(-3, 0, 6), 6, 1, 2)
	objects = append(objects, &object2)
	floor := entities.CreatePlaneObject(rl.NewVector3(-25, 0, -25), 50, 50, types.DirY)
	objects = append(objects, &floor)
	ceiling := entities.CreatePlaneObject(rl.NewVector3(-25, 3, -25), 50, 50, types.DirYminus)
	objects = append(objects, &ceiling)
	// points := []rl.Vector2{
	// 	rl.NewVector2(-10, -10),
	// 	rl.NewVector2(-10, 10),
	// 	rl.NewVector2(10, 10),
	// 	rl.NewVector2(10, -10),
	// 	rl.NewVector2(-10, -10)}
	// objects = append(objects, entities.CreateRoomWallsFromPoints(points, 0, 3)...)
	changes := []entities.Change{
		{Value: 20, Axis: types.DirX},
		{Value: 5, Axis: types.DirZ},
		{Value: 5, Axis: types.DirX},
		{Value: -5, Axis: types.DirZ},
		{Value: 10, Axis: types.DirX},
		{Value: 30, Axis: types.DirZ},
		{Value: -10, Axis: types.DirX},
		{Value: -20, Axis: types.DirZ},
		{Value: -5, Axis: types.DirX},
		{Value: 10, Axis: types.DirZ},
		{Value: -20, Axis: types.DirX},
		{Value: -10, Axis: types.DirZ},
		{Value: -5, Axis: types.DirX},
		{Value: 20, Axis: types.DirZ},
		{Value: -10, Axis: types.DirX},
		{Value: -30, Axis: types.DirZ},
		{Value: 10, Axis: types.DirX},
		{Value: 5, Axis: types.DirZ},
		{Value: 5, Axis: types.DirX},
		{Value: -5, Axis: types.DirZ},
	}
	objects = append(objects, entities.CreateRoomWallsFromChanges(rl.NewVector3(-10, 0, -10), changes, 3)...)
	pointObject := entities.CreateCubeObject(rl.Vector3{}, 0.1, 0.1, 0.1)

	// conn.Write([]byte("hello"))
	rl.HideCursor()
	centerx := rl.GetScreenWidth() / 2
	centery := rl.GetScreenHeight() / 2
	cameraRotationx := float32(-math.Pi / 2)
	cameraRotationy := float32(-math.Pi / 2)
	gravity := float32(0.005)
	isOnFloor := false
	velocity := rl.Vector3{}
	rl.SetMousePosition(centerx, centery)

	// moving := &player
	// waspressed := false
	for !rl.WindowShouldClose() {
		deltaMouse := rl.GetMousePosition()

		cameraRotationx += (deltaMouse.X - float32(centerx)) / 100 * config.CameraSensivity
		cameraRotationy -= (deltaMouse.Y - float32(centery)) / 100 * config.CameraSensivity
		if cameraRotationy > config.CameraLockMax {
			cameraRotationy = config.CameraLockMax
		} else if cameraRotationy < config.CameraLockMin {
			cameraRotationy = config.CameraLockMin
		}
		// println(cameraRotationy)

		// input := ""
		velocity.X = 0
		velocity.Z = 0
		movement := rl.Vector2{}
		if rl.IsKeyDown(rl.KeyW) {
			movement = rl.Vector2Add(movement, rl.NewVector2(0.1, 0))
		}
		if rl.IsKeyDown(rl.KeyS) {
			movement = rl.Vector2Add(movement, rl.NewVector2(-0.1, 0))
		}
		if rl.IsKeyDown(rl.KeyA) {
			movement = rl.Vector2Add(movement, rl.NewVector2(0, -0.1))
		}
		if rl.IsKeyDown(rl.KeyD) {
			movement = rl.Vector2Add(movement, rl.NewVector2(0, 0.1))
		}

		if rl.IsKeyDown(rl.KeySpace) && isOnFloor {
			velocity.Y = 0.2
		}
		movement = rl.Vector2Normalize(movement)
		movement = rl.Vector2Rotate(movement, cameraRotationx)
		movement = rl.Vector2Scale(movement, -0.1)
		velocity.Y -= gravity
		velocity = rl.Vector3Add(velocity, rl.NewVector3(movement.X, 0, movement.Y))

		// if rl.IsKeyDown(rl.KeyE) && !waspressed {
		// 	waspressed = true
		// 	if moving == &player {
		// 		moving = &player2
		// 	} else {
		// 		moving = &player
		// 	}
		// }
		// if rl.IsKeyReleased(rl.KeyE) {
		// 	waspressed = false
		// }

		player.Collider.AddPosition(velocity)
		isOnFloor = false
		for _, obj := range objects {
			if obj != nil {
				direction := player.Collider.PushbackFrom(obj.Collider)
				if direction == types.DirYminus {
					isOnFloor = true
					velocity.Y = 0
				} else if direction == types.DirY {
					velocity.Y = 0
				}
			}
		}

		target := rl.Vector3{X: float32(math.Sin(cameraRotationy) * math.Cos(cameraRotationx)),
			Z: float32(math.Sin(cameraRotationy) * math.Sin(cameraRotationx)),
			Y: float32(math.Cos(cameraRotationy))}
		target = rl.Vector3Normalize(target)

		camera.Position = rl.Vector3Add(player.Collider.GetPosition(), rl.NewVector3(0, 0.5, 0))
		playerRay := entities.Ray{Origin: camera.Position, Direction: target}
		target = rl.Vector3Add(target, camera.Position)
		camera.Target = target
		rl.SetMousePosition(centerx, centery)

		// point, _ := playerRay.GetCollisionPoint(objects[3].Collider)
		// fmt.Println(objects[3].Collider)
		// fmt.Println(player.Collider.GetPosition())
		// if input != "" {
		// 	_, err = conn.Write([]byte(input))
		// 	if err != nil {
		// 		fmt.Println("Send error:", err)
		// 	}
		// }

		// println(player.Collider.GetPosition().Y)
		// println(types.Xminus)
		// println(-types.X)

		if rl.IsKeyDown(rl.KeyUp) {
			camera.Target.Y += 0.1
		}
		if rl.IsKeyDown(rl.KeyDown) {
			camera.Target.Y -= 0.1
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			camera.Target.X -= 0.1
		}
		if rl.IsKeyDown(rl.KeyRight) {
			camera.Target.X += 0.1
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
				if plane, ok := obj.Collider.(*entities.PlaneCollider); ok {
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
		// rl.DrawGrid(10, 1.0)

		rl.EndMode3D()
		rl.DrawText("Collision demo", 10, 10, 20, rl.Black)
		rl.EndDrawing()
	}
}
