package client

import (
	"context"
	"fmt"

	"github.com/PawelZabc/ProjektZespolowy/internal/config"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type App struct {
	config    config.ClientConfig
	network   *Network
	gameState *GameState
	camera    *Camera
	input     *Input
	renderer  *Renderer
	running   bool
}

func NewApp(cfg config.ClientConfig) *App {
	return &App{
		config:  cfg,
		running: false,
	}
}

func (a *App) Run() error {
	a.initWindow() // init raylib

	defer a.cleanup() // end with cleanup

	// init components like input, network, state, etc..
	if err := a.initComponents(); err != nil {
		return fmt.Errorf("failed to initialize components: %w", err)
	}

	// goroutine to reveive messages and data from server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go a.network.StartReceiving(ctx)

	// initial connection to server
	if err := a.network.SendHello(); err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}

	a.running = true
	a.gameLoop()

	return nil
}

func (a *App) initWindow() {
	rl.InitWindow(int32(a.config.WindowWidth), int32(a.config.WindowHeight), a.config.WindowTitle)
	rl.SetTargetFPS(int32(a.config.TargetFPS))
}

func (a *App) initComponents() error {
	var err error

	a.gameState = NewGameState()
	a.network, err = NewNetwork(a.config.ServerIP, a.config.ServerPort, a.gameState)
	if err != nil {
		return fmt.Errorf("network initialization failed: %w", err)
	}
	a.camera = NewCamera()
	a.input = NewInput()
	a.renderer = NewRenderer(a.config.DebugMode)

	return nil
}

// simple game loop
func (a *App) gameLoop() {
	for !rl.WindowShouldClose() && a.running {
		a.update()
		a.render()
	}
}

// Handles all game logic updates
func (a *App) update() {
	inputData := a.input.ProcessInput(a.camera.GetRotationX(), a.camera.GetRotationY())

	centerX := a.config.WindowWidth / 2
	centerY := a.config.WindowHeight / 2
	a.camera.Update(centerX, centerY, a.input.IsMouseLocked())
	a.camera.UpdatePosition(a.gameState.GetPlayerPosition())

	if err := a.network.SendInput(inputData); err != nil {
		fmt.Printf("Failed to send input: %v\n", err)
	}

	// updating player "cursor"
	if a.config.DebugMode {
		playerRay := a.camera.GetPlayerCameraRay()
		a.gameState.UpdateRayCollision(playerRay)
	}

}

func (a *App) render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	rl.BeginMode3D(a.camera.GetCamera())
	a.renderer.RenderWorld(a.gameState)
	rl.EndMode3D()

	a.renderer.RenderUI()

	rl.EndDrawing()
}

func (a *App) cleanup() {
	if a.network != nil {
		a.network.Close()
	}
	rl.CloseWindow()
}

func (a *App) Stop() {
	a.running = false
}
