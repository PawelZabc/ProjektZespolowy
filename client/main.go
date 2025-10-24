package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"net"
	"time"

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
	serverIP := flag.String("ip", "127.0.0.1", "Server IP address")
	// serverIP := "10.230.125.200" //"127.0.0.1" for local or ipconfig to check lan network
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
		Position:   rl.NewVector3(0.0, 2.0, 4.0),
		Target:     rl.NewVector3(0.0, 1.0, 0.0),
		Up:         rl.NewVector3(0.0, 1.0, 0.0),
		Fovy:       45.0,
		Projection: rl.CameraPerspective,
	}

	model := rl.LoadModel("assets/barrel2.glb")
	defer rl.UnloadModel(model)
	// animations := rl.LoadModelAnimations("assets/dude.glb")
	// defer rl.UnloadModelAnimations(animations)

	shader := rl.LoadShader("lighting.vs", "lighting.fs")
	defer rl.UnloadShader(shader)

	model.Materials.Shader = shader

	lightDirLoc := rl.GetShaderLocation(shader, "lightDir")
	baseColorLoc := rl.GetShaderLocation(shader, "baseColor")
	ambientColorLoc := rl.GetShaderLocation(shader, "ambientColor")

	lightDir := []float32{0.0, -1.0, -1.0}
	rl.SetShaderValue(shader, lightDirLoc, lightDir, rl.ShaderUniformVec3)

	rl.SetShaderValue(shader, baseColorLoc, []float32{1.0, 0.3, 0.2, 1.0}, rl.ShaderUniformVec4)
	rl.SetShaderValue(shader, ambientColorLoc, []float32{0.2, 0.2, 0.2, 1.0}, rl.ShaderUniformVec4)

	var position rl.Vector3

	go func() {
		buffer := make([]byte, 1024)
		for {
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, _, err := conn.ReadFromUDP(buffer)
			if err != nil {
				continue
			}
			// var pos Position
			var data Data
			err = json.Unmarshal(buffer[:n], &data)
			// if len(animations) > 0 && animations[0].FrameCount > 0 && model.BoneCount > 0 {
			// 	rl.UpdateModelAnimation(model, animations[0], data.Frame%animations[0].FrameCount)
			// }
			// println(data.Frame % animations[0].FrameCount)
			if err == nil {
				position = rl.NewVector3(data.X, data.Y, data.Z)
			}
		}
	}()

	rl.HideCursor()
	centerx := rl.GetScreenWidth() / 2
	centery := rl.GetScreenHeight() / 2
	cameraRotationx := float32(0)
	cameraRotationy := float32(0)
	conn.Write([]byte("hello"))
	for !rl.WindowShouldClose() {
		input := ""
		if rl.IsKeyDown(rl.KeyKp8) {
			input += "W"
		}
		if rl.IsKeyDown(rl.KeyKp4) {
			input += "A"
		}
		if rl.IsKeyDown(rl.KeyKp5) {
			input += "S"
		}
		if rl.IsKeyDown(rl.KeyKp6) {
			input += "D"
		}

		if input != "" {
			_, err = conn.Write([]byte(input))
			if err != nil {
				fmt.Println("Send error:", err)
			}
		}

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

		if rl.IsKeyDown(rl.KeyA) {
			camera.Position.X -= 0.1
		}
		if rl.IsKeyDown(rl.KeyD) {
			camera.Position.X += 0.1
		}
		if rl.IsKeyDown(rl.KeyW) {
			camera.Position.Z -= 0.1
		}
		if rl.IsKeyDown(rl.KeyS) {
			camera.Position.Z += 0.1
		}

		deltaMouse := rl.GetMousePosition()

		cameraRotationx -= (deltaMouse.X - float32(centerx)) / 100
		cameraRotationy -= (deltaMouse.Y - float32(centery)) / 100

		target := rl.Vector3{X: float32(math.Sin(float64(cameraRotationx))),
			Z: float32(math.Cos(float64(cameraRotationx))) + float32(math.Sin(float64(cameraRotationy))),
			Y: float32(math.Cos(float64(cameraRotationy)))}
		// target = rl.Vector3Normalize(target)
		target = rl.Vector3Add(target, camera.Position)
		camera.Target = target
		rl.SetMousePosition(centerx, centery)

		println(float32(math.Sin(float64(cameraRotationy))), ",", float32(math.Cos(float64(cameraRotationy))))
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawModel(model, position, 1.0, rl.White)
		rl.DrawGrid(10, 1.0)

		rl.EndMode3D()
		rl.DrawText("Use WASD to move the object", 10, 10, 20, rl.Black)
		rl.DrawText(fmt.Sprintf("Position: X=%.2f  Y=%.2f  Z=%.2f", position.X, position.Y, position.Z), 10, 40, 20, rl.DarkBlue)
		rl.EndDrawing()
	}
}
