package core

import rl "github.com/gen2brain/raylib-go/raylib"

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
