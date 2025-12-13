package levels

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/entities"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ClientRoom struct {
	Objects       []*entities.Object //objects and walls in that room
	SharedObjects []*entities.Object //objects shared with visible rooms
	VisibleRooms  []*ClientRoom      //rooms visible that need to be rendered while in the room
}

// shared struct - JSON TODO
type RoomTWO struct {
	Objects       []*ObjectTWO //objects and walls in that room
	SharedObjects []*ObjectTWO //objects shared with visible rooms
	Id            uint16
	VisibleRooms  []*RoomTWO           //rooms visible that need to be rendered while in the room
	RoomChangers  []colliders.Collider //Colliders determining whenroom is changed
}

// shared struct -JSON TODO
type ObjectTWO struct {
	Colliders []colliders.Collider // TODO: possibly to delete if position provided
	DrawPoint rl.Vector3
	Model     string
	// Texture   string
}
