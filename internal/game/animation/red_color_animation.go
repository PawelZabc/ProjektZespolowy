package animation

import (
	"log"

	"github.com/PawelZabc/ProjektZespolowy/internal/game/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type RedColorAnimation struct {
	frame int
	max   int
}

func NewRedColorAnimation() *RedColorAnimation {
	return &RedColorAnimation{
		frame: 0,
		max:   32,
	}
}

func (r *RedColorAnimation) Apply(renderable core.Renderable) {
	notRed := min((8 * r.frame), 255)
	log.Print(notRed)
	renderable.SetColor(rl.NewColor(255, uint8(notRed), uint8(notRed), 255))
}

func (r *RedColorAnimation) Finished() bool {
	return r.frame >= r.max
}

func (r *RedColorAnimation) NextFrame() {
	r.frame++
}

func (r *RedColorAnimation) Reset() {
	r.frame = 0
}

var _ Animation = (*RedColorAnimation)(nil)
