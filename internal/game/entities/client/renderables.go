package client

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type BasicRenderable struct {
	Model    rl.Model  // old object
	Color    rl.Color  // old object
	Shader   rl.Shader // old object
	Offset   rl.Vector3
	Rotation float32 // old Actor
}

// GetColor implements Renderable.
func (s *BasicRenderable) GetColor() rl.Color {
	return s.Color
}

// GetModel implements Renderable.
func (s *BasicRenderable) GetModel() rl.Model {
	return s.Model
}

// GetOffset implements Renderable.
func (s *BasicRenderable) GetOffset() rl.Vector3 {
	return s.Offset
}

// GetRotation implements Renderable.
func (s *BasicRenderable) GetRotation() float32 {
	return s.Rotation
}

// GetShader implements Renderable.
func (s *BasicRenderable) GetShader() rl.Shader {
	return s.Shader
}

func (s BasicRenderable) Render(position rl.Vector3) {
	rl.DrawModelEx(s.Model, position, rl.NewVector3(0, 1, 0), s.Rotation, rl.Vector3One(), s.Color)
}

// SetColor implements Renderable.
func (s *BasicRenderable) SetColor(color rl.Color) {
	s.Color = color
}

// SetModel implements Renderable.
func (s *BasicRenderable) SetModel(model rl.Model) {
	s.Model = model
}

// SetOffset implements Renderable.
func (s *BasicRenderable) SetOffset(offset rl.Vector3) {
	s.Offset = offset
}

// SetRotation implements Renderable.
func (s *BasicRenderable) SetRotation(rotation float32) {
	s.Rotation = rotation
}

// SetShader implements Renderable.
func (s *BasicRenderable) SetShader(shader rl.Shader) {
	s.Shader = shader
}

var _ core.Renderable = (*BasicRenderable)(nil)
