package client

import rl "github.com/gen2brain/raylib-go/raylib"

// func (r Renderable) Draw(position rl.Vector3) {
// 	rl.DrawModelEx(r.Model, position, rl.NewVector3(0, 1, 0), r.Rotation, rl.Vector3One(), r.Color)
// }

type Renderable interface {
	GetModel() rl.Model
	SetModel(model rl.Model)
	GetColor() rl.Color
	SetColor(color rl.Color)
	GetShader() rl.Shader
	SetShader(shader rl.Shader)
	GetOffset() rl.Vector3
	SetOffset(offset rl.Vector3)
	GetRotation() float32
	SetRotation(rotation float32)
	Render(position rl.Vector3)
}

// func (r Renderable) Draw(position rl.Vector3) {
// 	rl.DrawModelEx(r.Model, position, rl.NewVector3(0, 1, 0), r.Rotation, rl.Vector3One(), r.Color)
// }
