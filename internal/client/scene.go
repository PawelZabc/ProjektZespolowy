package client

type SceneType int

const (
	SceneMainMenu SceneType = iota
	SceneGame
	SceneGameOver
)

type Scene interface {
	Update() SceneType
	Render()
	OnEnter()
	OnExit()
}
