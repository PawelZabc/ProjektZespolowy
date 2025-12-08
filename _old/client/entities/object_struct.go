package entities

import (
	s_types "github.com/PawelZabc/ProjektZespolowy/shared/_types"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Object struct {
	Collider s_types.Collider
	Model    rl.Model
}

func (o *Object) AddPosition(vec rl.Vector3) {
	o.Collider.SetPosition(rl.Vector3Add(o.Collider.GetPosition(), vec))
}
