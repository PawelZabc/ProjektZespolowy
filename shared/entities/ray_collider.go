package entities

import (
	types "github.com/PawelZabc/ProjektZespolowy/shared/_types"
	math "github.com/chewxy/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Ray struct {
	Origin    rl.Vector3
	Direction rl.Vector3
}

func (r *Ray) GetCollisionPoint(collider types.Collider) (*rl.Vector3, float32) {
	if plane, ok := collider.(*PlaneCollider); ok {
		return r.GetCollisionPointWithPlane(*plane)
	} else if cube, ok := collider.(*CubeCollider); ok {
		return r.GetCollisionPointWithCube(*cube)
	} else if cylinder, ok := collider.(*CylinderCollider); ok {
		return r.GetCollisionPointWithCylinder(*cylinder)
	}

	return nil, 0
}

func (r *Ray) GetCollisionPointWithCube(cube CubeCollider) (*rl.Vector3, float32) {
	var points = []*rl.Vector3{nil, nil, nil}
	var lengths = []float32{0, 0, 0}
	//X---------------------------------------------------------------------------
	if r.Direction.X != 0 { // if the direction X is 0, the it will never interact with the X plane, also removes the risk of dividing by 0
		distanceX := cube.GetPosition().X - r.Origin.X
		distanceX2 := cube.GetPosition().X + cube.SizeX - r.Origin.X
		if math.Abs(distanceX) > math.Abs(distanceX2) {
			distanceX = distanceX2
		} //only check for wall closer to the origin
		stepsX := distanceX / r.Direction.X //calculate scalar steps to get intersecting point with the X plane
		if stepsX > 0 {                     //if steps is negative the direction is opposite to the plane and there wont be an intersecting point
			point := rl.Vector3Scale(r.Direction, stepsX) //get the instersetion point with the X plane
			lengths[0] = rl.Vector3Length(point)
			point = rl.Vector3Add(point, r.Origin)
			points[0] = &point
			pointX := point
			if !(pointX.Z >= cube.Position.Z && pointX.Z <= cube.Position.Z+cube.SizeZ && pointX.Y >= cube.Position.Y && pointX.Y <= cube.Position.Y+cube.SizeY) {
				points[0] = nil // check if point is inside the bounds of the wall, if not set it to nil to invalidate it
			}
		}
	}
	//Y---------------------------------------------------------------------------
	if r.Direction.Y != 0 {
		distanceY := cube.GetPosition().Y - r.Origin.Y
		distanceY2 := cube.GetPosition().Y + cube.SizeY - r.Origin.Y
		if math.Abs(distanceY) > math.Abs(distanceY2) {
			distanceY = distanceY2
		}
		stepsY := distanceY / r.Direction.Y
		if stepsY > 0 {
			point := rl.Vector3Scale(r.Direction, stepsY)
			lengths[1] = rl.Vector3Length(point)
			point = rl.Vector3Add(point, r.Origin)
			points[1] = &point
			pointY := point
			if !(pointY.X >= cube.Position.X && pointY.X <= cube.Position.X+cube.SizeX && pointY.Z >= cube.Position.Z && pointY.Z <= cube.Position.Z+cube.SizeZ) {
				points[1] = nil
			}
		}
	}
	//Z---------------------------------------------------------------------------
	if r.Direction.Z != 0 {
		distanceZ := cube.GetPosition().Z - r.Origin.Z
		distanceZ2 := cube.GetPosition().Z + cube.SizeZ - r.Origin.Z
		if math.Abs(distanceZ) > math.Abs(distanceZ2) {
			distanceZ = distanceZ2
		}
		stepsZ := distanceZ / r.Direction.Z
		if stepsZ > 0 {
			point := rl.Vector3Scale(r.Direction, stepsZ)
			lengths[2] = rl.Vector3Length(point)
			point = rl.Vector3Add(point, r.Origin)
			points[2] = &point
			pointZ := point
			if !(pointZ.X >= cube.Position.X && pointZ.X <= cube.Position.X+cube.SizeX && pointZ.Y >= cube.Position.Y && pointZ.Y <= cube.Position.Y+cube.SizeY) {
				points[2] = nil
			}
		}
	}
	bestLength := float32(0)
	var bestPoint *rl.Vector3 = nil
	for i, el := range points {
		if el != nil {
			if bestLength == 0 || lengths[i] < bestLength {
				bestLength = lengths[i]
				bestPoint = points[i]
			}

		}
	}
	return bestPoint, bestLength
}

