package animation

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type GLBAnimation struct {
	anim       rl.ModelAnimation
	frame      int
	frameCount int32
}

func NewGLBAnimation(anim rl.ModelAnimation) *GLBAnimation {
	return &GLBAnimation{
		anim:       anim,
		frame:      0,
		frameCount: anim.FrameCount,
	}
}

func (a *GLBAnimation) NextFrame() {
	a.frame++
}

func (a *GLBAnimation) Apply(r core.Renderable) {
	rl.UpdateModelAnimation(r.GetModel(), a.anim, int32(a.frame%int(a.frameCount)))
}

func (a *GLBAnimation) Finished() bool {
	return a.frame >= int(a.frameCount)
}

func (a *GLBAnimation) Reset() {
	a.frame = 0
}

var _ Animation = (*GLBAnimation)(nil)
