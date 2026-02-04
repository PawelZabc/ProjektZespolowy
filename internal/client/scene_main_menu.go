package client

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type MainMenuScene struct {
	app           *App
	serverBrowser *ServerBrowser
	refreshTimer  float32
}

func NewMainMenuScene(app *App) *MainMenuScene {
	return &MainMenuScene{
		app:           app,
		serverBrowser: NewServerBrowser(app.config.ServerIP, app.config.ServerPort),
	}
}

func (s *MainMenuScene) OnEnter() {
	rl.ShowCursor()
	s.serverBrowser.Refresh()
}

func (s *MainMenuScene) OnExit() {
	rl.HideCursor()
}

func (s *MainMenuScene) Update() SceneType {
	// refresh player list every 3 seconds
	s.refreshTimer += rl.GetFrameTime()
	if s.refreshTimer >= 3.0 {
		s.refreshTimer = 0
		s.serverBrowser.Refresh()
	}

	if rl.IsKeyPressed(rl.KeyEnter) {
		return SceneGame
	}

	if rl.IsKeyPressed(rl.KeyEscape) {
		s.app.Stop()
	}

	return SceneMainMenu
}

func (s *MainMenuScene) Render() {
	screenWidth := int32(s.app.config.WindowWidth)
	screenHeight := int32(s.app.config.WindowHeight)
	centerX := screenWidth / 2

	rl.BeginDrawing()
	rl.ClearBackground(rl.NewColor(30, 30, 40, 255))

	// Title
	title := s.app.config.WindowTitle
	titleWidth := rl.MeasureText(title, 60)
	rl.DrawText(title, centerX-titleWidth/2, 100, 60, rl.White)

	// Press ENTER to start
	startText := "Press ENTER to start"
	startWidth := rl.MeasureText(startText, 24)
	rl.DrawText(startText, centerX-startWidth/2, 200, 24, rl.Green)

	// ESC to exit
	exitText := "ESC - Exit"
	exitWidth := rl.MeasureText(exitText, 16)
	rl.DrawText(exitText, centerX-exitWidth/2, 240, 16, rl.Gray)

	// Bottom: Server status and player list
	bottomY := screenHeight - 120

	// Server status line
	var statusText string
	var statusColor rl.Color
	if s.serverBrowser.Online {
		statusText = fmt.Sprintf("Server %s:%d - Online", s.app.config.ServerIP, s.app.config.ServerPort)
		statusColor = rl.Green
	} else {
		statusText = fmt.Sprintf("Server %s:%d - Offline", s.app.config.ServerIP, s.app.config.ServerPort)
		statusColor = rl.Red
	}
	rl.DrawText(statusText, 20, bottomY, 16, statusColor)

	// Player list
	if s.serverBrowser.Online {
		playerCount := len(s.serverBrowser.Players)
		if playerCount == 0 {
			rl.DrawText("No players online", 20, bottomY+22, 14, rl.Gray)
		} else {
			playerText := fmt.Sprintf("Players online: %d", playerCount)
			rl.DrawText(playerText, 20, bottomY+22, 14, rl.LightGray)

			// List player IDs (compact)
			x := int32(20)
			for i, p := range s.serverBrowser.Players {
				if i >= 8 {
					rl.DrawText("...", x, bottomY+44, 14, rl.Gray)
					break
				}
				idText := fmt.Sprintf("#%d", p.ID)
				rl.DrawText(idText, x, bottomY+44, 14, rl.LightGray)
				x += int32(rl.MeasureText(idText, 14)) + 15
			}
		}
	}

	rl.EndDrawing()
}
