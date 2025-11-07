package main

import rl "github.com/gen2brain/raylib-go/raylib"

func LoadModel(name string) rl.Model {
	model := rl.LoadModel("assets/" + name + ".glb")

	shader := rl.LoadShader("lighting.vs", "lighting.fs")

	model.Materials.Shader = shader

	lightDirLoc := rl.GetShaderLocation(shader, "lightDir")
	baseColorLoc := rl.GetShaderLocation(shader, "baseColor")
	ambientColorLoc := rl.GetShaderLocation(shader, "ambientColor")

	lightDir := []float32{0.0, -1.0, -1.0}
	rl.SetShaderValue(shader, lightDirLoc, lightDir, rl.ShaderUniformVec3)

	rl.SetShaderValue(shader, baseColorLoc, []float32{1.0, 0.3, 0.2, 1.0}, rl.ShaderUniformVec4)
	rl.SetShaderValue(shader, ambientColorLoc, []float32{0.2, 0.2, 0.2, 1.0}, rl.ShaderUniformVec4)
	return model
}

func createCylinderObject(position rl.Vector3, radius float32, height float32) Object {
	object := Object{Collider: &CylinderCollider{
		Position: position,
		Radius:   radius,
		Height:   height,
	}, Model: LoadModel("cylinder"),
	}
	object.Model.Transform = rl.MatrixScale(radius, height, radius)
	return object
}
func createCubeObject(position rl.Vector3, sizeX float32, sizeY float32, sizeZ float32) Object {
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
