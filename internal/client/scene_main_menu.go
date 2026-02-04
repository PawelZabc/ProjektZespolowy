package client

import rl "github.com/gen2brain/raylib-go/raylib"

type MainMenuScene struct {
	app *App
}

func NewMainMenuScene(app *App) *MainMenuScene {
	return &MainMenuScene{app: app}
}

func (s *MainMenuScene) OnEnter() {
	rl.ShowCursor()
}

func (s *MainMenuScene) OnExit() {
	rl.HideCursor()
}

func (s *MainMenuScene) Update() SceneType {
	if rl.IsKeyPressed(rl.KeyEnter) {
		return SceneGame
	}

	if rl.IsKeyPressed(rl.KeyEscape) {
		s.app.Stop()
	}

	return SceneMainMenu
}

func (s *MainMenuScene) Render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.DarkGray)

	title := "G-Game"
	titleWidth := rl.MeasureText(title, 40)
	rl.DrawText(title, (int32(s.app.config.WindowWidth)-titleWidth)/2, 150, 40, rl.White)

	startText := "Press ENTER to Start"
	startWidth := rl.MeasureText(startText, 20)
	rl.DrawText(startText, (int32(s.app.config.WindowWidth)-startWidth)/2, 300, 20, rl.LightGray)

	exitText := "Press ESC to Exit"
	exitWidth := rl.MeasureText(exitText, 20)
	rl.DrawText(exitText, (int32(s.app.config.WindowWidth)-exitWidth)/2, 340, 20, rl.LightGray)

	rl.EndDrawing()
}
