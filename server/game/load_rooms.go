package game

import (
	s_types "github.com/PawelZabc/ProjektZespolowy/shared/_types"
	leveldata "github.com/PawelZabc/ProjektZespolowy/shared/level_data"
)

func LoadRooms() []Room {
	rooms := make([]Room, 0, 10)
	roomShared := leveldata.Room1
	objects := make([]s_types.Collider, 0, len(roomShared.Objects))
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

type Room struct {
	Colliders    []s_types.Collider //objects and walls in that room
	RoomChangers []s_types.Collider
}
