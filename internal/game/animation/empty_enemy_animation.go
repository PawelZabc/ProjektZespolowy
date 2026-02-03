package animation

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type EmptyAnimation struct{}

func NewEmptyAnimation() *EmptyAnimation {
	return &EmptyAnimation{}
}

func (e *EmptyAnimation) NextFrame() {}

func (e *EmptyAnimation) Apply(r core.Renderable) {
	r.SetColor(rl.White)
}

func (e *EmptyAnimation) Finished() bool {
	return false
}

func (e *EmptyAnimation) Reset() {}

var _ Animation = (*EmptyAnimation)(nil)
