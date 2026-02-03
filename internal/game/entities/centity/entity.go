package client

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type CEntity struct {
	Renderable Renderable
	Position   rl.Vector3
	Collider   *colliders.Collider
	// IsInteractable bool
}

func (e CEntity) Render() {
	e.Renderable.Render(e.Position)
}

func RenderEntities(entities []*CEntity) {
	for _, e := range entities {
		e.Render()
	}
}
