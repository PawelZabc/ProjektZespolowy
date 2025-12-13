package levels

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	"github.com/PawelZabc/ProjektZespolowy/internal/shared"
	"github.com/chewxy/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var Room1 = CreateRoom1()

func CreateRoom1() RoomTWO {
	room := CreateRoomFromChanges(BasicLevel, rl.NewVector3(-10, 0, -10), 3)
	//add objects
	cylinder := CreateObjectFromCollider(colliders.NewCylinderCollider(rl.NewVector3(2, 0, 0), 0.5, 1))
	// cylinder.Model =
	room.Objects = append(room.Objects, cylinder)

	return room
}

func CreateRoomFromChanges(changes []Change, start rl.Vector3, height float32) RoomTWO {
	wallColliders := CreateRoomWallsFromChanges(start, changes, height)                   //create walls
	floor, ceiling := MakeFloorAndCeilingForWalls(wallColliders, start.Y, start.Y+height) //create floor and ceiling for the walls
	wallColliders = append(wallColliders, floor)                                          //add them to colliders
	wallColliders = append(wallColliders, ceiling)
	object := CreateObjectFromColliders(wallColliders, rl.NewVector3(0, 0, 0)) // convert to objects
	object.Model = "room.glb"

	return RoomTWO{
		Objects: []*ObjectTWO{object},
	}
}

func MakeFloorAndCeilingForWalls(wallColliders []colliders.Collider, startY float32, endY float32) (colliders.Collider, colliders.Collider) { // create planes that span that cover every wall
	var minVector, maxVector = wallColliders[0].GetPosition(), wallColliders[0].GetPosition()
	for _, wall := range wallColliders { //look for smallest x and z, and biggest x and z
		if wall.GetPosition().X < minVector.X {
			minVector.X = wall.GetPosition().X
		} else if wall.GetPosition().X > maxVector.X {
			maxVector.X = wall.GetPosition().X
		}
		if wall.GetPosition().Z < minVector.Z {
			minVector.Z = wall.GetPosition().Z
		} else if wall.GetPosition().Z > maxVector.Z {
			maxVector.Z = wall.GetPosition().Z
		}
	}
	planeWidth := math32.Abs(minVector.X - maxVector.X) //calculate the size of the plane
	planeHeight := math32.Abs(minVector.Z - maxVector.Z)
	minVector.Y = startY
	floor := colliders.NewPlaneCollider(minVector, planeWidth, planeHeight, shared.DirY)
	minVector.Y = minVector.Y + endY
	ceiling := colliders.NewPlaneCollider(minVector, planeWidth, planeHeight, -shared.DirY)
	return floor, ceiling
}

func CreateObjectsFromColliders(colliders []colliders.Collider) []*ObjectTWO {
	objects := make([]*ObjectTWO, 0, len(colliders))
	for _, collider := range colliders {
		objects = append(objects, CreateObjectFromCollider(collider))
	}
	return objects

}
func CreateObjectFromColliders(colliders []colliders.Collider, drawPoint rl.Vector3) *ObjectTWO {
	object := ObjectTWO{
		DrawPoint: drawPoint,
		Colliders: colliders,
	}
	return &object

}

func CreateObjectFromCollider(collider colliders.Collider) *ObjectTWO {
	object := ObjectTWO{
		DrawPoint: collider.GetPosition(),
		Colliders: []colliders.Collider{collider},
	}
	return &object
}
