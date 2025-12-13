package levels

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Room struct {
	Objects       []*Object //objects and walls in that room
	SharedObjects []*Object //objects shared with visible rooms
	VisibleRooms  []*Room   //rooms visible that need to be rendered while in the room
}



type RoomTWO struct {
	Objects       []*ObjectTWO //objects and walls in that room
	SharedObjects []*ObjectTWO //objects shared with visible rooms
	Id            uint16
	VisibleRooms  []*RoomTWO          //rooms visible that need to be rendered while in the room
	RoomChangers  []colliders.Collider //Colliders determining whenroom is changed
}

type ObjectTWO struct {
	Colliders []colliders.Collider
	DrawPoint rl.Vector3
	Model     string
	// Texture   string
}