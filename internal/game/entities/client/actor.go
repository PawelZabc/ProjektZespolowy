package client

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/animation"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/state"
)

type Actor struct {
	Entity           CEntity
	AnimationHandler animation.AnimationHandler
	State            state.State

	animationBase animation.AnimationMap
	lastState     state.State
}

func NewActor(entity CEntity, animationHandler animation.AnimationHandler) *Actor {
	return &Actor{
		Entity:           entity,
		AnimationHandler: animationHandler,
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

	a.AnimationHandler.Update(a.Entity.Renderable)
	a.Entity.Render()
}
