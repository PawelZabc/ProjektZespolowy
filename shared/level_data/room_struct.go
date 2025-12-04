package leveldata

import (
	types "github.com/PawelZabc/ProjektZespolowy/shared/_types"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Room struct {
	Objects      []*Object
	Id           uint16
	VisibleRooms []*Room
}

type Object struct {
	Colliders []types.Collider
	Model     rl.Model
}
