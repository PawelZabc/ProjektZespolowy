package client

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/animation"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/state"
)

type Actor struct {
	CEntity
	AnimationHandler animation.AnimationHandler
	State            state.State

	animationBase animation.AnimationMap
	lastState     state.State
}

func NewActor(entity CEntity) *Actor {
	return &Actor{
		CEntity:          entity,
		AnimationHandler: animation.AnimationHandler{},
	}
}

func (a *Actor) Render() {
	// TODO: state walking is per 20 ticks while attacking - to discover

	if a.State != a.lastState {
		if animationForState, ok := a.animationBase[a.State]; ok {
			a.AnimationHandler.SetAnimation(animationForState)
		}
		a.lastState = a.State
	}

	a.AnimationHandler.Update(a.Renderable)
	a.CEntity.Render()
}

func RenderActorsMap[T comparable](actors map[T]*Actor) {
	for _, actor := range actors {
		actor.Render()
	}
}
