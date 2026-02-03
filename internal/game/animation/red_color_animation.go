package animation

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type RedColorAnimation struct {
	frame int
	max   int
}

func NewRedColorAnimation(max int) *RedColorAnimation {
	return &RedColorAnimation{
		frame: 0,
		max:   max,
	}
}

func (r *RedColorAnimation) Apply(renderable core.Renderable) {
	notRed := max(255 - (8 * r.frame), 0)
	renderable.SetColor(rl.NewColor(255, uint8(notRed), uint8(notRed), 255))
}

func (r *RedColorAnimation) Finished() bool {
	return r.frame >= r.max
}

func (r *RedColorAnimation) NextFrame() {
	r.frame++
}

var _ Animation = (*RedColorAnimation)(nil)
