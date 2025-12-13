package physics

import (
	math "github.com/chewxy/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func GetVector2XZ(vec rl.Vector3) rl.Vector2 {
	return rl.NewVector2(vec.X, vec.Z)
}

func GetVector3FromXZ(vec rl.Vector2) rl.Vector3 {
	return rl.NewVector3(vec.X, 0, vec.Y)
}

func GetRotationX(vector rl.Vector2) float32 {
	return rl.Vector2Angle(rl.NewVector2(1, 0), vector)*rl.Rad2deg + 180
}

func GetCloserWallOnSameAxis(origin, position1, position2 float32) float32 {
	distance1 := position1 - origin
	distance2 := position2 - origin
	if math.Abs(distance1) < math.Abs(distance2) {
		return distance1
	} else {
		return distance2
	}
}
