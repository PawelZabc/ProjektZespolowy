package animation

import "github.com/PawelZabc/ProjektZespolowy/internal/game/state"

type AnimationMap map[state.State]Animation

func NewEnemyAnimationMap() AnimationMap {
	return AnimationMap{
		state.Attacking: NewRedColorAnimation(),
		state.Walking:   NewEmptyAnimation(),
	}
}
