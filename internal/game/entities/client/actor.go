package client

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/animation"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/state"
)

type Actor struct {
	Entity           CEntity
	AnimationHandler animation.AnimationHandler
	State            state.State
}

func NewActor(entity CEntity, animationHandler animation.AnimationHandler) *Actor {
	return &Actor{
		Entity:           entity,
		AnimationHandler: animationHandler,
	}
}

func (a *Actor) Render() {
	a.AnimationHandler.Update(a.Entity.Renderable)
	a.Entity.Render()
}
