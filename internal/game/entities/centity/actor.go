package client

type Actor struct {
	Entity           CEntity
	AnimationHandler AnimationHandler
}

func (a *Actor) Render() {
	a.AnimationHandler.Update(a.Entity.Renderable.GetModel())
	a.Entity.Render()
}
