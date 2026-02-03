package animation

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/core"
)

type AnimationHandler struct {
	current Animation
}

func (a *AnimationHandler) SetAnimation(animation Animation) {
	a.current = animation
}

func (a *AnimationHandler) Update(r core.Renderable) {
	if a.current == nil {
		return
	}

	a.current.NextFrame()
	a.current.Apply(r)

	if a.current.Finished() {
		a.current = nil
	}
}

// func (a *AnimationHandler) Update(model rl.Model) {
// 	a.Frame = (a.Frame + 1) % a.ModelAnimations[a.Animation].FrameCount
// 	rl.UpdateModelAnimation(model, a.ModelAnimations[a.Animation], a.Frame)
// }
