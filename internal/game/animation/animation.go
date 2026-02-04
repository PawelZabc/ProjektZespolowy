package animation

import "github.com/PawelZabc/ProjektZespolowy/internal/game/core"


type Animation interface {
	NextFrame()
	Apply(r core.Renderable)
	Finished() bool
	Reset()
}
