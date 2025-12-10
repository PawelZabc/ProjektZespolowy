package game

import (
	"github.com/PawelZabc/ProjektZespolowy/client/assets"
	s_types "github.com/PawelZabc/ProjektZespolowy/shared/_types"
	s_entities "github.com/PawelZabc/ProjektZespolowy/shared/entities"
	leveldata "github.com/PawelZabc/ProjektZespolowy/shared/level_data"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadRooms() []Room {
	rooms := make([]Room, 0, 10)
	roomShared := leveldata.Room1
	objects := make([]*Object, 0, len(roomShared.Objects))
	for _, object := range roomShared.Objects {
		objects = append(objects, ConvertObjectSharedToClient(object))
	}
	room := Room{
		Objects: objects,
	}

	rooms = append(rooms, room)

	return rooms
}

type Room struct {
	Objects       []*Object //objects and walls in that room
	SharedObjects []*Object //objects shared with visible rooms
	VisibleRooms  []*Room   //rooms visible that need to be rendered while in the room
}

type Object struct {
	Colliders []s_types.Collider
	DrawPoint rl.Vector3
	Model     rl.Model
	Color     rl.Color
}

func ConvertObjectSharedToClient(object *leveldata.Object) *Object {
	model := rl.Model{}
	drawPoint := object.DrawPoint
	if object.Model != "" {
		model2, _ := assets.GlobalManager.LoadModel(object.Model)
		model = model2.Data
	} else { //if the model is empty use the collider as the model
		model = NewModelFromCollider(object.Colliders[0])
		drawPoint = object.Colliders[0].GetPosition()
	}

	objectClient := Object{
		Colliders: object.Colliders,
		Model:     model,
		// DrawPoint: object.DrawPoint,
		DrawPoint: drawPoint,
		Color:     GetColorFromCollider(object.Colliders[0]),
	}

	return &objectClient

}

func GetColorFromCollider(collider s_types.Collider) rl.Color {
	// if plane, ok := collider.(*s_entities.PlaneCollider); ok { //check if the collider is a plane
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

func DrawRoom(room *Room) {
	DrawObjects(room.Objects)
	DrawObjects(room.SharedObjects)
	for _, room2 := range room.VisibleRooms {
		DrawObjects(room2.Objects)
	}
}

func DrawObjects(objects []*Object) {
	for _, object := range objects {
		rl.DrawModel(object.Model, object.DrawPoint, 1, object.Color)
	}
}

func NewModelFromCollider(collider s_types.Collider) rl.Model {

	switch c := collider.(type) {
	case *s_entities.CubeCollider:
		return NewModelFromCubeCollider(c)
	case *s_entities.CylinderCollider:
		return NewModelFromCylinderCollider(c)
	case *s_entities.PlaneCollider:
		return NewModelFromPlaneCollider(c)
	default:
		return rl.Model{}
	}

}

func NewModelFromCubeCollider(collider *s_entities.CubeCollider) rl.Model {
	modelData, _ := assets.GlobalManager.LoadModel(assets.ModelCube)
	model := modelData.Data
	model.Transform = rl.MatrixScale(collider.SizeX, collider.SizeY, collider.SizeZ)

	return model
}

func NewModelFromCylinderCollider(collider *s_entities.CylinderCollider) rl.Model {
	modelData, _ := assets.GlobalManager.LoadModel(assets.ModelCylinder)
	model := modelData.Data
	model.Transform = rl.MatrixScale(collider.Radius, collider.Height, collider.Radius)

	return model
}

func NewModelFromPlaneCollider(collider *s_entities.PlaneCollider) rl.Model {
	modelData, _ := assets.GlobalManager.LoadModel(assets.ModelCube)
	model := modelData.Data
	switch collider.Direction {
	case s_types.DirX, s_types.DirXminus:
		{
			model.Transform = rl.MatrixScale(0.01, collider.Height, collider.Width)
		}
	case s_types.DirY, s_types.DirYminus:
		{
			model.Transform = rl.MatrixScale(collider.Width, 0.01, collider.Height)
		}
	case s_types.DirZ, s_types.DirZminus:
		{
			model.Transform = rl.MatrixScale(collider.Width, collider.Height, 0.01)
		}
	}

	return model
}
