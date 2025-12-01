package types

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Collider interface {
	CollidesWith(Collider) bool
	GetPosition() rl.Vector3
	SetPosition(rl.Vector3)
	AddPosition(rl.Vector3)
	PushbackFrom(Collider) Direction
}

type Change struct {
	Value float32
	Axis  Direction
	Skip  bool
}
