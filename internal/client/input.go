package client

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/input"
	"github.com/PawelZabc/ProjektZespolowy/internal/protocol"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Input struct {
	mouseLocked bool
	justClicked bool
	inputBuffer []input.PlayerAction
}

func NewInput() *Input {
	rl.HideCursor()
	return &Input{
		mouseLocked: false,
		justClicked: false,
		inputBuffer: make([]input.PlayerAction, 0, 5),
	}
}

// Gather input and rotation ans parse it to send to the server
func (i *Input) ProcessInput(rotationX, rotationY float32) protocol.ClientData {

	i.inputBuffer = i.inputBuffer[:0] // resetting buffer without reallocation
	i.handleMouseLockToggle()

	// TODO: add keys in keymap
	if rl.IsKeyDown(rl.KeyW) {
		i.inputBuffer = append(i.inputBuffer, input.MoveForward)
	}
	if rl.IsKeyDown(rl.KeyS) {
		i.inputBuffer = append(i.inputBuffer, input.MoveBackward)
	}
	if rl.IsKeyDown(rl.KeyA) {
		i.inputBuffer = append(i.inputBuffer, input.MoveLeft)
	}
	if rl.IsKeyDown(rl.KeyD) {
		i.inputBuffer = append(i.inputBuffer, input.MoveRight)
	}
	if rl.IsKeyDown(rl.KeySpace) {
		i.inputBuffer = append(i.inputBuffer, input.Jump)
	}

	return protocol.ClientData{
		RotationX: rotationX,
		RotationY: rotationY,
		Inputs:    i.inputBuffer,
	}
}

// pressing R to lock and unlock mouse
func (i *Input) handleMouseLockToggle() {
	if rl.IsKeyDown(rl.KeyR) && !i.justClicked {
		i.mouseLocked = !i.mouseLocked
		i.justClicked = true

		if i.mouseLocked {
			rl.HideCursor()
		} else {
			rl.ShowCursor()
		}
	}

	if rl.IsKeyReleased(rl.KeyR) {
		i.justClicked = false
	}
}

// Returns whether the mouse is locked
func (i *Input) IsMouseLocked() bool {
	return i.mouseLocked
}
