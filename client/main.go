package main

import (
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
		Fovy:       45.0,
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

	player := createCylinderObject(rl.NewVector3(-2, 0, 0), 0.5, 1)

	object := createCylinderObject(rl.NewVector3(2, 0, 0), 0.5, 1)

	player2 := createCubeObject(rl.NewVector3(-2, 0, -3), 1, 1, 1)

	object2 := createCubeObject(rl.NewVector3(2, 0, -3), 1, 1, 1)

	// conn.Write([]byte("hello"))
	rl.HideCursor()
	centerx := rl.GetScreenWidth() / 2
	centery := rl.GetScreenHeight() / 2
	cameraRotationx := float32(0)
	cameraRotationy := float32(0)

	moving := &player
	waspressed := false
	for !rl.WindowShouldClose() {
		deltaMouse := rl.GetMousePosition()

		cameraRotationx += (deltaMouse.X - float32(centerx)) / 100
		cameraRotationy -= (deltaMouse.Y - float32(centery)) / 100

		target := rl.Vector3{X: float32(math.Sin(cameraRotationy) * math.Cos(cameraRotationx)),
			Z: float32(math.Sin(cameraRotationy) * math.Sin(cameraRotationx)),
			Y: float32(math.Cos(cameraRotationy))}
		target = rl.Vector3Normalize(target)
		target = rl.Vector3Add(target, camera.Position)
		camera.Target = target
		rl.SetMousePosition(centerx, centery)

		// input := ""
		velocity := rl.Vector3{}
		if rl.IsKeyDown(rl.KeyW) {
			velocity = rl.Vector3Add(velocity, rl.NewVector3(0, 0, -0.1))
		}
		if rl.IsKeyDown(rl.KeyS) {
			velocity = rl.Vector3Add(velocity, rl.NewVector3(0, 0, 0.1))
		}
		if rl.IsKeyDown(rl.KeyA) {
			velocity = rl.Vector3Add(velocity, rl.NewVector3(-0.1, 0, 0))
		}
		if rl.IsKeyDown(rl.KeyD) {
			velocity = rl.Vector3Add(velocity, rl.NewVector3(0.1, 0, 0))
		}
		if rl.IsKeyDown(rl.KeyLeftShift) {
			velocity = rl.Vector3Add(velocity, rl.NewVector3(0, -0.1, 0))
		}
		if rl.IsKeyDown(rl.KeySpace) {
			velocity = rl.Vector3Add(velocity, rl.NewVector3(0, 0.1, 0))
		}

		if rl.IsKeyDown(rl.KeyE) && !waspressed {
			waspressed = true
			if moving == &player {
				moving = &player2
			} else {
				moving = &player
			}
		}
		if rl.IsKeyReleased(rl.KeyE) {
			waspressed = false
		}

		moving.Collider.AddPosition(velocity)
		if moving.Collider.CollidesWith(object.Collider) {
			println("cylinder collision")
		}
		moving.Collider.PushbackFrom(object.Collider)
		if moving.Collider.CollidesWith(object2.Collider) {
			println("cube collision")
		}
		moving.Collider.PushbackFrom(object2.Collider)

		// if input != "" {
		// 	_, err = conn.Write([]byte(input))
		// 	if err != nil {
		// 		fmt.Println("Send error:", err)
		// 	}
		// }

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

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawModel(player.Model, player.Collider.GetPosition(), 1.0, rl.White)
		rl.DrawModel(object.Model, object.Collider.GetPosition(), 1.0, rl.White)
		rl.DrawModel(player2.Model, player2.Collider.GetPosition(), 1.0, rl.White)
		rl.DrawModel(object2.Model, object2.Collider.GetPosition(), 1.0, rl.White)
		rl.DrawGrid(10, 1.0)

		rl.EndMode3D()
		rl.DrawText("Collision demo", 10, 10, 20, rl.Black)
		rl.EndDrawing()
	}
}
