package entities

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/shared/types"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Object struct {
	Collider types.Collider
	Model    rl.Model
}

func (o *Object) AddPosition(vec rl.Vector3) {
	o.Collider.SetPosition(rl.Vector3Add(o.Collider.GetPosition(), vec))
}
