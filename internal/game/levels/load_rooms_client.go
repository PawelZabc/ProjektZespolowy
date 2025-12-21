package levels

// CLIENT LOAD ROOMS

import (
	"github.com/PawelZabc/ProjektZespolowy/assets"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/entities"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	"github.com/chewxy/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadRooms() []ClientRoom {
	rooms := make([]ClientRoom, 0, 10)
	roomShared := Room1
	objects := make([]*entities.Object, 0, len(roomShared.Objects))
	for _, object := range roomShared.Objects {
		objects = append(objects, ConvertObjectSharedToClient(object))
	}
	room := ClientRoom{
		Objects: objects,
	}

	rooms = append(rooms, room)

	return rooms
}

func ConvertObjectSharedToClient(object *ObjectTWO) *entities.Object {
	model := rl.Model{}
	drawPoint := object.DrawPoint
	if object.Model != "" {
		model2, _ := assets.GlobalManager.LoadModel(object.Model)
		model = model2.Data
	} else { //if the model is empty use the collider as the model
		model = NewModelFromCollider(object.Colliders[0])
		drawPoint = object.Colliders[0].GetPosition()
	}

	objectClient := entities.Object{
		Colliders: object.Colliders,
		Model:     model,
		// DrawPoint: object.DrawPoint,
		DrawPoint: drawPoint,
		Color:     GetColorFromCollider(object.Colliders[0]),
	}

	return &objectClient

}

func GetColorFromCollider(collider colliders.Collider) rl.Color {
	// if plane, ok := collider.(*colider.PlaneCollider); ok { //check if the collider is a plane
	// 	switch plane.Direction { //check which color to draw the plane as
	// 	case types.DirX:
	// 		{
	// 			return rl.Red
	// 		}
	// 	case types.DirY:
	// 		{
	// 			return rl.Orange
	// 		}
	// 	case types.DirYminus:
	// 		{
	// 			return rl.Green
	// 		}
	// 	case types.DirZ:
	// 		{
	// 			return rl.Yellow
	// 		}
	// 	}
	// }
	return rl.White //if its not a plane color white
}

func DrawRoom(room *ClientRoom) {
	DrawObjects(room.Objects)
	DrawObjects(room.SharedObjects)
	for _, room2 := range room.VisibleRooms {
		DrawObjects(room2.Objects)
	}
}

func DrawObjects(objects []*entities.Object) {
	for _, object := range objects {
		object.Draw()
	}
}

func NewModelFromCollider(collider colliders.Collider) rl.Model {

	switch c := collider.(type) {
	case *colliders.CubeCollider:
		return NewModelFromCubeCollider(c)
	case *colliders.CylinderCollider:
		return NewModelFromCylinderCollider(c)
	case *colliders.PlaneCollider:
		return NewModelFromPlaneCollider(c)
	default:
		return rl.Model{}
	}

}

func NewModelFromCubeCollider(collider *colliders.CubeCollider) rl.Model {
	modelData, _ := assets.GlobalManager.LoadModel(assets.ModelCube)
	model := modelData.Data
	model.Transform = rl.MatrixScale(collider.SizeX, collider.SizeY, collider.SizeZ)

	return model
}

func NewModelFromCylinderCollider(collider *colliders.CylinderCollider) rl.Model {
	modelData, _ := assets.GlobalManager.LoadModel(assets.ModelCylinder)
	model := modelData.Data
	model.Transform = rl.MatrixScale(collider.Radius, collider.Height, collider.Radius)

	return model
}

func NewModelFromPlaneCollider(collider *colliders.PlaneCollider) rl.Model {
	modelData, _ := assets.GlobalManager.LoadModel(assets.ModelCube)
	model := modelData.Data
	switch collider.Direction {
	case physics.DirX, physics.DirXminus:
		{
			model.Transform = rl.MatrixScale(0.01, collider.Height, collider.Width)
		}
	case physics.DirY, physics.DirYminus:
		{
			model.Transform = rl.MatrixScale(collider.Width, 0.01, collider.Height)
		}
	case physics.DirZ, physics.DirZminus:
		{
			model.Transform = rl.MatrixScale(collider.Width, collider.Height, 0.01)
		}
	}

	return model
}

func CreateRoomWallsFromChanges(StartPoint rl.Vector3, Changes []Change, Height float32) []colliders.Collider {

	count := 0
	for _, change := range Changes {
		if !change.Skip {
			count++
		}
	}
	walls := make([]colliders.Collider, len(Changes))
	skipped := 0
	for i, change := range Changes {
		if change.Axis == physics.DirX {
			change.Axis = physics.DirZ
		} else {
			change.Axis = physics.DirX
		}

		var object colliders.PlaneCollider
		if change.Value < 0 {
			if change.Axis == physics.DirX {
				StartPoint = rl.Vector3Add(StartPoint, rl.NewVector3(0, 0, change.Value))
			} else {
				StartPoint = rl.Vector3Add(StartPoint, rl.NewVector3(change.Value, 0, 0))
			}
		}
		if !change.Skip {
			object = *colliders.NewPlaneCollider(StartPoint, math32.Abs(change.Value), Height, change.Axis)
			walls[i-skipped] = &object
		} else {
			skipped += 1
		}
		if change.Value > 0 {
			if change.Axis == physics.DirX {
				StartPoint = rl.Vector3Add(StartPoint, rl.NewVector3(0, 0, change.Value))
			} else {
				StartPoint = rl.Vector3Add(StartPoint, rl.NewVector3(change.Value, 0, 0))
			}
		}

	}

	return walls
}
