package shared

type PlayerAction uint8

const (
	MoveForward PlayerAction = iota
	MoveBackward
	MoveLeft
	MoveRight
	Jump
)
