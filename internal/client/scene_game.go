package client

import (
	"context"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameScene struct {
	app         *App
	ctx         context.Context
	cancel      context.CancelFunc
	initialized bool
}

func NewGameScene(app *App) *GameScene {
	return &GameScene{app: app}
}

func (s *GameScene) OnEnter() {
	var err error

	err = s.app.initComponents()
	if err != nil {
		fmt.Printf("Failed to initialize components: %v\n", err)
		return
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())
	go s.app.network.StartReceiving(s.ctx)

	if err := s.app.network.SendHello(); err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		return
	}

	s.initialized = true
	rl.HideCursor()
}

func (s *GameScene) OnExit() {
	if s.cancel != nil {
		s.cancel()
	}
	if s.app.network != nil {
		s.app.network.Close()
		s.app.network = nil
	}
	rl.ShowCursor()
}

func (s *GameScene) Update() SceneType {
	if !s.initialized {
		return SceneMainMenu
	}

	if s.app.gameState != nil && s.app.gameState.playerHp <= 0 {
		return SceneGameOver
	}

	if rl.IsKeyPressed(rl.KeyEscape) {
		return SceneMainMenu
	}

	s.app.update()

	return SceneGame
}

func (s *GameScene) Render() {
	if !s.initialized {
		return
	}
	s.app.render()
}
