package client

import rl "github.com/gen2brain/raylib-go/raylib"

type GameOverScene struct {
	app     *App
	victory bool
}

func NewGameOverScene(app *App, victory bool) *GameOverScene {
	return &GameOverScene{
		app:     app,
		victory: victory,
	}
}

func (s *GameOverScene) OnEnter() {
	rl.ShowCursor()
}

func (s *GameOverScene) OnExit() {}

func (s *GameOverScene) Update() SceneType {
	if rl.IsKeyPressed(rl.KeyEnter) {
		return SceneMainMenu
	}

	if rl.IsKeyPressed(rl.KeyEscape) {
		s.app.Stop()
	}

	return SceneGameOver
}

func (s *GameOverScene) Render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	var title string
	var color rl.Color

	if s.victory {
		title = "VICTORY!"
		color = rl.Green
	} else {
		title = "GAME OVER"
		color = rl.Red
	}

	titleWidth := rl.MeasureText(title, 60)
	rl.DrawText(title, (int32(s.app.config.WindowWidth)-titleWidth)/2, 200, 60, color)

	menuText := "Press ENTER for Main Menu"
	menuWidth := rl.MeasureText(menuText, 20)
	rl.DrawText(menuText, (int32(s.app.config.WindowWidth)-menuWidth)/2, 350, 20, rl.White)

	rl.EndDrawing()
}
