package entities

import (
	types "github.com/PawelZabc/ProjektZespolowy/client/_types"
	math "github.com/chewxy/math32"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type CylinderCollider struct {
	Position rl.Vector3
	Radius   float32
	Height   float32
}

func (c CylinderCollider) CollidesWith(c2 types.Collider) bool {
	if cylinder, ok := c2.(*CylinderCollider); ok {
		if rl.Vector2Distance(rl.Vector2{X: c.Position.X, Y: c.Position.Z},
			rl.Vector2{X: cylinder.Position.X, Y: cylinder.Position.Z}) < (c.Radius+cylinder.Radius) &&
			c.Position.Y <= cylinder.Position.Y+cylinder.Height && c.Position.Y+c.Height >= cylinder.Position.Y {
			return true
		}

	} else if cube, ok := c2.(*CubeCollider); ok {
		if rl.Vector2Distance(rl.Vector2{X: math.Min(cube.Position.X+cube.SizeX, math.Max(cube.Position.X, c.Position.X)),
			Y: math.Min(cube.Position.Z+cube.SizeZ, math.Max(cube.Position.Z, c.Position.Z))},
			rl.Vector2{X: c.Position.X, Y: c.Position.Z}) <= c.Radius &&
			cube.Position.Y <= c.Position.Y+c.Height && cube.Position.Y+cube.SizeY >= c.Position.Y {
			return true
		}
	}

	return false
}

func (c CylinderCollider) GetPosition() rl.Vector3 {
	return c.Position
}

func (c *CylinderCollider) SetPosition(vec rl.Vector3) {
	c.Position = vec
}

func (c *CylinderCollider) AddPosition(vec rl.Vector3) {
	c.Position = rl.Vector3Add(c.Position, vec)
}

func (c *CylinderCollider) PushbackFrom(c2 types.Collider) types.Direction {
	if cylinder, ok := c2.(*CylinderCollider); ok {
		return c.PushbackFromCylinder(cylinder)
	} else if cube, ok := c2.(*CubeCollider); ok {
		return c.PushbackFromCube(cube)
	} else if plane, ok := c2.(*PlaneCollider); ok {
		return c.PushbackFromPlane(plane)
	}

	return types.DirNone

}

func (c *CylinderCollider) PushbackFromCube(cube *CubeCollider) types.Direction {
	diffrence := rl.Vector2Subtract(rl.Vector2{X: c.Position.X, Y: c.Position.Z},
		rl.Vector2{X: math.Min(cube.Position.X+cube.SizeX, math.Max(cube.Position.X, c.Position.X)),
			Y: math.Min(cube.Position.Z+cube.SizeZ, math.Max(cube.Position.Z, c.Position.Z))})
	distanceXZ := rl.Vector2Length(diffrence) - (c.Radius)
	distanceY1 := c.Position.Y - (cube.Position.Y + cube.SizeY)
	distanceY2 := cube.Position.Y - (c.Position.Y + c.Height)

	if distanceXZ <= 0 && distanceY1 <= 0 && distanceY2 <= 0 {
		if distanceXZ > distanceY1 && distanceXZ > distanceY2 {
			forceXZ := rl.Vector2Scale(rl.Vector2Normalize(diffrence), -distanceXZ)
			c.Position = rl.Vector3Add(c.Position, rl.NewVector3(forceXZ.X, 0, forceXZ.Y))
			return types.DirXZ

		} else if distanceY1 > distanceY2 {
			c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, -distanceY1, 0))
			return -types.DirY
		} else {
			c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, distanceY2, 0))
			return types.DirY
		}
	} else {
		return types.DirNone
	}

}

func (c *CylinderCollider) PushbackFromPlane(plane *PlaneCollider) types.Direction {

	switch plane.Direction {
	case types.DirX, types.DirXminus, types.DirZ, types.DirZminus:
		{
			var diffrence rl.Vector2
			if plane.Direction == types.DirZ || plane.Direction == types.DirZminus {
				diffrence = rl.Vector2Subtract(rl.Vector2{X: c.Position.X, Y: c.Position.Z},
					rl.Vector2{X: math.Min(plane.Position.X+plane.Width, math.Max(plane.Position.X, c.Position.X)),
						Y: plane.Position.Z,
					})
			} else {
				diffrence = rl.Vector2Subtract(rl.Vector2{X: c.Position.X, Y: c.Position.Z},
					rl.Vector2{X: plane.Position.X,
						Y: math.Min(plane.Position.Z+plane.Width, math.Max(plane.Position.Z, c.Position.Z)),
					})
			}
			distanceXZ := rl.Vector2Length(diffrence) - (c.Radius)
			distanceY1 := c.Position.Y - (plane.Position.Y + plane.Height)
			distanceY2 := plane.Position.Y - (c.Position.Y + c.Height)

			if distanceXZ <= 0 && distanceY1 <= 0 && distanceY2 <= 0 {

				forceXZ := rl.Vector2Scale(rl.Vector2Normalize(diffrence), -distanceXZ)
				c.Position = rl.Vector3Add(c.Position, rl.NewVector3(forceXZ.X, 0, forceXZ.Y))
				return -plane.Direction
			}
		}
	case types.DirY, types.DirYminus:
		{
			diffrence := rl.Vector2Subtract(rl.Vector2{X: c.Position.X, Y: c.Position.Z},
				rl.Vector2{X: math.Min(plane.Position.X+plane.Width, math.Max(plane.Position.X, c.Position.X)),
					Y: math.Min(plane.Position.Z+plane.Height, math.Max(plane.Position.Z, c.Position.Z))})
			distanceXZ := rl.Vector2Length(diffrence) - (c.Radius)
			distanceY1 := c.Position.Y - plane.Position.Y
			distanceY2 := plane.Position.Y - (c.Position.Y + c.Height)
			if distanceXZ <= 0 && distanceY1 <= 0 && distanceY2 <= 0 {
				if plane.Direction == -types.DirY {
					c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, distanceY2, 0))
					return types.DirY
				} else {
					c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, -distanceY1, 0))
					return -types.DirY
				}
			}

		}
		// default:
		// 	{
		// 		return types.DirNone
		// 	}
	}
	return types.DirNone

}

func (c *CylinderCollider) PushbackFromCylinder(cylinder *CylinderCollider) types.Direction {
	diffrence := rl.Vector2Subtract(rl.Vector2{X: c.Position.X, Y: c.Position.Z},
		rl.Vector2{X: cylinder.Position.X, Y: cylinder.Position.Z})
	distanceXZ := rl.Vector2Length(diffrence) - (c.Radius + cylinder.Radius)
	distanceY1 := c.Position.Y - (cylinder.Position.Y + cylinder.Height)
	distanceY2 := cylinder.Position.Y - (c.Position.Y + c.Height)

	if distanceXZ <= 0 && distanceY1 <= 0 && distanceY2 <= 0 {
		if distanceXZ > distanceY1 && distanceXZ > distanceY2 {
			forceXZ := rl.Vector2Scale(rl.Vector2Normalize(diffrence), -distanceXZ)
			c.Position = rl.Vector3Add(c.Position, rl.NewVector3(forceXZ.X, 0, forceXZ.Y))
			return types.DirXZ

		} else if distanceY1 > distanceY2 {
			c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, -distanceY1, 0))
			return -types.DirY
		} else {
			c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, distanceY2, 0))
			return types.DirY
		}

	} else {
		return types.DirNone
	}

}
