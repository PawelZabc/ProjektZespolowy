package colliders

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics"
	"github.com/PawelZabc/ProjektZespolowy/internal/shared"
	math "github.com/chewxy/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Ray struct {
	Origin    rl.Vector3
	Direction rl.Vector3
}

func (r *Ray) GetCollisionPoint(collider Collider) (*rl.Vector3, float32) {
	if plane, ok := collider.(*PlaneCollider); ok {
		return r.GetCollisionPointWithPlane(*plane)
	} else if cube, ok := collider.(*CubeCollider); ok {
		return r.GetCollisionPointWithCube(*cube)
	} else if cylinder, ok := collider.(*CylinderCollider); ok {
		return r.GetCollisionPointWithCylinder(*cylinder)
	}

	return nil, 0
}

func (r *Ray) GetRotationX() float32 {
	return GetRotationX(GetVector2XZ(r.Direction))
}

func (r *Ray) GetCollisionPointWithAxis(distance float32, plane types.Direction) (*rl.Vector3, float32) {
	direction := float32(0)
	switch plane {
	case types.DirX:
		direction = r.Direction.X
	case types.DirY, types.DirYminus:
		direction = r.Direction.Y
	case types.DirZ:
		direction = r.Direction.Z
	}

	if direction != 0 { // if the direction X is 0, the it will never interact with the X plane, also removes the risk of dividing by 0//only check for wall closer to the origin
		steps := distance / direction //calculate scalar steps to get intersecting point with the X plane
		if steps > 0 {                //if steps is negative the direction is opposite to the plane and there wont be an intersecting point
			point := rl.Vector3Scale(r.Direction, steps) //get the instersetion point with the X plane
			length := rl.Vector3Length(point)
			point = rl.Vector3Add(point, r.Origin)
			return &point, length
		}
	}
	return nil, 0

}

func (r *Ray) GetCollisionPointWithCube(cube CubeCollider) (*rl.Vector3, float32) {
	var points = []*rl.Vector3{nil, nil, nil}
	var lengths = []float32{0, 0, 0}
	//X---------------------------------------------------------------------------
	points[0], lengths[0] = r.GetCollisionPointWithAxis(GetCloserWallOnSameAxis(r.Origin.X, cube.Position.X, cube.Position.X+cube.SizeX), types.DirX)
	if !(points[0].Z >= cube.Position.Z && points[0].Z <= cube.Position.Z+cube.SizeZ && points[0].Y >= cube.Position.Y && points[0].Y <= cube.Position.Y+cube.SizeY) {
		points[0] = nil // check if point is inside the bounds of the wall, if not set it to nil to invalidate it
	}
	//Y---------------------------------------------------------------------------
	points[1], lengths[1] = r.GetCollisionPointWithAxis(GetCloserWallOnSameAxis(r.Origin.Y, cube.Position.Y, cube.Position.Y+cube.SizeY), types.DirY)
	if !(points[1].Z >= cube.Position.Z && points[1].Z <= cube.Position.Z+cube.SizeZ && points[1].X >= cube.Position.X && points[1].X <= cube.Position.X+cube.SizeX) {
		points[1] = nil
	}
	//Z---------------------------------------------------------------------------
	points[2], lengths[2] = r.GetCollisionPointWithAxis(GetCloserWallOnSameAxis(r.Origin.Z, cube.Position.Z, cube.Position.Z+cube.SizeZ), types.DirZ)
	if !(points[2].X >= cube.Position.X && points[2].X <= cube.Position.X+cube.SizeX && points[2].Y >= cube.Position.Y && points[2].Y <= cube.Position.Y+cube.SizeY) {
		points[2] = nil
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
	rayStart := GetVector2XZ(r.Origin) //source for math: https://youtu.be/ebzlMOw79Yw?si=I1rRq7fPx9mPEyjk
	circleCenter := GetVector2XZ(cylinder.Position)
	s := rl.Vector2Subtract(rayStart, circleCenter)
	a := rl.Vector2DotProduct(GetVector2XZ(r.Direction), GetVector2XZ(r.Direction))
	b := rl.Vector2DotProduct(s, GetVector2XZ(r.Direction))
	c := rl.Vector2DotProduct(s, s) - (cylinder.Radius * cylinder.Radius)
	h := b*b - (a * c)
	if h >= 0 {
		h = math.Sqrt(h)
		t := (-b - h) / a
		if t >= 0 {
			point := rl.Vector3Scale(r.Direction, t)
			point = rl.Vector3Add(point, r.Origin)
			if point.Y >= cylinder.Position.Y && point.Y <= (cylinder.Position.Y+cylinder.Height) {
				return &point, t
			}
		}
	}
	point, length := r.GetCollisionPointWithAxis(GetCloserWallOnSameAxis(r.Origin.Y, cylinder.Position.Y, cylinder.Position.Y+cylinder.Height), types.DirY)
	if point != nil && rl.Vector2Length(rl.Vector2Subtract(GetVector2XZ(*point), GetVector2XZ(cylinder.Position))) <= cylinder.Radius {
		return point, length
	}
	return nil, 0
}

func (r *Ray) GetCollisionPointWithPlane(plane PlaneCollider) (*rl.Vector3, float32) {
	switch plane.Direction {
	case types.DirX:
		point, length := r.GetCollisionPointWithAxis(plane.GetPosition().X-r.Origin.X, plane.Direction)
		if point != nil && point.Z >= plane.Position.Z && point.Z <= plane.Position.Z+plane.Width && point.Y >= plane.Position.Y && point.Y <= plane.Position.Y+plane.Height {
			return point, length
		}
	case types.DirY, types.DirYminus:
		point, length := r.GetCollisionPointWithAxis(plane.GetPosition().Y-r.Origin.Y, plane.Direction)
		if point != nil && point.X >= plane.Position.X && point.X <= plane.Position.X+plane.Width && point.Z >= plane.Position.Z && point.Z <= plane.Position.Z+plane.Height {
			return point, length
		}
	case types.DirZ:
		point, length := r.GetCollisionPointWithAxis(plane.GetPosition().Z-r.Origin.Z, plane.Direction)
		if point != nil && point.X >= plane.Position.X && point.X <= plane.Position.X+plane.Width && point.Y >= plane.Position.Y && point.Y <= plane.Position.Y+plane.Height {
			return point, length
		}
	}
	return nil, 0
}
