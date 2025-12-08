package leveldata

import (
	types "github.com/PawelZabc/ProjektZespolowy/shared/_types"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Room struct {
	Objects       []*Object //objects and walls in that room
	SharedObjects []*Object //objects shared with visible rooms
	Id            uint16
	VisibleRooms  []*Room          //rooms visible that need to be rendered while in the room
	RoomChangers  []types.Collider //Colliders determining whenroom is changed
}

type Object struct {
	Colliders []types.Collider
	DrawPoint rl.Vector3
	Model     string
}
