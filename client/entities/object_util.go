package entities

import (
	"fmt"

	math "github.com/chewxy/math32"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/PawelZabc/ProjektZespolowy/client/assets"
	s_types "github.com/PawelZabc/ProjektZespolowy/shared/_types"
	s_entities "github.com/PawelZabc/ProjektZespolowy/shared/entities"
)

// func LoadModel(name string) rl.Model {
// 	model, _ := assets.GlobalManager.LoadModel(name + ".glb")

// 	// shader, err := assets.GlobalManager.LoadShader(assets.ShaderLightingVS, assets.ShaderLightingFS)
// 	// if err != nil {
// 	// 	fmt.Println("Error with loading shader")
// 	// }

// 	// model.Data.Materials.Shader = shader.Data

// 	// lightDirLoc := rl.GetShaderLocation(shader.Data, "lightDir")
// 	// baseColorLoc := rl.GetShaderLocation(shader.Data, "baseColor")
// 	// ambientColorLoc := rl.GetShaderLocation(shader.Data, "ambientColor")

// 	// lightDir := []float32{0.0, -1.0, -1.0}
// 	// rl.SetShaderValue(shader.Data, lightDirLoc, lightDir, rl.ShaderUniformVec3)

// 	// rl.SetShaderValue(shader.Data, baseColorLoc, []float32{1.0, 0.3, 0.2, 1.0}, rl.ShaderUniformVec4)
// 	// rl.SetShaderValue(shader.Data, ambientColorLoc, []float32{0.2, 0.2, 0.2, 1.0}, rl.ShaderUniformVec4)
// 	return model.Data
// }

func CreateCylinderObject(position rl.Vector3, radius float32, height float32) Object {
	model, _ := assets.GlobalManager.LoadModel(assets.ModelCylinder)
	object := Object{Collider: &s_entities.CylinderCollider{
		Position: position,
		Radius:   radius,
		Height:   height,
	}, Model: model.Data,
	}
	object.Model.Transform = rl.MatrixScale(radius, height, radius)
	return object
}
func CreateCubeObject(position rl.Vector3, sizeX float32, sizeY float32, sizeZ float32) Object {
	model, _ := assets.GlobalManager.LoadModel(assets.ModelCube)
	object := Object{Collider: &s_entities.CubeCollider{
		Position: position,
		SizeX:    sizeX,
		SizeY:    sizeY,
		SizeZ:    sizeZ,
	}, Model: model.Data,
	}
	object.Model.Transform = rl.MatrixScale(sizeX, sizeY, sizeZ)
	return object
}
func CreatePlaneObject(position rl.Vector3, Width float32, Height float32, Direction s_types.Direction) Object {
	model, _ := assets.GlobalManager.LoadModel(assets.ModelCube)
	object := Object{Collider: &s_entities.PlaneCollider{
		Position:  position,
		Width:     Width,
		Height:    Height,
		Direction: Direction,
	}, Model: model.Data,
	}
	switch Direction {
	case s_types.DirX, s_types.DirXminus:
		{
			object.Model.Transform = rl.MatrixScale(0.01, Height, Width)
		}
	case s_types.DirY, s_types.DirYminus:
		{
			object.Model.Transform = rl.MatrixScale(Width, 0.01, Height)
		}
	case s_types.DirZ, s_types.DirZminus:
		{
			object.Model.Transform = rl.MatrixScale(Width, Height, 0.01)
		}
	}
	// object.Model.Transform = rl.MatrixScale(Width, 0.01, Height)
	return object
}

func CreateRoomWallsFromPoint(Points []rl.Vector2, StartHeight float32, Height float32) []*Object {
	if len(Points) < 2 {
		return nil
	}
	walls := make([]*Object, len(Points)-1)
	for i := 1; i < len(Points); i++ {
		point1 := Points[i-1]
		point2 := Points[i]
		diffrence := rl.Vector2Subtract(point2, point1)
		direction := s_types.DirNone
		Width := float32(0)
		if diffrence.X != 0 {
			direction = s_types.DirZ

			Width = diffrence.X
			if diffrence.X < 0 {
				point1 = point2
				Width = -diffrence.X
			}
		} else {
			direction = s_types.DirX
			Width = diffrence.Y
			if diffrence.Y < 0 {
				point1 = point2
				Width = -diffrence.Y
			}
		}
		// fmt.Println(point1, diffrence)
		position := rl.NewVector3(point1.X, StartHeight, point1.Y)
		object := CreatePlaneObject(position, Width, Height, direction)
		walls[i-1] = &object
		// walls[i-1] = CreateWall(Points[i-1], Points[i])
	}
	fmt.Println(walls)
	return walls

}

type Change struct {
	Value float32
	Axis  s_types.Direction
	Skip  bool
}

func CreateRoomWallsFromChanges(StartPoint rl.Vector3, Changes []Change, Height float32) []*Object {

	count := 0
	for _, change := range Changes {
		if !change.Skip {
			count++
		}
	}
	walls := make([]*Object, len(Changes))
	skipped := 0
	for i, change := range Changes {
		if change.Axis == s_types.DirX {
			change.Axis = s_types.DirZ
		} else {
			change.Axis = s_types.DirX
		}

		var object Object
		if change.Value < 0 {
			if change.Axis == s_types.DirX {
				StartPoint = rl.Vector3Add(StartPoint, rl.NewVector3(0, 0, change.Value))
			} else {
				StartPoint = rl.Vector3Add(StartPoint, rl.NewVector3(change.Value, 0, 0))
			}
		}
		if !change.Skip {
			object = CreatePlaneObject(StartPoint, math.Abs(change.Value), Height, change.Axis)
			walls[i-skipped] = &object
		} else {
			skipped += 1
		}
		if change.Value > 0 {
			if change.Axis == s_types.DirX {
				StartPoint = rl.Vector3Add(StartPoint, rl.NewVector3(0, 0, change.Value))
			} else {
				StartPoint = rl.Vector3Add(StartPoint, rl.NewVector3(change.Value, 0, 0))
			}
		}
		fmt.Println(object.Collider)

	}

	return walls
}
