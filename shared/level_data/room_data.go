package leveldata

import (
	types "github.com/PawelZabc/ProjektZespolowy/shared/_types"
	"github.com/PawelZabc/ProjektZespolowy/shared/entities"
	"github.com/chewxy/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var Room1 = CreateRoom1()

func CreateRoom1() Room {
	room := CreateRoomFromChanges(Changes, rl.NewVector3(-10, 0, -10), 3)
	//add objects
	cylinder := CreateObjectFromCollider(entities.NewCylinderCollider(rl.NewVector3(2, 0, 0), 0.5, 1))
	// cylinder.Model =
	room.Objects = append(room.Objects, cylinder)

	return room
}

func CreateRoomFromChanges(changes []types.Change, start rl.Vector3, height float32) Room {
	wallColliders := entities.CreateRoomWallsFromChanges(start, changes, height)          //create walls
	floor, ceiling := MakeFloorAndCeilingForWalls(wallColliders, start.Y, start.Y+height) //create floor and ceiling for the walls
	wallColliders = append(wallColliders, floor)                                          //add them to colliders
	wallColliders = append(wallColliders, ceiling)
	object := CreateObjectFromColliders(wallColliders, rl.NewVector3(0, 0, 0)) // convert to objects
	object.Model = "room.glb"

	return Room{
		Objects: []*Object{object},
	}
}

func MakeFloorAndCeilingForWalls(wallColliders []types.Collider, startY float32, endY float32) (types.Collider, types.Collider) { // create planes that span that cover every wall
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
	floor := entities.NewPlaneCollider(minVector, planeWidth, planeHeight, types.DirY)
	minVector.Y = minVector.Y + endY
	ceiling := entities.NewPlaneCollider(minVector, planeWidth, planeHeight, -types.DirY)
	return floor, ceiling
}

func CreateObjectsFromColliders(colliders []types.Collider) []*Object {
	objects := make([]*Object, 0, len(colliders))
	for _, collider := range colliders {
		objects = append(objects, CreateObjectFromCollider(collider))
	}
	return objects

}
func CreateObjectFromColliders(colliders []types.Collider, drawPoint rl.Vector3) *Object {
	object := Object{
		DrawPoint: drawPoint,
		Colliders: colliders,
	}
	return &object

}

func CreateObjectFromCollider(collider types.Collider) *Object {
	object := Object{
		DrawPoint: collider.GetPosition(),
		Colliders: []types.Collider{collider},
	}
	return &object
}
