package levels

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics"
)

type Change struct {
	Value float32
	Axis  physics.Direction
	Skip  bool
}
