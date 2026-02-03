package client

import rl "github.com/gen2brain/raylib-go/raylib"

type StaticRenderable struct {
	Model    rl.Model
	Color    rl.Color
	Shader   rl.Shader
	Offset   rl.Vector3
	Rotation float32
}

// GetColor implements Renderable.
func (s *StaticRenderable) GetColor() rl.Color {
	return s.Color
}

// GetModel implements Renderable.
func (s *StaticRenderable) GetModel() rl.Model {
	return s.Model
}

// GetOffset implements Renderable.
func (s *StaticRenderable) GetOffset() rl.Vector3 {
	return s.Offset
}

// GetRotation implements Renderable.
func (s *StaticRenderable) GetRotation() float32 {
	return s.Rotation
}

// GetShader implements Renderable.
func (s *StaticRenderable) GetShader() rl.Shader {
	return s.Shader
}

func (s StaticRenderable) Render(position rl.Vector3) {
	rl.DrawModelEx(s.Model, position, rl.NewVector3(0, 1, 0), s.Rotation, rl.Vector3One(), s.Color)
}

// SetColor implements Renderable.
func (s *StaticRenderable) SetColor(color rl.Color) {
	s.Color = color
}

// SetModel implements Renderable.
func (s *StaticRenderable) SetModel(model rl.Model) {
	s.Model = model
}

// SetOffset implements Renderable.
func (s *StaticRenderable) SetOffset(offset rl.Vector3) {
	s.Offset = offset
}

// SetRotation implements Renderable.
func (s *StaticRenderable) SetRotation(rotation float32) {
	s.Rotation = rotation
}

// SetShader implements Renderable.
func (s *StaticRenderable) SetShader(shader rl.Shader) {
	s.Shader = shader
}

var _ Renderable = (*StaticRenderable)(nil)
