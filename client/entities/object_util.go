package entities

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/PawelZabc/ProjektZespolowy/client/assets"
)

func LoadModel(name string) rl.Model {
	model, _ := assets.GlobalManager.LoadModel(name + ".glb")

	shader, err := assets.GlobalManager.LoadShader(assets.ShaderLightingVS, assets.ShaderLightingFS)
	if err != nil {
		fmt.Println("Error with loading shader")
	}

	model.Data.Materials.Shader = shader.Data

	lightDirLoc := rl.GetShaderLocation(shader.Data, "lightDir")
	baseColorLoc := rl.GetShaderLocation(shader.Data, "baseColor")
	ambientColorLoc := rl.GetShaderLocation(shader.Data, "ambientColor")

	lightDir := []float32{0.0, -1.0, -1.0}
	rl.SetShaderValue(shader.Data, lightDirLoc, lightDir, rl.ShaderUniformVec3)

	rl.SetShaderValue(shader.Data, baseColorLoc, []float32{1.0, 0.3, 0.2, 1.0}, rl.ShaderUniformVec4)
	rl.SetShaderValue(shader.Data, ambientColorLoc, []float32{0.2, 0.2, 0.2, 1.0}, rl.ShaderUniformVec4)
	return model.Data
}

func CreateCylinderObject(position rl.Vector3, radius float32, height float32) Object {
	object := Object{Collider: &CylinderCollider{
		Position: position,
		Radius:   radius,
		Height:   height,
	}, Model: LoadModel("cylinder"),
	}
	object.Model.Transform = rl.MatrixScale(radius, height, radius)
	return object
}
func CreateCubeObject(position rl.Vector3, sizeX float32, sizeY float32, sizeZ float32) Object {
	object := Object{Collider: &CubeCollider{
		Position: position,
		SizeX:    sizeX,
		SizeY:    sizeY,
		SizeZ:    sizeZ,
	}, Model: LoadModel("cube"),
	}
	object.Model.Transform = rl.MatrixScale(sizeX, sizeY, sizeZ)
	return object
}