func (r *Ray) GetCollisionPointWithCylinder(cylinder CylinderCollider) (*rl.Vector3, float32) {
	distanceY := cylinder.GetPosition().Y - r.Origin.Y
	distanceY2 := cylinder.GetPosition().Y + cylinder.Height - r.Origin.Y
	if math.Abs(distanceY) > math.Abs(distanceY2) {
		distanceY = distanceY2
	}
	stepsY := distanceY / r.Direction.Y
	if stepsY > 0 {
		point := rl.Vector3Scale(r.Direction, stepsY)
		length := rl.Vector3Length(point)
		point = rl.Vector3Add(point, r.Origin)
		if rl.Vector2Length(rl.Vector2Subtract(GetVector2DXZ(point), GetVector2DXZ(cylinder.Position))) <= cylinder.Radius {
			return &point, length
		}
	}
	rayStart := GetVector2DXZ(r.Origin) //source for math: https://youtu.be/ebzlMOw79Yw?si=I1rRq7fPx9mPEyjk
	circleCenter := GetVector2DXZ(cylinder.Position)
	s := rl.Vector2Subtract(rayStart, circleCenter)
	a := rl.Vector2DotProduct(GetVector2DXZ(r.Direction), GetVector2DXZ(r.Direction))
	b := rl.Vector2DotProduct(s, GetVector2DXZ(r.Direction))
	c := rl.Vector2DotProduct(s, s) - (cylinder.Radius * cylinder.Radius)
	h := b*b - (a * c)
	if h < 0 {
		return nil, 0
	}
	h = math.Sqrt(h)
	t := (-b - h) / a
	if t < 0 {
		return nil, 0
	}
	point := rl.Vector3Scale(r.Direction, t)
	point = rl.Vector3Add(point, r.Origin)
	if point.Y >= cylinder.Position.Y && point.Y <= (cylinder.Position.Y+cylinder.Height) {
		return &point, t
	}
	return nil, 0
}

func (r *Ray) GetCollisionPointWithPlane(plane PlaneCollider) (*rl.Vector3, float32) {
	switch plane.Direction {
	case types.DirZ:
		{
			if r.Direction.Z == 0 {
				return nil, 0
			}
			distanceZ := plane.GetPosition().Z - r.Origin.Z
			steps := distanceZ / r.Direction.Z
			if steps < 0 {
				return nil, 0
			} else {
				point := rl.Vector3Scale(r.Direction, steps)
				length := rl.Vector3Length(point)
				point = rl.Vector3Add(point, r.Origin)

				if point.X >= plane.Position.X && point.X <= plane.Position.X+plane.Width && point.Y >= plane.Position.Y && point.Y <= plane.Position.Y+plane.Height {
					return &point, length
				}
			}
		}
	case types.DirX:
		{
			if r.Direction.X == 0 {
				return nil, 0
			}
			distanceX := plane.GetPosition().X - r.Origin.X
			steps := distanceX / r.Direction.X
			if steps < 0 {
				return nil, 0
			} else {
				point := rl.Vector3Scale(r.Direction, steps)
				length := rl.Vector3Length(point)
				point = rl.Vector3Add(point, r.Origin)

				if point.Z >= plane.Position.Z && point.Z <= plane.Position.Z+plane.Width && point.Y >= plane.Position.Y && point.Y <= plane.Position.Y+plane.Height {
					return &point, length
				}
			}
		}
	case types.DirY, types.DirYminus:
		{
			if r.Direction.Y == 0 {
				return nil, 0
			}
			distanceY := plane.GetPosition().Y - r.Origin.Y
			steps := distanceY / r.Direction.Y
			if steps < 0 {
				return nil, 0
			} else {
				point := rl.Vector3Scale(r.Direction, steps)
				length := rl.Vector3Length(point)
				point = rl.Vector3Add(point, r.Origin)

				if point.X >= plane.Position.X && point.X <= plane.Position.X+plane.Width && point.Z >= plane.Position.Z && point.Z <= plane.Position.Z+plane.Height {
					return &point, length
				}
			}
		}
	}
	return nil, 0
}
