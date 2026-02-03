package client

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/core"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type CEntity struct {
	core.Renderable
	Position rl.Vector3
	Collider colliders.Collider
	// IsInteractable bool
}

func (e CEntity) Render() {
	e.Renderable.Render(e.Position)
}

func (e *CEntity) SetPosition(p rl.Vector3) {
	e.Position = p
}

func RenderEntities(entities []*CEntity) {
	for _, e := range entities {
		e.Render()
	}
}
