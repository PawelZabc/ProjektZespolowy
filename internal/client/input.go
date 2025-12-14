package client

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/protocol"
	"github.com/PawelZabc/ProjektZespolowy/internal/shared"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Input struct {
	mouseLocked bool
	justClicked bool
}

func NewInput() *Input {
	rl.HideCursor()
	return &Input{
		mouseLocked: false,
		justClicked: false,
	}
}

// Gather input and rotation ans parse it to send to the server
func (i *Input) ProcessInput(rotationX, rotationY float32) protocol.ClientData {
	
	i.handleMouseLockToggle()

	inputs := make([]shared.PlayerAction, 0, 5)

	// TODO: add keys in keymap
	if rl.IsKeyDown(rl.KeyW) {
		inputs = append(inputs, shared.MoveForward)
	}
	if rl.IsKeyDown(rl.KeyS) {
		inputs = append(inputs, shared.MoveBackward)
	}
	if rl.IsKeyDown(rl.KeyA) {
		inputs = append(inputs, shared.MoveLeft)
	}
	if rl.IsKeyDown(rl.KeyD) {
		inputs = append(inputs, shared.MoveRight)
	}
	if rl.IsKeyDown(rl.KeySpace) {
		inputs = append(inputs, shared.Jump)
	}

	return protocol.ClientData{
		RotationX: rotationX,
		RotationY: rotationY,
		Inputs:    inputs,
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
