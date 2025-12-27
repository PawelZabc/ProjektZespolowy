package levels

import "github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"


type Room struct {
	Colliders    []colliders.Collider //objects and walls in that room
	RoomChangers []colliders.Collider
}


func ServerLoadRooms() []Room {
	rooms := make([]Room, 0, 10)
	roomShared := Room1
	objects := make([]colliders.Collider, 0, len(roomShared.Objects))
	for _, object := range roomShared.Objects {
		for _, collider := range object.Colliders {
			objects = append(objects, collider)
		}
	}
	for _, object := range roomShared.SharedObjects {
		for _, collider := range object.Colliders {
			objects = append(objects, collider)
		}
	}
	// changers := make([]s_types.Collider, 0, len(roomShared.RoomChangers))

	room := Room{
		Colliders:    objects,
		RoomChangers: roomShared.RoomChangers,
	}

	rooms = append(rooms, room)

	return rooms
}
