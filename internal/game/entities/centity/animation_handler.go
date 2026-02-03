package client

import rl "github.com/gen2brain/raylib-go/raylib"

type AnimationHandler struct {
	Frame           int32
	Animation       int16
	ModelAnimations []rl.ModelAnimation
}

func (a *AnimationHandler) SetAnimation(animation int16) {
	if a.Animation != animation {
		a.Frame = 0
		a.Animation = animation
	}

}

func (a *AnimationHandler) Update(model rl.Model) {
	a.Frame = (a.Frame + 1) % a.ModelAnimations[a.Animation].FrameCount
	rl.UpdateModelAnimation(model, a.ModelAnimations[a.Animation], a.Frame)
}

// type Animation struct {

// }
