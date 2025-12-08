package entities

import rl "github.com/gen2brain/raylib-go/raylib"

func GetVector2XZ(vec rl.Vector3) rl.Vector2 {
	return rl.NewVector2(vec.X, vec.Z)
}

func GetVector3FromXZ(vec rl.Vector2) rl.Vector3 {
	return rl.NewVector3(vec.X, 0, vec.Y)
}
