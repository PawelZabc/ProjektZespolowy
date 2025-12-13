package levels

import "github.com/PawelZabc/ProjektZespolowy/internal/shared"

type Change struct {
	Value float32
	Axis  shared.Direction
	Skip  bool
}
